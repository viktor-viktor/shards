package internal

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dalMocked struct {
	workers []workerData
	events  []eventsBatch
	err     error
}

func (d *dalMocked) saveWorker(data workerData) error {
	d.workers = append(d.workers, data)
	return d.err
}

func (d *dalMocked) saveEvent(batch eventsBatch) error {
	d.events = append(d.events, batch)
	return d.err
}

func (d dalMocked) getAllWorkers() ([]workerData, error) {
	return d.workers, d.err
}

func (d dalMocked) getWorker(string) (workerData, error) {
	if len(d.workers) > 0 {
		return d.workers[0], d.err
	}
	return workerData{}, d.err
}

func TestStartPools(t *testing.T) {
	p := StartPools(&dalMocked{})

	done := make(chan struct{})
	go func() {
		p.ShutdownPools()
		done <- struct{}{}
	}()

	select {
	case <-time.After(time.Second * 5):
		t.Error("Pools aren't closed despite closing channels")
	case <-done:
	}
}

func TestNewWorkerPool(t *testing.T) {
	testData := []struct {
		name           string
		timeout        time.Duration
		wait           time.Duration
		expectWorker   int
		valuesToSend   []event
		valuesExpected int
	}{
		{
			name:         "Should create 3 workers on start",
			timeout:      4 * time.Second,
			wait:         2 * time.Second,
			expectWorker: 1,
		},
		{
			name:         "Should recreate workers after timeout",
			timeout:      3 * time.Second,
			wait:         5 * time.Second,
			expectWorker: 2,
		},
		{
			name:           "Should save full batch when more than 5 events given",
			timeout:        10 * time.Second,
			wait:           time.Second,
			expectWorker:   1,
			valuesToSend:   make([]event, 6),
			valuesExpected: 6,
		},
		{
			name:           "Should save not full batch when shutdown quickly",
			timeout:        10 * time.Second,
			wait:           time.Second,
			expectWorker:   1,
			valuesToSend:   make([]event, 3),
			valuesExpected: 3,
		},
	}
	// enforce single worker for tests
	workersAmount = 1

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			// start pool
			wg := sync.WaitGroup{}
			mdb := dalMocked{}
			wg.Add(1)
			ch := newWorkerPool(0, &wg, &mdb)

			// set time before workers timeout
			workerTimeout = v.timeout
			// wait some time to let workers either start or start + restart
			time.Sleep(v.wait)

			// verify amount of workers created
			assert.Equal(t, v.expectWorker, len(mdb.workers), "Unexpected amount of workers created")
			for _, val := range v.valuesToSend {
				ch <- val
			}

			// verify that pool has shutdown
			done := make(chan struct{})
			go func() {
				close(ch)
				wg.Wait()
				done <- struct{}{}
			}()
			select {
			case <-time.After(time.Second * 5):
				t.Error("Pools aren't closed despite shutdown")
			case <-done:
			}

			// verify that all sent events are saved
			if v.valuesExpected != 0 {
				amount := 0
				for _, val := range mdb.events {
					amount += len(val.Events)
				}

				assert.Equal(t, v.valuesExpected, amount, "Unexpected amount of events received")
				if v.valuesExpected > 5 {
					assert.Equal(t, 5, len(mdb.events[0].Events), "First batch should be full (5) when more than 5 events processed")
				}
			}
		})
	}
}
