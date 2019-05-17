package gou

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateSlice(t *testing.T) {
	a := []string{"0", "1", "2"}
	st := ""

	ok, err := IterateSlice(a, 1, func(s string) { st += s })
	assert.False(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "120", st)
}

func TestIterateSlice2(t *testing.T) {
	a := []string{"0", "1", "2"}
	st := ""

	ok, err := IterateSlice(a, 2, func(i int, s string) bool { st += s; return i == 0 })
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "20", st)
}

func TestIterateSlice3(t *testing.T) {
	a := []string{"0", "1", "2"}
	st := ""

	ok, err := IterateSlice(a, 2, func(i int, s string) (bool, error) { st += s; return i == 0, errors.New("returned to head") })
	assert.True(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, "20", st)
}

func TestSliceContains(t *testing.T) {
	a := []string{"0", "1", "2"}
	assert.True(t, SliceContains(a, "0"))
	assert.False(t, SliceContains(a, "3"))
}
