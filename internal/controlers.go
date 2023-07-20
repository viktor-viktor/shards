package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
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

		p.Send(events)
		c.JSON(http.StatusOK, nil)
	}
}

// BuildWorkersController initializes using provided dal and returns a controller that returns all existing,
// or finished, workers
func BuildWorkersController(db dal) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("workers controller")
		workers := db.getAllWorkers()
		c.JSON(http.StatusOK, gin.H{
			"workers": workers,
		})
	}
}

// BuildSingleWorkerController initializes using provided dal and returns a controller searches for a single worker
func BuildSingleWorkerController(db dal) func(c *gin.Context) {
	return func(c *gin.Context) {
		strId := c.Param("id")
		Id, err := strconv.Atoi(strId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": fmt.Sprintf("invlaid worker id provided: %v", err),
			})
			return
		}

		wd := db.getWorker(Id)
		c.JSON(http.StatusOK, gin.H{
			"workerData": wd,
		})
	}
}
