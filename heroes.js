// <SCRIPT LANGUAGE="JavaScript">

// First we get the viewport height and we multiple it by 1% to get a value for a vh unit
let vh = window.innerHeight * 0.01;
// Then we set the value in the --vh custom property to the root of the document
document.documentElement.style.setProperty('--vh', `${vh}px`);

// We listen to the resize event
window.addEventListener('resize', () => {
  // We execute the same script as before
  let vh = window.innerHeight * 0.01;
  document.documentElement.style.setProperty('--vh', `${vh}px`);
});

var attrs = ["speed",
			 "health",
			 "stamina",
			 "might",
			 "willpower",
			 "knowledge",
			 "awareness"];

function toTitleCase(str) {
	return str.replace(/(?:^|\s)\w/g, function(match) {
		return match.toUpperCase();
	});
}

function trigger(select)
{
	showHideArrows($(select));
	showHideRows();
	doSearch();
	addJumps();
	loadSelects();
	/*
	showHideRows();
	addJumps();
	loadSelects();
	doSearch();
	*/
}

function onload()
{
	$("select").each(function() {
		showHideArrows($(this));
	});
	showHideRows();
	addJumps();
	loadSelects();
	// $("#heroTable").tablesorter();
	
	if ($("#heroTable").length) {
		$("#heroTable").tablesorter({
			sortList       : [[2,0]],
			// widgets        : ['zebra', 'columns'],
			usNumberFormat : true,
			sortReset      : false,
			sortRestart    : true
		});
	}
}

function doSearch() {
	var input = document.getElementById('search-input');
	if (input == null)
		return;
	var filter = input.value.toUpperCase();
	$("tr", "table").each(function(i, tr){
		var trHide = null
		// CLASSES
		$(tr).find('img.equipment,img.skill').each(function(j, img){
			if (trHide == null) {
				trHide = true
			}
			var alt = $(img).attr("alt");
			var text = $(img).attr("text");
			var ranged = $(img).attr("ranged");
			var traits = $(img).attr("traits");
			if ((filter == "") ||
				(typeof alt !== "undefined" && alt.toUpperCase().indexOf(filter) > -1) ||
				(typeof type !== "undefined" && type.toUpperCase().indexOf(filter) > -1) ||
				(typeof ranged !== "undefined" && ranged.toUpperCase().indexOf(filter) > -1) ||
				(typeof traits !== "undefined" && traits.toUpperCase().indexOf(filter) > -1)) {
				$(img).show();
				trHide = false
			} else {
				$(img).hide();
			}
		});
		// HEROES
		$(tr).find('td.ability,td.heroic').each(function(j, td){
			if (trHide == null) {
				trHide = true
			}
			var str = $(td).html();
			// first create an element and add the string as its HTML
			var container = $('<div>').html(str);
			// then use .replaceWith(function()) to modify the HTML structure
			container.find('img').replaceWith(function() { return this.alt; })
			// finally get the HTML back as string
			var text = container.html();
			
			if (filter == "") {
				trHide = false
			} else if (typeof text === "undefined") {
			} else if (text.toUpperCase().indexOf(filter) > -1) {
				trHide = false
			}
		});
		// OVERLORD / PLOT
		$(tr).find('div.cardContainer').each(function(j, div){
			if (trHide == null) {
				trHide = true
			}
			var alt = $(div).find("img").attr("alt");
			var type = $(div).find("img").attr("type");
			var text = $(div).find("img").attr("text");
			var exp = $(div).find("img").attr("exp");
			var num = $(div).find("img").attr("num");
			if ((filter == "") ||
				(typeof alt !== "undefined" && alt.toUpperCase().indexOf(filter) > -1) ||
				(typeof type !== "undefined" && type.toUpperCase().indexOf(filter) > -1) ||
				(typeof text !== "undefined" && text.toUpperCase().indexOf(filter) > -1) ||
				(typeof exp !== "undefined" && exp.toUpperCase().indexOf(filter) > -1) ||
				(typeof num !== "undefined" && num.toUpperCase().indexOf(filter) > -1)) {
				$(div).show();
				trHide = false
			} else {
				$(div).hide();
			}
		});
		if (trHide) {
			$(tr).hide();
		}
	});
}

