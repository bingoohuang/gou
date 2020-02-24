package lo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// nolint gomnd
func TestRotateFile(t *testing.T) {
	file, err := NewRotateFile("./var/logs/my.log",
		MaxBackups(3),
		TimeFormat("20160102150405"))
	_ = file.Close()

	assert.Nil(t, err)

	for i := 0; i < 5; i++ {
		_, _ = file.Write([]byte("hello"))
		time.Sleep(time.Second)
	}
}
