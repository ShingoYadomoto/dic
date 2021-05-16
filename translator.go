package main

type (
	Result struct {
		Origin string
		Dist   string
	}

	Translator interface {
		Translate(string) ([]Result, error)
	}
)
