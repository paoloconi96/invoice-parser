package invparser

import (
	"context"
	"github.com/Rhymond/go-money"
	"time"
)

type MockParser struct {
}

func (p MockParser) Parse(ctx context.Context, invoice *InputInvoice) Invoice {
	time.Sleep(250 * time.Millisecond)

	return Invoice{
		amount: money.New(100, "EUR"),
	}
}
