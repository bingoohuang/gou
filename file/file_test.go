package file

import (
	"os"
	"testing"

	"github.com/bingoohuang/gou/lang"
	"github.com/stretchr/testify/assert"
)

func TestWriteTime(t *testing.T) {
	_ = os.Remove("/tmp/a.time")

	time, err := ReadTime("/tmp/a.time", "2020-04-23 12:35:19")
	assert.Nil(t, err)
	assert.Equal(t, lang.ParseTime(TimeFormat, "2020-04-23 12:35:19"), time)

	err = WriteTime("/tmp/a.time", lang.ParseTime(TimeFormat, "2020-04-23 12:35:20"))
	assert.Nil(t, err)

	time, err = ReadTime("/tmp/a.time", "2020-04-23 12:35:19")
	assert.Nil(t, err)
	assert.Equal(t, lang.ParseTime(TimeFormat, "2020-04-23 12:35:20"), time)

	_ = os.Remove("/tmp/a.time")
}
