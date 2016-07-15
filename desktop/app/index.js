// In renderer process (web page).
const {ipcRenderer} = require('electron');
window.document.addEventListener("DOMContentLoaded", function(event){
    var start = window.document.getElementById("start")
	start.addEventListener('click',function(){
		console.log("Start services")
		ipcRenderer.send("starting","elasticsearch")
	});
    console.log("DOM fully loaded and parsed");
});
