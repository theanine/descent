package main

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/theanine/utils"
)

const plotHtml = "plot.html"
const plotImg = "plots/Plot_Card_Back.png"

type pCard struct {
	name    string
	img     string
	imgBack string // agents
	num     int64
	exp     string
	expImg  string
	cost    int64
	text    string
}

type plot struct {
	name        string
	arch        string
	description string
	cards       []pCard
	url         string
}

var plots []plot

func dumpPlots() {
	// fmt.Printf("Plots: %d\n", len(plots))
	// for _, p := range plots {
	// 	if p.description != "" {
	// 		fmt.Printf("[%s: %s]\n", p.archetype, p.description)
	// 	} else {
	// 		fmt.Printf("[%s]\n", p.archetype)
	// 	}
	// 	for _, e := range p.cards {
	// 		fmt.Printf("\t[%d] %dx %s\n", e.cost, e.qty, e.name)
	// 	}
	// 	fmt.Printf("\n")
	// }
}

func pImgRtoL(img string, name string) (string, error) {
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
	return "plots/" + strings.Replace(file, "Plot_Card_-_", "", -1), nil
}

func downloadPImages() {
	for _, p := range plots {
		for _, c := range p.cards {
			var conf utils.Config
			if c.img != "" {
				conf.Url = wikiUrl + c.img
				outfile, err := pImgRtoL(c.img, c.name)
				if err != nil {
					panic(err)
				}
				conf.Outfile = outfile
				if _, _, err := utils.Wget(conf); err != nil {
					if strings.Contains(conf.Url, "/revision/latest") {
						conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
					}
					if _, _, err := utils.Wget(conf); err != nil {
						panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
					}
				}
			}
		}
	}
}

