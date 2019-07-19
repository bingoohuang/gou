package rand

import (
	"fmt"
	"math"
	"strconv"

	"github.com/averagesecurityguy/random"
)

// Float64 returns a random float64
func Float64() float64 {
	i, _ := random.Int64()
	return float64(i) * 1.0842021724855043e-19
}

// IntN returns a random int
func IntN(n uint64) int {
	i, _ := random.Uint64Range(0, n)
	return int(i)
}

// Int64 returns a random int64
func Int64() int64 {
	i, _ := random.Int64()
	return i
}

// String returns a random string
func String(n uint64) string {
	s, _ := random.AlphaNum(n)
	return s
}

// NumString returns a random number string
func Num(n int) string {
	f := fmt.Sprintf("%%0%dd", n)
	s := fmt.Sprintf(f, Int())
	if len(s) > n {
		return s[0:n]
	}

	return s
}

// Int returns a random positive int
func Int() int {
	i, _ := random.Uint64Range(0, math.MaxInt32)
	return int(i)
}

// IntAsString returns a random positive int
func IntAsString() string {
	return strconv.Itoa(Int())
}
