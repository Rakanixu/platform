'use strict';
/* jshint -W079 */
/* globals self: false */

var window = self;

(function() {
  self.importScripts(
    '/bower_components/qwest/qwest.min.js',
    '/bower_components/filesize.js/lib/filesize.min.js'
  );

  self.makeScrollRequest = function(url, query, token) {
    return qwest.post(
      url,
      query,
      {
        dataType: 'json',
        cache: true,
        headers: {
          Accept: 'application/json',
          Authorization: 'JWT ' + token
        }
      }
    );
  };

  self.generateCSV = function(data) {
     var dataString;
     var csvContent = '';

     data.forEach(function(infoArray, index) {
       dataString = infoArray.join('\t');
       csvContent += index < data.length ? dataString + '\n' : dataString;
     });

     self.postMessage({
       href: 'data:attachment/csv,' + encodeURIComponent(csvContent)
     });
  };

  self.handleResults = function(xhr, response) {
    var data = [];

    if (self.exportAs === 'search') {
      data.push([
        'Full path',
        'Filename',
        'Category',
        'Extension',
        'Modified',
        'Accessed',
        'Size',
        'Permissions'
      ]);

      response.hits.hits.forEach(function(item) {
        data.push([
          item._source.metadata.fullpath,
          item._source.metadata.filename,
          item._source.metadata.doc_type,
          item._source.metadata.extension,
          item._source.metadata.modified,
          item._source.metadata.accessed,
          filesize(item._source.metadata.size),
          item._source.permissions.allow.join()
        ]);
      });
    }

    if (self.exportAs === 'logs') {
      response.hits.hits.forEach(function(item) {
        data.push([
          item._source.category,
          item._source.level,
          item._source.message,
          item._source.timestamp,
        ]);
      });
    }

    self.generateCSV(data);
  };

  self.requestCSVGeneration = function(e) {
    self.exportAs = e.data.type;

    self.makeScrollRequest(e.data.url, e.data.params, e.data.authToken)
      .then(self.handleResults);
  };

  // Main app will ask worker for CSV generation
  self.addEventListener('message', self.requestCSVGeneration, false);
}());
