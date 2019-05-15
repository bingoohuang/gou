package gou

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestMapKeys(t *testing.T) {
	keys := MapKeys(map[string]string{"k1": "v1", "k0": "v0"}).([]string)

	a := assert.New(t)
	sort.Strings(keys)
	a.Equal(keys, []string{"k0", "k1"})
}
