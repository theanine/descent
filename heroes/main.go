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

		var h hero
		h.name = strings.TrimSpace(elements.Eq(0).Text())
		h.archetype = strings.TrimSpace(elements.Eq(1).Text())
		h.expansion = strings.TrimSpace(elements.Eq(2).Text())
		h.description = strings.TrimSpace(elements.Eq(3).Text())

		aTag := elements.Eq(0).Find("a")
		if heroUrl, ok := aTag.Attr("href"); ok {
			doc, err := goquery.NewDocument("http://descent2e.wikia.com" + heroUrl)
			if err != nil {
				panic(fmt.Sprintf("error on parsing: %s", err))
			}

			characters := doc.Find("tbody")

			if characters.Length() > 0 {
				base := characters.First().Find("td")
				h.img, _ = base.Eq(0).Find("a").Attr("href")
				h.speed = utils.MustAtoi(base.Eq(2).Text())
				h.health = utils.MustAtoi(base.Eq(5).Text())
				h.ability = strings.TrimSpace(base.Eq(6).Text())
				h.stamina = utils.MustAtoi(base.Eq(8).Text())
				h.defense = strings.TrimSpace(base.Eq(10).Text())
				h.might = utils.MustAtoi(base.Eq(12).Text())
				h.knowledge = utils.MustAtoi(base.Eq(15).Text())
				h.heroic = strings.TrimSpace(base.Eq(16).Text())
				h.willpower = utils.MustAtoi(base.Eq(17).Text())
				h.awareness = utils.MustAtoi(base.Eq(20).Text())
				h.quote = strings.TrimSpace(base.Eq(21).Text())

				heroes = append(heroes, h)
				h.print()
			}

			if characters.Length() > 1 {
				var ckh hero
				ckh.name = h.name
				ckh.archetype = h.archetype
				ckh.expansion = "Second Edition Conversion Kit"
				ckh.description = h.description

				ck := characters.Eq(1).Find("td")
				ckh.img, _ = ck.Eq(0).Find("a").Attr("href")
				ckh.speed = utils.MustAtoi(ck.Eq(2).Text())
				ckh.health = utils.MustAtoi(ck.Eq(5).Text())
				ckh.ability = strings.TrimSpace(ck.Eq(6).Text())
				ckh.stamina = utils.MustAtoi(ck.Eq(8).Text())
				ckh.defense = strings.TrimSpace(ck.Eq(10).Text())
				ckh.might = utils.MustAtoi(ck.Eq(12).Text())
				ckh.knowledge = utils.MustAtoi(ck.Eq(15).Text())
				ckh.heroic = strings.TrimSpace(ck.Eq(16).Text())
				ckh.willpower = utils.MustAtoi(ck.Eq(17).Text())
				ckh.awareness = utils.MustAtoi(ck.Eq(20).Text())
				ckh.quote = strings.TrimSpace(ck.Eq(21).Text())

				heroes = append(heroes, ckh)
				ckh.print()
			}
		}
	})
}
