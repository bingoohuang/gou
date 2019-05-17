package gou

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChecksum(t *testing.T) {
	crc := Checksum([]byte("bigoohuang"))
	assert.Equal(t, "380372004", crc)
}
