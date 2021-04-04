package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
)

// TODO: refractor the registering of events, possibly exporting them as consts

// SpellFetcher is used for fetching data from a file or an API (or both). It
// can also cache data, so once data is fetche from an external API, it doesn't
// have to be readownloaded again. Fetching from API (if it exists) can be
// forced. Statuses and errors are reported via EventRegister.
type SpellFetcher struct {
	// identifier of the fetched data, could be per example "custom Spells"
	name string
	// path of the data in local storage where it could/should be
	// can be "" if we dont want to cache spells from an API
	local string
	// URL of the API from which the spells will be fetched
	// can be "" if the spells can be found only locally
	apiURL string
	// data that will be set when Fetch method is called. Afet the method
	// is called, this must not be nil
	data *Spells
	// optional EventRegister where events will be reported
	evtReg *EventRegister
}

func NewSpellFetcher(name, local, apiUrl string, e *EventRegister) *SpellFetcher {
	data := Spells{}
	return &SpellFetcher{name, local, apiUrl, &data, e}
}

// Fetch the spells from local storage if it exists, if not, fetch them from
// an API if it was passed in. If not, the final result will be an empty slice.
// This function was designed to be run in a separate goroutine. The fetched
// spells can be accesed in two ways: through a channel, or via a field in the
// SpellFetcher struct. If isForce is true, then it will always attempt to refetch
// the data.
// DISCLAIMER: no mutex was put for the safety of SpellFetcher.data because it
// is not expected that the field is acessed at once in any point of time
func (s *SpellFetcher) FetchSpells(spellCh chan<- Spells, isForce bool) {
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
			formatedErr := fmt.Sprintf("error while parsing json: %v", err)
			status := fmt.Sprintf("Could not parse %v from local file, check logs", s.name)
			s.evtReg.Register(EventErr, formatedErr, status)
			return
		}
		info := fmt.Sprintf("Loaded offline cache for %v", s.name)
		s.evtReg.Register(EventInfo, info, info)
		return
	}
	// fetch the data from the API if it exists and nothing is cached
	if s.apiURL != "" {
		if isForce {
			info := fmt.Sprintf("Refetching %v from remote API...", s.name)
			s.evtReg.Register(EventInfo, info, info)
		} else {
			info := fmt.Sprintf("Could not find %v locally, fetching from remote API...", s.name)
			s.evtReg.Register(EventInfo, info, info)
		}
		if err := fetchSpells(s.apiURL, s.data); err != nil {
			formatedErr := fmt.Sprintf("error while fetching api: %v", err)
			status := fmt.Sprintf("Could not fetch %v from remote API", s.name)
			s.evtReg.Register(EventErr, formatedErr, status)
		}
		info := fmt.Sprintf("Fetched online spells for %v", s.name)
		s.evtReg.Register(EventInfo, info, info)
		if s.local != "" {
			info := fmt.Sprintf("Caching %v...", s.name)
			s.evtReg.Register(EventInfo, info, info)
			go s.Cache()
		}
		return
	}
	info := fmt.Sprintf("Could not fetch nor find cached data for %s", s.name)
	s.evtReg.Register(EventInfo, info, "")
}

// Cache the data that is currently held in s.data. It will panic
// if s.data is nil which should never happen
func (s *SpellFetcher) Cache() {
	if s.data == nil {
		err := "no data to cache"
		s.evtReg.Register(EventErr, err, err)
	}

	// marshall the fetched spells in human-readable format
	cacheReady, err := json.MarshalIndent(s.data, "", "    ")
	if err != nil {
		formatedErr := fmt.Sprintf("error while marshaling %s: %s", s.name, err)
		status := fmt.Sprintf("Error while trying to cache %s, check logs", err)
		s.evtReg.Register(EventErr, formatedErr, status)
	}

	// make sure that the dir of the file is created
	if err := readyDir(path.Dir(s.local)); err != nil {
		formatedErr := fmt.Sprintf("error while creating dir: %s", err)
		status := fmt.Sprintf("Error while trying to create dir for %s, check logs", err)
		s.evtReg.Register(EventErr, formatedErr, status)
	}

	file, err := os.Create(s.local)
	if err != nil {
		formatedErr := fmt.Sprintf("error while creating file: %s", err)
		status := fmt.Sprintf("Error while trying to cache %s, check logs", err)
		s.evtReg.Register(EventErr, formatedErr, status)
	}
	defer file.Close()

	// cache the spells
	if _, err := file.Write(cacheReady); err != nil {
		formatedErr := fmt.Sprintf("error while creating file: %s", err)
		status := fmt.Sprintf("Error while trying to cache %s, check logs", err)
		s.evtReg.Register(EventErr, formatedErr, status)
	}
}
