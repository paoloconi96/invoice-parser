package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var clientId = os.Getenv("GOOGLE_AUTH_CLIENT_ID")

type LoginCredentials struct {
	Credential string `json:"credential"`
}

func login(rCtx *gin.Context) {
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

	fmt.Println(payload)
	rCtx.JSON(http.StatusNoContent, nil)
}
