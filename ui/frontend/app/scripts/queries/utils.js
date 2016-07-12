'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.utils = ElasticSearch.utils || {};

(function() {
  ElasticSearch.utils.getDateByMonthsAgo = function(months) {
    var now = new Date();

    return new Date(now.setMonth(now.getMonth() - months));
  };

  return ElasticSearch;
}(ElasticSearch));
