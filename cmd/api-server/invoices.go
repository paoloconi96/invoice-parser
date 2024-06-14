package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paoloconi96/invoice-parser/internal/invparser"
	"github.com/paoloconi96/invoice-parser/internal/lang"
	"mime/multipart"
	"net/http"
)

type Invoice struct {
	File *multipart.FileHeader `form:"file"`
}

var parserChannel = make(chan string)

func main() {
	go processInvoices(parserChannel, lang.MockDetector{}, invparser.MockParser{})

	router := gin.Default()
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

	var fileContent []byte
	_, err = file.Read(fileContent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidMessage)
		return
	}

	parserChannel <- string(fileContent)

	ctx.JSON(http.StatusOK, invoice)
}

func processInvoices(tasks <-chan string, detector lang.Detector, parser invparser.Parser) {
	for taskInput := range tasks {
		fmt.Println(detector.Detect(taskInput))
		fmt.Println(parser.Parse(taskInput))
	}
}
