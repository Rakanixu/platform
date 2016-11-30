'use strict';

var window = self;

(function() {
  var elapsedTime = 60 * 3; // We count in seconds, not milliseconds
  var day = 60 * 60 * 24;

  self.initScheduler = function(e) {
    console.log("initScheduler")
    self.datasources = e.data.datasources;

    self.checkTime();

    setInterval(function() {
      self.checkTime();
    }, 60000); // Check every minute
  };

  self.checkTime = function() {

    console.log("checkTime called" )
    self.initTime = parseInt((Date.now() / 1000), 10);
    self.pool = [];

    for (var i = 0; i < self.datasources.length; i++) {

      //
      if (self.datasources[i].last_scan > self.datasources[i].last_scan_started) {
        if (self.datasources[i].last_scan < self.initTime - day) {
          console.log("DS to be scan as was incorrect for a day")

          self.pool.push(self.datasources[i].id);
        }
      }


      if (self.datasources[i].last_scan < self.initTime - elapsedTime) {
        console.log("DS to be scan it pass elapsed time", self.datasources[i].last_scan)
        self.pool.push(self.datasources[i].id);
      }

      // last_scan may not be set it, as never finished
/*      if (self.datasources[i].last_scan === undefined) {
        // Check how much time ago scann was kicked off
        if (self.datasources[i].last_scan_started < self.initTime - elapsedTime * 2) {
          self.pool.push(self.datasources[i].id);
        }
      }*/
    }

    self.postMessage({
      scan_datasources: self.pool
    });
  };

  self.addEventListener('message', self.initScheduler, false);
}());
