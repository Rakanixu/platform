'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.accessedPastYearDashboard = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
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
              ],
              'must': [
                {
                  'range': {
                    'accessed': {
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime(),
                      'gt': ElasticSearch.utils.getDateByMonthsAgo(12).getTime()
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
          'sum': {
            'field': 'metadata.size'
          }
        }
      },
      'size': 0
    }
  });

  ElasticSearch.queries.accessedPastYear = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
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
              ],
              'must': [
                {
                  'range': {
                    'accessed': {
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime(),
                      'gt': ElasticSearch.utils.getDateByMonthsAgo(12).getTime()
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
            'size': 0
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
