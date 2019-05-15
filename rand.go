package gou

import (
	"fmt"
	"math"
	"strconv"

	"github.com/averagesecurityguy/random"
)

// RandomFloat64 returns a random float64
func RandomFloat64() float64 {
	i, _ := random.Int64()
	return float64(i) * 1.0842021724855043e-19
}

// RandomIntN returns a random int
func RandomIntN(n uint64) int {
	i, _ := random.Uint64Range(0, n)
	return int(i)
}

// RandomInt64 returns a random int64
func RandomInt64() int64 {
	i, _ := random.Int64()
	return i
}

// RandomString returns a random string
func RandomString(n uint64) string {
	s, _ := random.AlphaNum(n)
	return s
}

// RandomNumString returns a random number string
func RandomNum(n int) string {
	f := fmt.Sprintf("%%0%dd", n)
	s := fmt.Sprintf(f, RandomInt())
	if len(s) > n {
		return s[0:n]
	}

	return s
}

// RandomInt returns a random positive int
func RandomInt() int {
	i, _ := random.Uint64Range(0, math.MaxInt32)
	return int(i)
}

// RandomIntAsString returns a random positive int
func RandomIntAsString() string {
	return strconv.Itoa(RandomInt())
}
