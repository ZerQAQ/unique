
var dropdown_content = document.getElementById("dropdown-content");


for(let i = 1; i <= 6; i++){
	let dropdown_div_block = document.createElement("div");
	dropdown_div_block.innerHTML = 
	"<div>\
	<h1 id = \"nav-dropdown-h1\"> Unique </h1> \
	</div> \
	<div> \
		<a id = \"nav-dropdown-a-1\"> PM </a> \
		<a id = \"nav-dropdown-a-2\"> Game </a> \
		<a id = \"nav-dropdown-a-3\"> AI </a> \
		<a id = \"nav-dropdown-a-4\"> Design </a>\
	</div>\
	<div>\
		<a id = \"nav-dropdown-a-5\"> ios </a>\
		<a id = \"nav-dropdown-a-6\"> Lab </a>\
		<a id = \"nav-dropdown-a-7\"> Android </a>\
		<a id = \"nav-dropdown-a-8\"> Web </a>\
	</div>";

	dropdown_div_block.onmouseover = function(){
		let as = this.getElementsByTagName("a");
		for(let i = 0; i <= 7; i++){
			as[i].style.color = "#f0f0f0";
		}
		let hs = this.getElementsByTagName("h1");
		hs[0].style.color = "rgb(230, 230, 230, 0.7)";
	}

	dropdown_div_block.onmouseout = function(){
		let as = this.getElementsByTagName("a");
		for(let i = 0; i <= 7; i++){
			as[i].style.color = "#707070";
		}
		let hs = this.getElementsByTagName("h1");
		hs[0].style.color = "#a0a0a0";
	}

	dropdown_content.appendChild(dropdown_div_block);
}

var dropdown_btn = document.getElementById("sel");
var nav_elm = document.getElementById("nav");
var timer_add, timer_sub, h = 0, h_max = 20;
var	dropdown_H = parseInt(getComputedStyle(dropdown_content)["height"]);
var nav_H = parseInt(getComputedStyle(nav_elm)["height"]);
dropdown_content.style.height = 0;
dropdown_content.style.top = nav_H;

dropdown_btn.onmouseover = function(){
	clearInterval(timer_add);
	clearInterval(timer_sub);
	timer_add = setInterval(() => {
		if(h == h_max){
			clearInterval(timer_add);
			return;
		}
		h++;
		dropdown_content.style.height = dropdown_H * h / h_max;
	}, 10);
}

var leaving = 0;

function dropdown_leave(){
	clearInterval(timer_add);
	clearInterval(timer_sub);
	timer_sub = setInterval(() => {
		if(h == 0){
			clearInterval(timer_sub);
			leaving = 0;
			return;
		}
		h--;
		dropdown_content.style.height = dropdown_H * h / h_max;
	}, 10);
}

function mousemove(e){
	if(leaving) return;
	e = e || window.event;
	if(e.clientY >= dropdown_H + nav_H){
		dropdown_leave();
		leaving = 1;
	}
}

document.getElementById("topbtn").onmouseover = function(){
	this.src = "imgs/icons/top02.png";
}

document.getElementById("topbtn").onmouseover = function(){
	this.src = "imgs/icons/top02.png";
}