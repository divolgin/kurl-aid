package ui

import (
	"strings"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pkg/errors"
)

type PromptResult int

const (
	Yes PromptResult = iota
	No
)

func Prompt(msg string) (PromptResult, error) {
	if err := termui.Init(); err != nil {
		return No, errors.Wrap(err, "create terminal ui")
	}
	defer termui.Close()

	drawMessage(msg)
	drawFooter()

	for e := range termui.PollEvents() {
		if e.Type != termui.KeyboardEvent {
			continue
		}

		if strings.ToLower(e.ID) == "y" {
			return Yes, nil
		}

		if strings.ToLower(e.ID) == "n" {
			return No, nil
		}
	}

	// this will never happen
	return No, errors.New("event poll")
}

func drawMessage(msg string) {
	termWidth, termHeight := termui.TerminalDimensions()

	p := widgets.NewParagraph()
	p.Text = msg
	p.WrapText = true
	p.SetRect(0, 0, termWidth, termHeight-1)

	termui.Render(p)
}

func drawFooter() {
	termWidth, termHeight := termui.TerminalDimensions()

	instructions := widgets.NewParagraph()
	instructions.Text = "[y] yes    [n] no"
	instructions.Border = false

	left := 0
	right := termWidth
	top := termHeight - 1
	bottom := termHeight

	instructions.SetRect(left, top, right, bottom)
	termui.Render(instructions)
}
