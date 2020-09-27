package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// check if file/dir exists, true if does, false if doesnt
func checkFile(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	if !os.IsNotExist(err) && err != nil {
		panic(err)
	}
	return true
}

// check if dir exists, make it if it doesnt
func readyDir(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println(dir, "does not exist, creating...")
		os.MkdirAll(dir, 0755)
	}
	if !os.IsNotExist(err) && err != nil {
		panic(err)
	}
}

func loadJSONFromFile(file string, dest interface{}) {
	cachedFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	stuff, err := ioutil.ReadAll(cachedFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(stuff, dest)
	if err != nil {
		panic(err)
	}
}

func loadData(file string) (*[]Spell, error) {
	if !checkFile(file) {
		fmt.Println("spells.json doesnt exist")

		data, err := fetchSpells()
		if err != nil {
			return nil, err
		}

		cached, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			panic(err)
		}

		file, err := os.Create(file)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		file.Write(cached)
		return &data, nil
	}
	var data []Spell
	loadJSONFromFile(file, &data)

	return &data, nil
}

func mergeMultipleSources(s1 *[]Spell, s2 *[]Spell) *[]Spell {
	var final []Spell
	var i1, i2 int
	arr1 := *s1
	arr2 := *s2

	for {
		spell1 := arr1[i1]
		spell2 := arr2[i2]

		switch {
		case spell1.Index == spell2.Index:
			final = append(final, spell1)
			i1++
			i2++
		// spell1 comes before spell2
		case spell1.Index < spell2.Index:

		// spell1 comes after spell2
		case spell1.Index > spell2.Index:
		}
	}
}

func readyAllData() (*[]Spell, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	readyDir(config + "/banshie/cache")
	readyDir(config + "/banshie/local")

	dataAPI, err := loadData(config + "/banshie/cache/spells.json")
	if err != nil {
		return nil, err
	}

	if !checkFile(config + "/banshie/local/spells.json") {
		return dataAPI, nil
	}
	var userData *[]Spell
	loadJSONFromFile(config+"/banshie/local/spells.json", userData)

	allSpells := append(*dataAPI, *userData...)

	return &allSpells, nil
}
