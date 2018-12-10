var hero0cardId = 0;
var hero1cardId = 0;
var hero2cardId = 0;
var hero3cardId = 0;

function deleteCookie(cname) {
    // TODO: better way to delete cookies?
    document.cookie = cname + '=, expires=Thu, 01 Jan 1970 00:00:01 GMT, path=/';
}

function setCookie(cname, cvalue) {
    var curCookie = cname + "=" + cvalue + 
    ", expires=" + "Tue, 19 Jan 2038 03:14:07 UTC" + 
    ", path=" + "/";
    document.cookie = curCookie;

    var expires = "expires=Tue, 19 Jan 2038 03:14:07 UTC";
    document.cookie = cname + "=" + cvalue + "," + expires + ",path=/";
}

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            c = c.split(',')[0];
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function loadCookies() {
    var ca = document.cookie.split(';');
    for(var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        c = c.split(',')[0];
        var cname = c.split('=')[0];
        var cval = c.substring(cname.length+1, c.length);
        
        if (cname.substring(0, 6) == "health") {
            var div = "div#health" + cname.substring(6, 7);
            if (cval != "NaN" && cval != "undefined") {
                $(div).text(cval);
            }
            continue;
        }
        if (cname.substring(0, 7) == "fatigue") {
            var div = "div#fatigue" + cname.substring(7, 8);
            if (cval != "NaN" && cval != "undefined") {
                $(div).text(cval);
            }
            continue;
        }
        
        var arr = cname.split('.');
        if (arr[2] != "img") {
            continue;
        }

        var heroId = arr[0].split('hero')[1];
        var cardId = arr[1].split('card')[1];
        var img = cval;
        if (img == "") {
            continue;
        }
        
        var x=getCookie("hero" + heroId + ".card" + cardId + ".x");
        var y=getCookie("hero" + heroId + ".card" + cardId + ".y");
        var heroArea = document.getElementById("hero" + heroId);
        addCard(heroArea, img, x, y);
    }
}

function onload() {
    randombg();
    loadCookies();
}

function heroToCardIdAndInc(heroId) {
    switch (heroId) {
        case "hero0":
            var id = hero0cardId;
            hero0cardId++;
            return id;
        case "hero1":
            var id = hero1cardId;
            hero1cardId++;
            return id;
        case "hero2":
            var id = hero2cardId;
            hero2cardId++;
            return id;
        case "hero3":
            var id = hero3cardId;
            hero3cardId++;
            return id;
    }
}

function addCard(heroArea, img, x, y)
{
    var id = $(heroArea).attr("id");
    var cardId = heroToCardIdAndInc(id);
    $(heroArea).append("<img src=\"images/"+img+"\" id=\""+id+".card"+cardId+"\" class=\"cardImg\" style=\"position:absolute;top:" + y + ";left:" + x + ";\">");
    setCookie(id + ".card" + cardId + ".img", img);
    setCookie(id + ".card" + cardId + ".x", x);
    setCookie(id + ".card" + cardId + ".y", y);
    document.getElementById(id+".card"+cardId).addEventListener('mousedown', mouseDown, false);
    document.getElementById(id+".card"+cardId).addEventListener('oncontextmenu', removeImg, false);
    window.addEventListener('mouseup', mouseUp, false);
}

function setHealthStamina(heroArea, li)
{
    var health = $(li).find("a").attr("health");
    var stamina = $(li).find("a").attr("stamina");
    var id = $(heroArea).attr("id").slice(-1);
    var div = "div#health" + id;
    $(div).text(health);
    setCookie("health"+id, health);
    var div = "div#fatigue" + id;
    $(div).text(stamina);
    setCookie("fatigue"+id, stamina);
}

function search(input, e) {
    var enterPressed = false;
    if (e.keyCode === 13) {
        enterPressed = true;
    }
    var filter, ul, li, a, i, txtValue;
    ul = $(input).parent().find(".myuL");
    filter = input.value.toUpperCase();
    li = ul.find("li");
    var count = 0;
    for (i = 0; i < li.length; i++) {
        a = li[i].getElementsByTagName("a")[0];
        txtValue = a.textContent || a.innerText;
        if (txtValue.toUpperCase().indexOf(filter) > -1 && count < 8) {
            if (enterPressed) {
                var img = $(li[i]).find("a").attr("href");
                var heroArea = $(li[i]).parent().parent().find(".hero");
                addCard(heroArea, img, "10px", "10px");
                setHealthStamina(heroArea, li[i]);
                enterPressed = false;
            }
            li[i].style.display = "";
            count++;
        } else {
            li[i].style.display = "none";
        }
    }
}

