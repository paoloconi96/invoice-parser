package main

import (
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/paoloconi96/invoice-parser/internal/invparser"
	"mime/multipart"
	"net/http"
)

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
