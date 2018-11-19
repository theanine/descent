package main

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/theanine/utils"
)

const oldWiki = "https://descent2e.wikia.com/wiki/Overlord_Card"
const wikiUrl = "http://wiki.descent-community.org"
const overlordHtml = "overlord.html"
const overlordImg = "olcards/Overlord_Card_Back.png"

type card struct {
	qty       int64
	name      string
	img       string
	expansion string
	expImg    string
	typ       string
	cost      int
}

type overlord struct {
	archetype   string
	cards       []card
	description string
	url         string
}

var overlords []overlord

func dumpOverlords() {
	fmt.Printf("Overlords: %d\n", len(overlords))
	for _, o := range overlords {
		if o.description != "" {
			fmt.Printf("[%s: %s]\n", o.archetype, o.description)
		} else {
			fmt.Printf("[%s]\n", o.archetype)
		}
		for _, e := range o.cards {
			fmt.Printf("\t[%d] %dx %s\n", e.cost, e.qty, e.name)
		}
		fmt.Printf("\n")
	}
}

func oImgRtoL(img string, name string) (string, error) {
	parts := strings.Split(img, "/")
	if len(parts) <= 5 {
		return "", fmt.Errorf("unexpected class img string (<=5 parts): %s\n", img)
	}
	file, err := url.QueryUnescape(parts[5])
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(file)
	if ext == "" {
		ext = ".png"
	}
	return "olcards/" + strings.Replace(file, "Overlord_Card_-_", "", -1), nil
}

func downloadOImages() {
	for _, o := range overlords {
		for _, c := range o.cards {
			var conf utils.Config
			if c.img != "" {
				conf.Url = wikiUrl + c.img
				outfile, err := oImgRtoL(c.img, c.name)
				if err != nil {
					panic(err)
				}
				conf.Outfile = outfile
				if _, err := utils.Wget(conf); err != nil {
					if strings.Contains(conf.Url, "/revision/latest") {
						conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
					}
					if _, err := utils.Wget(conf); err != nil {
						panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
					}
				}
			}
		}
	}
}

