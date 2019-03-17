package custommiddleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	crypto "github.com/SermoDigital/jose/crypto"
	jws "github.com/SermoDigital/jose/jws"
	jwt "github.com/SermoDigital/jose/jwt"
	"github.com/dr-sungate/google-oauth-gateway/api/repository/entity"
	"github.com/dr-sungate/google-oauth-gateway/api/service/client"
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"github.com/dr-sungate/google-oauth-gateway/api/service/parser"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const HTTP_REQUEST_TIMEOUT = 300
const OPENID_CONTXTKEY = "openidclaims"
const (
	PUBLICKEY_ENDPOINT = "https://www.googleapis.com/oauth2/v3/certs"
	PUBLICKEY_METHOD   = "GET"
	PUBLICKEY_TTL      = 5
)
const PUBLICKEY_CACHE_PREFIX = "publickey:::"

type (
	OAuth2Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// BeforeFunc defines a function which is executed just before the middleware.
		BeforeFunc middleware.BeforeFunc

		SuccessHandler OpenIDSuccessHandler
		ErrorHandler   OpenIDErrorHandler

		PublicKeyEndPoint string
		PublicKeyTtl      int
		PublicKey         []*rsa.PublicKey

		SigningMethod *crypto.SigningMethodRSA

		ContextKey string

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		GoCacheClient *client.GoCacheClient
	}

	OpenIDSuccessHandler func(echo.Context)

	OpenIDErrorHandler func(error) error

	tokenExtractor func(echo.Context) (string, error)
)

var (
	ErrTokenMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed token")
)
var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJOAuth2Config = OAuth2Config{
		Skipper:           middleware.DefaultSkipper,
		SigningMethod:     crypto.SigningMethodRS256,
		ContextKey:        OPENID_CONTXTKEY,
		PublicKeyEndPoint: PUBLICKEY_ENDPOINT,
		PublicKeyTtl:      PUBLICKEY_TTL,
		TokenLookup:       "header:" + echo.HeaderAuthorization,
		AuthScheme:        "Bearer",
	}
)

func OAuth2() echo.MiddlewareFunc {
	c := DefaultJOAuth2Config
	return OAuth2WithConfig(c)
}

func OAuth2WithConfig(config OAuth2Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultJOAuth2Config.Skipper
	}
	if config.SigningMethod == nil {
		config.SigningMethod = DefaultJOAuth2Config.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJOAuth2Config.ContextKey
	}
	if config.PublicKeyEndPoint == "" {
		config.PublicKeyEndPoint = DefaultJOAuth2Config.PublicKeyEndPoint
	}
	if config.PublicKeyTtl == 0 {
		config.PublicKeyTtl = DefaultJOAuth2Config.PublicKeyTtl
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJOAuth2Config.TokenLookup
	}
	if config.AuthScheme == " " || config.AuthScheme == "-" {
		config.AuthScheme = ""
	} else if config.AuthScheme == "" {
		config.AuthScheme = DefaultJOAuth2Config.AuthScheme
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := tokenFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = tokenFromQuery(parts[1])
	case "cookie":
		extractor = tokenFromCookie(parts[1])
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}
			if encryptKeys, err := getPublicKey(c, config); err == nil {
				for _, encryptKey := range encryptKeys {
					config.PublicKey = append(config.PublicKey, encryptKey.PublicKey)
				}
			} else {
				log.Error(err)
			}

			auth, err := extractor(c)
			if err != nil {
				log.Error(err)
				return responseError(err, config)
			}
			token, err := jws.ParseJWT([]byte(auth))
			if err != nil {
				log.Error(err)
				return responseError(err, config)
			}
			var validateerr error
			for _, publickey := range config.PublicKey {
				if validateerr = token.Validate(publickey, config.SigningMethod); validateerr == nil {
					// Store user information from token into context.
					c.Set(config.ContextKey, token.Claims())
					log.Info(token.Claims())
					if config.SuccessHandler != nil {
						config.SuccessHandler(c)
					}
					return next(c)
				} else {
					log.Warn(validateerr)
				}
			}
			return responseError(validateerr, config)
		}
	}
}

