## Installation 
Requirements: `docker`, `docker-compose`
1. Clone the repo
2. Run from the root directory: `./scripts/start.sh`

Once the script finishes running you should be able to make calls to the server by `localhost:9753`

## Usage
Server provides 3 endpoints:  
#### POST  - `/events` .  
Takes a json body in the form of:
```json
[
  {
    "timestamp": "2019-10-04T01:52:07-03:00",
    "data": 12312
  }
]
```
You can generate request data here: https://json-generator.com/ using the following config:
```javascript
[
  '{{repeat(10000, 15000)}}',
  {
    "timestamp": '{{date(new Date(2014, 0, 1), new Date(), "YYYY-MM-ddThh:mm:ssZ")}}',
    "data": 12312
  }
]
```
`data` field is optional and can be any valid json.


#### GET - `/workers`  
Returns all workers that were created during the session.
Response looks like:
```json
{
    "workers": [
        {
            "id": 7965774568352684057,
            "shard_id": 0,
            "events_count": 1775,
            "created_at": "2023-07-20T17:06:14.831Z"
        }
    ]
}
```

#### GET - `/workers/:id`
Returns detailed data about one of the workers.
```json
{
    "workerData": {
        "id": 7965774568352684057,
        "shard_id": 0,
        "events_count": 1779,
        "created_at": "2023-07-20T17:06:14.831Z"
    }
}
```

## Testing  
Requirements: you need `go` binaries and `go tools` installed.  
To test run from the root directory: `./scripts/test.sh`  
You can add `-v` flag to get visualize test coverage.  