function gotFocus(input) {
    var id = $(input).parent().attr('id');
    $('#' + id + ' .myUL').show();
    $('#' + id + ' .myUL li').hide().slice(0, 8).show();
}

function lostFocus(x) {
    $(x).parent().find(".search-input").val("");
    $(x).parent().find(".myUL").hide();
}

$(function() {
  $("div.health").on("mousedown",function(e) {
    var mid = $(this).offset().left + ($(this).width() / 2);
    var val = parseInt($(this).text());
    if ($(this).text() == "" || $(this).text() == "NaN" || $(this).text() == "undefined") {
        val = 0;
    } else if (e.clientX < mid) {
        // left
        val--;
    } else {
        // right
        val++;
    }
    $(this).text(val);
    setCookie(this.id, val);
  });
});

$(function() {
  $("div.fatigue").on("mousedown",function(e) {
    var mid = $(this).offset().left + ($(this).width() / 2);
    var val = parseInt($(this).text());
    if ($(this).text() == "" || $(this).text() == "NaN" || $(this).text() == "undefined") {
        val = 0;
    } else if (e.clientX < mid) {
        // left
        val--;
    } else {
        // right
        val++;
    }
    $(this).text(val);
    setCookie(this.id, val);
  });
});

$(function() {
  $(".card").on("mousedown",function() {
    var img = $(this).find("a").attr("href");
    var heroArea = $(this).parent().parent().find(".hero");
    addCard(heroArea, img, "10px", "10px");
    setHealthStamina(heroArea, this);
  });
});

var offset = [0,0];
var isDown = false;
var image = "";

function mouseUp(e)
{
    if (isDown) {
        isDown = false;
        window.removeEventListener('mousemove', mouseMove, true);
        
        var win = $(window), width = win.width(), height = win.height();
        var quad = 0;
        if (e.clientX < width/2) {
            if (e.clientY < height/2) {
                // top-left
                quad = 0;
            } else {
                // bottom-left
                quad = 2;
            }
        } else {
            if (e.clientY < height/2) {
                // top-right
                quad = 1;
            } else {
                // bottom-right
                quad = 3;
            }
        }
        var id = $(image).parent().attr("id");
        var digit = id.substr(id.length - 1);
        if (digit != quad) {
            deleteCookie($(image).attr("id") + ".img");
            deleteCookie($(image).attr("id") + ".x");
            deleteCookie($(image).attr("id") + ".y");
            $(image).remove();
        } else {
            setCookie($(image).attr("id") + ".x", image.style.left);
            setCookie($(image).attr("id") + ".y", image.style.top);
        }
    }
}

function mouseDown(e)
{
    parent = $(this).parent();
    isDown = true;
    image = this;
    offset = [
        parseInt(image.style.left) - e.clientX,
        parseInt(image.style.top) - e.clientY,
    ];
    window.addEventListener('mousemove', mouseMove, true);
}

function mouseMove(e){
    event.preventDefault();  // Cancel the default action, if needed
    if (isDown) {
        image.style.position = 'absolute';
        image.style.left = (e.clientX + offset[0]) + 'px';
        image.style.top  = (e.clientY + offset[1]) + 'px';
    }
}

function removeImg(e){
    this.parentNode.removeChild(this);
}

function randombg(){
    var bigSize = ["url('backgrounds/back0.png')",
                 "url('backgrounds/back1.jpg')",
                 "url('backgrounds/back2.jpg')",
                 "url('backgrounds/back3.jpg')",
                 "url('backgrounds/back4.png')",
                 "url('backgrounds/back5.jpg')",
                 "url('backgrounds/back6.jpg')"];
    var random = Math.floor(Math.random() * bigSize.length) + 0;
    document.getElementById("back").style.backgroundImage=bigSize[random];
}
