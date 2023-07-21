package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// EventsController expects []events on the input and pushed them to respective pool to process
func BuildEventsController(p iPoolEntry) func(*gin.Context) {
	return func(c *gin.Context) {

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": fmt.Sprintf("failed to read the body. Error: %v", err),
			})
			return
		}

		events := make([]event, 0)
		if err := json.Unmarshal(data, &events); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": fmt.Sprintf("failed to unmarshal body. Error: %v", err),
			})
			return
		}

		singleShard := c.Query("single_shard")

		p.Send(events, singleShard)
		c.JSON(http.StatusOK, nil)
	}
}

// BuildWorkersController initializes using provided dal and returns a controller that returns all existing,
// or finished, workers
func BuildWorkersController(db dal) func(c *gin.Context) {
	return func(c *gin.Context) {
		workers, err := db.getAllWorkers()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"workers": workers,
		})
	}
}

// BuildSingleWorkerController initializes using provided dal and returns a controller searches for a single worker
func BuildSingleWorkerController(db dal) func(c *gin.Context) {
	return func(c *gin.Context) {
		Id := c.Param("id")
		wd, err := db.getWorker(Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": fmt.Sprintf("Error when looking for an item: %v", err),
			})
			return
		}

		if wd.Id == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": fmt.Sprintf("worker with such Id: %v isn't found", Id),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"workerData": wd,
		})
	}
}
