package internal

import (
	"bytes"
	"encoding/json"
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
