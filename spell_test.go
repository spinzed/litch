package main

import (
	"reflect"
	"strconv"
	"testing"
)

var SpellsJSON []string = []string{`{
    "count": 321,
    "next": "2",
    "previous": null,
    "results": [
        {
            "slug": "acid-arrow",
            "name": "Acid Arrow",
            "desc": "Green arrow",
            "page": "phb 259",
            "range": "90 feet",
            "components": "V, S, M",
            "ritual": "no",
            "concentration": "no",
            "casting_time": "1 action",
            "level": "2nd-level",
            "level_int": 2,
            "school": "Evocation",
            "dnd_class": "Druid, Wizard",
            "archetype": "Druid: Swamp",
            "circles": "Swamp",
            "document__slug": "wotc-srd"
        },
        {
            "slug": "acid-splash",
            "name": "Acid Splash",
            "desc": "Bubble",
            "higher_level": "",
            "page": "phb 211",
            "range": "60 feet",
            "components": "V, S",
            "material": "",
            "ritual": "no",
            "duration": "Instantaneous",
            "concentration": "no",
            "casting_time": "1 action",
            "level": "Cantrip",
            "level_int": 0,
            "school": "Conjuration",
            "dnd_class": "Sorcerer,Wizard",
            "archetype": "",
            "circles": "",
            "document__title": "Systems Reference Document",
            "document__license_url": "http://aaheee.com/ssdgi"
        }
	]
}
`, `{
    "count": 321,
    "next": "",
    "previous": "1",
    "results": [
        {
            "slug": "cone-of-cold",
            "name": "Cone of Cold",
            "desc": "Blast of air",
            "higher_level": "1d8",
            "components": "V, S, M",
            "material": "A small crystal or glass cone.",
            "ritual": "no",
            "duration": "Instantaneous",
            "casting_time": "1 action",
            "level": "5th-level",
            "level_int": 5,
            "school": "Evocation",
            "dnd_class": "Druid,   Sorcerer  , Wizard",
            "document__slug": "wotc-srd",
            "document__license_url": "http://eeee.com/aa"
        },
        {
            "slug": "confusion",
            "name": "Confusion",
            "desc": "Twists minds",
            "higher_level": "5 feet",
            "page": "phb 224",
            "material": "Three walnut shells.",
            "ritual": "no",
            "concentration": "yes",
            "casting_time": "1 action",
            "level": "4th-level",
            "level_int": 4,
            "dnd_class": "  Bard,   Druid  ",
            "archetype": "Cleric: Knowledge",
            "circles": "",
            "document__slug": "wotc-srd"
        }
	]
}
`}

func __fetchSpellsTest(index string) ([]byte, error) {
	ind, err := strconv.Atoi(index)
	if err != nil {
		return nil, err
	}
	bits := []byte(SpellsJSON[ind-1])
	return bits, nil
}

func TestFetchSpells(t *testing.T) {
	fetchFunc = __fetchSpellsTest
	var test = struct {
		url      string
		dataAPI  []string
		wantData *Spells
		wantErr  error
	}{
		"1",
		SpellsJSON,
		&ExampleSpells,
		nil,
	}

	var output Spells
	err := fetchSpells(test.url, &output)
	if err != nil {
		t.Fatalf("Got error: \"%s\", but expected nil", err)
	}
	if !reflect.DeepEqual(output, *test.wantData) {
		//t.Fatalf("Data not identical, expected:\n%v\nbut got:\n%v\n", *test.wantData, output)
		t.Fatalf("Data not identical, expected:\n%#v\nbut got:\n%#v\n", *test.wantData, output)
	}
}

