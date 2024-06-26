package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

type LoginCredentials struct {
	Credential string `json:"credential"`
}

func login(rCtx *gin.Context) {
	var clientId = os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	var store = sessions.NewFilesystemStore("data/session", []byte(os.Getenv("SESSION_KEY")))

	var loginCredentials LoginCredentials

	if err := rCtx.BindJSON(&loginCredentials); err != nil {
		log.Printf("%s", err)
		rCtx.JSON(http.StatusBadRequest, "Invalid format")
		return
	}

	payload, err := validator.Validate(ctx, loginCredentials.Credential, clientId)
	if err != nil {
		log.Printf("%s", err)
		rCtx.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	userAttr := payload.Claims
	session, err := store.New(rCtx.Request, "session-name")
	if err != nil {
		log.Printf("%s", err)
		rCtx.JSON(http.StatusInternalServerError, "Error authenticating")
		return
	}

	sessionVal := session.Values
	for index, val := range userAttr {
		sessionVal[index] = val
	}
	err = session.Save(rCtx.Request, rCtx.Writer)
	if err != nil {
		log.Printf("%s", err)
		rCtx.JSON(http.StatusInternalServerError, "Error authenticating")
		return
	}

	rCtx.JSON(http.StatusNoContent, nil)
}
