## Usage

### Create
```
[POST] http[domain:micro API port]/datasource/create
{
    "endpoint": {
        "url": "local://home/"
    }
}

{}
```

### Delete
```
[POST] http[domain:micro API port]/datasource/delete
{
    "id":"AVW2cMk9kEu9W4s33YrP"
}

{}
```

### Search
```
[POST] http[domain:micro API port]/datasource/search
{
    "query":"home",
    "offset": 0, 
    "limit": 1000
}

{}
```