function search() {
	showHideRows();
	doSearch();
	addJumps();
	loadSelects();
}

function showHideArrows(select)
{
	if (select.val() == "") {
		select.removeClass("blank");
		select.addClass("normal");
	} else {
		select.removeClass("normal");
		select.addClass("blank");
	}
}

function showHideRow(tr)
{
	var e = $("#selectExp");
	var k = $("#selectCK");
	var c = $("#selectClass");
	var d = $("#selectDefense");
	var o = $("#selectCards");
	
	$(tr).show();
	if (typeof $(tr).attr("class") === "undefined" || $(tr).attr("class") === "tablesorter-headerRow")
		return;
	if (c.length && c.val() != "" && !$(tr).hasClass(c.val()))
		$(tr).hide();
	if (d.length && d.val() != "" && !$(tr).hasClass(d.val()))
		$(tr).hide();
	if (k.length && k.val() != "" && !$(tr).hasClass(k.val()))
		$(tr).hide();
	if (e.length && e.val() != "" && !$(tr).hasClass(e.val()))
		$(tr).hide();
	
	$(tr).find('div.cardContainer').each(function(index, div){
		$(div).show()
		if (o.length && o.val() != "") {
			if (typeof $(div).attr("class") === "undefined")
				return;
			if (!$(div).hasClass(o.val()))
				$(div).hide();
		}
	});
	
	$.each(attrs, function (i, attr) {
		var s = $("#select"+toTitleCase(attr));
		if (s.length && s.val() != "" && $('td.'+attr, tr).text() != s.val())
			$(tr).hide();
	});

	$(tr).find(".cards div.cardContainer").each(function(index, e){removeJumps(e)});
	$(tr).find(".skill img.skill").each(function(index, e){removeJumps(e)});
	$(tr).find(".equipment img.equipment").each(function(index, e){removeJumps(e)});
}

function removeJumps(e)
{
	$(e).removeClass("jumpDown");
	$(e).removeClass("jumpUp");
	$(e).removeClass("jumpLeft");
	$(e).removeClass("jumpRight");
}

function specialJump(e, jump)
{
	$(e).addClass(jump);
}

function showHideRows()
{
	$("tr", "#heroTable").each(function(index, tr){showHideRow(tr)});
	$("tr", "#classTable").each(function(index, tr){showHideRow(tr)});
	$("tr", "#overlordTable").each(function(index, tr){showHideRow(tr)});
	$("tr", "#plotTable").each(function(index, tr){showHideRow(tr)});
}

function addJumps()
{
	var td = $("#classTable tbody").find("tr:visible:first").find(".skill");
	var card = td.find("img.skill");
	var cardWidth = Math.floor(td.width() / card.width());
	$("#classTable tbody").find("tr:visible:first").find(".skill img.skill").each(function(index, e){specialJump(e, "jumpDown")});
	$("#classTable tbody").find("tr:visible:last").find(".skill img.skill").each(function(index, e){specialJump(e, "jumpUp")});
	$("#classTable tbody").find("tr:visible").find(".skill img.skill:eq("+(cardWidth-1)+")").each(function(index, e){specialJump(e, "jumpLeft")});
	$("#classTable tbody").find("tr:visible:first").find(".equipment img.equipment").each(function(index, e){specialJump(e, "jumpDown")});
	$("#classTable tbody").find("tr:visible:last").find(".equipment img.equipment").each(function(index, e){specialJump(e, "jumpUp")});

	var td = $("#overlordTable tbody").find("tr:visible:first").find(".cards");
	var card = td.find("div.cardContainer");
	var cardWidth = Math.floor(td.width() / card.width());
	$("#overlordTable tbody").find("tr:visible:first").find(".cards div.cardContainer").each(function(index, e){specialJump(e, "jumpDown")});
	$("#overlordTable tbody").find("tr:visible:last").find(".cards div.cardContainer").each(function(index, e){specialJump(e, "jumpUp")});
	$("#overlordTable tbody").find("tr:visible").find(".cards div.cardContainer:visible:eq("+(cardWidth-1)+")").each(function(index, e){specialJump(e, "jumpLeft")});

	var td = $("#plotTable tbody").find("tr:visible:first").find(".cards");
	var card = td.find("div.cardContainer");
	var cardWidth = Math.floor(td.width() / card.width());
	$("#plotTable tbody").find("tr:visible:first").find(".cards div.cardContainer").each(function(index, e){specialJump(e, "jumpDown")});
	$("#plotTable tbody").find("tr:visible:last").find(".cards div.cardContainer").each(function(index, e){specialJump(e, "jumpUp")});
	$("#plotTable tbody").find("tr:visible").find(".cards div.cardContainer:visible:eq("+(cardWidth-1)+")").each(function(index, e){specialJump(e, "jumpLeft")});
}

function uniquesort(arr) {
	if (arr.length === 0)
		return arr;
	arr = arr.sort(function (a, b) { return a*1 - b*1; });
	var a = [];
	var l = arr.length;
	for(var i=0; i<l; i++) {
		for(var j=i+1; j<l; j++) {
			if (arr[i] === arr[j])
				j = ++i;
		}
		a.push(arr[i]);
	}
	return a;
}

function loadSelects()
{
	$.each(attrs, function (i, attr) {
		var selected = $('#select'+toTitleCase(attr)).val();
		var arr = [];
		$("tr", "#heroTable").each(function(index, tr) {
			if (typeof $(tr).attr("class") === "undefined" || $(tr).attr("class") === "tablesorter-headerRow")
				return;
			$('td.'+attr, tr).each(function(index, td) {
				if (selected != "" || $(tr).is(":visible"))
					arr.push($(td).text());
			});
		});
		arr = uniquesort(arr);
		$('#select'+toTitleCase(attr)).empty().append($('<option>', { value: "", text : "" }));
		$.each(arr, function (i, item) {
			$('#select'+toTitleCase(attr)).append($('<option>', { value: item, text : item }));
		});
		$('#select'+toTitleCase(attr)).val(selected);
	});
}

function colorizeCells(cells, red, orange, yellow, green, dkGreen) {
	for (var i=0, len=cells.length; i<len; i++) {
		val = parseInt(cells[i].innerHTML, 10)
		if (val <= red) {
			cells[i].style.backgroundColor = '#f8696b';
		} else if (val == orange) {
			cells[i].style.backgroundColor = '#fbaa77' ;
		} else if (val == yellow) {
			cells[i].style.backgroundColor = '#ffeb84' ;
		} else if (val == green) {
			cells[i].style.backgroundColor = '#cbdc81' ;
		} else if (val >= dkGreen) {
			cells[i].style.backgroundColor = '#63be7b' ;
		}
	}
}

colorizeCells(document.getElementsByClassName("speed"), 2, 3, 4, -1, 5)
colorizeCells(document.getElementsByClassName("health"), 8, -1, 10, 12, 14)
colorizeCells(document.getElementsByClassName("stamina"), 3, -1, 4, 5, 6)
colorizeCells(document.getElementsByClassName("might"), 1, 2, 3, 4, 5)
colorizeCells(document.getElementsByClassName("willpower"), 1, 2, 3, 4, 5)
colorizeCells(document.getElementsByClassName("knowledge"), 1, 2, 3, 4, 5)
colorizeCells(document.getElementsByClassName("awareness"), 1, 2, 3, 4, 5)

function myFunction() {
    var popup = document.getElementById("myPopup");
    popup.classList.toggle("show");
}

function rateSpeed(speed) {
	return (speed - 4);
}

function rateStamina(stamina) {
	return (stamina - 4);
}

function rateHP(hp) {
	return (hp - 12)/2;
}

