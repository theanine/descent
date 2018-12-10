package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

const consoleHtml = "console.html"

// type Item struct {
// 	Name      string `json:"name"`
// 	Points    int    `json:"points"`
// 	Traits    string `json:"traits"`
// 	Attack    string `json:"attack"`
// 	Equip     string `json:"equip"`
// 	Dice      string `json:"dice"`
// 	Rules     string `json:"rules"`
// 	Expansion string `json:"expansion"`
// 	Image     string `json:"image"`
// 	Xws       string `json:"xws"`

// 	// Shop Items
// 	Count int    `json:"count"`
// 	Act   string `json:"act"`
// 	Cost  int    `json:"cost"`

// 	// Class Items
// 	Archetype string `json:"archetype"`
// 	Class     string `json:"class"`
// }

var Kitems []Item

type Familiar struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	Speed     string `json:"speed"`
	Health    string `json:"health"`
	Defense   string `json:"defense"`
	Attack    string `json:"attack"`
	Dice      string `json:"type"`
	Traits    string `json:"traits"`
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	Xws       string `json:"xws"`
}

var Kfamiliars []Familiar

type Hero struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	Archetype string `json:"archetype"`
	Speed     int    `json:"speed"`
	Health    int    `json:"health"`
	Stamina   int    `json:"stamina"`
	Defense   string `json:"defense"`
	Willpower int    `json:"willpower"`
	Might     int    `json:"might"`
	Knowledge int    `json:"knowledge"`
	Awareness int    `json:"awareness"`
	Ability   string `json:"ability"`
	Feat      string `json:"feat"`
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	Xws       string `json:"xws"`
}

var Kheroes []Hero
var KheroMap map[string]*Hero

type Skill struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	Archetype string `json:"archetype"`
	Class     string `json:"class"`
	XP        int    `json:"xp cost"`
	Rules     string `json:"rules"`
	Fatigue   int    `json:"fatigue"`
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	Xws       string `json:"xws"`
}

var Kskills []Skill

type Condition struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	Count     int    `json:"count"`
	Traits    string `json:"traits"`
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	Xws       string `json:"xws"`
}

var Kconditions []Condition

type Treasure struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	Count     int    `json:"count"`
	Traits    string `json:"traits"`
	Ability   string `json:"ability"`
	Gold      int    `json:"gold"`
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	Xws       string `json:"xws"`
}

var Ktreasures []Treasure

var KimageMap map[string]string
var KnameList []string

func loadItems() error {
	var shopItems []Item
	dat1, err := ioutil.ReadFile("data/shop-items.js")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(dat1, &shopItems); err != nil {
		return err
	}
	Kitems = append(Kitems, shopItems...)

	dat2, err := ioutil.ReadFile("data/class-items.js")
	if err != nil {
		return err
	}
	var classItems []Item
	if err = json.Unmarshal(dat2, &classItems); err != nil {
		return err
	}
	Kitems = append(Kitems, classItems...)

	dat3, err := ioutil.ReadFile("data/relics.js")
	if err != nil {
		return err
	}
	var relics []Item
	if err = json.Unmarshal(dat3, &relics); err != nil {
		return err
	}
	Kitems = append(Kitems, relics...)
	return nil
}

func loadFamiliars() error {
	dat, err := ioutil.ReadFile("data/class-familiars.js")
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, &Kfamiliars)
}

func loadHeroes() error {
	KheroMap = make(map[string]*Hero)
	dat, err := ioutil.ReadFile("data/heroes.js")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dat, &Kheroes); err != nil {
		return err
	}
	for i, h := range Kheroes {
		KheroMap[h.Name] = &Kheroes[i]
	}
	return nil
}

func loadSkills() error {
	dat1, err := ioutil.ReadFile("data/class-skills.js")
	if err != nil {
		return err
	}
	var cSkills []Skill
	if err = json.Unmarshal(dat1, &cSkills); err != nil {
		return err
	}
	Kskills = append(Kskills, cSkills...)

	dat2, err := ioutil.ReadFile("data/hybrid-class-skills.js")
	if err != nil {
		return err
	}
	var hSkills []Skill
	if err = json.Unmarshal(dat2, &hSkills); err != nil {
		return err
	}
	Kskills = append(Kskills, hSkills...)
	return nil
}

func loadConditions() error {
	dat, err := ioutil.ReadFile("data/conditions.js")
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, &Kconditions)
}

func loadTreasures() error {
	dat, err := ioutil.ReadFile("data/search-deck.js")
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, &Ktreasures)
}

