package web

import (
	"regexp"
	"time"

	filter "github.com/seiflotfy/cuckoofilter"
)

const (
	maxIDLen = 10 // length 10 will represent 62^10 ~ 8.39e17 urls
)

var base62Regex = regexp.MustCompile(`^[0-9A-Za-z]{1,10}$`)

// checkExpiredAt checks if exp is in the future
func checkExpiredAt(exp time.Time) bool {
	return time.Now().UTC().Unix() < exp.UTC().Unix()
}

// checkValidID checks:
// 1. id length <= 10
// 2. match regexp of base62 members
func checkValidID(id string) bool {
	if len(id) > maxIDLen {
		return false
	}
	return base62Regex.MatchString(id)
}

// checkCuckoo checks if the ID is possible
func checkCuckoo(id string, f *filter.Filter) bool {
	return f.Lookup([]byte(id))
}
