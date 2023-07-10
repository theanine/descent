package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const classesHtml = "classes.html"
const startingSkillsInEquipmentSection = true

// type skill struct {
// 	class string
// 	name  string
// 	img   string
// 	xp    int
// 	cost  int
// 	text  string
// }

// type equipment struct {
// 	class  string
// 	name   string
// 	typ    string // "Familiar", "Item", "Skill"
// 	img    string
// 	xp     int
// 	cost   int
// 	text   string
// 	ranged string
// 	traits []string
// }

// type class struct {
// 	name        string
// 	img         string
// 	archetype   string
// 	hybrid      bool
// 	equipments  []equipment
// 	description string
// 	expansion   string
// 	expImg      string
// 	skills      []skill
// 	url         string
// }

// var hybridMap = make(map[string]bool)
// var classes []class

type Class struct {
	Type string
	Name string `json:"name"`
	// Points    int    `json:"points"`
	Archetype string `json:"archetype"`
	Class     string `json:"class"`
	Hybrid    bool

	Speed       string `json:"speed"`        // FAMILIARS
	Health      string `json:"health"`       // FAMILIARS
	Defense     string `json:"defense"`      // FAMILIARS, defense dice
	Attack      string `json:"attack"`       // FAMILIARS | ITEMS, attack type
	Dice        string `json:"dice"`         // FAMILIARS | ITEMS, attack dice
	Traits      string `json:"traits"`       // ITEMS, bold keywords
	Equip       string `json:"equip"`        // ITEMS, equipment type
	XpCost      int    `json:"xp cost"`      // SKILLS
	FatigueCost string `json:"fatigue cost"` // SKILLS

	Rules     string `json:"rules"` // card text
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	// Xws       string `json:"xws"`
}

var classCards = map[string]map[string][]Class{}

func (c *Class) dump() {
	if c.Type != "" {
		fmt.Printf("\tType: %s\n", c.Type)
	}
	if c.Name != "" {
		fmt.Printf("\tName: %s\n", c.Name)
	}
	if c.Archetype != "" {
		fmt.Printf("\tArchetype: %s\n", c.Archetype)
	}
	if c.Class != "" {
		fmt.Printf("\tClass: %s\n", c.Class)
	}
	fmt.Printf("\tHybrid: %v\n", c.Hybrid)
	if c.Type == "familiar" {
		if c.Speed != "" {
			fmt.Printf("\tSpeed: %s\n", c.Speed)
		}
		if c.Health != "" {
			fmt.Printf("\tHealth: %s\n", c.Health)
		}
		if c.Defense != "" {
			fmt.Printf("\tDefense: %s\n", c.Defense)
		}
	}
	if c.Type == "familiar" || c.Type == "item" {
		if c.Attack != "" {
			fmt.Printf("\tAttack: %s\n", c.Attack)
		}
		if c.Dice != "" {
			fmt.Printf("\tDice: %s\n", c.Dice)
		}
	}
	if c.Type == "item" {
		if c.Traits != "" {
			fmt.Printf("\tTraits: %s\n", c.Traits)
		}
		if c.Equip != "" {
			fmt.Printf("\tEquip: %s\n", c.Equip)
		}
	}
	if c.Type == "skill" {
		fmt.Printf("\tXpCost: %d\n", c.XpCost)
		if c.FatigueCost != "" {
			fmt.Printf("\tFatigueCost: %s\n", c.FatigueCost)
		}
	}
	if c.Rules != "" {
		fmt.Printf("\tRules: %s\n", c.Rules)
	}
	if c.Expansion != "" {
		fmt.Printf("\tExpansion: %s\n", c.Expansion)
	}
	if c.Image != "" {
		fmt.Printf("\tImage: %s\n", c.Image)
	}
	fmt.Println()

	// fmt.Printf("[%s] [%s] [%s] [%s] [%s]\n", c.expansion, c.Name, c.Archetype, c.equipment, c.Description)
	// fmt.Println(c)
	// fmt.Printf("%v|%s|%s|%s|%s|%s|%s|", c.hybrid, c.expansion, c.Name, c.Archetype, c.expImg, c.Image, c.Url)
	// for i, e := range c.Equip {
	// 	fmt.Printf("%v", e)
	// 	if i < len(c.Equip)-1 {
	// 		fmt.Printf(",")
	// 	}
	// }
	// for i, s := range c.skills {
	// 	fmt.Printf("%v", s)
	// 	if i < len(c.skills)-1 {
	// 		fmt.Printf(",")
	// 	}
	// }
	// fmt.Printf("|%s\n", c.Description)
}

// func (s *skill) dump() {
// 	fmt.Printf("[%s] [%d] [%d] [%s] %s\n", s.class, s.xp, s.cost, s.Name, s.text)
// }

var classes []Class

