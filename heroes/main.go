package main

import (
	"fmt"
	"os"
)

const version = "v2.0.0.181110"
const downloadEnabled = false

var archetypes = []string{
	"Healer",
	"Mage",
	"Scout",
	"Warrior",
}

var expansions = map[string]string{
	"Bonds of the Wild":             "BotW",
	"Crown of Destiny":              "CoD",
	"Crusade of the Forgotten":      "CotF",
	"Guardians of Deephall":         "GoD",
	"Labyrinth of Ruin":             "LoR",
	"Lair of the Wyrm":              "LotW",
	"Manor of Ravens":               "MoR",
	"Oath of the Outcast":           "OotO",
	"Raythen Lieutenant Pack":       "LP",
	"Second Edition Base Game":      "2E",
	"Second Edition Conversion Kit": "1E",
	"Serena Lieutenant Pack":        "LP",
	"Shadow of Nerekhall":           "SoN",
	"Shards of Everdark":            "SoE",
	"Stewards of the Secret":        "SotS",
	"The Trollfens":                 "TT",
	"Treaty of Champions":           "ToC",
	"Visions of Dawn":               "VoD",
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
