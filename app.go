package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App holds all references for App components of the App.App as well as all
// info regarding an App instance
type App struct {
	app              *tview.Application
	list             *tview.List
	input            *tview.InputField
	statusBox        *tview.TextView
	widebox          *WideBox
	wideboxFakeFocus bool
	inputMode        InputMode
	inputText        string
	spells           *Spells
	dataChan         chan Spells
	statusChan       chan string
	fetchLock        bool
	eventReg         *EventRegister
}

// Instantiate a new app ready to run
func newApp() *App {
	app := App{}

	defer func() {
		// perform an app cleanup and after that keep on packing
		if r := recover(); r != nil {
			app.Quit()
			panic(r)
		}
	}()

	ui := tview.NewApplication().EnableMouse(false)
	app.app = ui
	// sets the default color for tview primitives. A bit weird way to do
	// that but eh. Does not set the background color of some pritives like
	// input fields
	tview.Styles.PrimitiveBackgroundColor = DefaultBgColor

	// instantiate all parts of the UI
	app.list = getList()
	app.input = getInputField(app.setInputText)
	app.statusBox = getStatusBox()
	app.widebox = getWideBox()
	app.setInputMode(InputNormal)
	app.dataChan = make(chan Spells)
	app.statusChan = make(chan string)

	// set up channel loops in separate goroutines which wait for data and statuses.
	// It is important that they are set up before the the first data fetch.
	go app.waitForData()
	go app.waitForStatuses()

	// instantiate a logger
	l := NewLogger(LogFile)
	app.eventReg = NewEventRegister(l, app.statusChan)

	// make sure that spells are initialized if fetching goes wrong
	app.spells = new(Spells)
	app.FetchData(false)

	// set the global input handler
	app.app.SetInputCapture(app.handleInput)

	// bind all the UI elements to the app instance
	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(app.list, 0, 3, true).
			AddItem(app.widebox.grid, 0, 7, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(app.input, 0, 3, false).
			AddItem(app.statusBox, 0, 7, false), 1, 0, false)

	ui.SetRoot(root, true)

	ui.SetFocus(app.input)
	app.focusList()

	return &app
}

// Run the app
func (app *App) Run() {
	defer func() {
		// perform an app cleanup and after that keep on panicking
		if r := recover(); r != nil {
			app.Quit()
			panic(r)
		}
	}()

	app.app.Run()
}

// Quit the app aka do the cleanup
func (app *App) Quit() {
	app.eventReg.logger.Close()
}

// Waits for the data in the data channel. Stays always open
func (app *App) waitForData() {
	for v := range app.dataChan {
		app.spells = &v
		app.updateSpellList()
		// update the data lock
		app.fetchLock = false
		// for some reason, the screen isn't auto updated on the initial spell set
		// so it has to be so manually.
		app.app.Draw()
	}
}

// Waits for the statuses in the status channel. Stays always open
func (app *App) waitForStatuses() {
	for v := range app.statusChan {
		app.setStatus(v)
	}
}

// The main app global input handler
func (app *App) handleInput(event *tcell.EventKey) *tcell.EventKey {
	// if status exists, clear it, but only if data is not being fetched atm
	if app.Status() != "" && !app.fetchLock {
		app.setStatus("")
	}
	//fmt.Print(event)
	switch event.Key() {
	case tcell.KeyEnter:
		//ui.wideboxFakeFocus = true
		if spell := app.currentSelectedSpell(); spell != nil {
			app.widebox.SetSpell(spell)
		}
	case tcell.KeyUp, tcell.KeyCtrlK:
		if app.wideboxFakeFocus {
			app.widebox.ScrollUp()
			break
		}
		app.list.SetCurrentItem(app.list.GetCurrentItem() - 1)
	case tcell.KeyCtrlJ, tcell.KeyDown:
		if app.wideboxFakeFocus {
			app.widebox.ScrollDown()
			break
		}
		item := app.list.GetCurrentItem()
		if item >= app.list.GetItemCount()-1 {
			app.list.SetCurrentItem(0)
			break
		}
		app.list.SetCurrentItem(item + 1)
	// tcell.KeyCtrlBackspace doesn't exist for whatever reason
	case tcell.KeyCtrlD:
		app.input.SetText("")
	case tcell.KeyF5:
		var isForce bool
		if event.Modifiers() == tcell.ModShift {
			isForce = true
		}
		go app.FetchData(isForce)
	case tcell.KeyLeft, tcell.KeyCtrlH:
		app.focusList()
	case tcell.KeyRight, tcell.KeyCtrlL:
		app.focusWideBox()
	case tcell.KeyTab:
		app.switchFocus()
	case tcell.KeyESC:
		app.setInputMode(InputNormal)
	case tcell.KeyCtrlN:
		app.setInputMode(InputCommand)
	// this will run before its default behavior (closing the application)
	case tcell.KeyCtrlC:
		app.Quit()
	}
	return event
}

