package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Spell struct {
	Index         string
	Name          string
	Desc          string
	Range         string
	Components    []string
	Material      string
	Ritual        bool
	Duration      string
	Concentration bool
	CastingTime   string `json:"casting_time"`
	Level         int
	School        struct{ Name string }
	Classes       []struct{ Name string }
	Subclasses    []struct{ Name string }
}

type SpellAPI struct {
	Count    int
	Next     string
	previous string
	Results  []SpellTemp
}

type SpellTemp struct {
	Index         string `json:"slug"`
	Name          string
	Desc          string
	HigherLevel   string `json:"higher_level"`
	Range         string
	Components    string
	Material      string
	Ritual        string
	Duration      string
	Concentration string
	CastingTime   string `json:"casting_time"`
	Level         int    `json:"level_int"` // api has level and level_int
	School        string
	Class         string `json:"dnd_class"`
	Archeatype    string
	Circles       string
}

func fetchSpells() ([]Spell, error) {
	url := "https://api.open5e.com/spells/"
	allSpells := []Spell{}

	for url != "" {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("cannot read response:", err)
			return nil, err
		}

		var jsonResp SpellAPI
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			fmt.Println("cannot parse json:", err)
			return nil, err
		}

		url = jsonResp.Next

		standard := spellAPIToStandard(&jsonResp.Results)
		allSpells = append(allSpells, *standard...)
	}

	return allSpells, nil
}

func spellAPIToStandard(spells *[]SpellTemp) *[]Spell {
	a := []Spell{}

	for _, spell := range *spells {
		s := Spell{}

		s.Index = spell.Index
		s.Name = spell.Name
		s.Desc = spell.Desc
		s.Range = spell.Range
		s.Components = strings.Split(spell.Components, ", ")
		s.Material = spell.Material
		if spell.Ritual == "yes" {
			s.Ritual = true
		}
		s.Duration = spell.Duration
		if spell.Concentration == "yes" {
			s.Concentration = true
		}
		s.CastingTime = spell.CastingTime
		s.Level = spell.Level
		s.School = struct{ Name string }{spell.School}
		for _, class := range strings.Split(spell.Class, ", ") {
			s.Classes = append(s.Classes, struct{ Name string }{class})
		}

		a = append(a, s)
	}

	return &a
}
