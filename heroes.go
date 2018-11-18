package main

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/theanine/utils"
)

const heroesHtml = "heroes.html"

type hero struct {
	trClass     string
	name        string
	url         string
	archetype   string
	expansion   string
	description string
	img         string // html (may contain img tags)
	die         string
	expImg      string
	ck          bool
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
	backstory   string
}

var heroes []hero

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

func heroUniqSortExps() []string {
	expMap := make(map[string]struct{})
	for _, h := range heroes {
		expMap[expansions[strings.ToLower(h.expansion)]] = struct{}{}
	}
	var exps []string
	for exp := range expMap {
		exps = append(exps, exp)
	}
	sort.Strings(exps)
	return exps
}

func outputHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Journeys in Hero Selection</title>\n")
	fmt.Fprintf(w, "<meta name=\"description\" content=\"%s\">\n", `With over 100+ heroes to choose from, it's painful to choose a character.

For owners of Descent: Journeys in the Dark (Second Edition), this Hero Selector makes the decision that much easier for newbies, casuals, and veterans.

Send your heroes to get some Coufee and they'll be adventuring in no time!`)
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.widgets.min.js\"></script>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"heroes.css?version=%s\">\n", version)
	fmt.Fprintf(w, "<link rel=\"icon\" type=\"image/png\" href=\"etc/favicon.png\">\n")
	fmt.Fprintf(w, "</head><body onload=\"onload()\">\n")

	// table
	fmt.Fprintf(w, "<table id=\"heroTable\" class=\"tablesorter\"><thead class=\"heroes\"><tr>\n")
	fmt.Fprintf(w, "<th class=\"expansion\">Exp</th>\n")
	fmt.Fprintf(w, "<th class=\"hero\">Name</th>\n")
	fmt.Fprintf(w, "<th class=\"image\">Hero</th>\n")
	fmt.Fprintf(w, "<th class=\"num speed\"><img src=\"attributes/speed.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num health\"><img src=\"attributes/health.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num stamina\"><img src=\"attributes/fatigue.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num dice\"><img src=\"attributes/defense.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num might\"><img src=\"attributes/might.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num willpower\"><img src=\"attributes/willpower.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num knowledge\"><img src=\"attributes/knowledge.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"num awareness\"><img src=\"attributes/awareness.png\" class=\"header\"></th>\n")
	fmt.Fprintf(w, "<th class=\"ability\">Hero Ability</th>\n")
	fmt.Fprintf(w, "<th class=\"heroic\">Heroic Feat</th>\n")
	fmt.Fprintf(w, "</tr>\n\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<td class=\"expansion\"><div><select id=\"selectExp\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	exps := heroUniqSortExps()
	for _, exp := range exps {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", exp, exp)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"hero\"><div><select id=\"selectCK\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	fmt.Fprintf(w, "<option value=\"override-ck\" selected=\"selected\">Override CK</option>\n")
	fmt.Fprintf(w, "<option value=\"ck-only\">CK Only</option>\n")
	fmt.Fprintf(w, "<option value=\"no-ck\">No CK</option>\n")
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"image\"><div><select id=\"selectClass\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	for _, k := range archetypes {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(k), k)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num speed\"><div><select id=\"selectSpeed\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num health\"><div><select id=\"selectHealth\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num stamina\"><div><select id=\"selectStamina\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num dice\"><div><select id=\"selectDefense\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	fmt.Fprintf(w, "<option value=\"brown\">b</option>\n")
	fmt.Fprintf(w, "<option value=\"white\">W</option>\n")
	fmt.Fprintf(w, "<option value=\"black\">B</option>\n")
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num might\"><div><select id=\"selectMight\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num willpower\"><div><select id=\"selectWillpower\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num knowledge\"><div><select id=\"selectKnowledge\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"num awareness\"><div><select id=\"selectAwareness\" onchange=\"trigger(this)\"><option value=\"\"></option></select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"ability\"><div></div></td>\n")
	fmt.Fprintf(w, "<td class=\"heroic\"><div></div></td>\n")
	fmt.Fprintf(w, "</tr></thead><tbody class=\"heroes\">\n\n")
}

