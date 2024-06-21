package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/paoloconi96/invoice-parser/internal/invparser"
	"github.com/paoloconi96/invoice-parser/internal/websocket"
	"mime/multipart"
	"net/http"
	"os"
)

type InvoiceRequestBody struct {
	File *multipart.FileHeader `form:"file"`
}

var parserChannel = make(chan *invparser.InputInvoice)
var ctx = context.Background()
var hub *websocket.Hub

func main() {
	parser := invparser.NewGDocAiParser(
		ctx,
		os.Getenv("GOOGLE_CLOUD_LOCATION"),
		os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
		os.Getenv("GOOGLE_CLOUD_DOCUMENT_AI_PROCESSOR_ID"),
	)
	hub = websocket.NewHub()
	go hub.Run()

	go processInvoices(parserChannel, parser)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/v1/invoices", addInvoice)
	router.GET("/api/v1/websocket", func(ctx *gin.Context) {
		websocket.HandleWebsocket(ctx, hub)
	})

	router.Run("localhost:8000")
}

func addInvoice(ctx *gin.Context) {
	var invoiceRequestBody InvoiceRequestBody
	invalidMessage := "Invalid file provided"

	if err := ctx.Bind(&invoiceRequestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, invalidMessage)
		return
	}

	fileReader, err := invoiceRequestBody.File.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidMessage)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(fileReader)

	id, _ := gonanoid.New()
	invoice := invparser.InputInvoice{
		Id:         invparser.InvoiceId(id),
		FileReader: fileReader,
	}

	parserChannel <- &invoice

	ctx.JSON(http.StatusOK, invoice)
}

func processInvoices(tasks <-chan *invparser.InputInvoice, parser invparser.Parser) {
	for taskInput := range tasks {
		parser.Parse(ctx, taskInput)
		hub.Broadcast <- &taskInput.Id
	}
}
