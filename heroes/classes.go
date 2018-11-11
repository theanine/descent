package main

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/theanine/utils"
)

const classesHtml = "classes.html"

type skill struct {
	name string
	xp   int
	text string
	cost int
	img  string
}

type class struct {
	name        string
	img         string
	archetype   string
	equipment   []string
	description string
	expansion   string
	expImg      string
	skills      []skill
	url         string
}

var classes []class

func filterText(str string) bool {
	if len(strings.TrimSuffix(str, ".")) == 0 {
		return true
	}
	if strings.Contains(str, "is one of the class options for the") {
		return true
	}
	if strings.Contains(str, "It is one of the first Hybrid classes, first introduced in The Chains that Rust") {
		return true
	}
	if strings.Contains(str, "This class was first introduced in the") {
		return true
	}
	return false
}

func dumpClasses() {
	fmt.Printf("Classes: %d\n", len(classes))
	for _, c := range classes {
		// fmt.Printf("[%s] [%s] [%s] [%s] [%s]\n", c.expansion, c.name, c.archetype, c.equipment, c.description)
		fmt.Printf("%s|%s|%s|", c.expansion, c.name, c.archetype)
		for i, e := range c.equipment {
			fmt.Printf("%s", e)
			if i < len(c.equipment)-1 {
				fmt.Printf(",")
			}
		}
		fmt.Printf("|%s\n", c.description)
	}
}

func cImgRtoL(img string, name string) (string, error) {
	parts := strings.Split(img, "/")
	if len(parts) <= 7 {
		return "", fmt.Errorf("unexpected class img string (<7 parts): %s\n", img)
	}
	file, err := url.QueryUnescape(parts[7])
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(file)
	if ext == "" {
		ext = ".png"
	}
	return "classes/" + strings.Replace(name, " ", "_", -1) + ext, nil
}

func sImgRtoL(img string) (string, error) {
	parts := strings.Split(img, "/")
	if len(parts) <= 7 {
		return "", fmt.Errorf("unexpected skill img string (<7 parts): %s\n", img)
	}
	parts = strings.Split(parts[7], "_-_")
	if len(parts) <= 1 {
		return "", fmt.Errorf("unexpected skill img string (no '_-_'): %s\n", img)
	}
	file, err := url.QueryUnescape(parts[1])
	if err != nil {
		return "", err
	}
	return "skills/" + file, nil
}

func downloadCImages() {
	for _, c := range classes {
		var conf utils.Config
		if c.img != "" {
			conf.Url = c.img
			outfile, err := cImgRtoL(c.img, c.name)
			if err != nil {
				panic(err)
			}
			conf.Outfile = outfile
			if _, err := utils.Wget(conf); err != nil {
				conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
				if _, err := utils.Wget(conf); err != nil {
					panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
				}
			}
		}

		for _, s := range c.skills {
			if s.img != "" {
				conf.Url = s.img
				outfile, err := sImgRtoL(s.img)
				if err != nil {
					panic(err)
				}
				conf.Outfile = outfile
				if _, err := utils.Wget(conf); err != nil {
					conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
					if _, err := utils.Wget(conf); err != nil {
						panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
					}
				}
			}
		}
	}
}

