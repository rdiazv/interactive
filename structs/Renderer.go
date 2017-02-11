package structs

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type renderer struct {
	Options   []*Option
	Selection []interface{}
	LineIndex int
}

func NewRenderer() *renderer {
	return &renderer{
		LineIndex: 0,
		Selection: make([]interface{}, 0),
	}
}

func (r renderer) Init() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	r.Render()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop

			case termbox.KeySpace:
				value := r.Options[r.LineIndex].Value

				if inArray(r.Selection, value) {
					r.Selection = removeValue(r.Selection, value)
				} else {
					r.Selection = append(r.Selection, value)
				}

			case termbox.KeyArrowUp:
				if r.LineIndex > 0 {
					r.LineIndex -= 1
				}

			case termbox.KeyArrowDown:
				r.LineIndex += 1
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		r.Render()
	}
}

func (r renderer) Render() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	_, height := termbox.Size()

	for i := 0; i < min(height, len(r.Options))-1; i++ {
		var selectionCharacter string
		var checkedCharacter string

		if inArray(r.Selection, r.Options[i].Value) {
			selectionCharacter = "◉"
		} else {
			selectionCharacter = "◯"
		}

		if i == r.LineIndex {
			checkedCharacter = "❯"
		} else {
			checkedCharacter = " "
		}

		fmt.Println(checkedCharacter, selectionCharacter, r.Options[i].Text)
	}

	termbox.Flush()
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func inArray(array []interface{}, value interface{}) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func removeValue(array []interface{}, value interface{}) []interface{} {
	for i, item := range array {
		if item == value {
			return append(array[:i], array[i+1:]...)
		}
	}

	return array
}
