'use strict';

var window = self;

(function() {
  var elapsedTime = 60 * 60; // We count in seconds, not milliseconds
  var day = 60 * 60 * 24;

  self.initScheduler = function(e) {
    self.datasources = e.data.datasources;

    self.checkTime();

    setInterval(function() {
      self.checkTime();
    }, 60000); // Check every minute
  };

  self.checkTime = function() {
    self.initTime = parseInt((Date.now() / 1000), 10);
    self.pool = [];

    for (var i = 0; i < self.datasources.length; i++) {

      console.log(self.datasources[i])

      // Something wrong happened, kick off scan to fix it
      if (self.datasources[i].last_scan > self.datasources[i].last_scan_started) {
        if (self.datasources[i].last_scan < self.initTime - day) {
          self.pool.push(self.datasources[i].id);
        }
        
        // Scan was stopped, probably because a new deployment while crawler was running
        if (self.datasources[i].last_scan_started < (parseInt(new Date().getTime() / 1000) + elapsedTime)) {
          self.pool.push(self.datasources[i].id);
        }
      }

      // Las scan has to be prior elapsed time, but we have to know if last scan happen a few time ago
      // This avoid to kick off several scan while previous one is running
      if (self.datasources[i].last_scan < self.initTime - elapsedTime &&
        self.datasources[i].last_scan + elapsedTime > self.datasources[i].last_scan_started) {
        self.pool.push(self.datasources[i].id);
      }
    }

    // Send DS that requires to be scan
    self.postMessage({
      scan_datasources: self.pool
    });
  };

  self.addEventListener('message', self.initScheduler, false);
}());