func responseError(err error, config OAuth2Config) error {
	if config.ErrorHandler != nil {
		return config.ErrorHandler(err)
	}
	return &echo.HTTPError{
		Code:     http.StatusUnauthorized,
		Message:  "invalid or expired openid",
		Internal: err,
	}

}
func tokenFromHeader(header string, authScheme string) tokenExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrTokenMissing
	}
}
func tokenFromQuery(param string) tokenExtractor {
	return func(c echo.Context) (string, error) {
		token := c.QueryParam(param)
		if token == "" {
			return "", ErrTokenMissing
		}
		return token, nil
	}
}

func tokenFromCookie(name string) tokenExtractor {
	return func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", ErrTokenMissing
		}
		return cookie.Value, nil
	}
}

func getPublicKey(c echo.Context, config OAuth2Config) ([]entity.EncryptKey, error) {
	encryptkeys := make([]entity.EncryptKey, 0)
	parsedurl, err := url.Parse(config.PublicKeyEndPoint)
	if err != nil {
		return encryptkeys, err
	}
	if cacheddata := getGoCache(config.GoCacheClient, parsedurl); cacheddata != nil {
		log.Info("PublicKey is Get From Cache")
		encryptkeys = cacheddata.([]entity.EncryptKey)
		return encryptkeys, nil
	}
	responsedata, err := requestGet(parsedurl.String())
	if err != nil {
		log.Error(err)
		return encryptkeys, err
	}
	log.Debug(string(responsedata))
	keys, err := parser.GetJsonItems(responsedata, "/keys")
	if err != nil {
		log.Error(err)
		return encryptkeys, err
	}
	switch convertedkeys := keys.(type) {
	case []interface{}:
		var keydata []byte
		var err error
		for _, convertedkey := range convertedkeys {
			if keydata, err = parser.JsonToByte(convertedkey); err != nil {
				log.Error(err)
				return encryptkeys, err
			}
			jwk, err := parser.ConvertJSONWebKey(keydata)
			if err != nil {
				log.Error(err)
				return encryptkeys, err
			}
			publickey, err := parser.ToRSAPublic(&jwk)
			if err != nil {
				log.Error(err)
				return encryptkeys, err
			}
			encryptkey := entity.EncryptKey{}
			encryptkey.PublicKey = publickey
			encryptkeys = append(encryptkeys, encryptkey)
		}
		setGoCache(config.GoCacheClient, parsedurl, encryptkeys, config.PublicKeyTtl)
	default:
		err = errors.New("No JWK Data")
		return encryptkeys, err

	}
	return encryptkeys, nil
}

func requestGet(requesturi string) ([]byte, error) {
	log.Info(fmt.Sprintf("Request URI: %s", requesturi))
	resp, err := http.Get(requesturi)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	return ioutil.ReadAll(resp.Body)
}

func getGoCache(gocache *client.GoCacheClient, parsedurl *url.URL) interface{} {
	if gocache == nil {
		log.Info("No GoCache Client")
		return nil
	}
	if value, exists := gocache.Get(PUBLICKEY_CACHE_PREFIX + parsedurl.String()); exists {
		return value
	}
	return nil
}

func setGoCache(gocache *client.GoCacheClient, parsedurl *url.URL, publickey interface{}, ttl int) {
	if gocache == nil {
		log.Info("No GoCache Client")
		return
	}
	gocache.Set(PUBLICKEY_CACHE_PREFIX+parsedurl.String(), publickey, ttl)
}

func GetOpenIDUserIDKeyFromContext(c echo.Context, useridkey string) string {
	claims := c.Get(OPENID_CONTXTKEY)
	if useridkey == "" {
		useridkey = "sub"
	}
	switch claimsconverted := claims.(type) {
	case jwt.Claims:
		return claimsconverted.Get(useridkey).(string)
	default:
		return ""
	}

}
