package main

import (
	"fmt"
	"interactive/prompt"
	"interactive/structs"
)

func main() {
	options := make([]*structs.Option, 0)

	for i := 1; i <= 100; i++ {
		options = append(options, &structs.Option{
			Text:  fmt.Sprintf("Option %d", i),
			Value: i,
		})
	}

	prompt.Selection("Choose which tenants to install.", options)
}
