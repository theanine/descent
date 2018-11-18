package main

import (
	"fmt"
	"os"
)

const version = "v2.5.1.181116"
const downloadEnabled = false

var archetypes = []string{
	"Healer",
	"Mage",
	"Scout",
	"Warrior",
}

var expansions = map[string]string{
	"bonds of the wild":             "BotW",
	"crown of destiny":              "CoD",
	"crusade of the forgotten":      "CotF",
	"guardians of deephall":         "GoD",
	"labyrinth of ruin":             "LoR",
	"lair of the wyrm":              "LotW",
	"manor of ravens":               "MoR",
	"mists of bilehall":             "MoB",
	"oath of the outcast":           "OotO",
	"raythen lieutenant pack":       "LP",
	"second edition base game":      "2E",
	"second edition conversion kit": "1E",
	"serena lieutenant pack":        "LP",
	"shadow of nerekhall":           "SoN",
	"shards of everdark":            "SoE",
	"stewards of the secret":        "SotS",
	"the chains that rust":          "TCtR",
	"the trollfens":                 "TT",
	"treaty of champions":           "ToC",
	"visions of dawn":               "VoD",
}

func usage() {
	fmt.Println(`Usage: heroes <OPTION>

Options:
  -h, --heroes     Generate heroes.html
  -c, --classes    Generate classes.html`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 1 || len(os.Args) > 2 {
		usage()
	}
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--heroes" {
		heroesGen()
	}
	if len(os.Args) == 1 || os.Args[1] == "-c" || os.Args[1] == "--classes" {
		classesGen()
	}
}
