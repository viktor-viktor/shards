package internal

import "time"

// workerData keeps track of worker history
type workerData struct {
	Id          int `json:"id"`
	ShardId     int `json:"shard_id"`
	EventsCount int `json:"events_count"`
}

// event represents an input data to the endpoints
type event struct {
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// eventBatch represents a piece of data that is being saved by a worker
type eventsBatch struct {
	WorkerId int     `json:"workerId"`
	Events   []event `json:"events"`
}
