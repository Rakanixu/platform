// Handle Squirrel events for Windows immediately on start
const electron = require("electron");
const fs = require("fs");
// Module to control application life.
const {app} = electron;

const messages = require('./messages.js');
const windows = require('./windows.js');
const srvs = require('./services.js');

let services;
let win;


// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", function() {
    win = windows.createMainWindow(app);
    //win = windows.createWidgetWindow(app);

});
app.on("activate", () => {
    // On macOS it"s common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (win === null) {
        win = windows.createMainWindow(app);
    }
});

// Quit when all windows are closed.
app.on("window-all-closed", () => {
    // On macOS it is common for applications and their menu bar
    // to stay active until the user quits explicitly with Cmd + Q
    if (process.platform !== "darwin") {
        app.quit();
    }
});