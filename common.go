package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const wikiUrlOld = "https://descent2e.wikia.com/wiki"
const wikiUrl = "http://wiki.descent-community.org"

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

var expImgs = map[string]string{
	"base game":                "",
	"bonds of the wild":        "Bonds_of_the_Wild.svg",
	"crown of destiny":         "Crown_of_Destiny.svg",
	"crusade of the forgotten": "Crusade_of_the_Forgotten.svg",
	"dark elements":            "Dark_Elements.svg",
	"forgotten souls":          "Forgotten_Souls.svg",
	"guardians of deephall":    "Guardians_of_Deephall.svg",
	"labyrinth of ruin":        "Labyrinth_of_Ruin.svg",
	"lair of the wyrm":         "Lair_of_the_Wyrm.svg",
	"lieutenant pack":          "Lieutenant_Pack.png",
	"lost legends":             "",
	"manor of ravens":          "Manor_of_Ravens.svg",
	"mists of bilehall":        "Mists_of_Bilehall.svg",
	"nature's ire":             "Nature's_Ire.svg",
	"oath of the outcast":      "Oath_of_the_Outcast.svg",
	"sands of the past":        "",
	"shadow of nerekhall":      "Shadow_of_Nerekhall.svg",
	"shards of everdark":       "Shards_of_Everdark.svg",
	"stewards of the secret":   "Stewards_of_the_Secret.svg",
	"the chains that rust":     "The_Chains_that_Rust.svg",
	"the trollfens":            "The_Trollfens.svg",
	"treaty of champions":      "Treaty_of_Champions.svg",
	"visions of dawn":          "Visions_of_Dawn.svg",
}

func replaceIcon(src string) string {
	if strings.Contains(src, "Heart.png") || strings.Contains(src, "Heart.svg") {
		return "attributes/health.svg"
	} else if strings.Contains(src, "Fatigue.png") || strings.Contains(src, "Fatigue.svg") {
		return "attributes/fatigue.svg"
	} else if strings.Contains(src, "Surge.png") || strings.Contains(src, "Surge.svg") {
		return "attributes/surge.svg"
	} else if strings.Contains(src, "Shield.png") || strings.Contains(src, "Shield.svg") {
		return "attributes/defense.svg"
	} else if strings.Contains(src, "Action.png") || strings.Contains(src, "Action.svg") {
		return "attributes/action.svg"
	} else if strings.Contains(src, "Might.png") || strings.Contains(src, "Might.svg") {
		return "attributes/might.svg"
	} else if strings.Contains(src, "Willpower.png") || strings.Contains(src, "Willpower.svg") {
		return "attributes/willpower.svg"
	} else if strings.Contains(src, "Knowledge.png") || strings.Contains(src, "Knowledge.svg") {
		return "attributes/knowledge.svg"
	} else if strings.Contains(src, "Awareness.png") || strings.Contains(src, "Awareness.svg") {
		return "attributes/awareness.svg"
	}
	return src
}

func replaceIcons(s *goquery.Selection) *goquery.Selection {
	s.Find("img").Each(func(i int, img *goquery.Selection) {
		if src, ok := img.Attr("src"); ok {
			img.SetAttr("src", replaceIcon(src))
		}
		if src, ok := img.Attr("data-src"); ok {
			img.SetAttr("src", replaceIcon(src))
		}
		alt, ok1 := img.Attr("alt")
		src, ok2 := img.Attr("src")
		if ok1 && ok2 {
			img.ReplaceWithHtml(fmt.Sprintf("<img src=\"%s\" alt=\"%s\" class=\"icon\">", src, alt))
		}
	})
	return s
}

func replaceErrata(img *string) {
	parts := strings.Split(*img, "/")
	if len(parts) < 2 {
		return
	}
	file := "errata/" + parts[1]
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return
	}
	*img = file
}

func outputLFooter(w *bufio.Writer, class string) {
	fmt.Fprintf(w, "</tbody>\n")
	fmt.Fprintf(w, "<tfoot class=\"%s\"><tr><td class=\"donateArea\">\n", class)
	fmt.Fprintf(w, `<form action="https://www.paypal.com/cgi-bin/webscr" method="post" target="_top">
	<input type="hidden" name="business" value="GAGMA422DQE9J">
	<input type="hidden" name="cmd" value="_s-xclick">
	<input type="hidden" name="hosted_button_id" value="85ZEFVNEAXV3A">
	<input type="image" src="etc/donate-paypal.svg" border="0" name="submit" alt="PayPal" class="donate">
	<img alt="" border="0" src="https://www.paypalobjects.com/en_US/i/scr/pixel.gif" width="1" height="1">
	</form>`)
	fmt.Fprintf(w, `<div class="popup" onclick="myFunction()"><img src="etc/donate-bitcoin.svg" class="donate">
						<span class="popuptext" id="myPopup">Donations Address<br><br>
						<img src="etc/bitcoin.png" width=200px height=200px><br><br>
						1KiR9rZJgSrF8xN3f2A6ZKYGFYuv7oRYxn</span></div></td>`)
	fmt.Fprintf(w, "<td class=\"support\"><img src=\"etc/support.png\"></td>\n")
}

func outputRFooter(w *bufio.Writer, cols int) {
	fmt.Fprintf(w, "<td class=\"version\">%s</td>\n", version)
	for i := 3; i < cols; i++ {
		fmt.Fprintf(w, "<td></td>")
	}
	fmt.Fprintf(w, "</tr></tfoot></table>\n")
	fmt.Fprintf(w, "</body></html>\n")
	fmt.Fprintf(w, "<script type=\"text/javascript\" src=\"heroes.js?version=%s\"></script>\n", version)
}

func outputFooter(w *bufio.Writer, class string, cols int) {
	outputLFooter(w, class)
	outputRFooter(w, cols)
}
