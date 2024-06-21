package invparser

import (
	documentai "cloud.google.com/go/documentai/apiv1beta3"
	"cloud.google.com/go/documentai/apiv1beta3/documentaipb"
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"google.golang.org/api/option"
	"net/http"
	"strconv"
	"time"
)

type GDocAiParser struct {
	client      *documentai.DocumentProcessorClient
	location    string
	projectID   string
	processorID string
}

func (p *GDocAiParser) Close() error {
	return p.client.Close()
}

func (p *GDocAiParser) Parse(ctx context.Context, inputInvoice *InputInvoice) Invoice {
	var content []byte
	inputInvoice.FileReader.Read(content)

	req := &documentaipb.ProcessRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", p.projectID, p.location, p.processorID),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  content,
				MimeType: http.DetectContentType(content),
			},
		},
	}
	resp, err := p.client.ProcessDocument(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
	}

	document := resp.GetDocument()

	// TODO: Compute this using all the pages and the confidence level for each
	// TODO: Use an enum for this value
	language := document.GetPages()[0].DetectedLanguages[0].LanguageCode
	invoice := Invoice{
		language: language,
	}

	var amount int
	var currency string
	for _, entity := range document.GetEntities() {
		switch entity.Type {
		case "total_amount":
			amount, err = strconv.Atoi(entity.GetNormalizedValue().GetText())
			if err != nil {
				// TODO: Handle invalid amount
			}
		case "currency":
			currency = entity.GetNormalizedValue().GetText()
		case "due_date":
			date := entity.NormalizedValue.GetStructuredValue().(*documentaipb.Document_Entity_NormalizedValue_DateValue).DateValue
			invoice.date = time.Date(int(date.Year), time.Month(date.Month), int(date.Day), 0, 0, 0, 0, time.UTC)
		}
	}

	// TODO: This could produce an invalid amount
	invoice.amount = money.New(int64(amount), currency)

	return invoice
}

func NewGDocAiParser(ctx context.Context, location string, projectID string, processorID string) *GDocAiParser {
	client, err := documentai.NewDocumentProcessorClient(
		ctx,
		option.WithEndpoint(fmt.Sprintf("%s-documentai.googleapis.com:443", location)),
	)
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
	}

	return &GDocAiParser{
		client:      client,
		location:    location,
		projectID:   projectID,
		processorID: processorID,
	}
}