function rateDef(def) {
	switch (def) {
		case 1: // brown
			return -1;
		case 2: // grey
			return 0;
		case 3: // black
			return 1;
	}
}

function rateAttr(arch, might, know, will, aware) {
	var total = might + know + will + aware;
	if (might == 3 && know == 3 && will == 3 && aware == 3) {
		return 1;
	}
	if (total != 11) {
		return null;	// error
	}
	if (arch == "W" && might < 3) {
		return null;	// error
	}
	if (arch == "M" && know < 3) {
		return null;	// error
	}
	if (arch == "S" && aware < 3) {
		return null;	// error
	}
	if (arch == "H" && will < 3) {
		return null;	// error
	}
	return 0;
}

function colorizeCells(cells, red, orange, yellow, green, dkGreen) {
	for (var i=0, len=cells.length; i<len; i++) {
		val = parseInt(cells[i].innerHTML, 10)
		if (val <= red) {
			cells[i].style.backgroundColor = '#f8696b';
		} else if (val == orange) {
			cells[i].style.backgroundColor = '#fbaa77' ;
		} else if (val == yellow) {
			cells[i].style.backgroundColor = '#ffeb84' ;
		} else if (val == green) {
			cells[i].style.backgroundColor = '#cbdc81' ;
		} else if (val >= dkGreen) {
			cells[i].style.backgroundColor = '#63be7b' ;
		}
	}
}

function colorizeAbilityFeat() {
	// var vals = [];
	$('tbody.heroes > tr').each(function(index, tr) {
		var total = 0;
		var speed = parseInt($(tr).find("td.speed").text(), 10);
		total += rateSpeed(speed);
		var hp = parseInt($(tr).find("td.health").text(), 10);
		total += rateHP(hp);
		var stamina = parseInt($(tr).find("td.stamina").text(), 10);
		total += rateStamina(stamina);
		var def = parseInt($(tr).find("td.dice > div.die").text(), 10);
		total += rateDef(def);
		var name = $(tr).find("td.hero").text();
		var arch = $(tr).find("td.image div.divImage").text();
		var might = parseInt($(tr).find("td.might").text(), 10);
		var know = parseInt($(tr).find("td.knowledge").text(), 10);
		var will = parseInt($(tr).find("td.willpower").text(), 10);
		var aware = parseInt($(tr).find("td.awareness").text(), 10);
		total += rateAttr(arch, might, know, will, aware);
		switch (total) {
			case 2:
				$(tr).find("td.ability").css('backgroundColor', '#f8696b');
				$(tr).find("td.heroic").css('backgroundColor', '#f8696b');
				break;
			case 1:
				$(tr).find("td.ability").css('backgroundColor', '#fbaa77');
				$(tr).find("td.heroic").css('backgroundColor', '#fbaa77');
				break;
			case 0:
				$(tr).find("td.ability").css('backgroundColor', '#FFFFFF');
				$(tr).find("td.heroic").css('backgroundColor', '#FFFFFF');
				break;
			case -1:
				$(tr).find("td.ability").css('backgroundColor', '#FFEB84');
				$(tr).find("td.heroic").css('backgroundColor', '#FFEB84');
				break;
			case -2:
				$(tr).find("td.ability").css('backgroundColor', '#CBDC81');
				$(tr).find("td.heroic").css('backgroundColor', '#CBDC81');
				break;
			case -3:
				$(tr).find("td.ability").css('backgroundColor', '#97CD7E');
				$(tr).find("td.heroic").css('backgroundColor', '#97CD7E');
				break;
			case -5:
				$(tr).find("td.ability").css('backgroundColor', '#63BE7B');
				$(tr).find("td.heroic").css('backgroundColor', '#63BE7B');
				break;
		}
	});
	// vals = uniquesort(vals);
	// alert(vals);
	// dkG  G  dkY  Y  W   O  R
	// [-5, -3, -2, -1, 0, 1, 2];
}
colorizeAbilityFeat();

// </SCRIPT>
