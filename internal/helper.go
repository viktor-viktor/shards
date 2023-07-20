package internal

import (
	"os"
	"strconv"
)

func getEnvInt(key string, def int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return def
	}
	return val
}
