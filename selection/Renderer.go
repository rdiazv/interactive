package selection

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"interactive/helper"
)

type renderer struct {
	Question   string
	Options    []*Option
	Selection  []interface{}
	LineIndex  int
	LineOffset int
}

func NewRenderer(question string, options []*Option) *renderer {
	return &renderer{
		LineIndex:  0,
		LineOffset: 0,
		Selection:  make([]interface{}, 0),
		Options:    options,
		Question:   question,
	}
}

func (r *renderer) Init() ([]interface{}, bool) {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	r.Render()

	canceled := false

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				canceled = true
				break mainloop

			case termbox.KeyEnter:
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

	if canceled {
		return nil, true
	}

	return r.Selection, false
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
	return helper.Min(height-2, len(r.Options))
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

	caretColor := color.New(color.FgHiGreen).SprintFunc()
	keyColor := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	selectionColor := color.New(color.FgGreen).SprintFunc()
	questionColor := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	optionColor := color.New(color.FgYellow).SprintFunc()
	topHintColor := color.New(color.Reset).SprintFunc()
	bottomHintColor := color.New(color.Faint).SprintFunc()

	fmt.Println()

	if len(r.Selection) > 0 {
		fmt.Println(
			selectionColor("?"),
			questionColor(r.Question),
			topHintColor("(Press"),
			keyColor("<enter>"),
			topHintColor("to confirm)"))
	} else {
		fmt.Println(
			selectionColor("?"),
			questionColor(r.Question),
			topHintColor("(Press"),
			keyColor("<space>"),
			topHintColor("to select)"))
	}

	for i := 0; i < height; i++ {
		var selectionCharacter string
		var checkedCharacter string

		if r.IsSelected(r.Options[i+r.LineOffset].Value) {
			selectionCharacter = selectionColor("◉")
		} else {
			selectionCharacter = "◯"
		}

		if i == r.LineIndex {
			checkedCharacter = caretColor("❯")
		} else {
			checkedCharacter = " "
		}

		fmt.Println(checkedCharacter, selectionCharacter, optionColor(r.Options[i+r.LineOffset].Text))
	}

	for i := 0; i < totalHeight-height-2; i++ {
		fmt.Println()
	}

	fmt.Print(bottomHintColor("(Move up and down to reveal more choices)"))

	termbox.Flush()
}
