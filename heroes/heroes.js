// <SCRIPT LANGUAGE="JavaScript">

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

function showHideRows()
{
	var e = $("#selectExp").val();
	var k = $("#selectCK").val();
	var c = $("#selectClass").val();
	var d = $("#selectDefense").val();
	$("tr", "#heroTable").each(function(index, tr){
		$(tr).show();
		if (typeof $(tr).attr("class") === "undefined")
			return;
		if (c != "" && !$(tr).hasClass(c))
			$(tr).hide();
		if (d != "" && !$(tr).hasClass(d))
			$(tr).hide();
		if (k != "" && !$(tr).hasClass(k))
			$(tr).hide();
		if (e != "" && !$(tr).hasClass(e))
			$(tr).hide();
		
		$.each(attrs, function (i, attr) {
			var s = $("#select"+toTitleCase(attr)).val();
			if (s != "" && $('td.'+attr, tr).text() != s)
				$(tr).hide();
		});
	});
	
	loadSelects();
}

function loadSelects()
{
	$.each(attrs, function (i, attr) {
		var selected = $('#select'+toTitleCase(attr)).val();
		var arr = [];
		$("tr", "#heroTable").each(function(index, tr) {
			$('td.'+attr, tr).each(function(index, td) {
				if (selected != "" || $(tr).is(":visible"))
					arr.push($(td).text());
			});
		});
		arr = jQuery.uniqueSort(arr).sort(function(a,b){
			return a-b;
		});
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
// </SCRIPT>
