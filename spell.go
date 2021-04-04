package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Spells []Spell

// These methods make Spells type satisfy the sort.Sort interface
// It may not be a good idea to sort with raw values, but not pointer
// values, but we shall see is it a bottleneck when benchmarks are written
func (s Spells) Len() int           { return len(s) }
func (s Spells) Less(i, j int) bool { return s[i].Index < s[j].Index }
func (s Spells) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Spell struct {
	Index         string
	Name          string
	Desc          string
	HigherLevel   string
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
	Count int
	Next  string
	//previous string
	Results []SpellTemp
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
	Archetype     string
	Circles       string
}

func fetchSpells(url string, data *Spells) error {
	allSpells := Spells{}

	for url != "" {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("cannot read response:", err)
			return err
		}

		var jsonResp SpellAPI
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			fmt.Println("cannot parse json:", err)
			return err
		}

		url = jsonResp.Next

		standard := spellAPIToStandard(&jsonResp.Results)
		allSpells = append(allSpells, *standard...)
	}

	*data = allSpells
	return nil
}

func spellAPIToStandard(spells *[]SpellTemp) *Spells {
	a := Spells{}

	for _, spell := range *spells {
		s := Spell{}

		s.Index = spell.Index
		s.Name = spell.Name
		s.Desc = spell.Desc
		s.HigherLevel = spell.HigherLevel
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
			if class != "Ritual Caster" {
				s.Classes = append(s.Classes, struct{ Name string }{class})
			}
		}

		if spell.Archetype != "" {
			for _, arch := range strings.Split(spell.Archetype, "<br/> ") {
				parts := strings.Split(arch, ": ")
				subclass := parts[0] + " (" + parts[1] + ")"
				s.Subclasses = append(s.Subclasses, struct{ Name string }{subclass})
			}
		}
		if spell.Circles != "" {
			for _, circle := range strings.Split(spell.Circles, ", ") {
				subclass := "Druid (" + circle + ")"
				// There is an inconsistency in the api where the circle is sometimes
				// contained in the archeatypes, this fixes that
				var contained bool
				for _, subcl := range s.Subclasses {
					if subcl.Name == subclass {
						contained = true
						break
					}
				}
				if contained {
					continue
				}
				s.Subclasses = append(s.Subclasses, struct{ Name string }{subclass})
			}
		}

		a = append(a, s)
	}

	return &a
}
