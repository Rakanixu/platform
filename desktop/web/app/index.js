// In renderer process (web page).
const electron = require('electron');
const {
    ipcRenderer
} = electron;

document.addEventListener('DOMContentLoaded', boot);

function boot() {

    console.log("DOM fully loaded and parsed");
    // check for Geolocation support
    if (navigator.geolocation) {
        console.log('Geolocation is supported!');
    } else {
        console.log('Geolocation is not supported for this Browser/OS version yet.');
    }
    var startPos;
    var geoSuccess = function(position) {
        startPos = position;
        console.log(position);
    };
    navigator.geolocation.getCurrentPosition(geoSuccess);

    //ipcRenderer.send("starting", "elasticsearch");
    //ipcRenderer.send("disks-message");
}



ipcRenderer.on('update-message', function(event, method) {
//    alert(method);
});
/*ipcRenderer.on('disks-message', function(event, disks) {
    console.log("Got disks", disks);
});*/
