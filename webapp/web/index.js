// In renderer process (web page).
const electron = require('electron');
const {ipcRenderer} = electron;


function boot() {

    console.log("DOM fully loaded and parsed");
    ipcRenderer.send("disks-message","ping");
}

document.addEventListener('DOMContentLoaded', boot);
//document.addEventListener("DOMContentLoaded", function(event){
//    var start = window.document.getElementById("start")
//	start.addEventListener('click',function(){
//		console.log("Start services")
//		ipcRenderer.send("starting","elasticsearch")
//	});
//
//
//});

ipcRenderer.on('update-message', function(event, method) {
    alert(method);
});
ipcRenderer.on('disks-message',function(event,disks){
   console.log("Got disks",disks);
});
