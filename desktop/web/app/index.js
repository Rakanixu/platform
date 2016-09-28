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
      endpoint: 'http://10.17.57.130:8082/rpc',
      web: 'http://10.17.57.130:8082',
      srvs:{
          config: {
              srv: 'com.kazoup.srv.config',
              setFlags: 'Config.SetFlags',
              status: 'Config.Status'
          },
          crawler: {
              srv: 'com.kazoup.srv.crawler',
              search: 'Crawler.Search'
          },
          datasource: {
              srv: 'com.kazoup.srv.datasource',
              create: 'Datasource.Create',
              delete: 'Datasource.Delete',
              search: 'Datasource.Search',
              scan: 'Datasource.Scan'
          },
          db: {
              srv: 'com.kazoup.srv.db',
              create: 'DB.Create',
              createIndexWithSettings: 'DB.CreateIndexWithSettings',
              delete: 'DB.Delete',
              putMappingFromJSON: 'DB.PutMappingFromJSON',
              read: 'DB.Read',
              search: 'DB.Search',
              status: 'DB.Status',
              update: 'DB.Update'
          },
          flag: {
              srv: 'com.kazoup.srv.flag',
              create: 'Flag.Create',
              delete: 'Flag.Delete',
              flip: 'Flag.Flip',
              list: 'Flag.List',
              read: 'Flag.Read'
          },
          search: {
              srv: 'com.kazoup.srv.search',
              search: 'Search.Search'
          }
      }
  };
}());
