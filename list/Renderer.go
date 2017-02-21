package list

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"reflect"
	"strings"
)

type renderer struct {
	Question   *Question
	Selection  interface{}
	LineIndex  int
	LineOffset int
}

func Ask(question *Question) (interface{}, bool) {
	r := &renderer{
		LineIndex:  0,
		LineOffset: 0,
		Selection:  nil,
		Question:   question,
	}

	return r.Init()
}

func (r *renderer) Init() (interface{}, bool) {
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
				r.ToggleSelection(r.Question.Choices[r.LineIndex+r.LineOffset].Value)

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
	return r.Selection == value
}

func (r *renderer) ToggleSelection(value interface{}) {
	if r.Selection == value {
		r.Selection = nil
	} else {
		r.Selection = value
	}
}

func (r *renderer) GetUsableHeight() int {
	_, height := termbox.Size()

	if len(r.Question.Choices) <= height-1 {
		return len(r.Question.Choices)
	} else {
		return height - 2
	}
}

func (r *renderer) GetMessage() string {
	typeOf := reflect.TypeOf(r.Question.Message)
	valueOf := reflect.ValueOf(r.Question.Message)

	switch typeOf.Kind() {
	case reflect.Func:
		in := []reflect.Value{
			reflect.ValueOf(r.Selection),
		}

		return valueOf.Call(in)[0].String()
	case reflect.String:
		return valueOf.String()
	}

	return ""
}

func (r *renderer) Move(lines int) {
	height := r.GetUsableHeight()
	middle := height / 2
	movingDown := lines > 0
	movingUp := lines < 0

	if r.LineIndex == middle {
		if movingUp && r.LineOffset > 0 {
			r.LineOffset += lines
		} else if movingDown && len(r.Question.Choices)-r.LineOffset > height {
			r.LineOffset += lines
		} else {
			r.LineIndex += lines
		}
	} else if movingUp && r.LineIndex > 0 {
		r.LineIndex += lines
	} else if movingDown && r.LineIndex+r.LineOffset+1 < len(r.Question.Choices) {
		r.LineIndex += lines
	}
}

func (r *renderer) Println(str ...string) {
	fmt.Print("\n", strings.Join(str, " "))
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

	if r.Selection != nil {
		r.Println(
			selectionColor("?"),
			questionColor(r.GetMessage()),
			topHintColor("(Press"),
			keyColor("<enter>"),
			topHintColor("to confirm)"))
	} else {
		r.Println(
			selectionColor("?"),
			questionColor(r.GetMessage()),
			topHintColor("(Press"),
			keyColor("<space>"),
			topHintColor("to select)"))
	}

	for i := 0; i < height; i++ {
		var selectionCharacter string
		var checkedCharacter string

		if r.IsSelected(r.Question.Choices[i+r.LineOffset].Value) {
			selectionCharacter = selectionColor("◉")
		} else {
			selectionCharacter = "◯"
		}

		if i == r.LineIndex {
			checkedCharacter = caretColor("❯")
		} else {
			checkedCharacter = " "
		}

		r.Println(checkedCharacter, selectionCharacter, optionColor(r.Question.Choices[i+r.LineOffset].Text))
	}

	hasScroll := len(r.Question.Choices) > height

	if hasScroll {
		for i := 0; i < totalHeight-height-2; i++ {
			r.Println()
		}

		r.Println(bottomHintColor("(Move up and down to reveal more choices)"))
	} else {
		for i := 0; i < totalHeight-height-1; i++ {
			r.Println()
		}
	}

	termbox.Flush()
}
