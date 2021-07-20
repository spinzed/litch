# Litch
A terminal app for viewing D&amp;D 5e spells. It extensively makes use of concurrency and caching.

## What does this app do?
It fetches spells from the public API [https://api.open5e.com](https://api.open5e.com) and **caches them for offline use**.
There are other public D&amp;D 5e APIs, but I chose this one because it was easiest to fetch all data from once.
(which is a priority so when spells are once downloaded, they can be viewed offline). The spells are cached in JSON format at the following location:
`$HOME/.config/litch/cache` on Linux, `$HOME/Library/Application Support/litch/cache` on macOS and `%AppData%/litch/cache` on Windows.  

This app also supports adding custom spells. Once they are detected and loaded, they will be mixed with the ones fetched from the API.
To add custom spells, they must be added to the following location: `$HOME/.config/litch/local/spells.json` on Linux,
`$HOME/Library/Application Support/litch/local/spells.json` on macOS and `%AppData%/litch/local/spells.json` on Windows.  

Format of the spells in the spells.json file (no fields are required except the name of the spell):  
[spel format will be added soon]

## How do I run this?
1) Install [Golang](https://golang.org/)
2) `git clone https://github.com/spinzed/litch.git`
3) `make` on Linux or `go run .` on any system.
    - to uninstall on Linux: `make uninstall`

## Aditional Notes
If you have any suggestions, please let me know. Although I planned to use this myself, I doubt that I'll be able to do so in foreseeable future.
With that in mind, I hope that this app can help someone else with their DMing (:
