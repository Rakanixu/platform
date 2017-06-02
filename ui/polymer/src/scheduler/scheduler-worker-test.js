// Requirments, Node.js and npm package manager
// npm install unit.js
// npm install mocha
// npm install mock-browser
// npm install workerjs

var test = require('unit.js');
var MockBrowser = require('mock-browser').mocks.MockBrowser;
var Worker = require('workerjs');

describe('scheduler-worker.js', function(){
  // Mock browser window object
  self = MockBrowser.createWindow();

  it('scheduler', function(done){
    var w = new Worker(__dirname + '/scheduler-worker.js', true);
    var result;
    var expectedIds = ["1", "4", "5"];

    w.addEventListener('message', function (e) {
      test.array(expectedIds).is(e.data.scan_datasources);

      w.terminate();
      w = undefined;

      done();
    });

    w.postMessage({
      datasources: [
        {
          "id": "1",
          "last_scan": 0,
          "last_scan_started": 0
        }, {
          "id": "2",
          "last_scan": Math.ceil(new Date().getTime() / 1000) - (60 * 60),
          "last_scan_started": Math.ceil(new Date().getTime() / 1000)
        }, {
          "id": "3",
          "last_scan": Math.ceil(new Date().getTime() / 1000),
          "last_scan_started": Math.ceil(new Date().getTime() / 1000)
        }, {
          "id": "4",
          "crawler_running": true,
          "last_scan": Math.ceil(new Date().getTime() / 1000) - ((60 * 60) * 2),
          "last_scan_started": Math.ceil(new Date().getTime() / 1000) - (60 * 60)
        }, {
          "id": "5",
          "crawler_running": true,
          "last_scan": 0,
          "last_scan_started": Math.ceil(new Date().getTime() / 1000) - (60 * 60)
        }, {
          "id": "6",
          "crawler_running": false,
          "last_scan": 0,
          "last_scan_started": Math.ceil(new Date().getTime() / 1000) - (60 * 60)
        }
      ]
    });
  });
});
