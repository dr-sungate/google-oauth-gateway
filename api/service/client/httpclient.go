package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
)

var version = "1.0"
var userAgent = fmt.Sprintf("Gateway GoClient/%s (%s)", version, runtime.Version())

const HTTP_STATUS_OK int = 200
const HTTP_DEFAULT_TIMEOUT = 300
const HTTP_DEFAULT_IDLECONN_TIMEOUT = 90
const HTTP_DEFAULT_TLSHANDSHAKE_TIMEOUT = 30
const MAX_IDLE_CONNS = 0
const MAX_IDLE_CONNS_PER_HOST = 300000
const MAX_CONNS_PER_HOST = 300000

const HTTP_GET = "GET"
const HTTP_POST = "POST"
const HTTP_PUT = "PUT"
const HTTP_PATCH = "PATCH"
const HTTP_DELETE = "DELETE"

const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMETextCSV                          = "text/csv"
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
	MIMEApplicationGraphql               = "application/graphql"
)

const (
	charsetUTF8 = "charset=UTF-8"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"
)

type HttpClient struct {
	Client             *http.Client
	Username, Password string
	CustomUserAgent    string
}

func NewHTTPClient(username, password string, timeout int, sslverifyflg bool, customuseragent string) *HttpClient {
	http.DefaultTransport.(*http.Transport).MaxIdleConns = MAX_IDLE_CONNS
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = MAX_IDLE_CONNS_PER_HOST
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = MAX_CONNS_PER_HOST
	http.DefaultTransport.(*http.Transport).DisableKeepAlives = false
	http.DefaultTransport.(*http.Transport).IdleConnTimeout = HTTP_DEFAULT_IDLECONN_TIMEOUT * time.Second
	http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout = HTTP_DEFAULT_TLSHANDSHAKE_TIMEOUT * time.Second
	http.DefaultTransport.(*http.Transport).ExpectContinueTimeout = 1 * time.Second
	http.DefaultTransport.(*http.Transport).DisableCompression = true
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	I := &HttpClient{
		Client:          &http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: transport},
		Username:        username,
		Password:        password,
		CustomUserAgent: customuseragent,
	}
	return I
}

func (c *HttpClient) NewRequest(ctx context.Context, method, requesturi string, body io.Reader, requesttype string) (*http.Request, error) {
	if _, err := url.ParseRequestURI(requesturi); err != nil {
		log.Error(fmt.Sprintf("failed to parse url: %s", requesturi))
		log.Error(err)
		return nil, err
	}
	log.Info("RequestURI : " + requesturi)
	log.Info("RequestMethod : " + method)
	log.Info("RequestType : " + requesttype)
	log.Debug(fmt.Sprintf("RequestData : %v", body))

	req, err := http.NewRequest(method, requesturi, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Close = true

	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	if strings.HasPrefix(requesttype, MIMEApplicationGraphql) || strings.HasPrefix(requesttype, MIMEApplicationJSON) || strings.HasPrefix(requesttype, MIMEApplicationForm) || strings.HasPrefix(requesttype, MIMEMultipartForm) || strings.HasPrefix(requesttype, MIMEApplicationXML) || strings.HasPrefix(requesttype, MIMETextXML) {
		req.Header.Set(HeaderContentType, requesttype)
	}
	if c.CustomUserAgent != "" {
		req.Header.Set("User-Agent", c.CustomUserAgent)
	} else {
		req.Header.Set("User-Agent", userAgent)
	}

	return req, nil
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

func (c *HttpClient) PostMultipart(ctx context.Context, method, requesturi string, params map[string][]string, files map[string][]*multipart.FileHeader, xrequestid string, requestheadermap map[string]string) (*http.Response, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for name, fileheaderlist := range files {
		for _, fileheader := range fileheaderlist {
			part, err := writer.CreateFormFile(name, filepath.Base(fileheader.Filename))
			if err != nil {
				return nil, err
			}
			f, err := fileheader.Open()
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(part, f); err != nil {
				return nil, err
			}
			f.Close()
		}
	}
	for key, vallist := range params {
		for _, val := range vallist {
			if err := writer.WriteField(key, val); err != nil {
				return nil, err
			}
		}
	}
	contenttype := writer.FormDataContentType()

	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := c.NewRequest(ctx, method, requesturi, &buf, contenttype)
	if err != nil {
		return nil, err
	}
	req.Header.Set(HeaderContentType, contenttype)
	for key, value := range requestheadermap {
		req.Header.Set(key, value)
	}
	log.Info(req.Header)
	return c.Do(req)
}

func (c *HttpClient) RequestMultipartWithData(ctx context.Context, method, requesturi string, postdata *bytes.Buffer, requesttype string, xrequestid string, requestheadermap map[string]string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, method, requesturi, postdata, requesttype)
	if err != nil {
		return nil, err
	}
	for key, value := range requestheadermap {
		req.Header.Set(key, value)
	}
	log.Info(req.Header)
	return c.Do(req)
}
