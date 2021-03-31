package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
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
	data   *Spells
}

func NewSpellFetcher(name, local, apiUrl string) *SpellFetcher {
	data := Spells{}
	return &SpellFetcher{name, local, apiUrl, &data}
}

// Fetch the spells from local storage if it exists, if not, fetch them from
// an API if it was passed in. If not, the final result will be an empty slice.
// This function was designed to be run in a separate goroutine. The fetched
// spells can be accesed in two ways: through a channel, or via a field in the
// SpellFetcher struct. If isForce is true, then it will always attempt to refetch
// the data.
// DISCLAIMER: no mutex was put for the safety of SpellFetcher.data because it
// is not expected that the field is acessed at once in any point of time
func (s *SpellFetcher) FetchSpells(spellCh chan<- Spells, statusCh chan<- string, isForce bool) {
	defer func() {
		// Sort and pass the data through the channel when the function ends.
		// If local storage nor API url weren't passed to the SpellFetched,
		// it will always pass an empty slice
		if !sort.IsSorted(s.data) {
			sort.Sort(s.data)
		}
		spellCh <- *s.data
	}()

	// fetch the files if they exist, nothing is cached already and force refetch
	// wasn't specified
	if s.local != "" && checkFile(s.local) && !isForce {
		if err := loadJSONFromFile(s.local, s.data); err != nil {
			fmt.Printf("error while parsing json: %v", err)
			statusCh <- fmt.Sprintf("Could not parse %v from local file", s.name)
			return
		}
		statusCh <- fmt.Sprintf("Loaded offline cache for %v", s.name)
		return
	}
	// fetch the data from the API if it exists and nothing is cached
	if s.apiUrl != "" {
		if isForce {
			statusCh <- fmt.Sprintf("Refetching %v from remote API...", s.name)
		} else {
			statusCh <- fmt.Sprintf("Could not find %v locally, fetching from remote API...", s.name)
		}
		if err := fetchSpells(s.apiUrl, s.data); err != nil {
			fmt.Printf("error while fetching api: %v", err)
			statusCh <- fmt.Sprintf("Could not fetch %v from remote API", s.name)
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

// Cache the data that is currently held in s.data. It will panic
// if s.data is nil which should never happen
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
