package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestMergeMultipleSources(t *testing.T) {
	var tests = []struct {
		source1 Spells
		source2 Spells
		want    Spells
	}{
		{ExampleSpells[:2], ExampleSpells[2:], ExampleSpells},
		{ExampleSpells[2:], ExampleSpells[:2], ExampleSpells},
		{Spells{ExampleSpells[0], ExampleSpells[2]}, Spells{ExampleSpells[1], ExampleSpells[3]}, ExampleSpells},
		{Spells{ExampleSpells[0], ExampleSpells[3]}, Spells{ExampleSpells[1], ExampleSpells[2]}, ExampleSpells},
	}

	for _, test := range tests {
		//t.Log(test.source1)
		//t.Log(test.source2)
		output := mergeMultipleSources(&test.source1, &test.source2)
		if !reflect.DeepEqual(*output, test.want) {
			var o, w []string
			for _, v := range *output {
				o = append(o, v.Index)
			}
			for _, v := range test.want {
				w = append(w, v.Index)
			}
			t.Errorf("Unexpected result.\nhave:\"%v\"\nwant:\"%v\"", strings.Join(o, ","), strings.Join(w, ","))
		}
	}
}