func loadCJson() error {
	var familiars []Class
	dat1, err := ioutil.ReadFile("../d2e-master/data/class-familiars.js")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(dat1, &familiars); err != nil {
		return err
	}
	for i := range familiars {
		familiars[i].Type = "familiar"
	}
	classes = append(classes, familiars...)

	dat2, err := ioutil.ReadFile("../d2e-master/data/class-items.js")
	if err != nil {
		return err
	}
	var items []Class
	if err = json.Unmarshal(dat2, &items); err != nil {
		return err
	}
	for i := range items {
		items[i].Type = "item"
	}
	classes = append(classes, items...)

	dat3, err := ioutil.ReadFile("../d2e-master/data/class-skills.js")
	if err != nil {
		return err
	}
	var skills []Class
	if err = json.Unmarshal(dat3, &skills); err != nil {
		return err
	}
	for i := range skills {
		skills[i].Type = "skill"
	}
	classes = append(classes, skills...)

	for _, c := range classes {
		if classCards[c.Class] == nil {
			classCards[c.Class] = make(map[string][]Class)
		}
		classCards[c.Class][c.Name] = append(classCards[c.Class][c.Name], c)
	}
	return nil
}

func classesGen() {
	if err := loadCJson(); err != nil {
		panic(err)
	}

	// for _, line := range strings.Split(string(dat), "\n") {
	// 	arr := strings.Split(line, ": ")
	// 	if len(arr) != 2 {
	// 		continue
	// 	}
	// 	class := arr[0]
	// 	desc := arr[1]
	// 	descMap[class] = desc
	// }

	dumpClasses()
	// dumpSkills()

	f, err := os.Create(classesHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// if downloadEnabled {
	// 	downloadCImages()
	// }
	fixClasses()
	outputCTable(w)
}

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

// func testClass(class string) {
// 	c := genClass("http://wiki.descent-community.org/"+class, class)
// 	c.expImg = ""
// 	imgFile := "expansions/" + strings.Replace(c.expansion, " ", "_", -1) + ".svg"
// 	if _, err := os.Stat(imgFile); !os.IsNotExist(err) {
// 		c.expImg = fmt.Sprintf("<img src=\"%s\" class=\"expansion\">", imgFile)
// 	} else if abbr, ok := expansions[strings.ToLower(c.expansion)]; ok {
// 		c.expImg = abbr
// 	}
// 	c.dump()
// }

// func (c *class) dump() {
// 	// fmt.Printf("[%s] [%s] [%s] [%s] [%s]\n", c.expansion, c.Name, c.Archetype, c.equipment, c.Description)
// 	fmt.Printf("%v|%s|%s|%s|%s|%s|%s|", c.hybrid, c.expansion, c.Name, c.Archetype, c.expImg, c.Image, c.Url)
// 	for i, e := range c.Equip {
// 		fmt.Printf("%v", e)
// 		if i < len(c.Equip)-1 {
// 			fmt.Printf(",")
// 		}
// 	}
// 	for i, s := range c.skills {
// 		fmt.Printf("%v", s)
// 		if i < len(c.skills)-1 {
// 			fmt.Printf(",")
// 		}
// 	}
// 	fmt.Printf("|%s\n", c.Description)
// }

// func (s *skill) dump() {
// 	fmt.Printf("[%s] [%d] [%d] [%s] %s\n", s.class, s.xp, s.cost, s.Name, s.text)
// }

func dumpClasses() {
	fmt.Printf("Classes: %d\n", len(classes))
	for _, c := range classes {
		c.dump()
	}
}

// func dumpSkills() {
// 	for _, c := range classes {
// 		for _, s := range c.skills {
// 			s.dump()
// 		}
// 	}
// }

func cImgRtoL(img string, name string) (string, error) {
	parts := strings.Split(img, "/")
	if len(parts) <= 5 {
		return "", fmt.Errorf("unexpected class img string (<5 parts): %s\n", img)
	}
	var part string
	if parts[3] == "thumb" {
		part = parts[6]
	} else {
		part = parts[5]
	}
	file, err := url.QueryUnescape(part)
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
	if len(parts) <= 5 {
		return "", fmt.Errorf("unexpected skill img string (<5 parts): %s\n", img)
	}
	if parts[3] == "thumb" {
		parts = strings.Split(parts[6], "_-_")
	} else {
		parts = strings.Split(parts[5], "_-_")
	}
	if len(parts) <= 1 {
		return "", fmt.Errorf("unexpected skill img string (no '_-_'): %s\n", img)
	}
	file, err := url.QueryUnescape(parts[1])
	if err != nil {
		return "", err
	}
	return "skills/" + file, nil
}

func eImgRtoL(img string) (string, error) {
	parts := strings.Split(img, "/")
	if len(parts) <= 5 {
		return "", fmt.Errorf("unexpected skill img string (<5 parts): %s\n", img)
	}
	if parts[3] == "thumb" {
		parts = strings.Split(parts[6], "_-_")
	} else {
		parts = strings.Split(parts[5], "_-_")
	}
	if len(parts) <= 1 {
		return "", fmt.Errorf("unexpected skill img string (no '_-_'): %s\n", img)
	}
	file, err := url.QueryUnescape(parts[1])
	if err != nil {
		return "", err
	}
	return "equipment/" + file, nil
}

// func downloadCImages() {
// 	for _, c := range classes {
// 		var conf utils.Config
// 		if c.Image != "" {
// 			conf.Url = c.Image
// 			outfile, err := cImgRtoL(c.Image, c.Name)
// 			if err != nil {
// 				panic(err)
// 			}
// 			conf.Outfile = outfile
// 			if _, _, err := utils.Wget(conf); err != nil {
// 				conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
// 				if _, _, err := utils.Wget(conf); err != nil {
// 					panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
// 				}
// 			}
// 		}

// 		for _, s := range c.skills {
// 			if s.Image != "" {
// 				conf.Url = s.Image
// 				outfile, err := sImgRtoL(s.Image)
// 				if err != nil {
// 					panic(err)
// 				}
// 				conf.Outfile = outfile
// 				if _, _, err := utils.Wget(conf); err != nil {
// 					conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
// 					if _, _, err := utils.Wget(conf); err != nil {
// 						panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
// 					}
// 				}
// 			}
// 		}

// 		for _, e := range c.Equip {
// 			if e.Image != "" {
// 				conf.Url = e.Image
// 				outfile, err := eImgRtoL(e.Image)
// 				if err != nil {
// 					panic(err)
// 				}
// 				conf.Outfile = outfile
// 				if _, _, err := utils.Wget(conf); err != nil {
// 					conf.Url = strings.Split(conf.Url, "/revision/latest")[0]
// 					if _, _, err := utils.Wget(conf); err != nil {
// 						panic(fmt.Sprintf("%s: %s\n", conf.Url, err))
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func iconsToText(s *goquery.Selection) *goquery.Selection {
	s.Find("img.icon").Each(func(i int, img *goquery.Selection) {
		if src, ok := img.Attr("src"); ok {
			img.ReplaceWithHtml(strings.Split(src, ".")[0])
		}
	})
	s.Find("span").Each(func(i int, span *goquery.Selection) {
		text := strings.TrimSpace(span.Text())
		if text == "𝄞" {
			span.ReplaceWithHtml("Melody")
		} else if text == "𝄢" {
			span.ReplaceWithHtml("Harmony")
		}
	})
	s.Find("br").ReplaceWithHtml(" ")
	s.Find("i").Each(func(i int, iTag *goquery.Selection) {
		iTag.ReplaceWithHtml(iTag.Text())
	})
	s.Find("b").Each(func(i int, bTag *goquery.Selection) {
		bTag.ReplaceWithHtml(bTag.Text())
	})
	s.Find("a").Each(func(i int, aTag *goquery.Selection) {
		if title, ok := aTag.Attr("title"); ok {
			aTag.ReplaceWithHtml(title)
		}
	})
	return s
}

func tdToSkill(td *goquery.Selection) string {
	td = replaceIcons(td)
	td = iconsToText(td)
	// s, _ := td.Html()
	s := td.Text()
	skill := strings.Join(strings.Fields(s), " ")
	if skill == "" {
		skill = strings.TrimSpace(td.Text())
	}
	return skill
}

// func genClass(url string, name string) Class {
// 	meta := Class{}
// 	meta.Name = name
// 	meta.Url = url

// 	doc, err := goquery.NewDocument(meta.Url)
// 	if err != nil {
// 		panic(fmt.Sprintf("error on parsing: %s", err))
// 	}

// 	// meta.Archetype = strings.TrimSuffix(strings.TrimSpace(doc.Find(".caption").Eq(1).Text()), ".")
// 	doc.Find(".dablink").Remove()
// 	doc.Find("div.notice").Remove()
// 	content := doc.Find("#mw-content-text")
// 	imageFound := false
// 	if src, ok := content.Find(".thumbimage").Eq(0).Attr("srcset"); ok {
// 		set := strings.Split(src, ", ")
// 		if len(set) > 0 {
// 			tmp := set[len(set)-1]
// 			meta.Image = strings.Split(tmp, " ")[0]
// 			if meta.Image != "" {
// 				imageFound = true
// 			}
// 		}
// 	}
// 	if !imageFound {
// 		src, ok := content.Find(".thumbimage").Eq(0).Attr("data-src")
// 		if ok {
// 			meta.Image = strings.Split(src, "/scale-to-width-down")[0]
// 		}
// 	}

// 	content.Find(".mw-gallery-traditional").Remove()
// 	content.Find("div.thumb").Remove()

// 	wikiCount := 0
// 	meta.hybrid = false
// 	content.Find("a").Each(func(i int, a *goquery.Selection) {
// 		text := strings.TrimSpace(a.Text())
// 		if len(text) == 0 {
// 			return
// 		}
// 		if text == "Heroes" {
// 			return
// 		}
// 		if text == "Hybrid" {
// 			meta.hybrid = true
// 			return
// 		}
// 		if _, ok := a.Attr("href"); ok {
// 			wikiCount++
// 			if wikiCount == 1 {
// 				meta.Archetype = text
// 				return
// 			}
// 			if wikiCount == 2 {
// 				meta.expansion = text
// 				return
// 			}
// 		}
// 	})

// 	classTable := content.Find(".wikitable")
// 	classTable.Find("tr").Each(func(i int, s1 *goquery.Selection) {
// 		elements := s1.Find("td")
// 		if elements.Length() == 0 {
// 			return
// 		}

// 		sMeta := skill{}
// 		sMeta.Name = strings.TrimSpace(elements.Eq(0).Text())
// 		sMeta.xp, _ = strconv.Atoi(strings.TrimSpace(elements.Eq(1).Text()))
// 		sMeta.text = tdToSkill(elements.Eq(2))
// 		sMeta.cost, _ = strconv.Atoi(strings.TrimSpace(elements.Eq(3).Text()))

// 		aTag := elements.Eq(0).Find("a")
// 		if skillUrl, ok := aTag.Attr("href"); ok {
// 			skillUrl = wikiUrl + skillUrl
// 			doc, err := goquery.NewDocument(skillUrl)
// 			if err != nil {
// 				panic(fmt.Sprintf("error on parsing: %s", err))
// 			}

// 			img := doc.Find(".wikitable").Find("a").First().Find("img").Eq(0)
// 			imageFound := false
// 			if src, ok := img.Attr("srcset"); ok {
// 				set := strings.Split(src, ", ")
// 				if len(set) > 0 {
// 					tmp := set[len(set)-1]
// 					sMeta.Image = strings.Split(tmp, " ")[0]
// 					if sMeta.Image != "" {
// 						imageFound = true
// 					}
// 				}
// 			}
// 			if !imageFound {
// 				if src, ok := img.Attr("src"); ok {
// 					sMeta.Image = strings.Split(src, "/scale-to-width-down")[0]
// 					if src[:10] == "data:image" {
// 						if src, ok := img.Attr("data-src"); ok {
// 							sMeta.Image = strings.Split(src, "/scale-to-width-down")[0]
// 						}
// 					}
// 				}
// 			}
// 		}

// 		sMeta.class = meta.Name
// 		meta.skills = append(meta.skills, sMeta)
// 	})

// 	content.Find(".wikitable").Remove()
// 	text := strings.TrimSpace(content.Text())
// 	if strings.Contains(strings.ToLower(text), "starting equipment") {
// 		content.Find("a").Each(func(i int, a *goquery.Selection) {
// 			if href, ok := a.Attr("href"); ok {
// 				doc, err := goquery.NewDocument(wikiUrl + href)
// 				if err != nil {
// 					panic(fmt.Sprintf("error on parsing: %s", err))
// 				}
// 				e := equipment{}
// 				e.Name = strings.TrimSpace(doc.Find(".wikitable").Find("div").Eq(0).Text())
// 				for i, s := range meta.skills {
// 					if s.Name == e.Name {
// 						if startingSkillsInEquipmentSection {
// 							meta.skills = append(meta.skills[:i], meta.skills[i+1:]...)
// 							break
// 						} else {
// 							return
// 						}
// 					}
// 				}
// 				e.typ = strings.TrimSpace(doc.Find(".wikitable").Find("tr").Eq(7).Find("td").Eq(1).Text())
// 				if e.typ == "Item" {
// 					doc.Find(".wikitable").Find("tr").Each(func(i int, tr *goquery.Selection) {
// 						td := strings.TrimSpace(tr.Find("td").Eq(0).Text())
// 						if td == "Range:" {
// 							e.ranged = strings.TrimSpace(tr.Find("td").Eq(1).Text())
// 						}
// 						if td == "Trait:" || td == "Traits:" {
// 							traits := strings.TrimSpace(tr.Find("td").Eq(1).Text())
// 							for _, trait := range strings.Split(traits, ", ") {
// 								e.traits = append(e.traits, strings.TrimSpace(trait))
// 							}
// 						}
// 					})
// 				}
// 				e.xp = 0
// 				if e.typ == "Skill" {
// 					val := strings.TrimSpace(doc.Find(".wikitable").Find("tr").Eq(8).Find("td").Eq(1).Text())
// 					if val != "" {
// 						if e.cost, err = strconv.Atoi(val); err != nil {
// 							panic(err)
// 						}
// 					}
// 				} else {
// 					e.cost = 0
// 				}
// 				img := doc.Find(".wikitable").Find("img").Eq(0)
// 				imageFound := false
// 				if src, ok := img.Attr("srcset"); ok {
// 					set := strings.Split(src, ", ")
// 					if len(set) > 0 {
// 						tmp := set[len(set)-1]
// 						e.Image = strings.Split(tmp, " ")[0]
// 						if e.Image != "" {
// 							imageFound = true
// 						}
// 					}
// 				}
// 				if !imageFound {
// 					if src, ok := img.Attr("src"); ok {
// 						imageFound = true
// 						e.Image = strings.Split(src, "/scale-to-width-down")[0]
// 						if e.Image[:10] == "data:image" {
// 							if src, ok := img.Attr("data-src"); ok {
// 								e.Image = strings.Split(src, "/scale-to-width-down")[0]
// 							}
// 						}
// 					}
// 				}
// 				doc.Find(".wikitable").Remove()
// 				e.text = tdToSkill(doc.Find("table"))
// 				if imageFound && e.typ != "" && e.Name != "" {
// 					e.class = meta.Name
// 					meta.Equip = append(meta.Equip, e)
// 				}
// 			}
// 			return
// 		})
// 	}
// 	// NOTE: this is destructive to the doc
// 	// text := strings.TrimSpace(content.Children().Remove().End().Text())
// 	// if filterText(text) {
// 	// 	text = ""
// 	// }
// 	// if len(text) > len(meta.Description) {
// 	// 	meta.Description = text
// 	// }

// 	url = fmt.Sprintf(wikiUrl+"/File:Back_-_%s.png", strings.Replace(meta.Name, " ", "_", -1))
// 	doc, err = goquery.NewDocument(url)
// 	if err != nil {
// 		panic(fmt.Sprintf("error on parsing: %s", err))
// 	}
// 	src, ok := doc.Find(".fullImageLink").Find("img").Eq(0).Attr("src")
// 	if ok {
// 		meta.Image = wikiUrl + "/" + src
// 	}

// 	// OLD METHOD:
// 	// content.Find("h2").Remove()
// 	// content.Find(".wikitable").Remove()
// 	// paras := strings.Split(content.Text(), "\n")
// 	// for _, p := range paras {
// 	// 	text := strings.TrimSpace(p)
// 	// 	if text != "" {
// 	// 		if meta.Description != "" {
// 	// 			meta.Description += "\n"
// 	// 		}
// 	// 		meta.Description += text
// 	// 	}
// 	// }
// 	// fmt.Printf("%s: [%s]\n", meta.Name, meta.Description)
// 	meta.Description = descMap[meta.Name]
// 	hybridMap[meta.Name] = meta.hybrid
// 	return meta
// }

var descMap = make(map[string]string)

func loadDescriptions() error {
	dat, err := ioutil.ReadFile("./classes/descriptions.txt")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(dat), "\n") {
		arr := strings.Split(line, ": ")
		if len(arr) != 2 {
			continue
		}
		class := arr[0]
		desc := arr[1]
		descMap[class] = desc
	}
	return nil
}

