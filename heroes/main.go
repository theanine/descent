package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/theanine/utils"
)

var archetypes = map[string]struct{}{
	"Warrior": {},
	"Healer":  {},
	"Scout":   {},
	"Mage":    {},
}

var expansions = map[string]struct{}{
	"Bonds of the Wild":             {},
	"Crown of Destiny":              {},
	"Crusade of the Forgotten":      {},
	"Guardians of Deephall":         {},
	"Labyrinth of Ruin":             {},
	"Lair of the Wyrm":              {},
	"Manor of Ravens":               {},
	"Oath of the Outcast":           {},
	"Raythen Lieutenant Pack":       {},
	"Second Edition Base Game":      {},
	"Second Edition Conversion Kit": {},
	"Serena Lieutenant Pack":        {},
	"Shadow of Nerekhall":           {},
	"Shards of Everdark":            {},
	"Stewards of the Secret":        {},
	"The Trollfens":                 {},
	"Treaty of Champions":           {},
	"Visions of Dawn":               {},
}

type hero struct {
	name        string
	url         string
	archetype   string
	expansion   string
	description string
	img         string // html (may contain img tags)
	speed       int
	health      int
	stamina     int
	defense     string // html (may contain img tags)
	might       int
	knowledge   int
	willpower   int
	awareness   int
	ability     string // html (may contain img tags)
	heroic      string // html (may contain img tags)
	quote       string
}

var heroes []hero

func usage() {
	panic("Usage: heroes")
}

func (h *hero) print() {
	fmt.Printf(`
[%s]
%s
%s
%s
%s
%d
%d
%d
%s
%d
%d
%d
%d
%s
%s
%s
`, h.name, h.archetype, h.expansion, h.description, h.img, h.speed, h.health, h.stamina, h.defense,
		h.might, h.knowledge, h.willpower, h.awareness, h.ability, h.heroic, h.quote)
}

func downloadImages() {
	for _, h := range heroes {
		var conf utils.Config
		conf.Url = h.img
		conf.Outfile = strings.Split(h.img, "/")[7]
		utils.Wget(conf)
	}
}

