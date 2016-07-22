// In renderer process (web page).
const electron = require('electron');
const {ipcRenderer} = electron;


function boot() {

    console.log("DOM fully loaded and parsed");
// check for Geolocation support
	if (navigator.geolocation) {
  		console.log('Geolocation is supported!');
	}
	else {
	  console.log('Geolocation is not supported for this Browser/OS version yet.');
	}
	var startPos;
  var geoSuccess = function(position) {
    	startPos = position;
	console.log(position);
    	//document.getElementById('startLat').innerHTML = startPos.coords.latitude;
    	//document.getElementById('startLon').innerHTML = startPos.coords.longitude;
  	};
 	 navigator.geolocation.getCurrentPosition(geoSuccess);
    //ipcRenderer.send("disks-message","ping");
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
