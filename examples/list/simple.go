package main

import (
	"fmt"
	"github.com/rdiazv/interactive/list"
)

func main() {
	options := make([]*list.Option, 0)

	for i := 1; i <= 20; i++ {
		options = append(options, &list.Option{
			Text:  fmt.Sprintf("Option %d", i),
			Value: i,
		})
	}

	value, canceled := list.Ask(&list.Question{
		Message: "Choose an option.",
		Choices: options,
	})

	if canceled {
		fmt.Println("Canceled!")
	} else {
		fmt.Println("Selected value:", value)
	}
}
