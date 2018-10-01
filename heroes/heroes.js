// <SCRIPT LANGUAGE="JavaScript">

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
