// Windows to be manage by electron main process
const electron = require("electron");
const os = require("os");
const menubar = require('menubar')
const {BrowserWindow} = electron;
const {autoUpdater} = electron;

// Keep a global reference of the window object, if you don"t, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;
let auth;
let widget;


module.exports = (function() {
  let updateFeed = "";
  let feedURL = "https://protected-reaches-10740.herokuapp.com";
  let isDevelopment = process.env.NODE_ENV === "development";
  let version;
  
  var _getMainWindow = function() {
    return win;
  };

  var _createMainWindow = function(app) {
    // Create the browser window.
    version = app.getVersion();
    win = new BrowserWindow({
      width: 1024,
      height: 768,
      webPreferences: {
        webSecurity: false
      }
    });

    // and load the index.html of the app.
    win.loadURL(`file://${__dirname}/index-electron.html`);

    // Open the DevTools.
    if (isDevelopment){
      win.webContents.openDevTools();
    }

    // Don"t use auto-updater if we are in development
    if (!isDevelopment) {
      if (os.platform() === "darwin") {
        updateFeed = `${feedURL}/update?version=${version}&platform=osx`;
      } else if (os.platform() === "win32") {
        updateFeed = `${feedURL}/update/win32/${version}`;
      }

      autoUpdater.addListener("update-available", function(event) {
        //console.log("A new update is available");
        if (win) {
          win.webContents.send("update-message", "update-available");
        }
      });

      autoUpdater.addListener("update-downloaded", function(event, releaseNotes, releaseName, releaseDate, updateURL) {
        console.log("A new update is ready to install", '${releaseName} is downloaded and will be automatically installed on Quit');
        if (win) {
          win.webContents.send("update-message", "update-downloaded");
        }
      });

      autoUpdater.addListener("error", function(error) {
        //         console.log(error);
        if (win) {
          win.webContents.send("update-message", "update-error");
        }
      });

      autoUpdater.addListener("checking-for-update", function(event) {
        //    console.log("Checking for update");
        if (win) {
          win.webContents.send("update-message", "checking-for-update");
        }
      });

      autoUpdater.addListener("update-not-available", function() {
        //  console.log("Update not available");
        if (win) {
          win.webContents.send("update-message", "update-not-available");
        }
      });

      autoUpdater.setFeedURL(updateFeed);
    }

    // Emitted when the window is closed.
    win.on("closed", () => {
      // Dereference the window object, usually you would store windows
      // in an array if your app supports multi windows, this is the time
      // when you should delete the corresponding element.
      win = null;
    });

    if (!isDevelopment) {
      win.webContents.on("did-frame-finish-load", function() {
        console.log("Checking for updates: " + updateFeed);
        autoUpdater.setFeedURL(updateFeed)
        autoUpdater.checkForUpdates();
      });
    }
  };

  var _createWidgetWindow = function(app) {
    mb = menubar({
      index: `file://${__dirname}/index-electron.html`
    });

    mb.on('ready', function ready () {
      console.log('app is ready')
    });

    return mb;
  };

  var _createAuthWindow =  function(event, arg) {
    auth = new BrowserWindow({
      parent: win,
      modal:false,
      width:420,
      height: 590,
      frame: true,
      webPreferences: {
        nodeIntegration: false
      }
    });

    auth.webContents.session.clearStorageData(function() {
      auth.loadURL(arg.url)
    });

    auth.on("closed", () => {
      // Let main window has the focus again
      event.sender.send('auth-callback-message', arg.token);
    });

    auth.show();
  }

  return {
    getMainWindow: _getMainWindow,
    createMainWindow: _createMainWindow,
    createWidgetWindow: _createWidgetWindow,
    createAuthWindow: _createAuthWindow
  }
}());
