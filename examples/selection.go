package main

import (
	"fmt"
	"interactive/selection"
)

func main() {
	options := make([]*selection.Option, 0)

	for i := 1; i <= 20; i++ {
		options = append(options, &selection.Option{
			Text:  fmt.Sprintf("Option %d", i),
			Value: i,
		})
	}

	values, canceled := selection.Ask("Choose which tenants to install.", options)

	if canceled {
		fmt.Println("Canceled!")
	} else {
		fmt.Println("Values:", values)
	}
}
