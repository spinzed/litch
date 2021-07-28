package main

import "fmt"

func (app *App) FetchData(isForce bool) {
    defer func () {
        if r := recover(); r != nil {
            app.eventReg.Register(EventErr, fmt.Sprintf("Could not load spells: %v", r), "Unforseen error while fetching spells, check logs")
        }
    }()
	// app.fetchLock keeps track whether data is being fetched at the moment.
	// it will only be fetched if that isn't happening already. The actual
	// fetching is in the separate method so that no routines are created
	// unnecesarily, although that is of a little importance tbh.
	if !app.fetchLock {
		go app.fetchAllData(isForce)
		// update the lock. It will be released when data is received
		// through app.dataChan channel in app.waitForData method
		app.fetchLock = true
	}
}

func (app *App) fetchAllData(isForce bool) {
	defer func() {
		app.eventReg.Register(EventInfo, "Loaded spells sucessfully", "Done")
	}()

	app.eventReg.Register(EventInfo, "Started loading spells...", "Loading spells...")

	// check if these exist, make them it they dont
	readyDir(CacheDir)
	readyDir(LocalDir)

	tempSpellChan := make(chan Spells, 2)

	custom := NewSpellFetcher("custom spells", LocalDir+"/spells.json", "", app.eventReg)
	api := NewSpellFetcher("remote spells", CacheDir+"/spells.json", "https://api.open5e.com/spells/", app.eventReg)

	custom.FetchSpells(tempSpellChan, isForce)
	api.FetchSpells(tempSpellChan, isForce)

	// Synchronise fetching of custom and API spells
	<-tempSpellChan
	<-tempSpellChan

	app.eventReg.Register(EventInfo, "Merging spells...", "")
	allSpells := mergeMultipleSources(custom.data, api.data)

	app.dataChan <- *allSpells
}

// Merge two spell lists alphabetically. It assumes that both lists are already
// ordered alphabetically. Spells in s1 have priority against those in s2,
// meaning that if spells with same indices occur in s1 and s2, the one in s1
// will end up in the final list.
func mergeMultipleSources(s1 *Spells, s2 *Spells) *Spells {
	arr1, arr2 := *s1, *s2

	if len(arr1) == 0 {
		return &arr2
	}
	if len(arr2) == 0 {
		return &arr1
	}

	var final Spells
	var i1, i2 int
	// lengths are precalculated for a small performance gain
	len1, len2 := len(arr1), len(arr2)

	for i1 < len1 || i2 < len2 {
		var spell1, spell2 Spell
		if i1 < len1 {
			spell1 = arr1[i1]
		}
		if i2 < len2 {
			spell2 = arr2[i2]
		}

		// the conditions are ugly but hey, they work
		switch {
		case spell1.Index == spell2.Index:
			final = append(final, spell1)
			i1++
			i2++
		// spell1 comes before spell2
		case spell1.Index != "" && (spell1.Index < spell2.Index || spell2.Index == ""):
			final = append(final, spell1)
			i1++

		// spell1 comes after spell2
		case spell2.Index != "" && (spell1.Index > spell2.Index || spell1.Index == ""):
			final = append(final, spell2)
			i2++
		}
	}

	return &final
}
