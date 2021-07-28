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

Format of the spells in the spells.json file (**no fields are required except the index field**):  
```jsonc
[
    {
        // index is looked when sorting spells, it is most commonly lowercase spell
        // name joined with dashes, but doesn't have to be
        "index": "magical-spell",
        "name": "Magical Spell",
        "desc": [ // description of the spell. It is an array of strings and each string is a paragraph
            "This spell does magical magic.",
            "The spell ends when you say \"hocus pocus\".",
            "If you cast it again, you will spawn a truckload of kittens."
        ],
        "higher_level": [ // it is also an array of strings like description
            "When you casting this spell using a spell slot of 3rd level and blah blah."
        ],
        "range": "60 ft",
        "components": [ // array of components. Items be either "V", "S" or "M"
            "V",
            "S",
            "M"
        ],
        "material": "Cat fur and 100gp worth of cat food", // only valid when material component is present
        "ritual": false,
        "duration": "1 hour",
        "concentration": false,
        "casting_time": "1 action",
        "level": 2, // must be an integer. For cantrips, the levels is 0
        "school": { 
            "name": "Abjuration",
        },
        "classes": [
            {
                "name": "Cleric",
            }
        ],
        "subclasses": [
            {
                "name": "Lore",
            }
        ],
    },
    {
        "index": "slighty-less-magical-spell"
        // ... 
        // rest of the spell(s)
        // ...
    }
] // EOF
```

## How do I run this?
1) Install [Golang](https://golang.org/)
2) `git clone https://github.com/spinzed/litch.git`
3) `go run .`
    - to install it, you can run `go install` but make sure that `$GOPATH\bin` is in your PATH
    - on Linux, you can run `make` to install the app and `make uninstall` to uninstall it

## Aditional Notes
If you have any suggestions, please let me know. Although I planned to use this myself, I doubt that I'll be able to do so in foreseeable future.
With that in mind, I hope that this app can help someone else with their DMing (:
