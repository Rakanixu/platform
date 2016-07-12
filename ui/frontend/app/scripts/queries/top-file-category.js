'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.topFileCategoryDashboard = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
            'bool': {
              'must': [
                {
                  'term': {
                    'exists_on_disk': true
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
        'top_categories': {
          'terms': {
            'field': 'doc_type',
            'order': {
              'total_size': 'desc'
            },
            'size': 1
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

  ElasticSearch.queries.topFileCategory = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
            'range': {
              'modified': {
                'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
              }
            }
          }
        }
      },
      'aggs': {
        'top_categories': {
          'terms': {
            'field': 'doc_type',
            'order': {
              'total_size': 'desc'
            },
            'size': 10
          },
          'aggs': {
            'total_size': {
              'sum': {
                'field': 'metadata.size'
              }
            },
            'stats': {
              'stats': {
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
