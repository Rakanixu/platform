'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.archiveSizeDashboard = JSON.stringify({
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
                  'exists': {
                    'field': 'archive_complete'
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

  ElasticSearch.queries.archiveSize = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
            'bool': {
              'must': [
                {
                  'prefix': {
                    'dirpath': '//'
                  }
                },
                {
                  'range': {
                    'modified': {
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
                    }
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
        }
      },
      'aggs': {
        'archive_size': {
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
