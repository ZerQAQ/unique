let btnn = 3;
let btnsize = 10;
let W = btnn * btnsize * 1.5;
let H = btnsize;
var st = btnn;
var lst = btnn - 1;
addbtn();
var pic_timer = setInterval(change_pics, 2000);
var change_pic_timer;

function addbtn(){
	let oelm = document.getElementById("showpics");
	let oW = parseInt(getComputedStyle(oelm)["width"]);
	let oH = parseInt(getComputedStyle(oelm)["height"]);
	let elm = document.getElementById("triger");
	elm.style.cssText =
		"position: absolute;" +
		"width:" + W.toString() + ";" +
		"height:" + H.toString() + ";" +
		"top:" + (oH - H * 2 + 80).toString() + ";" +
		"left:" + (oW / 2 - (W / 2)).toString() + "px;";

	for(let i = 1; i <= btnn; i++){
		let btnelm = document.createElement("img");
		btnelm.src = "imgs/icons/greybtn.png";
		btnelm.id = "btn" + i.toString();
		btnelm.style.cssText = "margin: 0px " + (btnsize / 4).toString() + 
		"px 0px " + (btnsize / 4).toString() + "px; width: " + H.toString() + ";";
		btnelm.style.cursor = "pointer";
		btnelm.onclick = function(){
			let ind = parseInt(this.id.slice(3));
			lst = st;
			st = ind;
			change();
			clearInterval(pic_timer);
			pic_timer = setInterval(change_pics, 2000);
		}
		elm.appendChild(btnelm);
	}
	let btnelm = document.getElementById("btn" + st.toString());
	btnelm.src = "imgs/icons/orangebtn.png"


	var lbt_div = document.getElementById("pic-lbtn")
	lbt_div.style.cssText = "position: absolute;" +
	"top: " + (80 + oH / 2 - 45 / 2).toString() + "px;" +
	"left: 0px;" +
	"opacity: 0.5;" +
	"z-index: 500;";
	var lbt_img = document.createElement("img");
	lbt_img.src = "imgs/icons/l.png";
	lbt_img.onmousedown = function(){
		clearInterval(pic_timer);
		lst = st;
		if(st == 1) st = btnn;
		else st--;
		change();
		pic_timer = setInterval(change_pics, 2000);
	}
	lbt_div.appendChild(lbt_img);

	var rbt_div = document.getElementById("pic-rbtn")
	rbt_div.style.cssText = "position: absolute;" +
	"top: " + (80 + oH / 2 - 45 / 2).toString() + "px;" +
	"left: " + (oW - 45).toString() +"px;" +
	"opacity: 0.5;" +
	"z-index: 500;";
	var rbt_img = document.createElement("img");
	rbt_img.src = "imgs/icons/r.png";
	rbt_img.onmousedown = function(){
		console.log("rd");
		clearInterval(pic_timer);
		lst = st;
		if(st == btnn) st = 1;
		else st++;
		change();
		pic_timer = setInterval(change_pics, 2000);
	}
	rbt_div.appendChild(rbt_img);
}

function delbtn(){
	document.getElementById("triger").innerHTML = " ";
	document.getElementById("pic-lbtn").innerHTML = " ";
	document.getElementById("pic-rbtn").innerHTML = " ";
}

window.onresize = () => {
	delbtn();
	addbtn();
}

function change(){
	clearInterval(change_pic_timer);
	let obtnelm = document.getElementById("btn" + lst.toString());
	obtnelm.src = "imgs/icons/greybtn.png";
	let btnelm = document.getElementById("btn" + st.toString());
	btnelm.src = "imgs/icons/orangebtn.png";
	let picelm1 = document.getElementById("picture1");
	let picelm2 = document.getElementById("picture2");
	picelm2.src = picelm1.src;
	picelm1.src = "imgs/pics/" + st.toString() + ".jpg";
	picelm1.style.opacity = 0;
	picelm2.style.opacity = 1;
	let maxv = 50, v = maxv;
	change_pic_timer = setInterval(() => {
		v--;
		picelm1.style.opacity = (maxv - v) / maxv;
		picelm2.style.opacity = v / maxv;
		if(v == 0){
			clearInterval(change_pic_timer);
		}
	}, 4);
}

function change_pics(){
	lst = st;
	if(st == btnn) st = 1;
	else st++;
	change();
}