func outputTable(w *bufio.Writer) {
	outputHeader(w)

	for _, h := range heroes {
		exp := expansions[strings.ToLower(h.expansion)]
		fmt.Fprintf(w, "<tr class=\"%s %s %s %s\" style=\"display:none;\">\n", strings.ToLower(h.archetype), h.die, h.trClass, exp)
		fmt.Fprintf(w, "<td class=\"expansion\">%s</td>\n", h.expImg)
		fmt.Fprintf(w, "<td class=\"hero\"><a href=\"%s\">%s</a></td>\n", h.url, h.name)
		// fmt.Fprintf(w, "<td class=\"image\"><img src=\"%s\"></td>\n", image+".png")
		fmt.Fprintf(w, "<td class=\"image\">")
		if h.backstory == "" {
			fmt.Fprintf(w, "<span title=\"%s\" class=\"quote\">", h.quote)
		} else {
			fmt.Fprintf(w, "<span title=\"%s\n\n%s\" class=\"quote\">", h.quote, h.backstory)
		}
		fmt.Fprintf(w, "<div class=\"divImage\">")
		fmt.Fprintf(w, "<img src=\"%s\" class=\"hero\">", h.img)
		fmt.Fprintf(w, "<img src=\"%s\" class=\"archetype\">", "classes/"+strings.ToLower(h.archetype)+".png")
		fmt.Fprintf(w, "%c</div>", h.archetype[0])
		fmt.Fprintf(w, "</span></td>\n")
		fmt.Fprintf(w, "<td class=\"num speed\">%d</td>\n", h.speed)
		fmt.Fprintf(w, "<td class=\"num health\">%d</td>\n", h.health)
		fmt.Fprintf(w, "<td class=\"num stamina\">%d</td>\n", h.stamina)
		die := 0
		if h.die == "brown" {
			die = 1
		} else if h.die == "white" {
			die = 2
		} else if h.die == "black" {
			die = 3
		} else {
			panic(h.die)
		}
		fmt.Fprintf(w, "<td class=\"num dice\"><img src=\"%s\" class=\"die\">", "attributes/"+h.die+"die.png")
		fmt.Fprintf(w, "<div class=\"die\">%d</div></td>\n", die)
		fmt.Fprintf(w, "<td class=\"num might\">%d</td>\n", h.might)
		fmt.Fprintf(w, "<td class=\"num willpower\">%d</td>\n", h.willpower)
		fmt.Fprintf(w, "<td class=\"num knowledge\">%d</td>\n", h.knowledge)
		fmt.Fprintf(w, "<td class=\"num awareness\">%d</td>\n", h.awareness)
		fmt.Fprintf(w, "<td class=\"ability\">%s</td>\n", h.ability)
		fmt.Fprintf(w, "<td class=\"heroic\">%s</td>\n", h.heroic)
		fmt.Fprintf(w, "</tr>\n\n")
	}

	outputFooter(w)
	w.Flush()
}

func outputFooter(w *bufio.Writer) {
	fmt.Fprintf(w, "</tbody>\n")
	fmt.Fprintf(w, "<tfoot class=\"heroes\"><tr><td class=\"donateArea\">\n")
	fmt.Fprintf(w, `<form action="https://www.paypal.com/cgi-bin/webscr" method="post" target="_top">
	<input type="hidden" name="cmd" value="_s-xclick">
	<input type="hidden" name="hosted_button_id" value="85ZEFVNEAXV3A">
	<input type="image" src="etc/donate-paypal.svg" border="0" name="submit" alt="PayPal" class="donate">
	<img alt="" border="0" src="https://www.paypalobjects.com/en_US/i/scr/pixel.gif" width="1" height="1">
	</form>`)
	fmt.Fprintf(w, `<div class="popup" onclick="myFunction()"><img src="etc/donate-bitcoin.svg" class="donate">
						<span class="popuptext" id="myPopup">Donations Address<br><br>
						<img src="etc/bitcoin.png" width=200px height=200px><br><br>
						3Q6y5d5c43Lj9maDr8dcZyXUFqxPcbBiEv</span></div>`)
	fmt.Fprintf(w, "</td><td class=\"fees\">Server Fees: $55.80/yr")
	fmt.Fprintf(w, "</td><td class=\"version\">%s</td>\n", version)
	fmt.Fprintf(w, "<td></td><td></td><td></td><td></td><td></td>\n")
	fmt.Fprintf(w, "<td></td><td></td><td></td><td></td><td></td>\n")
	fmt.Fprintf(w, "</tr></tfoot>\n")
	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "</body></html>\n")
	fmt.Fprintf(w, "<script type=\"text/javascript\" src=\"heroes.js?version=%s\"></script>\n", version)
}

