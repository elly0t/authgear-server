package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/oursky/ourd/oderr"
	"net/http"

	"github.com/oursky/ourd/authtoken"
	"github.com/oursky/ourd/oddb"
	"github.com/oursky/ourd/router"
)

type apiKeyValidatonPreprocessor struct {
	Key     string
	AppName string
}

func (p apiKeyValidatonPreprocessor) Preprocess(payload *router.Payload, response *router.Response) int {
	apiKey := payload.APIKey()
	if apiKey != p.Key {
		log.Debugf("Invalid APIKEY: %v", apiKey)
		response.Err = oderr.NewFmt(oderr.CannotVerifyAPIKey, "Cannot verify api key: %v", apiKey)
		return http.StatusUnauthorized
	}

	payload.AppName = p.AppName

	return http.StatusOK
}

type connPreprocessor struct {
	DBOpener func(string, string, string) (oddb.Conn, error)
	DBImpl   string
	Option   string
}

func (p connPreprocessor) Preprocess(payload *router.Payload, response *router.Response) int {
	log.Debugf("Opening DBConn: {%v %v %v}", p.DBImpl, payload.AppName, p.Option)

	conn, err := p.DBOpener(p.DBImpl, payload.AppName, p.Option)
	if err != nil {
		response.Err = err
		return http.StatusServiceUnavailable
	}
	payload.DBConn = conn

	log.Debugf("Get DB OK")

	return http.StatusOK
}

type tokenStorePreprocessor struct {
	authtoken.Store
}

func (p tokenStorePreprocessor) Preprocess(payload *router.Payload, response *router.Response) int {
	payload.TokenStore = p.Store
	return http.StatusOK
}

func authenticateUser(payload *router.Payload, response *router.Response) int {
	tokenString := payload.AccessToken()
	if tokenString == "" { // no access token, leave it unauthenticated
		return http.StatusOK
	}

	store := payload.TokenStore
	token := authtoken.Token{}

	if err := store.Get(tokenString, &token); err != nil {
		response.Err = err
		return http.StatusUnauthorized
	}

	payload.AppName = token.AppName
	payload.UserInfoID = token.UserInfoID

	return http.StatusOK
}

func injectUserIfPresent(payload *router.Payload, response *router.Response) int {
	if payload.UserInfoID == "" {
		log.Debugln("injectUser: empty UserInfoID, skipping")
		return http.StatusOK
	}

	conn := payload.DBConn
	userinfo := oddb.UserInfo{}
	if err := conn.GetUser(payload.UserInfoID, &userinfo); err != nil {
		log.Errorf("Cannot find UserInfo.ID = %#v\n", payload.UserInfoID)
		response.Err = err
		return http.StatusInternalServerError
	}

	payload.UserInfo = &userinfo

	return http.StatusOK
}

func injectDatabase(payload *router.Payload, response *router.Response) int {
	conn := payload.DBConn

	databaseID, _ := payload.Data["database_id"].(string)
	switch databaseID {
	case "_public":
		payload.Database = conn.PublicDB()
	case "_private":
		if payload.UserInfo != nil {
			payload.Database = conn.PrivateDB(payload.UserInfo.ID)
		} else {
			response.Err = errors.New("Authentication is needed for private DB access")
			return http.StatusUnauthorized
		}
	}

	return http.StatusOK
}
