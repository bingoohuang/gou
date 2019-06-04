package gou

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	keys := MapKeys(map[string]string{"k1": "v1", "k0": "v0"}).([]string)

	a := assert.New(t)
	sort.Strings(keys)
	a.Equal(keys, []string{"k0", "k1"})
}
