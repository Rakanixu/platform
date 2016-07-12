'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.nextYearSizeDashboard = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query': {
        'constant_score': {
          'filter': {
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
              ]
            }
          }
        }
      },
      'aggs': {
        'buckets': {
          'range': {
            'field': 'modified',
            'ranges': [
              {
                'key' : 'all_data',
                'from': 0,
                'to': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
              },
              {
                'key': 'all_data_except_last_1_year',
                'from': 0,
                'to': ElasticSearch.utils.getDateByMonthsAgo(12).getTime()
              },
              {
                'key': 'all_data_except_last_2_year',
                'from': 0,
                'to': ElasticSearch.utils.getDateByMonthsAgo(24).getTime()
              },
              {
                'key': 'all_data_except_last_3_year',
                'from': 0,
                'to': ElasticSearch.utils.getDateByMonthsAgo(36).getTime()
              }
            ]
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

  ElasticSearch.queries.nextYearSize = JSON.stringify({
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
                  'prefix': {
                    'dirpath': '//'
                  }
                }
              ],
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
