package client

import (
	"context"
	"errors"
	"fmt"
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"golang.org/x/oauth2"
	v2 "google.golang.org/api/oauth2/v2"
	"os"
	"time"
)

const HTTP_REQUEST_TIMEOUT_DEFAULT = 120

var HTTP_REQUEST_TIMEOUT time.Duration = HTTP_REQUEST_TIMEOUT_DEFAULT

const (
	authorizeEndpoint      = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint          = "https://www.googleapis.com/oauth2/v4/token"
	AUTH_KEY_ACCESS_TOKEN  = "access_token"
	AUTH_KEY_REFRESH_TOKEN = "refresh_token"
	AUTH_KEY_ID_TOKEN      = "id_token"
	AUTH_KEY_ID_EXPIRED    = "expired"
	AUTH_KEY_ID            = "id"
	AUTH_KEY_EMAIL         = "email"
)

type GoogleOAuth2Client struct {
	Config oauth2.Config
}

func NewGoogleOAuth2Client() *GoogleOAuth2Client {
	scopes := []string{"openid", "email", "profile"}
	log.Debug(fmt.Sprintf("ClientId: %s", os.Getenv("GOOGLE_CLIENTID")))
	log.Debug(fmt.Sprintf("ClientSecret : %s", os.Getenv("GOOGLE_CLIENTSECRET")))
	log.Debug(fmt.Sprintf("RedirectURL : %s", os.Getenv("GOOGLE_CALLBACKURL")))
	I := &GoogleOAuth2Client{
		Config: oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENTID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENTSECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  authorizeEndpoint,
				TokenURL: tokenEndpoint,
			},
			RedirectURL: os.Getenv("GOOGLE_CALLBACKURL"),
			Scopes:      scopes,
		},
	}
	return I
}

func (goc *GoogleOAuth2Client) GetAuthCodeUrl(state string) string {
	return goc.Config.AuthCodeURL(state)
}

func (goc *GoogleOAuth2Client) Callback(code string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_REQUEST_TIMEOUT*time.Second)
	defer cancel()
	authedmap := make(map[string]string, 0)
	token, err := goc.Config.Exchange(ctx, code)
	if err != nil {
		log.Error(err)
		return authedmap, err
	}
	if token.Valid() == false {
		err := errors.New("Invaild Token")
		log.Error(err)
		return authedmap, err
	}
	authedmap[AUTH_KEY_ACCESS_TOKEN] = token.AccessToken
	authedmap[AUTH_KEY_REFRESH_TOKEN] = token.RefreshToken
	authedmap[AUTH_KEY_ID_TOKEN] = fmt.Sprintf("%v", token.Extra("id_token"))
	authedmap[AUTH_KEY_ID_EXPIRED] = token.Expiry.Format(time.RFC3339)
	service, _ := v2.New(goc.Config.Client(ctx, token))
	tokenInfo, _ := service.Tokeninfo().AccessToken(token.AccessToken).Context(ctx).Do()

	authedmap[AUTH_KEY_ID] = tokenInfo.UserId
	authedmap[AUTH_KEY_EMAIL] = tokenInfo.Email
	return authedmap, nil
}
