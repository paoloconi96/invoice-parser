package main

import (
	"fmt"
	"os"
)

const invoicePath = "data/invoice-1.pdf"

func main() {
	file, err := os.Open(invoicePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Println(float64(stats.Size())/1024, "KB")

}