func fixHeroes() {
	for i := range heroes {
		// h.ck
		heroes[i].ck = false
		if heroes[i].expansion == "Second Edition Conversion Kit" {
			heroes[i].ck = true
		}

		// h.trClass
		if heroes[i].ck {
			heroes[i].trClass = "ck-only"
		} else {
			heroes[i].trClass = "no-ck"
		}
		if i == 0 || heroes[i].name != heroes[i-1].name {
			heroes[i].trClass += " override-ck"
		}

		// h.img
		heroes[i].img = "heroes-small/" + strings.Replace(heroes[i].name, " the ", " The ", -1)
		heroes[i].img = strings.Replace(heroes[i].img, " of ", " Of ", -1)
		heroes[i].img = strings.Replace(heroes[i].img, " and ", " And ", -1)
		heroes[i].img = strings.Replace(heroes[i].img, " ", "", -1)
		if heroes[i].ck {
			heroes[i].img += "CK"
		}
		heroes[i].img += ".png"

		// h.name
		if heroes[i].ck {
			heroes[i].name += " (CK)"
		}

		// h.die
		heroes[i].die = strings.ToLower(heroes[i].defense)
		if heroes[i].die == "1 gray" || heroes[i].die == "1 grey" {
			heroes[i].die = "white"
		} else if heroes[i].die == "1 black" {
			heroes[i].die = "black"
		} else if heroes[i].die == "1 brown" {
			heroes[i].die = "brown"
		} else {
			panic(heroes[i].defense)
		}

		// h.expImg
		heroes[i].expImg = ""
		imgFile := "expansions/" + strings.Replace(heroes[i].expansion, " ", "_", -1) + ".svg"
		if strings.Contains(heroes[i].expansion, "Lieutenant Pack") {
			imgFile = "expansions/Lieutenant_Pack.png"
		}
		if _, err := os.Stat(imgFile); !os.IsNotExist(err) {
			heroes[i].expImg = fmt.Sprintf("<img src=\"%s\" class=\"expansion\">", imgFile)
		} else if heroes[i].expansion == "Second Edition Base Game" {
			heroes[i].expImg = "2E"
		} else if heroes[i].expansion == "Second Edition Conversion Kit" {
			heroes[i].expImg = "1E"
		}

		// h.quote
		heroes[i].quote = html.EscapeString(heroes[i].quote)
		heroes[i].quote = strings.Replace(heroes[i].quote, "&#34;", "", -1)
		heroes[i].quote = strings.Replace(heroes[i].quote, "“", "", -1)
		heroes[i].quote = strings.Replace(heroes[i].quote, "”", "", -1)

		// h.backstory
		heroes[i].backstory = html.EscapeString(heroes[i].backstory)
		heroes[i].backstory = strings.Replace(heroes[i].backstory, "&#34;", "", -1)
		heroes[i].backstory = strings.Replace(heroes[i].backstory, "“", "", -1)
		heroes[i].backstory = strings.Replace(heroes[i].backstory, "”", "", -1)
	}
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

func heroesGen() {
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

			foundBio := false
			doc.Find("#mw-content-text").Children().Each(func(i int, c *goquery.Selection) {
				if !foundBio && goquery.NodeName(c) == "h2" {
					if strings.TrimSpace(c.Find("span").Eq(0).Text()) == "Biography" {
						foundBio = true
					}
				} else if foundBio && goquery.NodeName(c) == "p" {
					foundBio = false
					meta.backstory = strings.TrimSpace(c.Text())
				}
			})

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
				h.backstory = meta.backstory
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
				h.backstory = meta.backstory
				heroes = append(heroes, h)
			}
		}
	})

	f, err := os.Create(heroesHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if downloadEnabled {
		downloadImages()
	}
	fixHeroes()
	outputTable(w)
}