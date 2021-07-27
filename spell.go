package main

import (
	"encoding/json"
	"fmt"
	"io"
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

// Modified during testing
var fetchFunc func(string) ([]byte, error) = __fetchSpellsAPI

func fetchSpells(url string, data *Spells) error {
	allSpells := Spells{}

	for url != "" {
		body, err := fetchFunc(url)
		if err != nil {
			fmt.Println("cannot read response:", err)
			return err
		}

		var jsonResp SpellAPI
		if err := json.Unmarshal(body, &jsonResp); err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				return fmt.Errorf("Cannot parse json: %s, at byte offset %d", err, e.Offset)
			}
			return fmt.Errorf("Cannot parse json: %s", err)
		}

		url = jsonResp.Next

		standard := spellAPIToStandard(&jsonResp.Results)
		allSpells = append(allSpells, *standard...)
	}

	*data = allSpells
	return nil
}

func __fetchSpellsAPI(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
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
		if spell.Components != "" {
			s.Components = strings.Split(spell.Components, ",")
			for i, c := range s.Components {
				s.Components[i] = strings.TrimSpace(c)
			}
		}
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
		for _, class := range strings.Split(spell.Class, ",") {
			class = strings.TrimSpace(class)
			if class != "Ritual Caster" {
				s.Classes = append(s.Classes, struct{ Name string }{class})
			}
		}

		if spell.Archetype != "" {
			for _, arch := range strings.Split(spell.Archetype, ",") {
				arch = strings.TrimSpace(arch)
				parts := strings.Split(arch, ": ")
                // if there are at least 2 parts, continue, if there are not,
                // there is a problem in parsing without a doubt
                if len(parts) >= 2 {
                    subclass := parts[0] + " (" + parts[1] + ")"
                    s.Subclasses = append(s.Subclasses, struct{ Name string }{subclass})
                }
			}
		}
		if spell.Circles != "" {
			for _, circle := range strings.Split(spell.Circles, ",") {
				circle = strings.TrimSpace(circle)
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
