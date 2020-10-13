package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func loadAllData(channel chan []Spell) {
	defer close(channel)

	config, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	readyDir(config + "/banshie/cache")
	readyDir(config + "/banshie/local")

	dataAPI, err := loadData(config + "/banshie/cache/spells.json")
	if err != nil {
		// should be handled via errchannel
		panic(err)
	}

	if !checkFile(config + "/banshie/local/spells.json") {
		channel <- *dataAPI
		return
	}
	fmt.Println("local spells.json detected")
	var userData []Spell
	err = loadJSONFromFile(config+"/banshie/local/spells.json", &userData)
	if err != nil {
		// should be handled via errchannel
		panic(err)
	}

	allSpells := mergeMultipleSources(&userData, dataAPI)

	channel <- *allSpells
}

func loadData(fileStr string) (*[]Spell, error) {
	// if the file exists, load it
	if checkFile(fileStr) {
		var data []Spell
		if err := loadJSONFromFile(fileStr, &data); err != nil {
			return nil, err
		}

		return &data, nil
	}

	fmt.Println("spells.json doesnt exist")

	data, err := fetchSpells()
	if err != nil {
		return nil, err
	}

	cached, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileStr)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(cached)
	return &data, nil
}

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

func loadJSONFromFile(file string, dest interface{}) error {
	cachedFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	stuff, err := ioutil.ReadAll(cachedFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(stuff, dest)
	if _, iserr := err.(*json.SyntaxError); iserr {
		return err
	}
	if err != nil {
		panic(err)
	}

	return nil
}

func mergeMultipleSources(s1 *[]Spell, s2 *[]Spell) *[]Spell {
	arr1, arr2 := *s1, *s2

	if arr1 == nil {
		return &arr2
	}
	if arr2 == nil {
		return &arr1
	}

	var final []Spell
	var i1, i2 int
	// len is expensive so it is calculated since it doesnt change
	len1, len2 := len(arr1), len(arr2)

	for i1 < len1 && i2 < len2 {
		spell1 := arr1[len1-1]
		spell2 := arr2[len2-1]
		if i1 < len1 {
			spell1 = arr1[i1]
		}
		if i2 < len2 {
			spell2 = arr2[i2]
		}

		switch {
		case spell1.Index == spell2.Index:
			final = append(final, spell1)
			i1++
			i2++
		// spell1 comes before spell2
		case spell1.Index < spell2.Index:
			final = append(final, spell1)
			i1++

		// spell1 comes after spell2
		case spell1.Index > spell2.Index:
			final = append(final, spell2)
			i2++
		}
	}

	return &final
}
