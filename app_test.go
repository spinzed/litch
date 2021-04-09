package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var ExampleSpells = Spells{{Index: "acid-arrow", Name: "Acid Arrow", Desc: "Green arrow", HigherLevel: "", Range: "90 feet", Components: []string{"V", "S", "M"}, Material: "", Ritual: false, Duration: "", Concentration: false, CastingTime: "1 action", Level: 2, School: struct{ Name string }{Name: "Evocation"}, Classes: []struct{ Name string }{{Name: "Druid"}, {Name: "Wizard"}}, Subclasses: []struct{ Name string }{{Name: "Druid (Swamp)"}}}, {Index: "acid-splash", Name: "Acid Splash", Desc: "Bubble", HigherLevel: "", Range: "60 feet", Components: []string{"V", "S"}, Material: "", Ritual: false, Duration: "Instantaneous", Concentration: false, CastingTime: "1 action", Level: 0, School: struct{ Name string }{Name: "Conjuration"}, Classes: []struct{ Name string }{{Name: "Sorcerer"}, {Name: "Wizard"}}, Subclasses: []struct{ Name string }(nil)}, Spell{Index: "cone-of-cold", Name: "Cone of Cold", Desc: "Blast of air", HigherLevel: "1d8", Range: "", Components: []string{"V", "S", "M"}, Material: "A small crystal or glass cone.", Ritual: false, Duration: "Instantaneous", Concentration: false, CastingTime: "1 action", Level: 5, School: struct{ Name string }{Name: "Evocation"}, Classes: []struct{ Name string }{{Name: "Druid"}, {Name: "Sorcerer"}, {Name: "Wizard"}}, Subclasses: []struct{ Name string }(nil)}, Spell{Index: "confusion", Name: "Confusion", Desc: "Twists minds", HigherLevel: "5 feet", Range: "", Components: []string(nil), Material: "Three walnut shells.", Ritual: false, Duration: "", Concentration: true, CastingTime: "1 action", Level: 4, School: struct{ Name string }{Name: ""}, Classes: []struct{ Name string }{{Name: "Bard"}, {Name: "Druid"}}, Subclasses: []struct{ Name string }{{Name: "Cleric (Knowledge)"}}}}

var AppTest = newApp()

func TestStatus(t *testing.T) {
	var tests = []struct {
		mode InputMode
	}{
		{InputNormal},
		{InputCommand},
	}

	for _, test := range tests {
		AppTest.setInputMode(test.mode)
		if AppTest.inputMode != test.mode {
			t.Errorf("Unexpected mode, expected: \"%v\", but got \"%v\"", test.mode, AppTest.inputMode)
		}
	}
}

func TestSetInputMode(t *testing.T) {
	var tests = []struct {
		have InputMode
		want InputMode
		err  error
	}{
		{InputCommand, InputCommand, nil},
		{InputNormal, InputNormal, nil},
		{InputMode(3), -1, errors.New("Selected invalid mode: 3")},
	}

	for _, test := range tests {
		if test.want == -1 {
			test.want = AppTest.inputMode
		}
		err := AppTest.setInputMode(test.have)

		if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", test.err) {
			t.Errorf("Unexpected error, expected: \"%v\", but got \"%v\"", test.err, err)
		}

		if AppTest.inputMode != test.want {
			t.Errorf("Unexpected input mode, expected: \"%v\", but got \"%v\"", test.want, AppTest.inputMode)
		}
	}
}

func TestSetSpells(t *testing.T) {
	var tests = []struct {
		spells *Spells
		substr string
		want   []string
	}{
		// 2 Acid Arrow, 0 Acid Splash, 5 Cone of Cold, 4 Confusion
		{&ExampleSpells, "Splash", []string{fmt.Sprintf("0 Acid %sSplash%s", HlghtSubstr, HlghtNormal)}},
		{&ExampleSpells, "Acid", []string{fmt.Sprintf("2 %sAcid%s Arrow", HlghtSubstr, HlghtNormal), fmt.Sprintf("0 %sAcid%s Splash", HlghtSubstr, HlghtNormal)}},
		{&ExampleSpells, "co", []string{fmt.Sprintf("5 %sCo%sne of %sCo%sld", HlghtSubstr, HlghtNormal, HlghtSubstr, HlghtNormal), fmt.Sprintf("4 %sCo%snfusion", HlghtSubstr, HlghtNormal)}},
		{&ExampleSpells, "Coon", []string(nil)},
	}

	for _, test := range tests {
		AppTest.spells = test.spells
		AppTest.setInputText(test.substr)
		output := AppTest.setSpells()
		//t.Log(output)
		if !reflect.DeepEqual(output, &test.want) {
			t.Errorf("Unexpected result.\nhave: \"%#v\"\nwant: \"%#v\"", *output, test.want)
		}
	}
}

func TestHighlight(t *testing.T) {
	var tests = []struct {
		str    string
		substr string
		want   string
	}{
		{"Aura of Light", "Light", fmt.Sprintf("Aura of %sLight%s", HlghtSubstr, HlghtNormal)},
		{"Light of Aura", "Light ", fmt.Sprintf("%sLight %sof Aura", HlghtSubstr, HlghtNormal)},
		{"Testing of Testing Test", "Meow", "Testing of Testing Test"},
		{"Aura of Aura", "Aura", fmt.Sprintf("%sAura%s of %sAura%s", HlghtSubstr, HlghtNormal, HlghtSubstr, HlghtNormal)},
		{"Dog Meo Dog", "Meow", "Dog Meo Dog"},
		{"Work Meo work", "Meow", "Work Meo work"},
		{"Seagull is a boring animal", "Boring", fmt.Sprintf("Seagull is a %sboring%s animal", HlghtSubstr, HlghtNormal)},
		{"Seagull is a boring animal", "BoRiNg", fmt.Sprintf("Seagull is a %sboring%s animal", HlghtSubstr, HlghtNormal)},
		{"Woof", "", "Woof"},
	}
	for _, test := range tests {
		output := highlight(test.str, test.substr)
		//t.Log(output)
		if output != test.want {
			t.Errorf("Unexpected result, expected \"%s\", but got \"%s\"", test.want, output)
		}
	}
}
