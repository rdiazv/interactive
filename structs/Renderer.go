package structs

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"interactive/helper"
)

type renderer struct {
	Options    []*Option
	Selection  []interface{}
	LineIndex  int
	LineOffset int
}

func NewRenderer() *renderer {
	return &renderer{
		LineIndex:  0,
		LineOffset: 0,
		Selection:  make([]interface{}, 0),
	}
}

func (r *renderer) Init() {
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
				r.ToggleSelection(r.Options[r.LineIndex].Value)

			case termbox.KeyArrowUp:
				r.Move(-1)

			case termbox.KeyArrowDown:
				r.Move(1)
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		r.Render()
	}
}

func (r *renderer) IsSelected(value interface{}) bool {
	return helper.InArray(r.Selection, value)
}

func (r *renderer) ToggleSelection(value interface{}) {
	if r.IsSelected(value) {
		r.Selection = helper.RemoveFromArray(r.Selection, value)
	} else {
		r.Selection = append(r.Selection, value)
	}
}

func (r *renderer) Move(lines int) {
	r.LineIndex += lines
}

func (r *renderer) Render() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	_, height := termbox.Size()

	for i := 0; i < helper.Min(height, len(r.Options))-1; i++ {
		var selectionCharacter string
		var checkedCharacter string

		if r.IsSelected(r.Options[i].Value) {
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
