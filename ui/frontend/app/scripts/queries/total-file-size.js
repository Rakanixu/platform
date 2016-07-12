'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.totalFileSizeAndCountDashboard = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
            'bool':{
              'must':[
                {
                  'bool':{
                    'should':[
                      {
                        'term':{
                          'exists_on_disk':true
                        }
                      },
                      {
                        'exists':{
                          'field':'archive_complete'
                        }
                      }
                    ],
                    'must':[
                    ]
                  }
                },
                {
                  'range': {
                    'modified': {
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
                    }
                  }
                }
              ]
            }
          }
        }
      },
      'aggs': {
        'total_size': {
          'stats': {
            'field': 'metadata.size'
          }
        }
      },
      'size': 0
    }
  });

  ElasticSearch.queries.totalFileSize = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
            'bool': {
              'must': [
                {
                  'range': {
                    'modified': {
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
                    }
                  }
                },
                {
                  'bool': {
                    'should': [
                      {
                        'term': {
                          'exists_on_disk': true
                        }
                      },
                      {
                        'exists': {
                          'field': 'archive_complete'
                        }
                      }
                    ]
                  }
                }
              ]
            }
          }
        }
      },
      'aggs': {
        'total_size': {
          'date_histogram' : {
            'field' : 'modified',
            'interval' : 'year'
          },
          'aggs': {
            'total_size': {
              'sum': {
                'field': 'metadata.size'
              }
            }
          }
        }
      },
      'size': 0
    }
  });

  return ElasticSearch;
}(ElasticSearch));
