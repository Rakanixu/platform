const electron = require("electron");
const {ipcMain} = require("electron");
const os = require("os");
const dialog = electron.dialog;
const windows = require('./windows.js');

// Messages to communicate between node.js and IPC render process
module.exports = (function() {
  function openFolderWindow(event, arg) {
    dialog.showOpenDialog(win,{properties: ['openDirectory']},function(args){
      event.sender.send("add-folder", args);
    });
  }

  function focusMainWindow(event, arg) {
    event.sender.send("auth-callback-message", {});
  }

  ipcMain.on("starting", (event, arg) => {

  });

  ipcMain.on("disks-message", (event, arg) => {
    console.log("Got request for disks");
    drivelist.list(function(error, disks) {
      if (error) throw error;
      event.sender.send("disks-message", disks);
    });
  });

  ipcMain.on("auth-message", windows.createAuthWindow);

  ipcMain.on("auth-callback-message", focusMainWindow);

  ipcMain.on("open-folder",openFolderWindow);

  ipcMain.on("home-dir-message", (event, arg) => {
    event.sender.send("home-dir-message", os.homedir());
  });

  ipcMain.on("user-info-message", (event, arg) => {
    event.sender.send("user-info-message", {
      user: os.userInfo(),
      hostname: os.hostname()
    });
  });

  ipcMain.on("open-file-message", (event, arg) => {
    if (arg.url.indexOf("http") == -1){
      shell.openItem(arg.url);
    } else {
      shell.openExternal(arg.url);
    }
  });
}());
