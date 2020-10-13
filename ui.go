package main

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type UI struct {
	app              *tview.Application
	list             *tview.List
	input            *tview.InputField
	statusBox        *tview.TextView
	widebox          *WideBox
	wideboxFakeFocus bool
	inputText        string
	spells           *[]Spell
	dataChan         chan []Spell
}

func newApp() *UI {
	ui := UI{}

	app := tview.NewApplication().EnableMouse(false)
	ui.app = app

	ui.list = getList()
	ui.input = getInputField(ui.setInputText)
	ui.statusBox = tview.NewTextView().SetTextAlign(tview.AlignRight)
	ui.statusBox.Box.SetBorder(false)
	ui.widebox = getWideBox()
	ui.dataChan = make(chan []Spell)
	go ui.waitForData()
	go loadAllData(ui.dataChan)

	ui.app.SetInputCapture(ui.handleInput)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(ui.list, 0, 3, true).
			AddItem(ui.widebox.grid, 0, 7, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(ui.input, 0, 3, false).
			AddItem(ui.statusBox, 0, 7, false), 1, 0, false)

	app.SetRoot(root, true)

	app.SetFocus(ui.input)

	return &ui
}

func (ui *UI) Run() {
	ui.app.Run()
}

func (ui *UI) waitForData() {
	for v := range ui.dataChan {
		ui.setSpells(&v)
	}
}

func (ui *UI) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		//ui.wideboxFakeFocus = true
		ui.widebox.SetSpell(ui.currentSelectedSpell())
	case tcell.KeyUp, tcell.KeyCtrlK:
		if ui.wideboxFakeFocus {
			ui.widebox.ScrollUp()
			break
		}
		ui.list.SetCurrentItem(ui.list.GetCurrentItem() - 1)
	case tcell.KeyCtrlJ, tcell.KeyDown:
		if ui.wideboxFakeFocus {
			ui.widebox.ScrollDown()
			break
		}
		item := ui.list.GetCurrentItem()
		if item >= ui.list.GetItemCount()-1 {
			ui.list.SetCurrentItem(0)
			break
		}
		ui.list.SetCurrentItem(item + 1)
	case tcell.KeyLeft, tcell.KeyCtrlH:
		ui.wideboxFakeFocus = false
	case tcell.KeyRight, tcell.KeyCtrlL:
		ui.wideboxFakeFocus = true
	case tcell.KeyEsc:
		ui.wideboxFakeFocus = !ui.wideboxFakeFocus

	}
	return event
}

func (ui *UI) setSpells(spells *[]Spell) {
	ui.spells = spells

	for i, s := range *spells {
		ui.list.AddItem(s.Name, strconv.Itoa(i), 0, nil)
	}
	// for some reason, the screen isnt auto updated so it has to be so manually
	// a fix without invoking this function below is prefered
	ui.app.Draw()
}

func (ui UI) currentSelectedSpell() *Spell {
	if ui.list.GetItemCount() < 1 {
		return nil
	}
	_, s := ui.list.GetItemText(ui.list.GetCurrentItem())
	index, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	spells := *ui.spells
	return &spells[index]
}

func (ui *UI) setInputText(text string) {
	ui.wideboxFakeFocus = false
	ui.inputText = text
	ui.list.Clear()
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

func getList() *tview.List {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true)

	return list
}

func getInputField(processInput func(string)) *tview.InputField {
	input := tview.NewInputField().SetLabel(">>> ")
	input.SetChangedFunc(processInput)

	return input
}
