package main

// shop.html
// Name	Exp	Gold	Qty	Equip	AttackType	Dice	Traits	Text
const itemsHtml = "items.html"

type Item struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	Traits    string `json:"traits"`
	Attack    string `json:"attack"`
	Equip     string `json:"equip"`
	Dice      string `json:"dice"`
	Rules     string `json:"rules"`
	Expansion string `json:"expansion"`
	Image     string `json:"image"`
	Xws       string `json:"xws"`

	// Shop Items
	Count int    `json:"count"`
	Act   string `json:"act"`
	Cost  int    `json:"cost"`

	// Class Items
	Archetype string `json:"archetype"`
	Class     string `json:"class"`
}

var items []Item

/*
func loadXML() error {
	var shopItems []Item
	dat1, err := ioutil.ReadFile("../d2e-master/data/shop-items.js")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(dat1, &shopItems); err != nil {
		return err
	}
	items = append(items, shopItems...)

	dat2, err := ioutil.ReadFile("../d2e-master/data/class-items.js")
	if err != nil {
		return err
	}
	var classItems []Item
	if err = json.Unmarshal(dat2, &classItems); err != nil {
		return err
	}
	items = append(items, classItems...)

	dat3, err := ioutil.ReadFile("../d2e-master/data/relics.js")
	if err != nil {
		return err
	}
	var relics []Item
	if err = json.Unmarshal(dat3, &relics); err != nil {
		return err
	}
	items = append(items, relics...)
	return nil
}

func itemsGen() {
	if err := loadXML(); err != nil {
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

	// dumpClasses()
	// dumpSkills()

	f, err := os.Create(itemsHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if downloadEnabled {
		downloadIImages()
	}
	fixItems()
	outputITable(w)
}

func fixItems() {
}

func downloadIImages() {
}

func outputIHeader(w *bufio.Writer) {
	fmt.Fprintf(w, "<html><head>\n")
	fmt.Fprintf(w, "<title>Coufee: Journeys in Item Selection</title>\n")
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

func outputITableRow(w *bufio.Writer, c1 class, c2 *class) {
	arch1 := strings.ToLower(c1.archetype)
	var arch2 string
	if c1.hybrid {
		arch2 = strings.ToLower(c1.archetype) + "/" + strings.ToLower(c2.archetype)
	}
	exp := expansions[strings.ToLower(c1.expansion)]
	fmt.Fprintf(w, "<tr class=\"%s %s %s\" style=\"display:none;\">\n", arch1, arch2, exp)
	fmt.Fprintf(w, "<td class=\"expansion\">%s</td>\n", c1.expImg)
	fmt.Fprintf(w, "<td class=\"class\">")
	if c1.hybrid {
		fmt.Fprintf(w, "<span title=\"%s\">", html.EscapeString(c1.description+"\n"+c2.description))
	} else {
		fmt.Fprintf(w, "<span title=\"%s\">", html.EscapeString(c1.description))
	}
	fmt.Fprintf(w, "<a href=\"%s\" class=\"class\">", c1.url)
	fmt.Fprintf(w, "<div class=\"divImage\">")
	exists := true
	if _, err := os.Stat(c1.img); os.IsNotExist(err) {
		exists = false
		arch := strings.ToLower(c1.archetype)
		if c1.hybrid {
			if arch == "mage" {
				c1.img = "classes/generic_mage_warrior.png"
			} else if arch == "warrior" {
				c1.img = "classes/generic_warrior_mage.png"
			} else if arch == "scout" {
				c1.img = "classes/generic_scout_healer.png"
			} else if arch == "healer" {
				c1.img = "classes/generic_healer_scout.png"
			}
		} else {
			c1.img = "classes/generic_" + arch + ".png"
		}
	}
	fmt.Fprintf(w, "<img src=\"%s\" class=\"class\">", c1.img)
	if !exists {
		if !c1.hybrid {
			fmt.Fprintf(w, "<div class=\"className\">%s</div>", c1.name)
		} else {
			fmt.Fprintf(w, "<div class=\"className hybrid\">%s<br>%s</div>", c1.name, c2.name)
		}
	}
	if len(c1.archetype) > 0 {
		fmt.Fprintf(w, "%c", c1.archetype[0])
	}
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</a></span></td>\n")

	sort.Slice(c1.equipments, func(i, j int) bool {
		// Item -> Skill -> Familiar
		if c1.equipments[i].typ == "Item" {
			return true
		}
		if c1.equipments[i].typ == "Skill" {
			if c1.equipments[j].typ == "Item" {
				return false
			} else {
				return true
			}
		}
		if c1.equipments[i].typ == "Familiar" {
			return false
		}
		return false
	})
	if c1.hybrid {
		sort.Slice(c2.equipments, func(i, j int) bool {
			// Item -> Skill -> Familiar
			if c2.equipments[i].typ == "Item" {
				return true
			}
			if c2.equipments[i].typ == "Skill" {
				if c2.equipments[j].typ == "Item" {
					return false
				} else {
					return true
				}
			}
			if c2.equipments[i].typ == "Familiar" {
				return false
			}
			return false
		})
	}
	fmt.Fprintf(w, "<td class=\"equipment\">")
	for _, e := range c1.equipments {
		if e.img == "" {
			e.img = "skills/blank.png"
		}
		// fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\">", e.img, e.typ)
		t := strings.Join(e.traits, " ")
		fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\" alt=\"%s\" text=\"%s\" cost=\"%d\" xp=\"%d\" ranged=\"%s\" traits=\"%s\">", e.img, e.typ, e.name, e.text, e.cost, e.xp, e.ranged, t)
	}
	if c1.hybrid {
		for _, e := range c2.equipments {
			if e.img == "" {
				e.img = "skills/blank.png"
			}
			// fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\">", e.img, e.typ)
			t := strings.Join(e.traits, " ")
			fmt.Fprintf(w, "<img src=\"%s\" class=\"equipment e%s\" alt=\"%s\" text=\"%s\" cost=\"%d\" xp=\"%d\" ranged=\"%s\" traits=\"%s\">", e.img, e.typ, e.name, e.text, e.cost, e.xp, e.ranged, t)
		}
	}
	fmt.Fprintf(w, "</td>")
	fmt.Fprintf(w, "<td class=\"skill\">")

	var skillPool []skill
	for _, s := range c1.skills {
		if c1.hybrid {
			if s.xp == 0 {
				continue
			}
		}
		skillPool = append(skillPool, s)
	}
	if c1.hybrid {
		for _, s := range c2.skills {
			if s.xp == 3 {
				continue
			}
			skillPool = append(skillPool, s)
		}
	}

	sort.Slice(skillPool, func(i, j int) bool {
		if skillPool[i].xp < skillPool[j].xp {
			return true
		}
		if skillPool[i].xp > skillPool[j].xp {
			return false
		}
		if hybridMap[skillPool[i].class] && !hybridMap[skillPool[j].class] {
			return true
		}
		if !hybridMap[skillPool[i].class] && hybridMap[skillPool[j].class] {
			return false
		}
		return skillPool[i].name < skillPool[j].name
	})
	for _, s := range skillPool {
		img := s.img
		if img == "" {
			img = "skills/blank.png"
		}
		fmt.Fprintf(w, "<img src=\"%s\" class=\"skill\" alt=\"%s\" text=\"%s\" cost=\"%d\" xp=\"%d\">", img, s.name, s.text, s.cost, s.xp)
	}
	fmt.Fprintf(w, "</td></tr>\n\n")
}

func outputITable(w *bufio.Writer) {
	for _, i := range items {
		fmt.Fprintf(w, "<img src=\"%s\">", "../d2e-master/images/"+i.Image)
	}

	// outputIHeader(w)

	// for _, c := range classes {
	// 	if !c.hybrid {
	// 		outputITableRow(w, c, nil)
	// 		continue
	// 	}

	// 	a := strings.ToLower(c.archetype)
	// 	for _, c2 := range classes {
	// 		if c2.hybrid {
	// 			continue
	// 		}
	// 		a2 := strings.ToLower(c2.archetype)
	// 		if a == "mage" || a2 == "mage" {
	// 			if a == "warrior" || a2 == "warrior" {
	// 				outputITableRow(w, c, &c2)
	// 			}
	// 		}
	// 		if a == "scout" || a2 == "scout" {
	// 			if a == "healer" || a2 == "healer" {
	// 				outputITableRow(w, c, &c2)
	// 			}
	// 		}
	// 	}
	// }

	// outputIFooter(w)
	w.Flush()
}

func outputIFooter(w *bufio.Writer) {
	outputFooter(w, "classes", 4)
}
*/
