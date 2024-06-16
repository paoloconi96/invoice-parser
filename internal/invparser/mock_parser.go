package invparser

import (
	"context"
	"github.com/Rhymond/go-money"
)

type MockParser struct {
}

func (p MockParser) Parse(ctx context.Context, content *[]byte) Invoice {
	return Invoice{
		amount: money.New(100, "EUR"),
	}
}
