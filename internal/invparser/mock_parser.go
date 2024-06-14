package invparser

type MockParser struct {
}

func (p MockParser) Parse(content string) Invoice {
	return Invoice{
		amount: "100",
	}
}
