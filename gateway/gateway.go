package main

import (
	"github.com/dr-sungate/google-oauth-gateway/api/handler"
	"github.com/dr-sungate/google-oauth-gateway/api/service/client"
	"github.com/dr-sungate/google-oauth-gateway/api/service/custommiddleware"
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"

	"net/http"
	_ "net/http/pprof"
	"runtime"
)

const DEFAULT_PORT = "8080"

func main() {
	//############## 計測モード ###############
	if os.Getenv("VERIFY_MODE") == "enable" {
		runtime.SetBlockProfileRate(1)
		go func() {
			log.Error("", http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}
	//#####################################
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	if os.Getenv("VERIFY_MODE") == "enable" {
		e.Debug = true
	}

	e.GET("/oauth2/authorize", handler.GoogleOauth2Handler{}.Authorize)
	e.GET("/oauth2/callback", handler.GoogleOauth2Handler{}.Callback)

	oauth2config := custommiddleware.OAuth2Config{
		GoCacheClient: client.NewGoCacheClient(custommiddleware.DefaultJOAuth2Config.PublicKeyTtl),
	}
	oauth2group := e.Group("/api/v1")
	oauth2group.Use(custommiddleware.OAuth2WithConfig(oauth2config))

	oauth2group.GET("/user/:id", handler.Users{}.GetUsers)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	e.Logger.Fatal(e.Start(":" + port))
}
