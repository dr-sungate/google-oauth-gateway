package handler

import (
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Users struct {
}

func (us Users) GetUsers(c echo.Context) error {
	log.Warn(c.QueryParams())
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}
