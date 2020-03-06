
let web_groups = document.getElementById("group-web");

for(let i = 1; i <= 2; i++){
	let web_groups_block_div = document.createElement("div");
	web_groups_block_div.innerHTML = "\
	<img src = \"imgs/groups/web.png\">\
	<div class = \"animation1\"> \
		<p id = \"animation1-text\"> Uniuqe Web </p>\
	</div>";
	web_groups_block_div.className = "group-block";

	let web_groups_timer_add, web_groups_timer_sub, web_groups_h = 0, web_groups_maxh = 30;

	web_groups_block_div.onmouseover = function(){
		clearInterval(web_groups_timer_add);
		clearInterval(web_groups_timer_sub);
		let anima = this.getElementsByTagName("div")[0];
		let H = parseInt(getComputedStyle(this)["height"]);
		web_groups_timer_add = setInterval(() => {
			if(web_groups_h == web_groups_maxh){
				clearInterval(web_groups_timer_add);
				return;
			}
			web_groups_h++;
			anima.style.height = H * web_groups_h / web_groups_maxh;
		}, 3);
	}

	web_groups_block_div.onmouseout = function(){
		clearInterval(web_groups_timer_add);
		clearInterval(web_groups_timer_sub);
		let anima = this.getElementsByTagName("div")[0];
		let H = parseInt(getComputedStyle(this)["height"]);
		web_groups_timer_add = setInterval(() => {
			if(web_groups_h == 0){
				clearInterval(web_groups_timer_add);
				return;
			}
			web_groups_h--;
			anima.style.height = H * web_groups_h / web_groups_maxh;
		}, 5);
	}

	web_groups.appendChild(web_groups_block_div);
}

var groups_name = new Array("design", "game", "ios", "lab", "pm", "Android");
var groups_add_timer = new Array(), groups_sub_timer = new Array();


for(let i = 0; i < 6; i++){
	let groups_div = document.createElement("div");
	groups_add_timer[i] = new Array();
	groups_sub_timer[i] = new Array();
	groups_div.className = "groups";
	for(let j = 0; j < 4; j++){
		let groups_block_div = document.createElement("div");
		groups_block_div.className = "group-block";
		groups_block_div.innerHTML = "\
			<img src = \"imgs/groups/" + groups_name[i] + ".png\">\
			<div style = \"position: absolute; bottom: 0px; background: rgb(254,102,0); z-index: -100\">\
			</div>\
			<div>\
				<p> Unique-" + groups_name[i] + "</p>\
				<h1> 659 </h1>\
				<img src = \"imgs/icons/works_clickgray.png\">\
			</div>";
		
		let h = 0, maxh = 20;
		
		groups_block_div.onmouseover = function(){
			clearInterval(groups_add_timer[i][j]);
			clearInterval(groups_sub_timer[i][j]);
			let anima = this.getElementsByTagName("div")[0];
			let p = this.getElementsByTagName("p")[0];
			let h1 = this.getElementsByTagName("h1")[0];
			let img = this.getElementsByTagName("img")[1];
			let bimg = this.getElementsByTagName("img")[0];
			let H = parseInt(getComputedStyle(this).height) - parseInt(getComputedStyle(bimg).height);
			p.style.color = "#fff";
			h1.style.color = "#fff";
			img.style.display = "none";
			
			groups_add_timer[i][j] = setInterval(() => {
				if(h == maxh){
					clearInterval(groups_add_timer[i][j]);
					return;
				}
				h++;
				anima.style.height = H * h / maxh;
			}, 2);
		}

		groups_block_div.onmouseout = function(){
			clearInterval(groups_add_timer[i][j]);
			clearInterval(groups_sub_timer[i][j]);
			let anima = this.getElementsByTagName("div")[0];
			let p = this.getElementsByTagName("p")[0];
			let h1 = this.getElementsByTagName("h1")[0];
			let img = this.getElementsByTagName("img")[1];
			let bimg = this.getElementsByTagName("img")[0];
			let H = parseInt(getComputedStyle(this).height) - parseInt(getComputedStyle(bimg).height);
			p.style.color = "#888";
			h1.style.color = "rgb(254,102,0)";
			img.style.display = "block";

			groups_sub_timer[i][j] = setInterval(() => {
				if(h == 0){
					clearInterval(groups_sub_timer[i][j]);
					return;
				}
				h--;
				anima.style.height = H * h / maxh;
			}, 2);
		}

		groups_div.appendChild(groups_block_div);
	}
	document.body.appendChild(groups_div);
}

