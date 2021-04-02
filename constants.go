package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

type InputMode int

const (
	InputNormal InputMode = iota
	InputCommand
	DefaultBgColor = tcell.ColorDefault
)

var CacheDir string = fmt.Sprintf("%s/cache", ProjectDir)
var LocalDir string = fmt.Sprintf("%s/local", ProjectDir)

var ProjectDir string = func() string {
	config, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/litch", config)
}()