func TestSpellAPIToStandard(t *testing.T) {
	var tests = []struct {
		have []SpellTemp
		want Spells
	}{
		{
			[]SpellTemp{{Index: "acid-arrow", Name: "Acid Arrow", Desc: "Green arrow", HigherLevel: "", Range: "90 feet", Components: "V, S, M", Material: "", Ritual: "no", Duration: "", Concentration: "no", CastingTime: "1 action", Level: 2, School: "Evocation", Class: "Druid, Wizard", Archetype: "Druid: Swamp", Circles: "Swamp"}, {Index: "acid-splash", Name: "Acid Splash", Desc: "Bubble", HigherLevel: "", Range: "60 feet", Components: "V, S", Material: "", Ritual: "no", Duration: "Instantaneous", Concentration: "no", CastingTime: "1 action", Level: 0, School: "Conjuration", Class: "Sorcerer,Wizard", Archetype: "", Circles: ""}},
			Spells{Spell{Index: "acid-arrow", Name: "Acid Arrow", Desc: "Green arrow", HigherLevel: "", Range: "90 feet", Components: []string{"V", "S", "M"}, Material: "", Ritual: false, Duration: "", Concentration: false, CastingTime: "1 action", Level: 2, School: struct{ Name string }{Name: "Evocation"}, Classes: []struct{ Name string }{{Name: "Druid"}, {Name: "Wizard"}}, Subclasses: []struct{ Name string }{{Name: "Druid (Swamp)"}}}, Spell{Index: "acid-splash", Name: "Acid Splash", Desc: "Bubble", HigherLevel: "", Range: "60 feet", Components: []string{"V", "S"}, Material: "", Ritual: false, Duration: "Instantaneous", Concentration: false, CastingTime: "1 action", Level: 0, School: struct{ Name string }{Name: "Conjuration"}, Classes: []struct{ Name string }{{Name: "Sorcerer"}, {Name: "Wizard"}}, Subclasses: []struct{ Name string }(nil)}},
		},
		{
			[]SpellTemp{{Index: "cone-of-cold", Name: "Cone of Cold", Desc: "Blast of air", HigherLevel: "1d8", Range: "", Components: "V, S, M", Material: "A small crystal or glass cone.", Ritual: "no", Duration: "Instantaneous", Concentration: "", CastingTime: "1 action", Level: 5, School: "Evocation", Class: "Druid,   Sorcerer  , Wizard", Archetype: "", Circles: ""}, {Index: "confusion", Name: "Confusion", Desc: "Twists minds", HigherLevel: "5 feet", Range: "", Components: "", Material: "Three walnut shells.", Ritual: "no", Duration: "", Concentration: "yes", CastingTime: "1 action", Level: 4, School: "", Class: "  Bard,   Druid  ", Archetype: "Cleric: Knowledge", Circles: ""}},
			Spells{Spell{Index: "cone-of-cold", Name: "Cone of Cold", Desc: "Blast of air", HigherLevel: "1d8", Range: "", Components: []string{"V", "S", "M"}, Material: "A small crystal or glass cone.", Ritual: false, Duration: "Instantaneous", Concentration: false, CastingTime: "1 action", Level: 5, School: struct{ Name string }{Name: "Evocation"}, Classes: []struct{ Name string }{{Name: "Druid"}, {Name: "Sorcerer"}, {Name: "Wizard"}}, Subclasses: []struct{ Name string }(nil)}, Spell{Index: "confusion", Name: "Confusion", Desc: "Twists minds", HigherLevel: "5 feet", Range: "", Components: []string(nil), Material: "Three walnut shells.", Ritual: false, Duration: "", Concentration: true, CastingTime: "1 action", Level: 4, School: struct{ Name string }{Name: ""}, Classes: []struct{ Name string }{{Name: "Bard"}, {Name: "Druid"}}, Subclasses: []struct{ Name string }{{Name: "Cleric (Knowledge)"}}}},
		},
	}
	for _, test := range tests {
		result := spellAPIToStandard(&test.have)
		if !reflect.DeepEqual(*result, test.want) {
			//t.Fatalf("Data not identical, expected:\n%v\nbut got:\n%v\n", *test.wantData, output)
			t.Fatalf("Data not identical, expected:\n%#v\nbut got:\n%#v\n", test.want, *result)
		}
	}
}
