// Handle Squirrel events for Windows immediately on start
const electron = require("electron");
const fs = require("fs");
const drivelist = require('drivelist');
const remote = require("electron").remote;
const dialog = electron.dialog;
const {shell} = electron;
// Module to control application life.
const {
    app
} = electron;
const {
    ipcMain
} = require("electron");

const spawn = require("child_process").spawn
const {
    autoUpdater
} = electron;
const os = require("os");
// Module to create native browser window.
const {
    BrowserWindow
} = electron;

const messages = require('./messages.js');
const windows = require('./windows.js');

// Keep a global reference of the window object, if you don"t, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;
let auth;
let es;
let paths;
let services;
let running;
let isDevelopment = process.env.NODE_ENV === "development";
let feedURL = "https://protected-reaches-10740.herokuapp.com";
const version = app.getVersion();


// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", windows.createMainWindow);
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
            running[i].kill("SIGHUP");
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

function findDisks() {
    var results = []
    drivelist.list(function(error, disks) {
        if (error) throw error;
        results = disks

    });
    return results;
}