// func classesGen_old() {
// 	if err := loadDescriptions(); err != nil {
// 		panic(err)
// 	}

// 	doc, err := goquery.NewDocument(wikiUrl + "/Classes")
// 	if err != nil {
// 		panic(fmt.Sprintf("error on parsing: %s", err))
// 	}

// 	classMetadata := doc.Find("#mw-content-text").First()
// 	classMetadata.Find("li").Each(func(_ int, s1 *goquery.Selection) {
// 		classUrl, ok := s1.Find("a").Attr("href")
// 		if !ok {
// 			return
// 		}

// 		name := strings.TrimSpace(s1.Text())
// 		class := genClass(wikiUrl+classUrl, name)
// 		// fmt.Printf("%s => TEXT: [%s]\n", meta.Name, text)
// 		classes = append(classes, class)
// 	})

// 	// dumpClasses()
// 	// dumpSkills()

// 	f, err := os.Create(classesHtml)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	w := bufio.NewWriter(f)

// 	if downloadEnabled {
// 		downloadCImages()
// 	}
// 	fixClasses()
// 	outputCTable(w)
// }

func fixClasses() {
	// var err error
	sort.Slice(classes, func(i, j int) bool {
		if classes[i].Archetype < classes[j].Archetype {
			return true
		}
		if classes[i].Archetype > classes[j].Archetype {
			return false
		}
		if classes[i].Archetype == classes[j].Archetype {
			if !classes[i].Hybrid && classes[j].Hybrid {
				return true
			}
			if classes[i].Hybrid && !classes[j].Hybrid {
				return false
			}
		}
		if classes[i].Name < classes[j].Name {
			return true
		}
		if classes[i].Name > classes[j].Name {
			return false
		}
		return false
	})
	// for i, c := range classes {
	// 	c.expImg = ""
	// 	imgFile := "expansions/" + strings.Replace(c.Expansion, " ", "_", -1) + ".svg"
	// 	if _, err := os.Stat(imgFile); !os.IsNotExist(err) {
	// 		c.expImg = fmt.Sprintf("<img src=\"%s\" class=\"expansion\">", imgFile)
	// 	} else if abbr, ok := expansions[strings.ToLower(c.Expansion)]; ok {
	// 		c.expImg = abbr
	// 	}

	// 	if c.Image != "" {
	// 		if c.Image, err = cImgRtoL(c.Image, c.Name); err != nil {
	// 			panic(err)
	// 		}
	// 	}

	// 	for j, s := range c.Skills {
	// 		if s.Image != "" {
	// 			if c.Skills[j].Image, err = sImgRtoL(s.Image); err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 		replaceErrata(&classes[i].Skills[j].Image)
	// 	}

	// 	for j, e := range c.Equip {
	// 		if e.Image != "" {
	// 			if c.Equip[j].Image, err = eImgRtoL(e.Image); err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 		replaceErrata(&classes[i].Equip[j].Image)
	// 	}

	// 	classes[i] = c
	// }
}

