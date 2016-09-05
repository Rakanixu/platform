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

// Keep a global reference of the window object, if you don"t, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;
let auth;
let es;
let paths;
let services;
let running;
let updateFeed = "";
let isDevelopment = process.env.NODE_ENV === "development";
let feedURL = "https://protected-reaches-10740.herokuapp.com";
const version = app.getVersion();


function openFolderWindow(event, arg) {
    dialog.showOpenDialog(win,{properties: ['openDirectory']},function(args){
        event.sender.send("add-folder", args);
    });
}

function createAuthWindow(event, arg) {
    auth = new BrowserWindow({
        parent: win,
        modal:false,
        width:420,
        height: 590,
        frame: true,
        webPreferences: {
          nodeIntegration: true
        }
    });

    auth.webContents.session.clearStorageData(function() {
        auth.loadURL(arg)
    });

    auth.on("closed", () => {
        // Let main window has the focus again
        event.sender.send('auth-callback-message');
    });

    auth.show();
}

function focusMainWindow(event, arg) {
    event.sender.send("auth-callback-message", {});
}

function createWindow() {
    // Create the browser window.
    win = new BrowserWindow({
        width: 1024,
        height: 768,
        webPreferences: {
              webSecurity: false
        }
    });
    //Start micro services
    startServices()
    // and load the index.html of the app.
    win.loadURL(`file://${__dirname}/index-electron.html`);

    //win.loadURL("http://google.com");
    // Open the DevTools.
    win.webContents.openDevTools();

    // Don"t use auto-updater if we are in development
    if (!isDevelopment) {
        if (os.platform() === "darwin") {
            updateFeed = `${feedURL}/update?version=${version}&platform=osx`;

        } else if (os.platform() === "win32") {
            //updateFeed = "https://protected-reaches-10740.herokuapp.com/" + (os.arch() === "x64" ? "64" : "32");

            updateFeed = `${feedURL}/update/win32/${version}`;
        }

        autoUpdater.addListener("update-available", function(event) {
            console.log("A new update is available");
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
            console.log(error);
            if (win) {
                win.webContents.send("update-message", "update-error");
            }
        });
        autoUpdater.addListener("checking-for-update", function(event) {
            console.log("Checking for update");
            if (win) {
                win.webContents.send("update-message", "checking-for-update");
            }
        });
        autoUpdater.addListener("update-not-available", function() {
            console.log("Update not available");
            if (win) {
                win.webContents.send("update-message", "update-not-available");
            }
        });
        autoUpdater.setFeedURL(updateFeed);
    }



    // Detect drives
    //
    //console.log(findDisks());

    //win.webContents.send("disks-message", findDisks());
    //win.webContents.send("starting", "elasicsearch");
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
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", createWindow);

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

app.on("activate", () => {
    // On macOS it"s common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (win === null) {
        createWindow();
    }
});

// In this file you can include the rest of your app"s specific main process
// code. You can also put them in separate files and require them here.
app.on("asynchronous-message", (event, arg) => {
    console.log(arg); // prints "ping"
    event.sender.send("asynchronous-reply", "pong");
});

ipcMain.on("starting", (event, arg) => {
});

app.on("synchronous-message", (event, arg) => {
    console.log(arg); // prints "ping"
    event.returnValue = "pong";
});

ipcMain.on("disks-message", (event, arg) => {
    console.log("Got request for disks");
    drivelist.list(function(error, disks) {
        if (error) throw error;
        event.sender.send("disks-message", disks);
    });
});
ipcMain.on("auth-message",createAuthWindow);
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
    shell.openItem(arg.url);
});

function startService(path, args) {
    es = spawn(path, args, {
        wd: path,
        stdio: 'ignore'
    });
    //es.stdout.on("data", function(data) {
    //    console.log("stdout: " + data);
    //});
    es.on("close", function(code, signal) {
        console.log("stdout: es process terminated due to receipt of signal " + signal);
    });

    return es;
}

function discoverServices(platform, architecture) {
    url = "bin/" + platform + "/" + architecture + "/";
    services = []
    files = fs.readdirSync(url);

    for (var i in files) {
        services.push(url + files[i]);
    }
    return services;
}

function archConvert(arch) {

    switch (arch) {
        case "x64":
            arch = "amd64";
            break;
        case "ia32":
            arch = "386";
            break;
    }

    return arch;
}

function findDisks() {
    var results = []
    drivelist.list(function(error, disks) {
        if (error) throw error;
        results = disks

    });
    return results;
}

function startServices(){

    running = [];
    resourcesPath = "";
    //Set resources path depending if we running in development - needs to do it better
    if (!isDevelopment) {
        resourcesPath = process.resourcesPath;// Starting all required services FIXME:
        if (process.platform == "win32") {
            elastic = startService(resourcesPath + "/elasticsearch/bin/elasticsearch.bat", []);
            kazoup = startService(resourcesPath 		+ "/bin/windows/" + archConvert(process.arch) + "/kazoup.exe", ["desktop"])
        } else {
            elastic = startService(resourcesPath + "/elasticsearch/bin/elasticsearch", []);
            kazoup = startService(resourcesPath 		+ "/bin/" + process.platform + "/" + archConvert(process.arch) + "/kazoup", ["desktop"])
        }
        running.push(elastic)
        running.push(kazoup)

    } else {
        //TODO: stupid hack fixme should we point to desktop folder ?
        resourcesPath = __dirname + "/..";

    }

}
