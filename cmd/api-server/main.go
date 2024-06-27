package main

import (
	idtoken2 "cloud.google.com/go/auth/credentials/idtoken"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/paoloconi96/invoice-parser/internal/invparser"
	"github.com/paoloconi96/invoice-parser/internal/websocket"
	"mime/multipart"
)

type InvoiceRequestBody struct {
	File *multipart.FileHeader `form:"file"`
}

var parserChannel = make(chan *invparser.InputInvoice)
var ctx = context.Background()
var hub *websocket.Hub

// TODO: Handle err
var validator, _ = idtoken2.NewValidator(nil)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	//parser := invparser.NewGDocAiParser(
	//	ctx,
	//	os.Getenv("GOOGLE_CLOUD_LOCATION"),
	//	os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
	//	os.Getenv("GOOGLE_CLOUD_DOCUMENT_AI_PROCESSOR_ID"),
	//)
	parser := invparser.MockParser{}
	hub = websocket.NewHub()
	go hub.Run()

	go processInvoices(parserChannel, parser)

	router := gin.Default()

	// TODO: Change this config when the domain gets unified
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	router.POST("/api/v1/invoices", addInvoice)
	router.GET("/api/v1/websocket", func(ctx *gin.Context) {
		websocket.HandleWebsocket(ctx, hub)
	})
	router.POST("/api/v1/login", login)

	router.Run("localhost:8000")
}
