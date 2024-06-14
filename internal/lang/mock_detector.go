package lang

type MockDetector struct {
}

func (d MockDetector) Detect(input string) Iso639 {
	return IT
}