func classUniqSortExps() []string {
	expMap := make(map[string]struct{})
	for _, c := range classes {
		expMap[expansions[strings.ToLower(c.Expansion)]] = struct{}{}
	}
	var exps []string
	for exp := range expMap {
		exps = append(exps, exp)
	}
	sort.Strings(exps)
	return exps
}

func outputCHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Journeys in Class Selection</title>\n")
	fmt.Fprintf(w, "<meta name=\"description\" content=\"%s\">\n", `With over 25+ classes and 200+ skills to choose from, it's painful to choose a character.

For owners of Descent: Journeys in the Dark (Second Edition), this Class Selector makes the decision that much easier for newbies, casuals, and veterans.

Send your heroes to get some Coufee and they'll be adventuring in no time!`)
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.min.js\"></script>\n")
	fmt.Fprintf(w, "<script src=\"lib/tablesorter/jquery.tablesorter.widgets.min.js\"></script>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"heroes.css?version=%s\">\n", version)
	fmt.Fprintf(w, "<link rel=\"icon\" type=\"image/png\" href=\"etc/favicon.png\">\n")
	fmt.Fprintf(w, "</head><body onload=\"onload()\">\n")

	// table
	fmt.Fprintf(w, "<table id=\"classTable\" class=\"tablesorter\"><thead class=\"classes\"><tr style=\"width:100%\">\n")
	fmt.Fprintf(w, "<th class=\"expansion\">Exp</th>\n")
	fmt.Fprintf(w, "<th class=\"class\">Class</th>\n")
	fmt.Fprintf(w, "<th class=\"equipment\">Starting Equipment</th>\n")
	fmt.Fprintf(w, "<th class=\"skill\">Skill Tree</th>\n")
	fmt.Fprintf(w, "</tr>\n\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<td class=\"expansion\"><div><select id=\"selectExp\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	exps := classUniqSortExps()
	for _, exp := range exps {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", exp, exp)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"class\"><div><select id=\"selectClass\" onchange=\"trigger(this)\">\n")
	fmt.Fprintf(w, "<option value=\"\"></option>\n")
	for _, k := range archetypes {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(k), k)
	}
	for _, k := range archetypes {
		switch k {
		case "Mage":
			k += "/Warrior"
		case "Warrior":
			k += "/Mage"
		case "Scout":
			k += "/Healer"
		case "Healer":
			k += "/Scout"
		}
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>\n", strings.ToLower(k), k)
	}
	fmt.Fprintf(w, "</select></div></td>\n")
	fmt.Fprintf(w, "<td class=\"cSearch\" colspan=\"2\"><div class=\"search\">\n")
	fmt.Fprintf(w, "<input type=\"text\" class=\"search-input\" id=\"search-input\" name=\"search\" placeholder=\"Search\" onkeyup=\"search()\"/>\n")
	fmt.Fprintf(w, "<input type=\"submit\" class=\"search-submit\"/></div></td>\n")
	fmt.Fprintf(w, "</tr></thead><tbody class=\"classes\">\n\n")
}

