const electron = require('electron');
const ipc = require('ipc');
// Module to control application life.
const {app} = electron;
const {ipcMain} = require('electron');
const {os} = require('os')
const spawn = require("child_process").spawn

// Module to create native browser window.
const {BrowserWindow} = electron;

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;
let es;
function createWindow() {
  // Create the browser window.
  win = new BrowserWindow({width: 800, height: 600});
  //Start micro services

  // and load the index.html of the app.
  win.loadURL(`file://${__dirname}/app/index.html`);

  // Open the DevTools.
  win.webContents.openDevTools();
  win.webContents.send('starting', "elasicsearch");
  // Emitted when the window is closed.
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null;
  });
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', createWindow);

// Quit when all windows are closed.
app.on('window-all-closed', () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  //
  //es.kill('SIGHUP');
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (win === null) {
    createWindow();
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
app.on('asynchronous-message', (event, arg) => {
  console.log(arg);  // prints "ping"
  event.sender.send('asynchronous-reply', 'pong');
});

ipcMain.on('starting', (event, arg) => {
  console.log(os.platform());  // prints "ping"
 
  startServices()
  event.returnValue = "OK";
});

app.on('synchronous-message', (event, arg) => {
  console.log(arg);  // prints "ping"
  event.returnValue = 'pong';
});

function startServices(){


   var wd =  __dirname + "/bin"
   var es = spawn( __dirname + '/bin/auth-api_darwin_amd64', ['--registry=mdns'], {wd:__dirname});
   es.stdout.on('data', function (data) {
   console.log('stdout: ' + data);
   });
   es.on('close', function (code, signal) {
   console.log('stdout: es process terminated due to receipt of signal ' + signal);
   });
}
