# Policy API

Policy API for managing Kazoup data policies. This API exposes Policy srv over HTTP interface.

FQDN: go.micro.api.policy

## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service

```
$ go run main.go
```

#### Usage

#### Create
```
[POST] http[domain:micro API port]/policy/create
{
    "filter": {"bool":{"must":[{"prefix":{"metadata.dirpath":"//10.17.57.20/print$/color"}},{"range":{"accessed":{"gte":"now-26y/d","lte":"now-0y/d"}}}]}}, 
    "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D", 
    "is_archive_policy": true,
    "is_deletion_policy": false, 
    "name": "policy name"
}

{}
```


#### Read
```
[POST] http[domain:micro API port]/policy/read
{
    "name": "policy name"
}

{
  "created": "2016-06-29 17:29:28.256763573 +0100 BST",
  "created_by": "unknown",
  "filter": {
    "bool": {
      "must": [
        {
          "prefix": {
            "metadata.dirpath": "//10.17.57.20/print$/color"
          }
        },
        {
          "range": {
            "accessed": {
              "gte": "now-26y/d",
              "lte": "now-0y/d"
            }
          }
        }
      ]
    }
  },
  "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
  "is_archive_policy": true,
  "name": "policy name"
}
```

#### Delete
```
[POST] http[domain:micro API port]/policy/delete
{
    "name": "policy name"
}

{}
```

#### List
```
[POST] http[domain:micro API port]/policy/list
{}

{
  "result": [
    {
      "created": "2016-06-30 09:47:03.521376973 +0100 BST",
      "created_by": "unknown",
      "filter": {
        "bool": {
          "must": [
            {
              "prefix": {
                "metadata.dirpath": "//10.17.57.20/print$/color"
              }
            },
            {
              "range": {
                "accessed": {
                  "gte": "now-26y/d",
                  "lte": "now-0y/d"
                }
              }
            }
          ]
        }
      },
      "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
      "name": "policy name1"
    },
    {
      "created": "2016-06-30 09:46:58.224188769 +0100 BST",
      "created_by": "unknown",
      "filter": {
        "bool": {
          "must": [
            {
              "prefix": {
                "metadata.dirpath": "//10.17.57.20/print$/color"
              }
            },
            {
              "range": {
                "accessed": {
                  "gte": "now-26y/d",
                  "lte": "now-0y/d"
                }
              }
            }
          ]
        }
      },
      "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
      "name": "policy name2"
    },
    {
      "created": "2016-06-30 09:47:07.721637975 +0100 BST",
      "created_by": "unknown",
      "filter": {
        "bool": {
          "must": [
            {
              "prefix": {
                "metadata.dirpath": "//10.17.57.20/print$/color"
              }
            },
            {
              "range": {
                "accessed": {
                  "gte": "now-26y/d",
                  "lte": "now-0y/d"
                }
              }
            }
          ]
        }
      },
      "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
      "name": "policy name3"
    },
    {
      "created": "2016-06-30 09:58:40.889553602 +0100 BST",
      "created_by": "unknown",
      "filter": {
        "bool": {
          "must": [
            {
              "prefix": {
                "metadata.dirpath": "//10.17.57.20/print$/color"
              }
            },
            {
              "range": {
                "accessed": {
                  "gte": "now-26y/d",
                  "lte": "now-0y/d"
                }
              }
            }
          ]
        }
      },
      "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
      "name": "policy name4"
    }
  ]
}
```
