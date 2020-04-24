package file

import (
	"os"
	"testing"

	"github.com/bingoohuang/gou/lang"
	"github.com/stretchr/testify/assert"
)

func TestWriteTime(t *testing.T) {
	_ = os.Remove("/tmp/a.time")

	time, err := ReadTime("/tmp/a.time", "2020-04-23 12:35:19.000")
	assert.Nil(t, err)
	assert.Equal(t, lang.ParseTime(TimeFormat, "2020-04-23 12:35:19.000"), time)

	err = WriteTime("/tmp/a.time", lang.ParseTime(TimeFormat, "2020-04-23 12:35:20.000"))
	assert.Nil(t, err)

	time, err = ReadTime("/tmp/a.time", "2020-04-23 12:35:19.000")
	assert.Nil(t, err)
	assert.Equal(t, lang.ParseTime(TimeFormat, "2020-04-23 12:35:20.000"), time)

	_ = os.Remove("/tmp/a.time")
}
