package invparser

import (
	"time"
)

type Category uint8

const (
	Other Category = iota
)

type Product struct {
	name     string
	amount   string
	category Category
}

type Invoice struct {
	amount   string
	date     time.Time
	products []Product
}

type Parser interface {
	Parse(content string) Invoice
}
