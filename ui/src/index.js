// In renderer process (web page).
const electron = require('electron');
const {
    ipcRenderer
} = electron;

document.addEventListener('DOMContentLoaded', boot);

function boot() {

    console.log("DOM fully loaded and parsed");

    // check for Geolocation support
    if (navigator.geolocation) {
        console.log('Geolocation is supported!');
    } else {
        console.log('Geolocation is not supported for this Browser/OS version yet.');
    }
    var startPos;
    var geoSuccess = function(position) {
        startPos = position;
        console.log(position);
    };
    navigator.geolocation.getCurrentPosition(geoSuccess);
}

var Endpoints = (function() {
  return {
      endpoint: 'https://web.kazoup.io:8082/rpc',
      web: 'https://web.kazoup.io:8082',
      socket: 'wss://web.kazoup.io:8082',
      srvs:{
          crawler: {
              srv: 'com.kazoup.srv.crawler',
              search: 'Crawler.Search'
          },
          datasource: {
              srv: 'com.kazoup.srv.datasource',
              create: 'DataSource.Create',
              delete: 'DataSource.Delete',
              search: 'DataSource.Search',
              scan: 'DataSource.Scan'
          },
          db: {
              srv: 'com.kazoup.srv.db',
              create: 'DB.Create',
              createIndex: 'DB.CreateIndex',
              delete: 'DB.Delete',
              read: 'DB.Read',
              search: 'DB.Search',
              searchById: 'DB.SearchById',
              status: 'DB.Status',
              update: 'DB.Update'
          },
          search: {
              srv: 'com.kazoup.srv.search',
              search: 'Search.Search',
              aggregate: 'Search.Aggregate'
          },
          file: {
              srv: 'com.kazoup.srv.file',
              create: 'File.Create',
              delete: 'File.Delete',
              share: 'File.Share'
          }
      }
  };
}());
