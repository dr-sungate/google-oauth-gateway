package handler

import (
	"fmt"
	"github.com/dr-sungate/google-oauth-gateway/api/service/client"
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"github.com/dr-sungate/google-oauth-gateway/api/service/parser"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

type GoogleOauth2Handler struct {
}

func (goh GoogleOauth2Handler) Authorize(c echo.Context) error {
	state := uuid.NewV4().String()

	oauth2client := client.NewGoogleOAuth2Client()
	return c.Redirect(http.StatusFound, fmt.Sprintf("%s", oauth2client.GetAuthCodeUrl(state)))
}

func (goh GoogleOauth2Handler) Callback(c echo.Context) error {
	code := c.Request().FormValue("code")
	oauth2client := client.NewGoogleOAuth2Client()
	authmap, err := oauth2client.Callback(code)
	if err != nil {
		return err
	}
	log.Info(authmap)
	return c.JSONPretty(http.StatusOK, authmap, parser.MARSHAL_INDENT)
}
