package structs

import (
	"fmt"
	"github.com/fatih/color"
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
				r.ToggleSelection(r.Options[r.LineIndex+r.LineOffset].Value)

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

func (r *renderer) GetUsableHeight() int {
	_, height := termbox.Size()
	return helper.Min(height-1, len(r.Options))
}

func (r *renderer) Move(lines int) {
	height := r.GetUsableHeight()
	middle := height / 2
	movingDown := lines > 0
	movingUp := lines < 0

	if r.LineIndex == middle {
		if movingUp && r.LineOffset > 0 {
			r.LineOffset += lines
		} else if movingDown && len(r.Options)-r.LineOffset > height {
			r.LineOffset += lines
		} else {
			r.LineIndex += lines
		}
	} else if movingUp && r.LineIndex > 0 {
		r.LineIndex += lines
	} else if movingDown && r.LineIndex+r.LineOffset+1 < len(r.Options) {
		r.LineIndex += lines
	}
}

func (r *renderer) Render() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	height := r.GetUsableHeight()
	_, totalHeight := termbox.Size()

	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	for i := 0; i < height; i++ {
		var selectionCharacter string
		var checkedCharacter string

		if r.IsSelected(r.Options[i+r.LineOffset].Value) {
			selectionCharacter = green("◉")
		} else {
			selectionCharacter = "◯"
		}

		if i == r.LineIndex {
			checkedCharacter = yellow("❯")
		} else {
			checkedCharacter = " "
		}

		fmt.Println(checkedCharacter, selectionCharacter, r.Options[i+r.LineOffset].Text)
	}

	for i := 0; i < totalHeight-height-1; i++ {
		fmt.Println()
	}

	termbox.Flush()
}
