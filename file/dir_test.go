package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseDir(t *testing.T) {
	assert.Equal(t, "/a", BaseDir([]string{"/a/b", "/a/c"}))
	assert.Equal(t, "/", BaseDir([]string{"/a/b", "/b/c"}))
}
