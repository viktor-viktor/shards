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
You can set optional query parameter `single_shard` to `true` if you want all events  
to processed by a single random shard.


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
            //"finished_at" is optional
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
      //"finished_at" is optional
    }
}
```

## Testing  
Requirements: you need `go` binaries and `go tools` installed.  
To test run from the root directory: `./scripts/test.sh`  
You can add `-v` flag to get visualize test coverage.  

## Tasks addressed
 - the server has 5 shards pools with their own set of workers, with 3 workers by default. 
 - workers have a timeout of 2 minutes
 - workers write data in batches of 5 messages.
 - history about workers is available through `/workers` and `/workers/:id` endpoints
 - Testability: DAL and Pool provide interfaces for easier tests integrations. There is also a simple 
versions of DI. The system is small and benefits little from it, but it shows the general idea.
 - Shipping to the cloud. While I didn't ship it to any of the clouds, the app provides dockerfile,
and the general local setup in docker, which should limit cloud challenge to the cloud setup itself, 
and simply pulling + running the image.
 - Protecting on the public internet. I added and option to start service with https connection. I didn't
add basic or JWT authentication. Many additional improvements (such as rate limiting) are generally done 
on infrastructure/cloud level thus aren't added here.
 - Other small improvements. 
   - added graceful shutdown with minimizing data loss.
   - added scripts for easier tests and local setup.
   - added time fields to the workers to better visualize their history.
 