func loadJson() error {
	if err := loadItems(); err != nil {
		return err
	}
	if err := loadFamiliars(); err != nil {
		return err
	}
	if err := loadHeroes(); err != nil {
		return err
	}
	if err := loadSkills(); err != nil {
		return err
	}
	if err := loadConditions(); err != nil {
		return err
	}
	if err := loadTreasures(); err != nil {
		return err
	}
	return nil
}

func loadMap() {
	KimageMap = make(map[string]string)
	for _, i := range items {
		KimageMap[i.Name] = i.Image
		KnameList = append(KnameList, i.Name)
	}
	for _, f := range Kfamiliars {
		KimageMap[f.Name] = f.Image
		KnameList = append(KnameList, f.Name)
	}
	for _, h := range Kheroes {
		KimageMap[h.Name] = h.Image
		KnameList = append(KnameList, h.Name)
	}
	for _, s := range Kskills {
		KimageMap[s.Name] = s.Image
		KnameList = append(KnameList, s.Name)
	}
	for _, c := range Kconditions {
		KimageMap[c.Name] = c.Image
		KnameList = append(KnameList, c.Name)
	}
	for _, t := range Ktreasures {
		KimageMap[t.Name] = t.Image
		KnameList = append(KnameList, t.Name)
	}
	sort.Strings(KnameList)
}

func consoleGen() {
	if err := loadJson(); err != nil {
		log.Fatal(err)
	}
	loadMap()

	f, err := os.Create(consoleHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	outputKTable(w)
}

func outputKHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Overlord Dashboard</title>\n")
	fmt.Fprintf(w, "<meta name=\"description\" content=\"%s\">\n", `It's painful to track all skills, conditions, health, and fatigue for all 4 heroes.

For owners of Descent: Journeys in the Dark (Second Edition), this Overlord Dashboard makes it that much easier for newbies, casuals, and veterans.

Send your overlords to get some Coufee and they'll be adventuring in no time!`)
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"console.css?version=%s\">\n", version)
	fmt.Fprintf(w, "<script type=\"text/javascript\" src=\"console.js?version=%s\"></script>\n", version)
	fmt.Fprintf(w, "<link rel=\"icon\" type=\"image/png\" href=\"etc/favicon.png\">\n")
	fmt.Fprintf(w, "</head><body onload=\"onload()\">\n")
	fmt.Fprintf(w, "<div id=\"back\"></div>\n")
}

func outputKQuad(w *bufio.Writer, num int) {
	fmt.Fprintf(w, "<div id=\"search%d\" class=\"search\">\n", num)
	// fmt.Fprintf(w, "<input type=\"text\" class=\"search-input\" id=\"search-input\" name=\"search\" placeholder=\"Search\" onkeyup=\"search()\"/>\n")
	// fmt.Fprintf(w, "<input type=\"submit\" class=\"search-submit\"/></div></td>\n")
	fmt.Fprintf(w, "<input type=\"text\" class=\"search-input\" onkeyup=\"search(this, event)\" onfocus=\"gotFocus(this)\" onblur=\"lostFocus(this)\" placeholder=\"Search\" title=\"???\">\n")
	fmt.Fprintf(w, "<ul class=\"myUL\" style=\"display:none;\">\n")
	for _, name := range KnameList {
		var info string
		if h, ok := KheroMap[name]; ok {
			info = fmt.Sprintf("health=\"%d\" stamina=\"%d\"", h.Health, h.Stamina)
		}
		fmt.Fprintf(w, "<li class=\"card\"><a href=\"%s\" %s>%s</a></li>\n", KimageMap[name], info, name)
	}
	fmt.Fprintf(w, "</ul>\n")
	fmt.Fprintf(w, "<div id=\"hero%d\" class=\"hero\"></div>\n", num)
	fmt.Fprintf(w, "<div id=\"stats%d\" class=\"stats\">\n", num)
	fmt.Fprintf(w, "<img src=\"attributes/health-icon.png\" class=\"health\">\n")
	fmt.Fprintf(w, "<div id=\"health%d\" class=\"health\"></div>\n", num)
	fmt.Fprintf(w, "<img src=\"attributes/fatigue-icon.png\" class=\"fatigue\">\n")
	fmt.Fprintf(w, "<div id=\"fatigue%d\" class=\"fatigue\"></div>\n", num)
	fmt.Fprintf(w, "</div>\n")
	fmt.Fprintf(w, "</div>\n")
}

func outputKTable(w *bufio.Writer) {
	outputKHeader(w)

	for i := 0; i < 4; i++ {
		outputKQuad(w, i)
	}

	outputKFooter(w)
	w.Flush()
}

func outputKFooter(w *bufio.Writer) {
	fmt.Fprintf(w, "</body></html>")
}
