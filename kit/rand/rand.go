package rand

import (
	"math/rand"
	"time"
)

// Bool returns a randomly
// generated boolean value.
func Bool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%2 == 0
}

// Intn returns, as an int,
// a non-negative pseudo-random
// number in the half-open interval
// [0,n) from the default Source.
// It panics if n <= 0.
func Intn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
