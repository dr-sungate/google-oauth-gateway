package custommiddleware

import (
	"crypto/rsa"
	crypto "github.com/SermoDigital/jose/crypto"
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

const OPENID_CONTXTKEY = "openidclaims"

type (
	OAuth2Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// BeforeFunc defines a function which is executed just before the middleware.
		BeforeFunc middleware.BeforeFunc

		SuccessHandler OpenIDSuccessHandler
		ErrorHandler   OpenIDErrorHandler

		PublicKey *rsa.PublicKey

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
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: crypto.SigningMethodRS256,
		ContextKey:    OPENID_CONTXTKEY,
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
	}
)

func OAuth2(key interface{}) echo.MiddlewareFunc {
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
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJOAuth2Config.TokenLookup
	}
	if config.AuthScheme == "" && config.AuthScheme != " " {
		config.AuthScheme = DefaultJOAuth2Config.AuthScheme
	} else if config.AuthScheme != " " {
		config.AuthScheme = ""
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := tokenFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = tokenFromQuery(parts[1])
	case "cookie":
		extractor = tokenFromCookie(parts[1])
	}
	log.Info(extractor)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}
			return next(c)
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
