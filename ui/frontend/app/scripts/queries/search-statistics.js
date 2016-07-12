'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.searchStatistics = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query':{
        'constant_score':{
          'filter':{
            'bool':{
              'must':[
                {
                  'range':{
                    'timestamp':{
                      'gt':'now-1M/d'
                    }
                  }
                },
                {
                  'term':{
                    'category':'Search'
                  }
                }
              ]
            }
          }
        }
      },
      'size':0,
      'aggs':{
        'searches':{
          'date_histogram':{
            'field':'timestamp',
            'interval':'day',
            'min_doc_count': 0, // Ensure to retrieve empty buckets
            'extended_bounds': { // Ensure to retrieve empty for the whole given range
              'min': 'now-1M/d',
              'max': 'now-0M/d'
            }
          }
        }
      }
    }
  });

  ElasticSearch.queries.searchStatisticsCount = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query':{
        'constant_score':{
          'filter':{
            'bool':{
              'must':[
                {
                  'range':{
                    'timestamp':{
                      'gt':'now-1M/d'
                    }
                  }
                },
                {
                  'term':{
                    'category':'Search'
                  }
                }
              ]
            }
          }
        }
      },
      'size':0
    }
  });

  return ElasticSearch;
}(ElasticSearch));
