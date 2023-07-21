package internal

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

var (
	shards        = [5]chan event{}
	workerTimeout = 2 * time.Minute
	workersAmount = getEnvInt("WORKERS_MAX", 3)
)

// iPoolEntry main purpose to provide a chance of mocking pools
// its main implementation - poolEntry
type iPoolEntry interface {
	ShutdownPools()
	Send([]event, string)
}

// poolEntry implements iPoolEntry
// wg is used for a graceful shutdown
type poolEntry struct {
	wg *sync.WaitGroup
}

// workerStoppedSignal used to signal that worker is completed.
// if channelClosed is true - no new worker should be created
type workerStoppedSignal struct {
	channelClosed bool
}

// ensure interface
var _ iPoolEntry = &poolEntry{}

// StartPools start 5 pools on the server start
func StartPools(db dal) iPoolEntry {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		shards[i] = newWorkerPool(i, &wg, db)
	}
	return poolEntry{wg: &wg}
}

// ShutdownPools shuts pools down by closing all worker channels and waiting for them to complete
func (p poolEntry) ShutdownPools() {
	defer fmt.Println("Shutting pools completed")

	for i := range shards {
		close(shards[i])
	}
	p.wg.Wait()
}

func (poolEntry) Send(events []event, singleShard string) {
	if singleShard != "true" {
		for _, e := range events {
			number := shardNumber(e)
			shards[number] <- e
		}
	} else {
		number := shardNumber(event{Timestamp: time.Now()})
		for _, e := range events {
			shards[number] <- e
		}
	}
}

// newWorkerPool creates a new pool with 3 starting workers
// pool recreates a worker when it's finished
func newWorkerPool(sharId int, wg *sync.WaitGroup, db dal) chan event {
	workerStopped := make(chan workerStoppedSignal)
	messagesChan := make(chan event)

	for i := 0; i < workersAmount; i++ {
		go worker(sharId, messagesChan, workerStopped, db)
	}

	go func(wg *sync.WaitGroup) {
		defer fmt.Println("Pool is finished. Id: ", sharId)
		defer wg.Done()

		workersNumber := workersAmount
		for workersNumber > 0 {
			channelClosed := (<-workerStopped).channelClosed
			if !channelClosed {
				go worker(sharId, messagesChan, workerStopped, db)
			} else {
				workersNumber--
			}
		}
	}(wg)

	return messagesChan
}

// worker is the main unit of processing events.
// it saves data about itself to the database and stars listening for events
// until it either times out or events channel is closed
// in both cases it sends workerStoppedSignal to the pool
// events are saved in batches of 5 unless the worked is stopped than whatever exist is saved
func worker(shardId int, messages <-chan event, done chan<- workerStoppedSignal, db dal) {
	channelClosed := false
	workerId := uuid.New().String()
	defer func() {
		db.archiveWorker(workerData{Id: workerId})
		done <- workerStoppedSignal{channelClosed: channelClosed}
		fmt.Println("Worker is finished. ", shardId, workerId)
	}()

	db.createWorker(workerData{Id: workerId, ShardId: shardId})
	batch := make([]event, 0)
	run := true

	for run {
		select {
		case m, open := <-messages:
			run = open
			channelClosed = !open

			if open {
				batch = append(batch, m)
				if len(batch) == 5 {
					db.saveEvent(eventsBatch{Events: batch, WorkerId: workerId})
					batch = nil
				}
			} else {
				db.saveEvent(eventsBatch{Events: batch, WorkerId: workerId})
			}
		case <-time.After(workerTimeout):
			run = false
			db.saveEvent(eventsBatch{Events: batch, WorkerId: workerId})
		}
	}
}

// shardNumber simply returns a number of a pool to use based on a timestamp.
func shardNumber(e event) int {
	unix := e.Timestamp.Unix()

	return int(unix) % 5
}
