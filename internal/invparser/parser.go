package invparser

import (
	"context"
	"github.com/Rhymond/go-money"
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

type Invoice struct {
	amount   *money.Money
	date     time.Time
	products []Product
	language string
}

type Parser interface {
	Parse(ctx context.Context, content *[]byte) Invoice
}
