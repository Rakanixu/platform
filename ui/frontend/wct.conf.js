var path = require('path');

var ret = {
  'suites': [
    'app/test'
  ],
  'webserver': {
    'pathMappings': []
  },
  'plugins': {
    'local': {
      'browsers': ['chrome'/*, 'firefox'*/] // Only chrome and ff aresupported
    },
    'istanbul': {
      'dir': './app/test/coverage',
      'reporters': [
        'text-summary',
        'lcov'
      ],
      // Path format is important, inline scripts are being serialized
      // to get test coverage, if the whitelist (inlcude) does not match the path requested
      // by the browser when the tests are being executed, web-component-tester-instanbul
      // won't serialize the script, and thecoverage result will be green (yei!) but 0/0 (upss..)
      'include': [
        '/app/elements/kazoup-ajax/kazoup-ajax.html',
        '/app/elements/kazoup-analytics/kazoup-analytics.html',
        '/app/elements/kazoup-analytics-brush-slider/kazoup-analytics-brush-slider.html',
        '/app/elements/kazoup-analytics-data/kazoup-analytics-data.html',
        '/app/elements/kazoup-analytics-filters/kazoup-analytics-filters.html',
        '/app/elements/kazoup-analytics-filters-menu/kazoup-analytics-filters-menu.html',
        '/app/elements/kazoup-breadcrumbs/kazoup-breadcrumbs.html',
        '/app/elements/kazoup-checksum-icon/kazoup-checksum-icon.html',
        '/app/elements/kazoup-d3-charts/kazoup-d3-charts.html',
        '/app/elements/kazoup-d3-histogram/kazoup-d3-histogram.html',
        '/app/elements/kazoup-dashboard/kazoup-dashboard.html',
        '/app/elements/kazoup-docs/kazoup-docs.html',
        '/app/elements/kazoup-download-report/kazoup-download-report.html',
        '/app/elements/kazoup-empty-state/kazoup-empty-state.html',
        '/app/elements/kazoup-file-preview/kazoup-file-preview.html',
        '/app/elements/kazoup-filters-converter/kazoup-filters-converter.html',
        '/app/elements/kazoup-form/kazoup-form.html',
        '/app/elements/kazoup-login/kazoup-login.html',
        '/app/elements/kazoup-logs/kazoup-logs.html',
        '/app/elements/kazoup-menu/kazoup-menu.html',
        '/app/elements/kazoup-mini-chart/kazoup-mini-chart.html',
        '/app/elements/kazoup-mobile-analytics/kazoup-mobile-analytics.html',
        '/app/elements/kazoup-mobile-search/kazoup-mobile-search.html',
        '/app/elements/kazoup-policies-dialog/kazoup-policies-dialog.html',
        '/app/elements/kazoup-register-appliance/kazoup-register-appliance.html',
        '/app/elements/kazoup-register-demo/kazoup-register-demo.html',
        '/app/elements/kazoup-restore-appliance/kazoup-restore-appliance.html'
      ],
      'exclude': [
        '/app/bower_components/polymer/polymer.html'/*,
        '/app/bower_components/PinkySwear.js/pinkyswear.min.js'*/
      ]
    }
  }
};

var mapping = {};
var rootPath = (__dirname).split(path.sep).slice(-1)[0];

mapping['/' + rootPath  + '/app/bower_components'] = 'bower_components';

ret.webserver.pathMappings.push(mapping);

module.exports = ret;
