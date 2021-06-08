package env

import (
	"os"
	"strconv"
)

// MustGetString takes key and look up environment variables, fallback if not found
func MustGetString(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// MustGetInt resembles MustGetString but return int
func MustGetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return res
}