func outputCTableRow(w *bufio.Writer, c1 []Class, c2 *Class) {
	if c1[0].Expansion == "Sands Of The Past" {
		return
	}
	arch1 := strings.ToLower(c1[0].Archetype)
	var arch2 string
	// if c1.Hybrid {
	// 	arch2 = strings.ToLower(c1.Archetype) + "/" + strings.ToLower(c2.Archetype)
	// }
	exp := expansions[strings.ToLower(c1[0].Expansion)]
	fmt.Fprintf(w, "<tr class=\"%s %s %s\" style=\"display:none;\">\n", arch1, arch2, exp)
	fmt.Fprintf(w, "<td class=\"expansion\"><img src=\"expansions/%s\" class=\"expansion\"></td>\n",
		expImgs[strings.ToLower(c1[0].Expansion)])
	fmt.Fprintf(w, "<td class=\"class\">")
	desc := descMap[c1[0].Class]
	// if c1.Hybrid {
	// 	desc += "\n"
	// 	desc += descMap[c2[0].Class]
	// }
	fmt.Fprintf(w, "<span title=\"%s\">", html.EscapeString(desc))
	// fmt.Fprintf(w, "<a href=\"%s\" class=\"class\">", c1.Url)
	fmt.Fprintf(w, "<div class=\"divImage\">")
	// exists := true
	// if _, err := os.Stat(c1.Image); os.IsNotExist(err) {
	// 	exists = false
	// 	arch := strings.ToLower(c1.Archetype)
	// 	if c1.Hybrid {
	// 		if arch == "mage" {
	// 			c1.Image = "classes/generic_mage_warrior.png"
	// 		} else if arch == "warrior" {
	// 			c1.Image = "classes/generic_warrior_mage.png"
	// 		} else if arch == "scout" {
	// 			c1.Image = "classes/generic_scout_healer.png"
	// 		} else if arch == "healer" {
	// 			c1.Image = "classes/generic_healer_scout.png"
	// 		}
	// 	} else {
	// 		c1.Image = "classes/generic_" + arch + ".png"
	// 	}
	// }
	image := "classes/generic_" + arch1 + ".png"
	fmt.Fprintf(w, "<img src=\"%s\" class=\"class\">", image)
	fmt.Fprintf(w, "<div class=\"className\">%s</div>", c1[0].Class)
	// if !exists {
	// 	if !c1.Hybrid {
	// 		fmt.Fprintf(w, "<div class=\"className\">%s</div>", c1.Name)
	// 	} else {
	// 		fmt.Fprintf(w, "<div class=\"className hybrid\">%s<br>%s</div>", c1.Name, c2.Name)
	// 	}
	// }
	if len(c1[0].Archetype) > 0 {
		fmt.Fprintf(w, "%c", c1[0].Archetype[0])
	}
	fmt.Fprintf(w, "</div>")
	// fmt.Fprintf(w, "</a></span></td>\n")
	fmt.Fprintf(w, "</span></td>\n")

	sort.Slice(c1, func(i, j int) bool {
		// Item -> Skill -> Familiar
		if c1[i].Type == "item" {
			return true
		}
		if c1[i].Type == "skill" {
			if c1[j].Type == "item" {
				return false
			} else {
				return true
			}
		}
		if c1[i].Type == "familiar" {
			return false
		}
		return false
	})
	// if c1.Hybrid {
	// 	sort.Slice(c2.Equip, func(i, j int) bool {
	// 		// Item -> Skill -> Familiar
	// 		if c2.Equip[i].typ == "Item" {
	// 			return true
	// 		}
	// 		if c2.Equip[i].typ == "Skill" {
	// 			if c2.Equip[j].typ == "Item" {
	// 				return false
	// 			} else {
	// 				return true
	// 			}
	// 		}
	// 		if c2.Equip[i].typ == "Familiar" {
	// 			return false
	// 		}
	// 		return false
	// 	})
	// }
	fmt.Fprintf(w, "<td class=\"equipment\">")

	var elementalistCards []Class
	for _, e := range c1 {
		if e.XpCost == 9 {
			elementalistCards = append(elementalistCards, e)
		}
		if e.Type == "skill" && e.XpCost > 0 {
			continue
		}
		if e.Image == "" {
			e.Image = "skills/blank.png"
		}
		// fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\">", e.Image, e.typ)
		// t := strings.Join(e.Traits, " ")
		t := e.Traits
		// cardContainer
		if e.Name == "Elemental Focus" {
			fmt.Fprintf(w, "<div class=\"cardContainer\">")
			fmt.Fprintf(w, "<img src=\"images/%s\" class=\"equipment e%s\" alt=\"%s\" text=\"%s\" cost=\"%s\" xp=\"%d\" ranged=\"%s\" traits=\"%s\">", e.Image, e.Type, e.Name, e.Rules, e.FatigueCost, e.XpCost, e.Attack, t)
			for _, c := range elementalistCards {
				fmt.Fprintf(w, "<img src=\"images/%s\" class=\"agentBack equipment e%s\" alt=\"%s\" text=\"%s\" cost=\"%s\" xp=\"%d\" ranged=\"%s\" traits=\"%s\">", c.Image, c.Type, c.Name, c.Rules, c.FatigueCost, c.XpCost, c.Attack, t)
			}
			fmt.Fprintf(w, "</div>")
		} else {
			fmt.Fprintf(w, "<img src=\"images/%s\" class=\"equipment e%s\" alt=\"%s\" text=\"%s\" cost=\"%s\" xp=\"%d\" ranged=\"%s\" traits=\"%s\">", e.Image, e.Type, e.Name, e.Rules, e.FatigueCost, e.XpCost, e.Attack, t)
		}
	}
	// if c1.Hybrid {
	// 	for _, e := range c2.Equip {
	// 		if e.Image == "" {
	// 			e.Image = "skills/blank.png"
	// 		}
	// 		// fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\">", e.Image, e.typ)
	// 		t := strings.Join(e.traits, " ")
	// 		fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\" alt=\"%s\" text=\"%s\" cost=\"%d\" xp=\"%d\" ranged=\"%s\" traits=\"%s\">", e.Image, e.typ, e.Name, e.text, e.cost, e.xp, e.ranged, t)
	// 	}
	// }
	fmt.Fprintf(w, "</td>")
	fmt.Fprintf(w, "<td class=\"skill\">")

	// var skillPool []skill
	// for _, s := range c1.skills {
	// 	if c1.Hybrid {
	// 		if s.xp == 0 {
	// 			continue
	// 		}
	// 	}
	// 	skillPool = append(skillPool, s)
	// }
	// if c1.Hybrid {
	// 	for _, s := range c2.skills {
	// 		if s.xp == 3 {
	// 			continue
	// 		}
	// 		skillPool = append(skillPool, s)
	// 	}
	// }

	sort.Slice(c1, func(i, j int) bool {
		// Item -> Skill -> Familiar
		if c1[i].Type == "item" {
			return true
		}
		if c1[i].Type == "skill" {
			if c1[j].Type == "item" {
				return false
			} else {
				return true
			}
		}
		if c1[i].Type == "familiar" {
			return false
		}
		return false
	})
	sort.Slice(c1, func(i, j int) bool {
		if c1[i].XpCost < c1[j].XpCost {
			return true
		}
		if c1[i].XpCost > c1[j].XpCost {
			return false
		}
		// if hybridMap[c1[i].class] && !hybridMap[c1[j].class] {
		// 	return true
		// }
		// if !hybridMap[c1[i].class] && hybridMap[c1[j].class] {
		// 	return false
		// }
		return c1[i].Class < c1[j].Class
	})
	for _, s := range c1 {
		if s.Type != "skill" || s.XpCost == 0 || s.XpCost == 9 {
			continue
		}
		img := s.Image
		if img == "" {
			img = "skills/blank.png"
		}
		fmt.Fprintf(w, "<img src=\"images/%s\" class=\"skill\" alt=\"%s\" text=\"%s\" cost=\"%d\" xp=\"%d\">", img, s.Name, s.Rules, s.FatigueCost, s.XpCost)
	}
	fmt.Fprintf(w, "</td></tr>\n\n")
}

