package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockedPoolEntry struct {
	events []event
}

func (mockedPoolEntry) ShutdownPools() {}

func (m *mockedPoolEntry) Send(events []event) {
	m.events = append(m.events, events...)
}

func TestEventsController(t *testing.T) {
	validEvents, _ := json.Marshal([]event{{Timestamp: time.Now()}})

	testData := []struct {
		name           string
		body           io.ReadCloser
		expectedStatus int
	}{
		{
			name:           "should return 400 when failed to read body",
			body:           io.NopCloser(bytes.NewReader([]byte{'w'})),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 400 when failed to unmarshal events",
			body:           io.NopCloser(bytes.NewReader([]byte{'w'})),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 200 when processed events",
			body:           io.NopCloser(bytes.NewReader(validEvents)),
			expectedStatus: http.StatusOK,
		},
	}
	workersAmount = 1

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			// Create a test context with a mock request and response
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{Body: v.body}

			mockPool := &mockedPoolEntry{}
			BuildEventsController(mockPool)(ctx)
			assert.Equal(t, v.expectedStatus, ctx.Writer.Status(), "unexpected status code received")
		})
	}
}

func TestBuildWorkersController(t *testing.T) {
	testData := []struct {
		name           string
		workers        []workerData
		err            error
		expectedStatus int
	}{
		{
			name:           "should return 400 when dal returns error",
			err:            fmt.Errorf("error"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 200 when dal doesn't return error",
			err:            nil,
			workers:        []workerData{{Id: 123, EventsCount: 321}},
			expectedStatus: http.StatusOK,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			mockedDB := &dalMocked{workers: v.workers, err: v.err}
			BuildWorkersController(mockedDB)(ctx)
			assert.Equal(t, v.expectedStatus, ctx.Writer.Status(), "unexpected status code received")
		})
	}
}

func TestBuildSingleWorkerController(t *testing.T) {
	testData := []struct {
		name           string
		paramId        string
		workers        []workerData
		err            error
		expectedStatus int
	}{
		{
			name:           "Should return 400 when param isn't int",
			paramId:        "wow",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 400 when dal returns error",
			paramId:        "123",
			err:            fmt.Errorf("error"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 200 when dal doesn't return error",
			paramId:        "123",
			err:            nil,
			workers:        []workerData{{Id: 123, EventsCount: 321}},
			expectedStatus: http.StatusOK,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			engine := gin.Default()

			// register endpoint
			mockedDB := &dalMocked{workers: v.workers, err: v.err}
			engine.GET("/workers/:id", BuildSingleWorkerController(mockedDB))

			// mock request
			url := fmt.Sprintf("/workers/%v", v.paramId)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal("Error creating request:", err)
			}

			engine.ServeHTTP(writer, req)
			assert.Equal(t, v.expectedStatus, writer.Code, "unexpected status code received")
		})
	}
}
