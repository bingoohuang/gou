package ran

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint gomnd
func TestRand(t *testing.T) {
	num := Num(10)
	assert.True(t, len(num) == 10)
}
