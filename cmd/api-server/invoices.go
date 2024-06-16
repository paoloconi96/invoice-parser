package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/paoloconi96/invoice-parser/internal/invparser"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Invoice struct {
	File *multipart.FileHeader `form:"file"`
}

var parserChannel = make(chan *[]byte)
var ctx = context.Background()

func main() {
	parser := invparser.NewGDocAiParser(
		ctx,
		os.Getenv("GOOGLE_CLOUD_LOCATION"),
		os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
		os.Getenv("GOOGLE_CLOUD_DOCUMENT_AI_PROCESSOR_ID"),
	)
	go processInvoices(parserChannel, parser)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/v1/invoices", addInvoice)

	router.Run("localhost:8000")
}

func addInvoice(ctx *gin.Context) {
	var invoice Invoice
	invalidMessage := "Invalid file provided"

	if err := ctx.Bind(&invoice); err != nil {
		ctx.JSON(http.StatusBadRequest, invalidMessage)
		return
	}

	file, err := invoice.File.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidMessage)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidMessage)
		return
	}

	fmt.Println(string(fileContent))
	parserChannel <- &fileContent

	ctx.JSON(http.StatusOK, invoice)
}

func processInvoices(tasks <-chan *[]byte, parser invparser.Parser) {
	for taskInput := range tasks {
		fmt.Println(parser.Parse(ctx, taskInput))
	}
}
