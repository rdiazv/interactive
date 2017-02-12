package main

import (
	"fmt"
	"github.com/rdiazv/interactive/selection"
)

func main() {
	options := make([]*selection.Option, 0)

	for i := 1; i <= 20; i++ {
		options = append(options, &selection.Option{
			Text:  fmt.Sprintf("Option %d", i),
			Value: i,
		})
	}

	values, canceled := selection.Ask(&selection.Question{
		Message: func(selection []interface{}) string {
			if len(selection) == 0 {
				return "Choose some options."
			}

			return fmt.Sprintf("Your have chosen %d option(s).", len(selection))
		},
		Choices: options,
	})

	if canceled {
		fmt.Println("Canceled!")
	} else {
		fmt.Println("Values:", values)
	}
}
