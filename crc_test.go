package gou

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	crc := Checksum([]byte("bigoohuang"))
	assert.Equal(t, "380372004", crc)
}
