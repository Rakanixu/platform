'use strict';

// Imports map by app context
// window.appEnvironment is being defined on index.html and index-electron.html
(function(window) {
  var electronAppPath;

  try {
    if (electron !== undefined) {
      electronAppPath = electron.remote.app.getAppPath()
    }
  } catch(e) {}

  window.polymerApp = {
    schedulerWorker: '/src/scheduler/scheduler-worker.js',
    unknownImg: '/src/static/unknown-img.jpg',
    unknownVideo: '/src/static/unknown-video.png',
    menuBackground: '/src/static/menu-background.jpg'
  };

  window.electronApp = {
    schedulerWorker: electronAppPath + '/src/scheduler/scheduler-worker.js',
    unknownImg: electronAppPath + '/src/static/unknown-img.jpg',
    unknownVideo: electronAppPath + '/src/static/unknown-video.png',
    menuBackground: electronAppPath + '/src/static/menu-background.jpg'
  };

})(window);
