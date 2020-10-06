package main

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type UI struct {
	app       *tview.Application
	list      *tview.List
	input     *tview.InputField
	widebox   *WideBox
	inputText string
	spells    *[]Spell
}

func newApp(spells *[]Spell) *UI {
	ui := UI{}
	ui.spells = spells

	app := tview.NewApplication().EnableMouse(true)
	ui.app = app

	list := getList(spells)
	ui.list = list
	list.SetInputCapture(ui.handleListInput)
	input := getInputField(ui.setInputText)
	ui.input = input
	widebox := getWideBox()
	ui.widebox = widebox

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(list, 0, 3, true).
			AddItem(widebox.grid, 0, 7, false), 0, 1, false).
		AddItem(input, 1, 0, false)

	app.SetRoot(root, true)

	app.SetFocus(input)

	return &ui
}

func (ui *UI) Run() {
	ui.app.Run()
}

func (ui *UI) handleListInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		ui.widebox.SetSpell(ui.currentSelectedSpell())
	}
	return event
}

func (ui UI) currentSelectedSpell() *Spell {
	_, s := ui.list.GetItemText(ui.list.GetCurrentItem())
	index, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	spells := *ui.spells
	return &spells[index]
}

func (ui *UI) setInputText(text string) {
	ui.list.Clear()
	ui.inputText = text
	for i, s := range *ui.spells {
		lname := strings.ToLower(s.Name)
		linput := strings.ToLower(ui.inputText)

		if strings.Contains(lname, linput) {
			ui.list.AddItem(formatSpell(s.Name, text), strconv.Itoa(i), 0, nil)
		}
	}
}

// may not work properly with unicode
func formatSpell(name, pattern string) string {
	lname := strings.ToLower(name)
	linput := strings.ToLower(pattern)
	parts := strings.Split(lname, linput)
	pre := "[#ff0000]"
	post := "[white]"

	// precalculated lengths for small performance benefits
	prelen := len(pre)
	postlen := len(post)
	patternlen := len(pattern)

	var final string
	for i, w := range parts {
		startx := len(final)
		if i > 1 {
			startx -= (i - 1) * (prelen + postlen)
		}
		if i != 0 {
			final += pre + name[startx:startx+patternlen] + post
			startx += patternlen
		}
		final += name[startx : startx+len(w)]
	}
	return final
}

func getList(spells *[]Spell) *tview.List {
	list := tview.NewList().ShowSecondaryText(false)

	for i, s := range *spells {
		list.AddItem(s.Name, strconv.Itoa(i), 0, nil)
	}
	list.SetBorder(true)
	return list
}

func getInputField(processInput func(string)) *tview.InputField {
	input := tview.NewInputField().SetLabel(">>> ")
	input.SetChangedFunc(func(text string) {
		processInput(text)
	})
	return input
}