func printTable() {
	fmt.Println(`<html><head>`)
	fmt.Println(`<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>`)
	fmt.Println(`<link rel="stylesheet" type="text/css" href="heroes.css">`)
	fmt.Println(`</head><body><table id="heroTable"><thead><tr>`)
	fmt.Printf("<th class=\"expansion\">Expansion</th>\n")
	fmt.Printf("<th class=\"hero\">Hero</th>\n")
	fmt.Printf("<th class=\"image\">Image</th>\n")
	fmt.Printf("<th class=\"archetype\">Type</th>\n")
	fmt.Printf("<th class=\"num speed\"><img src=\"attributes/speed.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num health\"><img src=\"attributes/health.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num stamina\"><img src=\"attributes/fatigue.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num die\"><img src=\"attributes/defense.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num might\"><img src=\"attributes/might.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num willpower\"><img src=\"attributes/willpower.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num knowledge\"><img src=\"attributes/knowledge.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"num awareness\"><img src=\"attributes/awareness.png\" class=\"header\"></th>\n")
	fmt.Printf("<th class=\"ability\">Hero Ability</th>\n")
	fmt.Printf("<th class=\"heroic\">Heroic Feat</th>\n")
	fmt.Printf("<th class=\"quote\">Quote</th>\n")
	fmt.Println("</tr></thead><tbody>\n")
	fmt.Println(`<tr><td><td><td><td>`)
	fmt.Println(`<div id="selectType">`)
	fmt.Println(`<select onclick="showHideRows(this)">`)
	fmt.Println(`<option value=""></option>`)
	fmt.Println(`<option value="healer">Healer</option>`)
	fmt.Println(`<option value="mage">Mage</option>`)
	fmt.Println(`<option value="scout">Scout</option>`)
	fmt.Println(`<option value="warrior">Warrior</option>`)
	fmt.Println(`</select></div><td></tr>`)

	for _, h := range heroes {
		ck := false
		if h.expansion == "Second Edition Conversion Kit" {
			ck = true
		}
		image := "./heroes-small/" + strings.Replace(h.name, " the ", " The ", -1)
		image = strings.Replace(image, " of ", " Of ", -1)
		image = strings.Replace(image, " and ", " And ", -1)
		image = strings.Replace(image, " ", "", -1)
		if ck {
			image += "CK"
			h.name += " (CK)"
		}
		fmt.Println("<tr>")
		expansionImage := "expansions/" + strings.Replace(h.expansion, " ", "_", -1) + ".svg"
		if _, err := os.Stat(expansionImage); !os.IsNotExist(err) {
			fmt.Printf("<td class=\"expansion\"><img src=\"%s\"></td>\n", "expansions/"+strings.Replace(h.expansion, " ", "_", -1)+".svg")
		} else if h.expansion == "Second Edition Base Game" {
			fmt.Println("<td class=\"expansion\">2E</td>")
		} else if h.expansion == "Second Edition Conversion Kit" {
			fmt.Println("<td class=\"expansion\">1E</td>")
		} else {
			fmt.Println("<td class=\"expansion\"></td>")
		}
		fmt.Printf("<td class=\"hero\"><a href=\"%s\">%s</a></td>\n", h.url, h.name)
		fmt.Printf("<td class=\"image\"><img src=\"%s\"></td>\n", image+".png")
		fmt.Printf("<td class=\"archetype\"><img src=\"%s\" class=\"archetype\"></td>\n", "classes/"+strings.ToLower(h.archetype)+".png")
		fmt.Printf("<td class=\"num speed\">%d</td>\n", h.speed)
		fmt.Printf("<td class=\"num health\">%d</td>\n", h.health)
		fmt.Printf("<td class=\"num stamina\">%d</td>\n", h.stamina)

		die := strings.ToLower(h.defense)
		if die == "1 gray" || die == "1 grey" {
			die = "attributes/whitedie.png"
		} else if die == "1 black" {
			die = "attributes/blackdie.png"
		} else if die == "1 brown" {
			die = "attributes/browndie.png"
		} else {
			panic(h.defense)
		}
		fmt.Printf("<td class=\"num die\"><img src=\"%s\" class=\"die\"></td>\n", die)
		fmt.Printf("<td class=\"num might\">%d</td>\n", h.might)
		fmt.Printf("<td class=\"num willpower\">%d</td>\n", h.willpower)
		fmt.Printf("<td class=\"num knowledge\">%d</td>\n", h.knowledge)
		fmt.Printf("<td class=\"num awareness\">%d</td>\n", h.awareness)
		fmt.Printf("<td class=\"ability\">%s</td>\n", h.ability)
		fmt.Printf("<td class=\"heroic\">%s</td>\n", h.heroic)
		fmt.Printf("<td class=\"quote\">%s</td>\n", h.quote)
		fmt.Println("</tr>\n")
	}
	fmt.Println("</tbody></table></body></html>")
	fmt.Println(`<script type="text/javascript" src="heroes.js"></script>`)
}

// func replaceHearts(html string) string {
// 	// heart := `<a href="https://vignette.wikia.nocookie.net/descent2e/images/d/d9/Heart.png/revision/latest?cb=20121016005115" class="image image-thumbnail" title="Heart"><img src="https://vignette.wikia.nocookie.net/descent2e/images/d/d9/Heart.png/revision/latest/scale-to-width-down/15?cb=20121016005115" alt="Heart" class="" style="vertical-align: sub" data-image-key="Heart.png" data-image-name="Heart.png" width="15" height="14"/></a>`
// 	// newheart := `<img src="attributes/health.png" alt="Heart" class="" width="15" height="14" style="vertical-align: sub">`
// 	heart := `https://vignette.wikia.nocookie.net/descent2e/images/d/d9/Heart.png/revision/latest/scale-to-width-down/15?cb=20121016005115`
// 	newheart := `attributes/health.png`
// 	return strings.Replace(html, heart, newheart, -1)
// }

// func replaceFatigue(html string) string {
// 	// fatigue := `<a href="https://vignette.wikia.nocookie.net/descent2e/images/b/b4/Fatigue.png/revision/latest?cb=20121016005054" class="image image-thumbnail" title="Fatigue"><img src="https://vignette.wikia.nocookie.net/descent2e/images/b/b4/Fatigue.png/revision/latest/scale-to-width-down/10?cb=20121016005054" alt="Fatigue" class="" style="vertical-align: baseline" data-image-key="Fatigue.png" data-image-name="Fatigue.png" width="10" height="15"/></a>`
// 	// newfatigue := `<img src="attributes/fatigue.png" alt="Fatigue" class="" width="10" height="15" style="vertical-align: baseline">`
// 	fatigue := `https://vignette.wikia.nocookie.net/descent2e/images/b/b4/Fatigue.png/revision/latest/scale-to-width-down/10?cb=20121016005054`
// 	newfatigue := `attributes/fatigue.png`
// 	html = strings.Replace(html, fatigue, newfatigue, -1)
// 	fatigue = `https://vignette.wikia.nocookie.net/descent2e/images/b/b4/Fatigue.png/revision/latest?cb=20121016005054`
// 	return strings.Replace(html, fatigue, newfatigue, -1)
// }