func (app *App) InputMode() InputMode {
	return app.inputMode
}

func (app *App) setInputMode(mode InputMode) error {
	if mode == app.InputMode() {
		return nil
	}
	oldMode := app.inputMode
	app.inputMode = mode
	switch mode {
	case InputNormal:
		app.input.SetLabel("> ")
		return nil
	case InputCommand:
		app.input.SetLabel(": ")
		return nil
	}
	// by this point, if the mode was valid, the function would return, that
	// means that it is invalid and must be reverted
	app.inputMode = oldMode
	user := "Error while trying to switch modes, check logs"
	err := fmt.Errorf("Selected invalid mode: %d", mode)
	app.eventReg.Register(EventInfo, user, err.Error())
	return err
}

func (app *App) Status() string {
	// the bool signifies whether should the color tags be stripped off or not
	return app.statusBox.GetText(true)
}

func (app *App) setStatus(text string) {
	app.statusBox.SetText(text)
}

// Switches focus between the list on the left and the main content area to the right
func (app *App) switchFocus() {
	if app.wideboxFakeFocus {
		app.focusList()
		return
	}
	app.focusWideBox()
}

// Focuses the list
func (app *App) focusList() {
	app.wideboxFakeFocus = false
	app.list.SetBorderAttributes(tcell.AttrBold)
	app.widebox.grid.SetBorderAttributes(tcell.AttrNone)
}

// Focuses the main content area on the right
func (app *App) focusWideBox() {
	app.wideboxFakeFocus = true
	app.list.SetBorderAttributes(tcell.AttrNone)
	app.widebox.grid.SetBorderAttributes(tcell.AttrBold)
}

// Filters and sets the spells from app.spells and updates it on the screen.
// Does NOT update app.spells. The main task of the return value is for testing.
func (app *App) updateSpellList() *[]string {
	var items []string
	app.list.Clear()
	for i, s := range *app.spells {
		lname := strings.ToLower(s.Name)
		linput := strings.ToLower(app.inputText)

		if !strings.Contains(lname, linput) {
			continue
		}
		nameString := strconv.Itoa(s.Level) + " " + s.Name

		if s.Ritual || s.Concentration {
			_, _, w, _ := app.list.Box.GetInnerRect()
			padLen := w - len(nameString)
			padNum := 0
			if s.Concentration {
				padNum++
			}
			if s.Ritual {
				padNum++
			}

			if padLen >= 3 {
				nameString += strings.Repeat(" ", padLen-padNum)
				if s.Concentration {
					nameString += "C"
				}
				if s.Ritual {
					nameString += "R"
				}
			}
		}

		hlght := highlight(nameString, app.inputText)
		items = append(items, hlght)
		app.list.AddItem(hlght, strconv.Itoa(i), 0, nil)
	}
    // title shows how many spells are shown out of total
    app.list.SetTitle(fmt.Sprintf("%d/%d", len(items), len(*app.spells)))
	return &items
}

// Returns the current selected spell. Returns nil if there are no spells in the list
func (app App) currentSelectedSpell() *Spell {
	if app.list.GetItemCount() < 1 {
		return nil
	}
	_, s := app.list.GetItemText(app.list.GetCurrentItem())
	index, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	spells := *app.spells
	return &spells[index]
}

// Handler than should be ran on every text input change. Filters the spell list
// on text update.
func (app *App) setInputText(text string) {
	// focus the list on key input if the main content box happens to be focused atm
	// TODO: change behavior on different modes (normal, command...)
	app.focusList()
	app.inputText = text
	app.updateSpellList()
}

// Highlight a substring in a string regardless of it's capitalisation.
// May not work properly with unicode
func highlight(str, substr string) string {
	if substr == "" {
		return str
	}
	lname := strings.ToLower(str)
	linput := strings.ToLower(substr)
	parts := strings.Split(lname, linput)

	// precalculated lengths of color prefixes for a small performance gain
	prelen := len(HlghtSubstr)
	postlen := len(HlghtNormal)
	patternlen := len(substr)

	var final string
	for i, w := range parts {
		startx := len(final)
		if i > 1 {
			startx -= (i - 1) * (prelen + postlen)
		}
		if i != 0 {
			final += HlghtSubstr + str[startx:startx+patternlen] + HlghtNormal
			startx += patternlen
		}
		final += str[startx : startx+len(w)]
	}
	return final
}

// Returns a pointer to a new list element preconfigured for the app
func getList() *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorBlack).
		SetHighlightFullLine(true)
	list.SetBorder(true)
	return list
}

// Returns a pointer to a new input element preconfigured for the app
func getInputField(processInput func(string)) *tview.InputField {
	input := tview.NewInputField().
		SetLabel("> ")
	input.SetChangedFunc(processInput).SetFieldBackgroundColor(DefaultBgColor)
	return input
}

func getStatusBox() *tview.TextView {
	box := tview.NewTextView().SetTextAlign(tview.AlignRight)
	box.SetBorder(false)
	return box
}
