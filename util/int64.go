package util

import (
	"math/rand"
	"time"
)

// Int64MinMax TBD
func Int64MinMax(a, b int64) (min, max int64) {
	if a < b {
		return a, b
	}
	return b, a
}

// NextUniqID TBD
func NextUniqID() int64 {
	return (time.Now().Unix() << 32) | rand.Int63n(2147483646)
}

// PInt64HasChnaged TBD
func PInt64HasChnaged(old, new *int64) bool {
	fromValueToNil := old != nil && new == nil
	fromNilToValue := old == nil && new != nil
	valueChanged := (old != nil && new != nil && *old != *new)

	return fromValueToNil || fromNilToValue || valueChanged
}
