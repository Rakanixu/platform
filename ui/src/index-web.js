var Endpoints = (function() {
  return {
    endpoint: 'https://web.kazoup.io:8082/rpc',
    web: 'https://web.kazoup.io:8082',
    socket: 'wss://web.kazoup.io:8082',
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