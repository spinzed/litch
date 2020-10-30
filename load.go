package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Loads all spell data that it can fetch. Checks if API spells are cached
// and fetches them if available. If cached content doesn't exist, it fetches
// them from an API and caches them. Then it checks for user files at the
// specific location and merges them with API spells.
// KNOWN ISSUES:
// - if the spell parsing from local files fails, it will panic
// - if fetching from API is interrupted, it will panic
func loadAllData(spellChan chan []Spell, statusChan chan string) {
	defer func() {
		close(spellChan)
		statusChan <- "Done"
	}()
	statusChan <- "Loading spells..."

	config, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	readyDir(config + "/banshie/cache")
	readyDir(config + "/banshie/local")

	// get the spell API data, cached or fetched
	// TODO: make this run in a goroutine
	dataAPI, err := loadAPISpells(config+"/banshie/cache/spells.json", statusChan)
	if err != nil {
		// should be handled via errchannel
		panic(err)
	}

	// if local custom spells file doesnt exist, return only API spells
	if !checkFile(config + "/banshie/local/spells.json") {
		spellChan <- *dataAPI
		return
	}
	statusChan <- "Loading custom spells..."
	var userData []Spell
	err = loadJSONFromFile(config+"/banshie/local/spells.json", &userData)
	if err != nil {
		// should be handled via errchannel
		panic(err)
	}

	allSpells := mergeMultipleSources(&userData, dataAPI)

	spellChan <- *allSpells
}

// Load API spells, attempt to load cached data and if cache doesn't exist,
// fetch the spells from the API
func loadAPISpells(fileStr string, statusChan chan string) (*[]Spell, error) {
	// if the file exists, load it
	if checkFile(fileStr) {
		statusChan <- "Loading cached spells..."
		var data []Spell
		if err := loadJSONFromFile(fileStr, &data); err != nil {
			return nil, err
		}

		return &data, nil
	}

	statusChan <- "Fetching spells..."
	fmt.Println("spells.json doesnt exist")

	// fetch the spells
	data, err := fetchSpells()
	if err != nil {
		return nil, err
	}

	// marshall the fetched spells in human-readable format
	cached, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	// create the files where spells will be cached. It should not be created
	// because if it was, then spells would be read from it instead, so there
	// should be no error
	file, err := os.Create(fileStr)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// cache the spells
	file.Write(cached)
	return &data, nil
}

// Check if file/dir exists, true if does, false if doesnt
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

// Load JSON from file
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

// Merge two spell lists alphabetically. It assumes that both lists are already
// ordered alphabetically. Spells in s1 have priority against those in s2,
// meaning that if spells with same indices occur in s1 and s2, the one in s1
// will end up in the final list.
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
