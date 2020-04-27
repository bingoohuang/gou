package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseDir(t *testing.T) {
	assert.Equal(t, "/a", BaseDir([]string{"/a/b", "/a/c"}))
	assert.Equal(t, "/", BaseDir([]string{"/a/b", "/b/c"}))
}

func TestExists(t *testing.T) {
	assert.True(t, DoesExists("/tmp"))
	assert.False(t, DoesExists("/notexist"))
	assert.True(t, DoesNotExists("/notexist"))
	assert.False(t, DoesNotExists("/tmp"))
	assert.True(t, ExistsAsFile("dir.go"))
	assert.False(t, ExistsAsFile("."))
	assert.True(t, ExistsAsDir("."))
	assert.False(t, ExistsAsDir("dir.go"))
}