// func replaceSurges(html string) string {
// 	// surge := `<a href="/wiki/Surge" class="image image-thumbnail link-internal" title="Surge"><img src="https://vignette.wikia.nocookie.net/descent2e/images/1/1a/Surge.png/revision/latest/scale-to-width-down/16?cb=20121015120900" alt="Surge" class="" style="vertical-align: text-bottom" data-image-key="Surge.png" data-image-name="Surge.png" width="16" height="16"/></a>`
// 	// newsurge := `<img src="attributes/surge.png" alt="Surge" class="" width="16" height="16" style="vertical-align: text-bottom">`
// 	surge := `https://vignette.wikia.nocookie.net/descent2e/images/1/1a/Surge.png/revision/latest/scale-to-width-down/16?cb=20121015120900`
// 	newsurge := `attributes/surge.png`
// 	return strings.Replace(html, surge, newsurge, -1)
// }

// func replaceShields(html string) string {
// 	// shield := `<a href="/wiki/Shield" class="image image-thumbnail link-internal" title="Shield"><img src="https://vignette.wikia.nocookie.net/descent2e/images/1/1a/Shield.png/revision/latest/scale-to-width-down/16?cb=20121015120900" alt="Shield" class="" style="vertical-align: text-bottom" data-image-key="Shield.png" data-image-name="Shield.png" width="16" height="16"/></a>`
// 	// newshield := `<img src="attributes/shield.png" alt="Shield" class="" width="16" height="16" style="vertical-align: text-bottom">`
// 	shield := `https://vignette.wikia.nocookie.net/descent2e/images/1/1a/Shield.png/revision/latest/scale-to-width-down/16?cb=20121015120900`
// 	newshield := `attributes/shield.png`
// 	html = strings.Replace(html, shield, newshield, -1)
// 	shield = `https://vignette.wikia.nocookie.net/descent2e/images/c/cf/Shield.png/revision/latest?cb=20121021145103`
// 	newshield = `attributes/shield.png`
// 	return strings.Replace(html, shield, newshield, -1)
// }

// func replaceActions(html string) string {
// 	action := `https://vignette.wikia.nocookie.net/descent2e/images/c/c2/Action.png/revision/latest?cb=20121015121410`
// 	newaction := `attributes/action.png`
// 	return strings.Replace(html, action, newaction, -1)
// }

// func replaceWillpower(html string) string {
// 	willpower := `https://vignette.wikia.nocookie.net/descent2e/images/8/88/Willpower.png/revision/latest/scale-to-width-down/15?cb=20121013062622`
// 	newwillpower := `attributes/willpower.png`
// 	return strings.Replace(html, willpower, newwillpower, -1)
// }

// func replaceKnowledge(html string) string {
// 	knowledge := `https://vignette.wikia.nocookie.net/descent2e/images/a/ad/Knowledge.png/revision/latest/scale-to-width-down/32?cb=20121013062540`
// 	newknowledge := `attributes/knowledge.png`
// 	return strings.Replace(html, knowledge, newknowledge, -1)
// }

// func replaceAwareness(html string) string {
// 	awareness := `https://vignette.wikia.nocookie.net/descent2e/images/f/f5/Awareness.png/revision/latest/scale-to-width-down/20?cb=20121013062510`
// 	newawareness := `attributes/awareness.png`
// 	return strings.Replace(html, awareness, newawareness, -1)
// }

func iconHelper(src string, img *goquery.Selection) {
	if strings.Contains(src, "Heart.png") {
		img.SetAttr("src", "attributes/health.png")
	} else if strings.Contains(src, "Fatigue.png") {
		img.SetAttr("src", "attributes/fatigue.png")
	} else if strings.Contains(src, "Surge.png") {
		img.SetAttr("src", "attributes/surge.png")
	} else if strings.Contains(src, "Shield.png") {
		img.SetAttr("src", "attributes/defense.png")
	} else if strings.Contains(src, "Action.png") {
		img.SetAttr("src", "attributes/action.png")
	} else if strings.Contains(src, "Willpower.png") {
		img.SetAttr("src", "attributes/willpower.png")
	} else if strings.Contains(src, "Knowledge.png") {
		img.SetAttr("src", "attributes/knowledge.png")
	} else if strings.Contains(src, "Awareness.png") {
		img.SetAttr("src", "attributes/awareness.png")
	}
}

