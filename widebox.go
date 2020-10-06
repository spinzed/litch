package main

import (
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

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

func (b *WideBox) SetSpell(s *Spell) {
	b.SetName(s.Name)
	b.SetLevel(s.Level)
	b.SetRitual(s.Ritual)
	b.SetConentration(s.Concentration)
	b.SetClasses(s.Classes)
	b.SetCastingTime(s.CastingTime)
	b.SetRange(s.Range)
	b.SetComponents(s.Components)
	b.SetDuration(s.Duration)
	b.SetDescription(s.Desc, s.HigherLevel)
}

func (b *WideBox) SetName(s string) {
	b.namebox.SetText(s)
}

func (b *WideBox) SetLevel(lvl int) {
	if lvl == 0 {
		b.lvlbox.SetText("Cantrip")
	}
	b.lvlbox.SetText("Level " + strconv.Itoa(lvl))
}

func (b *WideBox) SetRitual(r bool) {
	b.ritualbox.SetText("")
	if r {
		b.ritualbox.SetText("Ritual")
	}
}

func (b *WideBox) SetConentration(c bool) {
	b.concentrbox.SetText("")
	if c {
		b.concentrbox.SetText("Concentration")
	}
}

func (b *WideBox) SetClasses(c []struct{ Name string }) {
	var names []string
	for _, n := range c {
		names = append(names, n.Name)
	}
	b.classbox.SetText(strings.Join(names, ", "))
}

func (b *WideBox) SetCastingTime(s string) {
	b.timecastbox.SetText(s)
}

func (b *WideBox) SetRange(s string) {
	b.rangebox.SetText(s)
}

func (b *WideBox) SetComponents(c []string) {
	b.componentbox.SetText(strings.Join(c, ", "))
}

func (b *WideBox) SetDuration(d string) {
	b.durationbox.SetText(d)
}

func (b *WideBox) SetDescription(d string, hl string) {
	b.descbox.SetText(d + "\n\n" + hl)
}

func getWideBox() *WideBox {
	box := WideBox{}

	namebox := tview.NewTextView()
	box.namebox = namebox
	lvlbox := tview.NewTextView()
	box.lvlbox = lvlbox
	ritualbox := tview.NewTextView()
	box.ritualbox = ritualbox
	concentrbox := tview.NewTextView()
	box.concentrbox = concentrbox
	classbox := tview.NewTextView()
	box.classbox = classbox
	timecastbox := tview.NewTextView()
	box.timecastbox = timecastbox
	rangebox := tview.NewTextView()
	box.rangebox = rangebox
	componentbox := tview.NewTextView()
	box.componentbox = componentbox
	durationbox := tview.NewTextView()
	box.durationbox = durationbox
	descbox := tview.NewTextView().SetWrap(true).SetWordWrap(true)
	box.descbox = descbox

	grid := tview.NewGrid().
		SetRows(1, 1, 1, 2, 2).
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
