package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvInt(t *testing.T) {
	testData := []struct {
		name     string
		envVar   string
		def      int
		expected int
	}{
		{
			name:     "should return default when can't convert env var to int",
			envVar:   "nope",
			def:      1,
			expected: 1,
		},
		{
			name:     "should return env var when it's valid",
			envVar:   "123",
			def:      2,
			expected: 123,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			os.Setenv("TEST", v.envVar)
			res := getEnvInt("TEST", v.def)
			assert.Equal(t, v.expected, res)
		})
	}
}
