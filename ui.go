package main

import "github.com/rivo/tview"

type UI struct {
	app *tview.Application
}

func newApp() *UI {
	ui := UI{}

	app := tview.NewApplication()
	ui.app = app

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox().SetBorder(true), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true), 0, 7, false), 0, 1, false).
		AddItem(tview.NewInputField().SetLabel(">>> "), 1, 0, true)

	app.SetRoot(root, true)

	return &ui
}

func (ui *UI) Run() {
	ui.app.Run()
}
