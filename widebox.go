package main

import (
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

// WideBox is the area where the spell data is shown.
type WideBox struct {
	grid         *tview.Grid
	namebox      *tview.TextView
	lvlbox       *tview.TextView
	ritualbox    *tview.TextView
	concentrbox  *tview.TextView
	classbox     *tview.TextView
	timecastbox  *tview.TextView
	rangebox     *tview.TextView
	componentbox *tview.TextView
	durationbox  *tview.TextView
	descbox      *tview.TextView
}

// Intialize a new widebox for the app
func getWideBox() *WideBox {
	newTextViewLeft := func() *tview.TextView {
		return tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(true).
			SetWordWrap(true)
	}
	newTextViewMid := func() *tview.TextView {
		return newTextViewLeft().SetTextAlign(tview.AlignCenter)
	}
	newTextViewRight := func() *tview.TextView {
		return newTextViewLeft().SetTextAlign(tview.AlignRight)
	}
	box := WideBox{}

	namebox := newTextViewLeft()
	box.namebox = namebox
	lvlbox := newTextViewLeft()
	box.lvlbox = lvlbox
	ritualbox := newTextViewRight()
	box.ritualbox = ritualbox
	concentrbox := newTextViewRight()
	box.concentrbox = concentrbox
	classbox := newTextViewLeft()
	box.classbox = classbox
	timecastbox := newTextViewMid()
	box.timecastbox = timecastbox
	rangebox := newTextViewMid()
	box.rangebox = rangebox
	componentbox := newTextViewMid()
	box.componentbox = componentbox
	durationbox := newTextViewMid()
	box.durationbox = durationbox
	descbox := newTextViewLeft()
	box.descbox = descbox

	grid := tview.NewGrid().
		SetRows(1, 1, 2, 3, 4).
		SetColumns(-1, -1).
		AddItem(namebox, 0, 0, 1, 1, 1, 1, false).
		AddItem(lvlbox, 1, 0, 1, 1, 1, 1, false).
		AddItem(ritualbox, 0, 1, 1, 1, 1, 1, false).
		AddItem(concentrbox, 1, 1, 1, 1, 1, 1, false).
		AddItem(classbox, 2, 0, 1, 2, 1, 1, false).
		AddItem(timecastbox, 3, 0, 1, 1, 1, 1, false).
		AddItem(rangebox, 3, 1, 1, 1, 1, 1, false).
		AddItem(componentbox, 4, 0, 1, 1, 1, 1, false).
		AddItem(durationbox, 4, 1, 1, 1, 1, 1, false).
		AddItem(descbox, 5, 0, 1, 2, 1, 1, false)
	grid.SetBorder(true)

	box.grid = grid
	return &box
}

// Scroll up the description
func (b *WideBox) ScrollUp() {
	r, c := b.descbox.GetScrollOffset()
	b.descbox.ScrollTo(r-1, c)
}

// Scroll down the description
func (b *WideBox) ScrollDown() {
	r, c := b.descbox.GetScrollOffset()
	b.descbox.ScrollTo(r+1, c)
}

// Sets spell data to be shown
func (b *WideBox) SetSpell(s *Spell) {
	if s == nil {
		// Level is set to -1 because cantrip is level 0
		s = &Spell{Level: -1}
	}
	b.SetName(s.Name)
	b.SetLevel(s.Level, s.School.Name)
	b.SetRitual(s.Ritual)
	b.SetConentration(s.Concentration)
	b.SetClasses(s.Classes, s.Subclasses)
	b.SetCastingTime(s.CastingTime)
	b.SetRange(s.Range)
	b.SetComponents(s.Components, s.Material)
	b.SetDuration(s.Duration)
	b.SetDescription(s.Desc, s.HigherLevel)

	b.descbox.ScrollToBeginning()
}

func (b *WideBox) SetName(s string) {
	b.namebox.SetText("[#ff5522::bu]" + s)
}

func (b *WideBox) SetLevel(lvl int, school string) {
	if lvl < 0 {
		b.lvlbox.SetText("")
		return
	}
	if lvl == 0 {
		b.lvlbox.SetText(school + " Cantrip")
		return
	}
	b.lvlbox.SetText("Level " + strconv.Itoa(lvl) + " " + school)
}

func (b *WideBox) SetRitual(r bool) {
	if r {
		b.ritualbox.SetText("[#00ff00]Ritual")
		return
	}
	b.ritualbox.SetText("")
}

func (b *WideBox) SetConentration(c bool) {
	if c {
		b.concentrbox.SetText("[yellow]Concentration")
		return
	}
	b.concentrbox.SetText("")
}

// Sets classes and subclasses
func (b *WideBox) SetClasses(c []struct{ Name string }, s []struct{ Name string }) {
	var names []string
	for _, n := range append(c, s...) {
		names = append(names, n.Name)
	}
	b.classbox.SetText(strings.Join(names, ", "))
}

func (b *WideBox) SetCastingTime(s string) {
	b.timecastbox.SetText("[orange]Casting Time[white]\n" + s)
}

func (b *WideBox) SetRange(s string) {
	b.rangebox.SetText("[orange]Range[white]\n" + s)
}

// Sets the spell components. If the material component is present, it will
// appear yellow and if it contains items that are worth x gp, it will appear red
func (b *WideBox) SetComponents(c []string, m string) {
	text := strings.Join(c, ", ")
	if m != "" {
		text += " (" + m + ")"
		split := strings.SplitN(text, "M", 2)
		color := "yellow"
		if strings.Contains(m, "gp") {
			color = "red"
		}
		text = split[0] + "[" + color + "]M[white]" + split[1]
	}
	b.componentbox.SetText("[orange]Components[white]\n" + text)
}

func (b *WideBox) SetDuration(d string) {
	b.durationbox.SetText("[orange]Duration[white]\n" + d)
}

// Sets the spell description and at higher levels description
func (b *WideBox) SetDescription(d string, hl string) {
	text := d
	if hl != "" {
		text += "\n\n[::b]At higher levels: [::-]" + hl
	}
	b.descbox.SetText(text)
}
