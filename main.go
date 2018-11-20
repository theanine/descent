package main

import (
	"fmt"
	"os"
)

const version = "v3.1.2.181119"
const downloadEnabled = false

var archetypes = []string{
	"Healer",
	"Mage",
	"Scout",
	"Warrior",
}

var expCodes = map[string]string{
	"DJ01": "Second Edition Base Game",
	"DJ02": "Second Edition Conversion Kit",
	"DJ03": "Lair of the Wyrm",
	"DJ04": "Labyrinth of Ruin",
	"DJ05": "The Trollfens",
	"DJ06": "UNKNOWN",
	"DJ07": "Shadow of Nerekhall",
	"DJ08": "Dice Pack",
	"DJ09": "Splig Lieutenant Pack",
	"DJ10": "Belthir Lieutenant Pack",
	"DJ11": "Zachareth Lieutenant Pack",
	"DJ12": "Alric Farrow Lieutenant Pack",
	"DJ13": "Merick Farrow Lieutenant Pack",
	"DJ14": "Eliza Farrow Lieutenant Pack",
	"DJ15": "Valyndra Lieutenant Pack",
	"DJ16": "Raythen Lieutenant Pack",
	"DJ17": "Serena Lieutenant Pack",
	"DJ18": "Ariad Lieutenant Pack",
	"DJ19": "Queen Ariad Lieutenant Pack",
	"DJ20": "Bol'Goreth Lieutenant Pack",
	"DJ21": "Manor of Ravens",
	"DJ22": "Rylan Olliven Lieutenant Pack",
	"DJ23": "Verminous Lieutenant Pack",
	"DJ24": "Tristayne Olliven Lieutenant Pack",
	"DJ25": "Gargan Mirklace Lieutenant Pack",
	"DJ26": "Oath of the Outcast",
	"DJ27": "Crown of Destiny",
	"DJ28": "Crusade of the Forgotten",
	"DJ29": "Guardians of Deephall",
	"DJ30": "Visions of Dawn",
	"DJ31": "Bonds of the Wild",
	"DJ32": "Treaty of Champions",
	"DJ33": "Stewards of the Secret",
	"DJ34": "Shards of Everdark",
	"DJ35": "Skarn Lieutenant Pack",
	"DJ36": "Forgotten Souls",
	"DJ37": "Nature's Ire",
	"DJ38": "Dark Elements",
	"DJ39": "Heirs of Blood",
	"DJ40": "Mists of Bilehall",
	"DJ41": "Kyndrithul Lieutenant Pack",
	"DJ42": "Zarihell Lieutenant Pack",
	"DJ43": "Ardus Ix'Erebus Lieutenant Pack",
	"DJ44": "The Chains that Rust",
}

var expansions = map[string]string{
	"alric farrow lieutenant pack":      "LP",
	"ardus ix'erebus lieutenant pack":   "LP",
	"ariad lieutenant pack":             "LP",
	"belthir lieutenant pack":           "LP",
	"bol'goreth lieutenant pack":        "LP",
	"bonds of the wild":                 "BotW",
	"crown of destiny":                  "CoD",
	"crusade of the forgotten":          "CotF",
	"dark elements":                     "DE",
	"dice pack":                         "DP",
	"eliza farrow lieutenant pack":      "LP",
	"forgotten souls":                   "FS",
	"gargan mirklace lieutenant pack":   "LP",
	"guardians of deephall":             "GoD",
	"heirs of blood":                    "HoB",
	"kyndrithul lieutenant pack":        "LP",
	"labyrinth of ruin":                 "LoR",
	"lair of the wyrm":                  "LotW",
	"manor of ravens":                   "MoR",
	"merick farrow lieutenant pack":     "LP",
	"mists of bilehall":                 "MoB",
	"nature's ire":                      "NI",
	"oath of the outcast":               "OotO",
	"queen ariad lieutenant pack":       "LP",
	"raythen lieutenant pack":           "LP",
	"rylan olliven lieutenant pack":     "LP",
	"second edition base game":          "2E",
	"second edition conversion kit":     "CK",
	"serena lieutenant pack":            "LP",
	"shadow of nerekhall":               "SoN",
	"shards of everdark":                "SoE",
	"skarn lieutenant pack":             "LP",
	"splig lieutenant pack":             "LP",
	"stewards of the secret":            "SotS",
	"the chains that rust":              "TCtR",
	"the trollfens":                     "TT",
	"treaty of champions":               "ToC",
	"tristayne olliven lieutenant pack": "LP",
	"unknown":                           "?",
	"valyndra lieutenant pack":          "LP",
	"verminous lieutenant pack":         "LP",
	"visions of dawn":                   "VoD",
	"zachareth lieutenant pack":         "LP",
	"zarihell lieutenant pack":          "LP",
}

func usage() {
	fmt.Println(`Usage: heroes <OPTION>

Options:
  -h, --heroes     Generate heroes.html
  -c, --classes    Generate classes.html
  -o, --overlord   Generate overlord.html`)
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
	if len(os.Args) == 1 || os.Args[1] == "-o" || os.Args[1] == "--overlord" {
		overlordGen()
	}
}
