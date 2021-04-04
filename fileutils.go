package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// This file contains functions that are wrappers around
// standard library file mainpulation function which do
// all error checking

// Check if file/dir exists,
func checkFile(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if !os.IsNotExist(err) && err != nil {
		panic(err)
	}
	return true
}

// Check if dir exists, make it if it doesn't.
func readyDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println(dir, "does not exist, creating...")
		os.MkdirAll(dir, 0755)
	}
	if !os.IsNotExist(err) && err != nil {
		return err
	}
	return nil
}

// Load JSON from file
func loadJSONFromFile(file string, dest interface{}) error {
	cachedFile, err := os.Open(file)
	if err != nil {
		return err
	}
	stuff, err := ioutil.ReadAll(cachedFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(stuff, dest); err != nil {
		return err
	}
	//if _, iserr := err.(*json.SyntaxError); iserr {
	//	return err
	//}

	return nil
}