func overlordGen() {
	doc, err := goquery.NewDocument(wikiUrl + "/Overlord_Card")
	if err != nil {
		panic(fmt.Sprintf("error on parsing: %s", err))
	}

	var decks []string

	overlordMetadata := doc.Find("#mw-content-text").First()
	headers := overlordMetadata.Find("h2")
	headers.Each(func(i int, s1 *goquery.Selection) {
		if s1.Find("span").Eq(0).Find("a").Length() > 0 {
			// fmt.Printf("[%s]\n", s1.Find("span").Eq(0).Find("a").Text())
			archetype := s1.Find("span").Eq(0).Find("a").Text()
			if archetype == "Basic" {
				archetype = "Basic I"
			}
			if archetype == "Overlord Reward" {
				archetype = "Reward"
			}
			decks = append(decks, archetype)
		}
		return
	})

	reQty := regexp.MustCompile("[^0-9]")
	reParens := regexp.MustCompile(`\((.*)\)`)
	cardLists := overlordMetadata.Find("ul").Slice(2, goquery.ToEnd)
	cardLists.Each(func(i int, u1 *goquery.Selection) {
		olClass := overlord{}
		olClass.archetype = decks[i]

		u1.Find("li").Each(func(j int, l1 *goquery.Selection) {
			olCard := card{}
			li := l1.RemoveClass("li")
			olCard.name = li.Text()
			if olCard.name[1:3] == "x " {
				olCard.name = olCard.name[3:]
			}
			olCard.name = reParens.ReplaceAllString(olCard.name, "")
			olCard.name = strings.TrimSpace(olCard.name)

			aTag := li.Find("a")
			if overlordUrl, ok := aTag.Attr("href"); ok {
				overlordUrl = wikiUrl + overlordUrl
				doc, err := goquery.NewDocument(overlordUrl)
				if err != nil {
					panic(fmt.Sprintf("error on parsing: %s", err))
				}

				typeFound := false
				costFound := false
				wikitable := doc.Find(".wikitable")
				wikitable.Find("td").Each(func(i int, t *goquery.Selection) {
					text := strings.TrimSpace(t.Text())
					if typeFound {
						olCard.typ = text
						typeFound = false
					}
					if text == "Type:" {
						typeFound = true
					}
					if costFound {
						olCard.cost = 0
						val := strings.Split(text, " XP")[0]
						if val != "-" {
							var err error
							olCard.cost, err = strconv.Atoi(val)
							if err != nil {
								log.Printf("%s\n", err)
							}
						}
						costFound = false
					}
					if text == "XP cost:" {
						costFound = true
					}
				})

				imageFound := false
				if src, ok := wikitable.Find("a").First().Find("img").Eq(0).Attr("srcset"); ok {
					set := strings.Split(src, ", ")
					if len(set) > 0 {
						tmp := set[len(set)-1]
						olCard.img = strings.Split(tmp, " ")[0]
						if olCard.img != "" {
							imageFound = true
						}
					}
				}
				if !imageFound {
					if src, ok := wikitable.Find("a").First().Find("img").Eq(0).Attr("src"); ok {
						if strings.Contains(src, "/thumb/") {
							olCard.img = strings.Replace(src[:strings.LastIndex(src, "/")], "/thumb/", "/", -1)
						} else {
							olCard.img = strings.Split(src, "/scale-to-width-down")[0]
							if src[:10] == "data:image" {
								if src, ok := wikitable.Find("a").First().Find("img").Eq(0).Attr("data-src"); ok {
									olCard.img = strings.Split(src, "/scale-to-width-down")[0]
								}
							}
						}
					}
				}
			}

			qty := strings.TrimSpace(l1.Children().Remove().End().Text())
			qty = reQty.ReplaceAllString(qty, "")
			if qty == "" {
				qty = "1"
			}
			qty = reQty.ReplaceAllString(qty, "")
			if olCard.qty, err = strconv.ParseInt(qty, 10, 64); err != nil {
				panic(err)
			}

			olClass.cards = append(olClass.cards, olCard)
		})
		overlords = append(overlords, olClass)
	})

	idx := 0
	headers.Each(func(i int, s1 *goquery.Selection) {
		if s1.Find("span").Eq(0).Find("a").Length() > 0 {
			aTag := s1.Find("span").Eq(0).Find("a")
			if url, ok := aTag.Attr("href"); ok {
				overlords[idx].url = wikiUrl + url
			}
			idx++
		}
	})
	// dumpOverlords()

	f, err := os.Create(overlordHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if downloadEnabled {
		downloadOImages()
	}
	fixOverlords()
	outputOTable(w)
}

func fixOverlords() {
	for i, o := range overlords {
		// c.cards
		sort.Slice(o.cards, func(i, j int) bool { return o.cards[i].cost < o.cards[j].cost })

		for j, c := range overlords[i].cards {
			// c.img
			if c.img != "" {
				var err error
				if overlords[i].cards[j].img, err = oImgRtoL(c.img, c.name); err != nil {
					panic(err)
				}
			}

			// c.expImg
			overlords[i].cards[j].expImg = ""
			imgFile := "expansions/" + strings.Replace(overlords[i].cards[j].expansion, " ", "_", -1) + ".svg"
			if strings.Contains(overlords[i].cards[j].expansion, "Lieutenant Pack") {
				imgFile = "expansions/Lieutenant_Pack.png"
			}
			if _, err := os.Stat(imgFile); !os.IsNotExist(err) {
				overlords[i].cards[j].expImg = fmt.Sprintf("<img src=\"%s\" class=\"expansion\">", imgFile)
			} else if strings.ToLower(overlords[i].cards[j].expansion) == "second edition base game" {
				overlords[i].cards[j].expImg = "2E"
			} else if strings.ToLower(overlords[i].cards[j].expansion) == "second edition conversion kit" {
				overlords[i].cards[j].expImg = "1E"
			}
		}
	}

}

func overlordUniqSortExps() []string {
	expMap := make(map[string]struct{})
	for _, o := range overlords {
		for _, c := range o.cards {
			expMap[expansions[strings.ToLower(c.expansion)]] = struct{}{}
		}
	}
	var exps []string
	for exp := range expMap {
		exps = append(exps, exp)
	}
	sort.Strings(exps)
	return exps
}

func overlordUniqSortArchs() []string {
	archMap := make(map[string]struct{})
	for _, o := range overlords {
		archMap[o.archetype] = struct{}{}
	}
	var archs []string
	for arch := range archMap {
		archs = append(archs, arch)
	}
	sort.Strings(archs)
	return archs
}

func overlordUniqSortTypes() []string {
	typMap := make(map[string]struct{})
	for _, o := range overlords {
		for _, c := range o.cards {
			typMap[c.typ] = struct{}{}
		}
	}
	var typs []string
	for typ := range typMap {
		typs = append(typs, typ)
	}
	sort.Strings(typs)
	return typs
}

func outputOHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Journeys in Overlord Selection</title>\n")
	fmt.Fprintf(w, "<meta name=\"description\" content=\"%s\">\n", `With over 10+ overlord decks and 100+ cards to choose from, it's painful to make a good deck.

For owners of Descent: Journeys in the Dark (Second Edition), this Overlord Selector makes the decision that much easier for newbies, casuals, and veterans.

Send your overlord to get some Coufee and they'll be adventuring in no time!`)
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.widgets.min.js\"></script>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"heroes.css?version=%s\">\n", version)
	fmt.Fprintf(w, "<link rel=\"icon\" type=\"image/png\" href=\"etc/favicon.png\">\n")
	fmt.Fprintf(w, "</head><body onload=\"onload()\">\n")

	// table
	fmt.Fprintf(w, "<table id=\"overlordTable\" class=\"tablesorter\"><thead class=\"overlord\"><tr>\n")
	fmt.Fprintf(w, "<th class=\"archetype\">Class</th>\n")
	fmt.Fprintf(w, "<th class=\"cards\">Cards</th>\n")
	fmt.Fprintf(w, "</tr>\n\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<td class=\"archetype\"><div><select id=\"selectClass\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	archs := overlordUniqSortArchs()
	for _, arch := range archs {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(arch), arch)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"cards\"><div><select id=\"selectCards\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	types := overlordUniqSortTypes()
	for _, typ := range types {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(typ), typ)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "</tr></thead><tbody class=\"overlord\">\n\n")
}

