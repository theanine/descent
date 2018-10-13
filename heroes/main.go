package main

import (
	"fmt"
	"html"
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
	fmt.Println(`</head><body onload="showHideRows()"><table id="heroTable"><thead><tr>`)
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
	fmt.Println("</tr></thead><tbody>\n")
	fmt.Println(`<tr><td>`)
	fmt.Println(`<td><div><select id="selectCK" onclick="showHideRows()">`)
	fmt.Println(`<option value=""></option>`)
	fmt.Println(`<option value="+ck">+CK</option>`)
	fmt.Println(`<option value="-ck">-CK</option>`)
	fmt.Println(`</select></div>`)
	fmt.Println(`<td>`)
	fmt.Println(`<td><div><select id="selectClass" onclick="showHideRows()">`)
	fmt.Println(`<option value=""></option>`)
	fmt.Println(`<option value="healer">Healer</option>`)
	fmt.Println(`<option value="mage">Mage</option>`)
	fmt.Println(`<option value="scout">Scout</option>`)
	fmt.Println(`<option value="warrior">Warrior</option>`)
	fmt.Println(`</select></div>`)
	fmt.Println(`<td><td><td>`)
	fmt.Println(`<td><div><select id="selectDefense" onclick="showHideRows()">`)
	fmt.Println(`<option value=""></option>`)
	fmt.Println(`<option value="brown">b</option>`)
	fmt.Println(`<option value="white">W</option>`)
	fmt.Println(`<option value="black">B</option>`)
	fmt.Println(`</select></div></tr>`)

	for i, h := range heroes {
		ckStr := ""
		ck := false
		if h.expansion == "Second Edition Conversion Kit" {
			ck = true
		}
		image := "./heroes-small/" + strings.Replace(h.name, " the ", " The ", -1)
		image = strings.Replace(image, " of ", " Of ", -1)
		image = strings.Replace(image, " and ", " And ", -1)
		image = strings.Replace(image, " ", "", -1)
		if ck {
			ckStr = "+ck"
			if h.name == heroes[i-1].name {
				ckStr += " ck+"
			}
			image += "CK"
			h.name += " (CK)"
		} else {
			ckStr = "-ck"
		}

		die := strings.ToLower(h.defense)
		if die == "1 gray" || die == "1 grey" {
			die = "white"
		} else if die == "1 black" {
			die = "black"
		} else if die == "1 brown" {
			die = "brown"
		} else {
			panic(h.defense)
		}

		fmt.Printf("<tr class=\"%s %s %s\">\n", strings.ToLower(h.archetype), die, ckStr)
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
		// fmt.Printf("<td class=\"image\"><img src=\"%s\"></td>\n", image+".png")
		fmt.Printf("<td class=\"image\">")
		fmt.Printf("<span title=\"%s\">", html.EscapeString(h.quote))
		fmt.Printf("<img src=\"%s\">", image+".png")
		fmt.Printf("</span></td>\n")
		fmt.Printf("<td class=\"archetype\"><img src=\"%s\" class=\"archetype\"></td>\n", "classes/"+strings.ToLower(h.archetype)+".png")
		fmt.Printf("<td class=\"num speed\">%d</td>\n", h.speed)
		fmt.Printf("<td class=\"num health\">%d</td>\n", h.health)
		fmt.Printf("<td class=\"num stamina\">%d</td>\n", h.stamina)
		fmt.Printf("<td class=\"num die\"><img src=\"%s\" class=\"die\"></td>\n", "attributes/"+die+"die.png")
		fmt.Printf("<td class=\"num might\">%d</td>\n", h.might)
		fmt.Printf("<td class=\"num willpower\">%d</td>\n", h.willpower)
		fmt.Printf("<td class=\"num knowledge\">%d</td>\n", h.knowledge)
		fmt.Printf("<td class=\"num awareness\">%d</td>\n", h.awareness)
		fmt.Printf("<td class=\"ability\">%s</td>\n", h.ability)
		fmt.Printf("<td class=\"heroic\">%s</td>\n", h.heroic)
		fmt.Println("</tr>\n")
	}
	fmt.Println("</tbody></table></body></html>")
	fmt.Println(`<script type="text/javascript" src="heroes.js"></script>`)
}

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
				ck := characters.Eq(1).Find("td")
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
