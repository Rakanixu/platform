'use strict';

window.ElasticSearch = window.ElasticSearch || {};
ElasticSearch.queries = ElasticSearch.queries || {};

(function() {
  ElasticSearch.queries.duplicatesDashboard = JSON.stringify({
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
                    'modified':{
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
                    }
                  }
                },
                {
                  'term': {
                    'content.unique': false
                  }
                },
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
                    ]
                  }
                }
              ]
            }
          }
        }
      },
      'aggs': {
        'known': {
          'terms': {
            'field': 'content.unique'
          },
          'aggs': {
            'total_size':{
              'sum':{
                'field':'metadata.size'
              }
            }
          }
        }
      },
      'size':0
    }
  });

  ElasticSearch.queries.duplicates = JSON.stringify({
    'index': 'files',
    'type': 'file',
    'query': {
      'query':{
        'filtered':{
          'filter':{
            'bool':{
              'must':[
                {
                  'range':{
                    'modified':{
                      'lt': ElasticSearch.utils.getDateByMonthsAgo(0).getTime()
                    }
                  }
                },
                {
                  'term': {
                    'content.unique': false
                  }
                },
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
                    ]
                  }
                }
              ]
            }
          }
        }
      },
      'aggs': {
        'doc_types': {
          'terms': {
            'field': 'doc_type',
            'size': 0,
            'order': {
              'total_size': 'desc'
            }
          },
          'aggs':{
            'total_size':{
              'sum':{
                'field':'metadata.size'
              }
            }
          }
        }
      },
      'size':0
    }
  });

  return ElasticSearch;
}(ElasticSearch));