func outputOTableRow(w *bufio.Writer, o overlord) {
	fmt.Fprintf(w, "<tr class=\"%s\" style=\"display:none;\">\n", strings.ToLower(o.archetype))
	fmt.Fprintf(w, "<td class=\"archetype\">")
	fmt.Fprintf(w, "<span title=\"%s\">", html.EscapeString(o.description))
	fmt.Fprintf(w, "<a href=\"%s\" class=\"overlord\">", o.url)
	fmt.Fprintf(w, "<div class=\"divImage\">")
	fmt.Fprintf(w, "<img src=\"%s\" class=\"overlord\">", overlordImg)
	fmt.Fprintf(w, "<div class=\"archetype\">%s</div>", o.archetype)
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</a></span></td>\n")

	fmt.Fprintf(w, "<td class=\"cards\">")
	for _, c := range o.cards {
		if c.img == "" || c.img == "olcards/Back_-_Overlord_Card.png" {
			// TODO: fix these
			continue
		}
		fmt.Fprintf(w, "<div class=\"cardContainer %s\">", strings.ToLower(c.typ))
		fmt.Fprintf(w, "<img src=\"%s\" alt=\"%s\" class=\"cards\">", c.img, c.name)
		if c.qty > 1 {
			fmt.Fprintf(w, "<div class=\"quantity\">%d</div>", c.qty)
		} else {
			fmt.Fprintf(w, "<div class=\"quantity\"></div>")
		}
		fmt.Fprintf(w, "</div>")
	}
	fmt.Fprintf(w, "</td></tr>\n\n")
}

func outputOTable(w *bufio.Writer) {
	outputOHeader(w)

	for _, o := range overlords {
		outputOTableRow(w, o)
	}

	outputOFooter(w)
	w.Flush()
}

func outputOFooter(w *bufio.Writer) {
	fmt.Fprintf(w, "</tbody>\n")
	fmt.Fprintf(w, "<tfoot class=\"overlord\"><tr><td class=\"donateArea\">\n")
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
						3Q6y5d5c43Lj9maDr8dcZyXUFqxPcbBiEv</span></div>`)
	fmt.Fprintf(w, "</td><td class=\"fees\">Server Fees: $55.80/yr")
	fmt.Fprintf(w, "</td><td class=\"version\">%s</td>\n", version)
	fmt.Fprintf(w, "</tr></tfoot>\n")
	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "</body></html>\n")
	fmt.Fprintf(w, "<script type=\"text/javascript\" src=\"heroes.js?version=%s\"></script>\n", version)
}
