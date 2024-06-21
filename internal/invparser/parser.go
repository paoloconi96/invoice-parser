package invparser

import (
	"context"
	"github.com/Rhymond/go-money"
	"io"
	"time"
)

type Category uint8

const (
	Other Category = iota
)

type Product struct {
	name     string
	amount   *money.Money
	category Category
}

type InvoiceId string

type InputInvoice struct {
	Id         InvoiceId `json:"id"`
	FileReader io.Reader `json:"-"`
}

type Invoice struct {
	amount   *money.Money
	date     time.Time
	products []Product
	language string
}

type Parser interface {
	Parse(ctx context.Context, invoice *InputInvoice) Invoice
}