func classesGen() {
	doc, err := goquery.NewDocument("https://descent2e.wikia.com/wiki/Class")
	if err != nil {
		panic(fmt.Sprintf("error on parsing: %s", err))
	}

	classMetadata := doc.Find(".mw-content-text").First()
	classMetadata.Find("li").Each(func(i int, s1 *goquery.Selection) {
		meta := class{}
		meta.name = strings.TrimSpace(s1.Text())

		classUrl, ok := s1.Find("a").Attr("href")
		if !ok {
			return
		}

		meta.url = "http://descent2e.wikia.com" + classUrl
		doc, err := goquery.NewDocument(meta.url)
		if err != nil {
			panic(fmt.Sprintf("error on parsing: %s", err))
		}

		// meta.archetype = strings.TrimSuffix(strings.TrimSpace(doc.Find(".caption").Eq(1).Text()), ".")
		doc.Find(".dablink").Remove()
		content := doc.Find("#mw-content-text")
		src, ok := content.Find(".thumbimage").Eq(0).Attr("data-src")
		if ok {
			meta.img = strings.Split(src, "/scale-to-width-down")[0]
		}

		wikiCount := 0
		content.Find("a").Each(func(i int, a *goquery.Selection) {
			text := strings.TrimSpace(a.Text())
			if len(text) == 0 {
				return
			}
			if strings.Contains(text, "stub") || strings.Contains(text, "expanding it or adding missing information") {
				return
			}
			if text == "Heroes" || text == "Hybrid" {
				return
			}
			if href, ok := a.Attr("href"); ok {
				if strings.HasPrefix(href, "/wiki/") {
					wikiCount++
					if wikiCount == 1 {
						meta.archetype = text
						return
					}
					if wikiCount == 2 {
						meta.expansion = text
						return
					}
				}
			}
		})

		classTable := content.Find(".wikitable")
		classTable.Find("tr").Each(func(i int, s1 *goquery.Selection) {
			elements := s1.Find("td")
			if elements.Length() == 0 {
				return
			}

			sMeta := skill{}
			sMeta.name = strings.TrimSpace(elements.Eq(0).Text())
			sMeta.xp, _ = strconv.Atoi(strings.TrimSpace(elements.Eq(1).Text()))
			sMeta.text = strings.TrimSpace(elements.Eq(2).Text())
			sMeta.cost, _ = strconv.Atoi(strings.TrimSpace(elements.Eq(3).Text()))

			aTag := elements.Eq(0).Find("a")
			if skillUrl, ok := aTag.Attr("href"); ok {
				skillUrl = "http://descent2e.wikia.com" + skillUrl
				doc, err := goquery.NewDocument(skillUrl)
				if err != nil {
					panic(fmt.Sprintf("error on parsing: %s", err))
				}

				if src, ok := doc.Find(".wikitable").Find("a").First().Find("img").Eq(0).Attr("src"); ok {
					sMeta.img = strings.Split(src, "/scale-to-width-down")[0]
					if src[:10] == "data:image" {
						if src, ok := doc.Find(".wikitable").Find("a").First().Find("img").Eq(0).Attr("data-src"); ok {
							sMeta.img = strings.Split(src, "/scale-to-width-down")[0]
						}
					}
				}
			}

			meta.skills = append(meta.skills, sMeta)
		})

		italics := true
		doc.Find("p").Each(func(i int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())
			if strings.Contains(strings.ToLower(text), "starting equipment") {
				p.Find("a").Each(func(i int, a *goquery.Selection) {
					if href, ok := a.Attr("href"); ok {
						href = "https://descent2e.wikia.com/wiki/" + strings.Split(href, "/wiki/")[1]
						meta.equipment = append(meta.equipment, href)
					}
					return
				})
				return
			}
			if p.HasClass("caption") {
				return
			}
			if filterText(text) {
				return
			}
			if len(text) < len(meta.description) {
				return
			}
			if p.Find("i").Length() > 0 && !italics {
				return
			}
			if p.Find("i").Length() == 0 {
				italics = false
			}

			meta.description = text
			// fmt.Printf("%d:%s => [%s]\n", i, meta.name, text)
		})
		// NOTE: this is destructive to the doc
		text := strings.TrimSpace(content.Children().Remove().End().Text())
		if filterText(text) {
			text = ""
		}
		if len(text) > len(meta.description) {
			meta.description = text
		}

		url := fmt.Sprintf("https://wiki.descent-community.org/File:Back_-_%s.png", strings.Replace(meta.name, " ", "_", -1))
		doc, err = goquery.NewDocument(url)
		if err != nil {
			panic(fmt.Sprintf("error on parsing: %s", err))
		}
		src, ok = doc.Find(".fullImageLink").Find("img").Eq(0).Attr("src")
		if ok {
			meta.img = "https://wiki.descent-community.org/" + src
		}

		// fmt.Printf("%s => TEXT: [%s]\n", meta.name, text)
		classes = append(classes, meta)
	})

	// dumpClasses()

	f, err := os.Create(classesHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if downloadEnabled {
		downloadCImages()
	}
	fixClasses()
	outputCTable(w)
}