func replaceIcons(td *goquery.Selection) *goquery.Selection {
	td.Find("img").Each(func(i int, img *goquery.Selection) {
		if src, ok := img.Attr("src"); ok {
			iconHelper(src, img)
		}
		if src, ok := img.Attr("data-src"); ok {
			iconHelper(src, img)
		}
	})
	// html = replaceHearts(html)
	// html = replaceFatigue(html)
	// html = replaceSurges(html)
	// html = replaceShields(html)
	// html = replaceActions(html)
	// html = replaceWillpower(html)
	// html = replaceKnowledge(html)
	// html = replaceAwareness(html)
	return td
}

func tdToAbility(td *goquery.Selection) string {
	td = replaceIcons(td)
	s, _ := td.Html()
	ability := strings.TrimSpace(s)
	if ability == "" {
		ability = strings.TrimSpace(td.Text())
	}
	return ability
}

func tdToHeroic(td *goquery.Selection) string {
	td = replaceIcons(td)
	heroic := strings.TrimSpace(td.Text())
	return heroic
}

func heroFromTd(td *goquery.Selection) hero {
	var h hero
	h.img, _ = td.Eq(0).Find("a").Attr("href")
	h.speed = utils.MustAtoi(td.Eq(2).Text())
	h.health = utils.MustAtoi(td.Eq(5).Text())
	h.ability = tdToAbility(td.Eq(6))
	h.stamina = utils.MustAtoi(td.Eq(8).Text())
	h.defense = strings.TrimSpace(td.Eq(10).Text())
	h.might = utils.MustAtoi(td.Eq(12).Text())
	h.knowledge = utils.MustAtoi(td.Eq(15).Text())
	h.heroic = tdToAbility(td.Eq(16))
	h.willpower = utils.MustAtoi(td.Eq(17).Text())
	h.awareness = utils.MustAtoi(td.Eq(20).Text())
	h.quote = strings.TrimSpace(td.Eq(21).Text())
	return h
}

func main() {
	if len(os.Args) != 1 {
		usage()
	}

	doc, err := goquery.NewDocument("http://descent2e.wikia.com/wiki/Hero")
	if err != nil {
		panic(fmt.Sprintf("error on parsing: %s", err))
	}

	heroesTable := doc.Find(".wikitable").First()
	heroesTable.Find("tr").Each(func(i int, s1 *goquery.Selection) {
		elements := s1.Find("td")

		var meta hero
		meta.name = strings.TrimSpace(elements.Eq(0).Text())
		meta.archetype = strings.TrimSpace(elements.Eq(1).Text())
		meta.expansion = strings.TrimSpace(elements.Eq(2).Text())
		meta.description = strings.TrimSpace(elements.Eq(3).Text())

		aTag := elements.Eq(0).Find("a")
		if heroUrl, ok := aTag.Attr("href"); ok {
			heroUrl = "http://descent2e.wikia.com" + heroUrl
			doc, err := goquery.NewDocument(heroUrl)
			if err != nil {
				panic(fmt.Sprintf("error on parsing: %s", err))
			}

			characters := doc.Find("tbody")

			characters.Find("a").Each(func(i int, a *goquery.Selection) {
				if class, ok := a.Attr("class"); ok {
					if strings.Contains(class, "image image-thumbnail") {
						s, _ := a.Html()
						a = a.ReplaceWithHtml(s)
					}
				}
				if href, ok := a.Attr("href"); ok {
					if strings.Contains(href, "/wiki/") {
						a.SetAttr("href", "http://descent2e.wikia.com"+href)
					}
				}
			})

			if characters.Length() > 0 {
				base := characters.First().Find("td")
				h := heroFromTd(base)
				h.url = heroUrl
				h.name = meta.name
				h.archetype = meta.archetype
				h.expansion = meta.expansion
				h.description = meta.description
				heroes = append(heroes, h)
			}

			if characters.Length() > 1 {
				ck := characters.First().Find("td")
				h := heroFromTd(ck)
				h.url = heroUrl
				h.name = meta.name
				h.archetype = meta.archetype
				h.expansion = "Second Edition Conversion Kit"
				h.description = meta.description
				heroes = append(heroes, h)
			}
		}
	})

	// for _, h := range heroes {
	// 	h.print()
	// }

	// downloadImages()

	printTable()
}