func plotGen() {
	doc, err := goquery.NewDocument(wikiUrl + "/Lieutenant_Pack")
	if err != nil {
		panic(fmt.Sprintf("error on parsing: %s", err))
	}

	doc.Find("ul").Eq(0).Find("a").Each(func(i int, a *goquery.Selection) {
		pClass := plot{}
		if href, ok := a.Attr("href"); ok {
			doc, err := goquery.NewDocument(wikiUrl + href)
			if err != nil {
				panic(err)
			}
			p := doc.Find("p").Last()
			pClass.name = strings.TrimSpace(p.Prev().Text())
			pClass.arch = strings.TrimSpace(strings.Replace(doc.Find("h1").Text(), " Lieutenant Pack", "", -1))
			descArr := strings.Split(strings.TrimSpace(p.Text()), ". ")
			if strings.Contains(descArr[len(descArr)-1], ":") {
				descArr = descArr[:len(descArr)-1]
			}
			descArr = descArr[1:]
			pClass.description = strings.Join(descArr, ". ")
			pClass.url = wikiUrl + href

			p.Find("sup").Remove()
			p.Next().Find("a").Each(func(i int, a *goquery.Selection) {
				pcard := pCard{}
				var exp string
				if href, ok := a.Attr("href"); ok {
					numFound := false
					costFound := false
					expFound := false
					doc, err := goquery.NewDocument(wikiUrl + href)
					if err != nil {
						panic(err)
					}
					wikitable := doc.Find(".wikitable")
					if wikitable.Length() == 0 {
						return
					}

					pcard.name = strings.TrimSpace(wikitable.Find("div").First().Text())

					wikitable.Find("td").Each(func(j int, t *goquery.Selection) {
						text := strings.TrimSpace(t.Text())
						if costFound {
							pcard.cost = 0
							val := strings.Split(text, " Threat")[0]
							if val != "" {
								if pcard.cost, err = strconv.ParseInt(val, 10, 64); err != nil {
									panic(err)
								}
							}
							costFound = false
						}
						if text == "Play Cost:" {
							costFound = true
						}
						if numFound {
							num := strings.Split(text, "/")[0]
							if num != "" {
								if pcard.num, err = strconv.ParseInt(num, 10, 64); err != nil {
									panic(err)
								}
							}
							numFound = false
						}
						if text == "Plot card number:" {
							numFound = true
						}
						if expFound {
							exp = strings.Split(text, " ")[0]
							if exp != "" {
								pcard.exp = expCodes[exp]
							}
							expFound = false
						}
						if text == "Expansion" {
							expFound = true
						}
					})

					doc.Find(".wikitable").Remove()
					table := doc.Find("table")
					href, _ := table.Find("a").First().Attr("href")
					pcard.text = tdToSkill(table)
					pClass.cards = append(pClass.cards, pcard)

					if strings.Contains(pcard.name, "Summon ") {
						doc, err := goquery.NewDocument(wikiUrl + href)
						if err != nil {
							panic(err)
						}
						pcard1 := pCard{
							name: strings.TrimSpace(doc.Find("h1").Text()),
							num:  int64(10),
							exp:  expCodes[exp],
							text: tdToSkill(doc.Find("#mw-content-text").Find("ul").Last()),
						}
						img := "agents/" + strings.Replace(pcard1.name, " (Agent)", "", -1)
						img = strings.Replace(img, " ", "_", -1)
						if _, err := os.Stat(img + "_I_Front.jpg"); !os.IsNotExist(err) {
							pcard1.img = img + "_I_Front.jpg"
						} else if _, err := os.Stat(img + "_I_Front.png"); !os.IsNotExist(err) {
							pcard1.img = img + "_I_Front.png"
						}
						if _, err := os.Stat(img + "_I_Back.jpg"); !os.IsNotExist(err) {
							pcard1.imgBack = img + "_I_Back.jpg"
						} else if _, err := os.Stat(img + "_I_Back.png"); !os.IsNotExist(err) {
							pcard1.imgBack = img + "_I_Back.png"
						}
						pClass.cards = append(pClass.cards, pcard1)

						pcard2 := pCard{
							name: strings.TrimSpace(doc.Find("h1").Text()),
							num:  int64(11),
							exp:  expCodes[exp],
							text: tdToSkill(doc.Find("#mw-content-text").Find("ul").Last()),
						}
						img = "agents/" + strings.Replace(pcard2.name, " (Agent)", "", -1)
						img = strings.Replace(img, " ", "_", -1)
						if _, err := os.Stat(img + "_II_Front.jpg"); !os.IsNotExist(err) {
							pcard2.img = img + "_II_Front.jpg"
						} else if _, err := os.Stat(img + "_II_Front.png"); !os.IsNotExist(err) {
							pcard2.img = img + "_II_Front.png"
						}
						if _, err := os.Stat(img + "_II_Back.jpg"); !os.IsNotExist(err) {
							pcard2.imgBack = img + "_II_Back.jpg"
						} else if _, err := os.Stat(img + "_II_Back.png"); !os.IsNotExist(err) {
							pcard2.imgBack = img + "_II_Back.png"
						}
						pClass.cards = append(pClass.cards, pcard2)
					}
				}
			})
			plots = append(plots, pClass)
		}
	})

	// dumpPlots()

	f, err := os.Create(plotHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if downloadEnabled {
		downloadPImages()
	}
	fixPlots()
	outputPTable(w)
}

func fixPlots() {
	for i, p := range plots {
		for j, c := range p.cards {
			// c.img
			if c.img == "" || c.img == "plots/Plot_Card_Back.png" {
				imgFile := "plots/" + strings.Replace(c.name, " ", "_", -1) + ".png"
				if _, err := os.Stat(imgFile); os.IsNotExist(err) {
					imgFile = "plots/" + strings.Replace(c.name, " ", "_", -1) + ".jpg"
					if _, err := os.Stat(imgFile); os.IsNotExist(err) {
						imgFile = "plots/Plot_Card_Back.png"
					}
				}
				plots[i].cards[j].img = imgFile
			}
			replaceErrata(&plots[i].cards[j].img)
			replaceErrata(&plots[i].cards[j].imgBack)

			// c.expImg
			plots[i].cards[j].expImg = ""
			imgFile := "expansions/" + strings.Replace(c.exp, " ", "_", -1) + ".svg"
			if _, err := os.Stat(imgFile); !os.IsNotExist(err) {
				plots[i].cards[j].expImg = fmt.Sprintf("<img src=\"%s\" class=\"expansion\">", imgFile)
			} else if abbr, ok := expansions[strings.ToLower(c.exp)]; ok {
				plots[i].cards[j].expImg = abbr
			}
		}
	}
}

func plotUniqSortArchs() []string {
	archMap := make(map[string]struct{})
	for _, p := range plots {
		archMap[p.arch] = struct{}{}
	}
	var archs []string
	for arch := range archMap {
		archs = append(archs, arch)
	}
	sort.Strings(archs)
	return archs
}

func outputPHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Journeys in Plot Selection</title>\n")
	fmt.Fprintf(w, "<meta name=\"description\" content=\"%s\">\n", `With 20 plot decks and 240 cards to choose from, it's painful to make a good plot deck.

For owners of Descent: Journeys in the Dark (Second Edition), this Plot Selector makes the decision that much easier for newbies, casuals, and veterans.

Send your overlord to get some Coufee and they'll be adventuring in no time!`)
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.widgets.min.js\"></script>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"heroes.css?version=%s\">\n", version)
	fmt.Fprintf(w, "<link rel=\"icon\" type=\"image/png\" href=\"etc/favicon.png\">\n")
	fmt.Fprintf(w, "</head><body onload=\"onload()\">\n")

	// table
	fmt.Fprintf(w, "<table id=\"plotTable\" class=\"tablesorter\"><thead class=\"plot\"><tr>\n")
	fmt.Fprintf(w, "<th class=\"archetype\">Class</th>\n")
	fmt.Fprintf(w, "<th class=\"pCards\">Cards</th>\n")
	fmt.Fprintf(w, "</tr>\n\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<td class=\"archetype\"><div><select id=\"selectClass\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	archs := plotUniqSortArchs()
	for _, arch := range archs {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(strings.Replace(arch, " ", "-", -1)), arch)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"pSearch\"><div class=\"search\">\n")
	fmt.Fprintf(w, "<input type=\"text\" class=\"search-input\" id=\"search-input\" name=\"search\" placeholder=\"Search\" onkeyup=\"search()\"/>\n")
	fmt.Fprintf(w, "<input type=\"submit\" class=\"search-submit\"/></div></td>\n")
	fmt.Fprintf(w, "</tr></thead><tbody class=\"plot\">\n\n")
}