func fixClasses() {
	for i := range classes {
		// c.expImg
		classes[i].expImg = ""
		imgFile := "expansions/" + strings.Replace(classes[i].expansion, " ", "_", -1) + ".svg"
		if strings.Contains(classes[i].expansion, "Lieutenant Pack") {
			imgFile = "expansions/Lieutenant_Pack.png"
		}
		if _, err := os.Stat(imgFile); !os.IsNotExist(err) {
			classes[i].expImg = fmt.Sprintf("<img src=\"%s\" class=\"expansion\">", imgFile)
		} else if strings.ToLower(classes[i].expansion) == "second edition base game" {
			classes[i].expImg = "2E"
		} else if strings.ToLower(classes[i].expansion) == "second edition conversion kit" {
			classes[i].expImg = "1E"
		}
	}

	var err error
	for i, c := range classes {
		if c.img != "" {
			if classes[i].img, err = cImgRtoL(c.img, c.name); err != nil {
				panic(err)
			}
		}
		for j, s := range c.skills {
			if s.img != "" {
				if classes[i].skills[j].img, err = sImgRtoL(s.img); err != nil {
					panic(err)
				}
			}
		}
	}
}

func outputCHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Journeys in Hero Selection</title>\n")
	fmt.Fprintf(w, "<meta name=\"description\" content=\"%s\">\n", `With over 100+ heroes to choose from, it's painful to choose a character.

For owners of Descent: Journeys in the Dark (Second Edition), this Hero Selector makes the decision that much easier for newbies, casuals, and veterans.

Send your heroes to get some Coufee and they'll be adventuring in no time!`)
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"jquery.tablesorter.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"jquery.tablesorter.widgets.min.js\"></script>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"heroes.css?version=%s\">\n", version)
	fmt.Fprintf(w, "<link rel=\"icon\" type=\"image/png\" href=\"etc/favicon.png\">\n")
	fmt.Fprintf(w, "</head><body>\n")

	// table
	fmt.Fprintf(w, "<table id=\"classTable\" class=\"tablesorter\"><thead class=\"classes\"><tr style=\"width:100%\">\n")
	fmt.Fprintf(w, "<th class=\"expansion\">Exp</th>\n")
	fmt.Fprintf(w, "<th class=\"class\">Class</th>\n")
	fmt.Fprintf(w, "<th class=\"skill\">Skills</th>\n")
	fmt.Fprintf(w, "</tr>\n\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<td class=\"expansion\"><div><select id=\"selectExp\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	exps := uniqueSortedExps()
	for _, exp := range exps {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", exp, exp)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"class\"><div><select id=\"selectClass\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	for _, k := range archetypes {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(k), k)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"skill\"><div></div></td>\n")
	fmt.Fprintf(w, "</tr></thead><tbody class=\"classes\">\n\n")
}

func outputCTable(w *bufio.Writer) {
	outputCHeader(w)

	for _, c := range classes {
		fmt.Fprintf(w, "<tr class=\"\" style=\"\">\n")
		fmt.Fprintf(w, "<td class=\"expansion\">%s</td>\n", c.expImg)
		fmt.Fprintf(w, "<td class=\"class\">")
		fmt.Fprintf(w, "<span title=\"%s\">", html.EscapeString(c.description))
		fmt.Fprintf(w, "<a href=\"%s\" class=\"class\">", c.url)
		fmt.Fprintf(w, "<div class=\"divImage\">")
		exists := true
		if _, err := os.Stat(c.img); os.IsNotExist(err) {
			exists = false
			c.img = "classes/generic_" + strings.ToLower(c.archetype) + ".png"
		}
		fmt.Fprintf(w, "<img src=\"%s\" class=\"class\">", c.img)
		if !exists {
			fmt.Fprintf(w, "<div class=\"className\">%s</div>", c.name)
		}
		fmt.Fprintf(w, "%c</div>", c.archetype[0])
		fmt.Fprintf(w, "</a></span></td>\n")

		// for _, e := range c.equipment {
		// 	// TODO: add these
		// }
		fmt.Fprintf(w, "<td class=\"skill\">")
		for _, s := range c.skills {
			img := s.img
			if img == "" {
				img = "skills/blank.png"
			}
			fmt.Fprintf(w, "<img src=\"%s\" class=\"skill\">", img)
		}
		fmt.Fprintf(w, "</td></tr>\n\n")
	}

	outputCFooter(w)
	w.Flush()
}

func outputCFooter(w *bufio.Writer) {
	fmt.Fprintf(w, "</tbody>\n")
	fmt.Fprintf(w, "<tfoot class=\"classes\"><tr><td class=\"donateArea\">\n")
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
	fmt.Fprintf(w, "</tr></tfoot>\n")
	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "</body></html>\n")
	fmt.Fprintf(w, "<script type=\"text/javascript\" src=\"heroes.js?version=%s\"></script>\n", version)
}