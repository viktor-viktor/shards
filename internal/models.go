package internal

import "time"

// workerData keeps track of worker history
type workerData struct {
	Id          string    `json:"id"`
	ShardId     int       `json:"shard_id"`
	EventsCount int       `json:"events_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// event represents an input data to the endpoints
type event struct {
	Timestamp time.Time   `json:"timestamp,required"`
	Data      interface{} `json:"data"`
}

// eventBatch represents a piece of data that is being saved by a worker
type eventsBatch struct {
	WorkerId string  `json:"worker_id"`
	Events   []event `json:"events"`
}