func outputPTableRow(w *bufio.Writer, p plot) {
	fmt.Fprintf(w, "<tr class=\"%s %s\" style=\"display:none;\">\n", strings.ToLower(p.name), strings.ToLower(strings.Replace(p.arch, " ", "-", -1)))
	fmt.Fprintf(w, "<td class=\"archetype\">")
	fmt.Fprintf(w, "<span title=\"%s\">", html.EscapeString(p.description))
	fmt.Fprintf(w, "<a href=\"%s\" class=\"plot\">", p.url)
	fmt.Fprintf(w, "<div class=\"divImage\">")
	fmt.Fprintf(w, "<img src=\"%s\" class=\"plot\" >", plotImg)
	fmt.Fprintf(w, "<div class=\"name\">%s</div>", p.name)
	fmt.Fprintf(w, "<div class=\"archetype\" style=\"display: none;\">%s</div>", p.arch)
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</a></span></td>\n")

	fmt.Fprintf(w, "<td class=\"cards\">")
	for _, c := range p.cards {
		// if c.img == "" || c.img == "plots/Back_-_Plot_Card.png" {
		// 	c.img = "plots/" + strings.Replace(c.name, " ", "_", -1) + ".png"
		// }
		fmt.Fprintf(w, "<div class=\"cardContainer\">")
		fmt.Fprintf(w, "<img src=\"%s\" class=\"cards\" exp=\"%s\" alt=\"%s\" num=\"%d\" cost=\"%d\" text=\"%s\">", c.img, c.exp, c.name, c.num, c.cost, c.text)
		if c.imgBack != "" {
			fmt.Fprintf(w, "<img src=\"%s\" class=\"cards agentBack\">", c.imgBack)
		}
		fmt.Fprintf(w, "</div>")
	}
	fmt.Fprintf(w, "</td></tr>\n\n")
}

func outputPTable(w *bufio.Writer) {
	outputPHeader(w)

	for _, p := range plots {
		outputPTableRow(w, p)
	}

	outputPFooter(w)
	w.Flush()
}

func outputPFooter(w *bufio.Writer) {
	outputFooter(w, "plot", 3)
}
