package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type SpellFetcher struct {
	// identifier of the fetched data, could be per example "custom Spells"
	name string
	// path of the data in local storage where it could/should be
	// can be "" if we dont want to cache spells from an API
	local string
	// URL of the API from which the spells will be fetched
	// can be "" if the spells can be found only locally
	apiUrl string
	data   *[]Spell
}

func NewSpellFetcher(name, local, apiUrl string) *SpellFetcher {
	data := []Spell{}
	return &SpellFetcher{name, local, apiUrl, &data}
}

// Fetch the spells from local storage if it exists, if not, fetch them from
// an API if it was passed in. If not, the final result will be an empty slice.
// This function was designed to be run in a separate goroutine. The fetched
// spells can be accesed in two ways: through a channel, or via a field in the
// SpellFetcher struct.
// DISCLAIMER: no mutex was put for the safety of SpellFetcher.data because it
// is not expected that the field is acessed at once in any point of time
func (s *SpellFetcher) FetchSpells(spellCh chan []Spell, statusCh chan string) {
	defer func() {
		spellCh <- *s.data
	}()

	// fetch the files locally if they exist
	if s.local != "" && checkFile(s.local) {
		if err := loadJSONFromFile(s.local, s.data); err != nil {
			log.Fatalf("error while parsing json: %v", err)
		}
		statusCh <- fmt.Sprintf("Loaded offline cache for %v", s.name)
		return
	}
	if s.apiUrl != "" {
		statusCh <- fmt.Sprintf("Could not find %v locally, fetching from API...", s.name)
		if err := fetchSpells(s.apiUrl, s.data); err != nil {
			log.Fatalf("error while fetching api: %v", err)
		}
		statusCh <- fmt.Sprintf("Fetched online spells for %v", s.name)
		if s.local != "" {
			statusCh <- fmt.Sprintf("Caching %v...", s.name)
			go s.Cache()
		}

		return
	}
	//statusCh <- fmt.Sprintf("Could not find any spells for %v", s.name)
}

func (s *SpellFetcher) Cache() {
	if s.data == nil {
		log.Fatalf("no data to cache")
		//errChan <- errors.New("no data to cache")
	}

	// marshall the fetched spells in human-readable format
	cacheReady, err := json.MarshalIndent(s.data, "", "    ")
	if err != nil {
		log.Fatalf("error while marshaling spells: %v", err)
		//errChan <- fmt.Errorf("error while marshaling spells: %v", err)
	}

	// make sure that the dir of the file is created
	readyDir(path.Dir(s.local))

	file, err := os.Create(s.local)
	if err != nil {
		log.Fatalf("error while creating file: %v", err)
		//errChan <- err
	}
	defer file.Close()

	// cache the spells
	file.Write(cacheReady)
}

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

	// check if these exist, make them it they dont
	readyDir(CacheDir)
	readyDir(LocalDir)

	tempSpellChan := make(chan []Spell, 2)

	custom := NewSpellFetcher("custom spells", LocalDir+"/spells.json", "")
	api := NewSpellFetcher("API spells", CacheDir+"/spells.json", "https://api.open5e.com/spells/")

	go custom.FetchSpells(tempSpellChan, statusChan)
	go api.FetchSpells(tempSpellChan, statusChan)

	<-tempSpellChan
	<-tempSpellChan

	allSpells := mergeMultipleSources(custom.data, api.data)

	spellChan <- *allSpells
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

	if arr1 == nil || len(arr1) == 0 {
		return &arr2
	}
	if arr2 == nil || len(arr2) == 0 {
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