// var classCards = map[string]map[string][]Class{}

func outputCTable(w *bufio.Writer) {
	outputCHeader(w)

	for _, v1 := range classCards {
		var cards []Class
		// if len(v1) > 2 || len(v1) < 1 {
		// 	panic("uh oh")
		// }
		for _, v2 := range v1 {
			if len(v2) == 1 {
				cards = append(cards, v2[0])
				continue
			}
			var card *Class
			for i := range v2 {
				if card == nil {
					card = &v2[i]
					continue
				}
				if strings.Contains(v2[i].Image, "errata.png") {
					card = &v2[i]
				}
			}
			cards = append(cards, *card)
		}
		outputCTableRow(w, cards, nil)
		// a := strings.ToLower(c.Archetype)
		// for _, c2 := range classes {
		// 	if c2.Hybrid {
		// 		continue
		// 	}
		// 	a2 := strings.ToLower(c2.Archetype)
		// 	if a == "mage" || a2 == "mage" {
		// 		if a == "warrior" || a2 == "warrior" {
		// 			outputCTableRow(w, c, &c2)
		// 		}
		// 	}
		// 	if a == "scout" || a2 == "scout" {
		// 		if a == "healer" || a2 == "healer" {
		// 			outputCTableRow(w, c, &c2)
		// 		}
		// 	}
		// }
	}

	// for _, c := range classes {
	// 	if !c.Hybrid {
	// 		outputCTableRow(w, c, nil)
	// 		continue
	// 	}

	// 	a := strings.ToLower(c.Archetype)
	// 	for _, c2 := range classes {
	// 		if c2.Hybrid {
	// 			continue
	// 		}
	// 		a2 := strings.ToLower(c2.Archetype)
	// 		if a == "mage" || a2 == "mage" {
	// 			if a == "warrior" || a2 == "warrior" {
	// 				outputCTableRow(w, c, &c2)
	// 			}
	// 		}
	// 		if a == "scout" || a2 == "scout" {
	// 			if a == "healer" || a2 == "healer" {
	// 				outputCTableRow(w, c, &c2)
	// 			}
	// 		}
	// 	}
	// }

	outputCFooter(w)
	w.Flush()
}

func outputCFooter(w *bufio.Writer) {
	outputFooter(w, "classes", 4)
}
