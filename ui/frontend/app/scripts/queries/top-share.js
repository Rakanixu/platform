'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.topShareDashboard = JSON.stringify({
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
        'top_shares': {
          'terms': {
            'field': 'sharepath',
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

  ElasticSearch.queries.topShare = JSON.stringify({
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
        'top_shares': {
          'terms': {
            'field': 'sharepath',
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
