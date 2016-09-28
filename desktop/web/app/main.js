// Handle Squirrel events for Windows immediately on start
const electron = require("electron");
const fs = require("fs");
// Module to control application life.
const {app} = electron;

const messages = require('./messages.js');
const windows = require('./windows.js');
const srvs = require('./services.js');

let paths;
let services;
let running;
let isDevelopment = process.env.NODE_ENV === "development";

const version = app.getVersion();

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", function() {
    //Start micro services
    srvs.startServices();

    setTimeout(function() {
        windows.createMainWindow();
    }, 2000);
});
app.on("activate", () => {
    // On macOS it"s common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (win === null) {
        window.createMainWindow();
    }
});

// Quit when all windows are closed.
app.on("window-all-closed", () => {
    // On macOS it is common for applications and their menu bar
    // to stay active until the user quits explicitly with Cmd + Q
    //
    //es.kill("SIGHUP");
    if (!isDevelopment) {
        for (var i = 0; i < running.length; i++) {
            srvs.getServices()[i].kill("SIGHUP");
            console.log("Killing " + i);
        }
    }

    if (process.platform !== "darwin") {
        app.quit();
    }
});

function discoverServices(platform, architecture) {
    url = "bin/" + platform + "/" + architecture + "/";
    services = []
    files = fs.readdirSync(url);

    for (var i in files) {
        services.push(url + files[i]);
    }
    return services;
}
