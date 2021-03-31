package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// LoadAllData is the entry point of the of the data loading sequence
//
// Loads all spell data that it can fetch. Checks if API spells are cached
// and fetches them if available. If cached content doesn't exist, it fetches
// them from an API and caches them. Then it checks for user files at the
// specific location and merges them with API spells.
// When fetching of the spells is done, the result is passed through the
// spellChan channel. StatusChan reports the current state of the data
// fetching. IsForce signalises whether cached results should be refetched
// KNOWN ISSUES:
// - if the spell parsing from local files fails, it will panic
// - if fetching from API is interrupted, it will panic
func fetchAllData(spellChan chan Spells, statusChan chan string, isForce bool) {
	defer func() {
		//close(spellChan)
		statusChan <- "Done"
	}()
	statusChan <- "Loading spells..."

	// check if these exist, make them it they dont
	readyDir(CacheDir)
	readyDir(LocalDir)

	tempSpellChan := make(chan Spells, 2)

	custom := NewSpellFetcher("custom spells", LocalDir+"/spells.json", "")
	api := NewSpellFetcher("remote spells", CacheDir+"/spells.json", "https://api.open5e.com/spells/")

	go custom.FetchSpells(tempSpellChan, statusChan, isForce)
	go api.FetchSpells(tempSpellChan, statusChan, isForce)

	// Synchronise fetching of custom and API spells
	<-tempSpellChan
	<-tempSpellChan

	allSpells := mergeMultipleSources(custom.data, api.data)

	spellChan <- *allSpells
}

// Check if file/dir exists, true if it does, false if it doesn't
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
func mergeMultipleSources(s1 *Spells, s2 *Spells) *Spells {
	arr1, arr2 := *s1, *s2

	if arr1 == nil || arr1.Len() == 0 {
		return &arr2
	}
	if arr2 == nil || arr2.Len() == 0 {
		return &arr1
	}

	var final Spells
	var i1, i2 int
	// lengths are precalculated for a small performance gain
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
