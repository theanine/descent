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
}

function onload()
{
	$("select").each(function() {
		showHideArrows($(this));
	});
	showHideRows();
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
	var td = $("#classTable tbody").find("tr:visible:first").find(".skill");
	var card = td.find("img.skill");
	var cardWidth = Math.floor(td.width() / card.width());
	$("#classTable tbody").find("tr:visible:first").find(".skill img.skill").each(function(index, e){specialJump(e, "jumpDown")});
	$("#classTable tbody").find("tr:visible:last").find(".skill img.skill").each(function(index, e){specialJump(e, "jumpUp")});
	$("#classTable tbody").find("tr:visible").find(".skill img.skill:eq("+(cardWidth-1)+")").each(function(index, e){specialJump(e, "jumpLeft")});
	$("#classTable tbody").find("tr:visible:first").find(".equipment img.equipment").each(function(index, e){specialJump(e, "jumpDown")});
	$("#classTable tbody").find("tr:visible:last").find(".equipment img.equipment").each(function(index, e){specialJump(e, "jumpUp")});

	$("tr", "#overlordTable").each(function(index, tr){showHideRow(tr)});
	var td = $("#overlordTable tbody").find("tr:visible:first").find(".cards");
	var card = td.find("div.cardContainer");
	var cardWidth = Math.floor(td.width() / card.width());
	$("#overlordTable tbody").find("tr:visible:first").find(".cards div.cardContainer").each(function(index, e){specialJump(e, "jumpDown")});
	$("#overlordTable tbody").find("tr:visible:last").find(".cards div.cardContainer").each(function(index, e){specialJump(e, "jumpUp")});
	$("#overlordTable tbody").find("tr:visible").find(".cards div.cardContainer:visible:eq("+(cardWidth-1)+")").each(function(index, e){specialJump(e, "jumpLeft")});

	loadSelects();
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
// </SCRIPT>
