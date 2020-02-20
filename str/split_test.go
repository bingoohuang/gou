package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitN(t *testing.T) {
	a := assert.New(t)
	a.Equal("bingoo", Split1("bingoo%huang", "%", true, false))

	x, y := Split2("bingoo%%huang", "%", true, true)
	a.Equal("bingoo", x)
	a.Equal("huang", y)

	x, y = Split2("bingoo%huang%xxx", "%", true, true)
	a.Equal("bingoo", x)
	a.Equal("huang", y)

	x, y = Split2("bingoo", "%", true, true)
	a.Equal("bingoo", x)
	a.Equal("", y)
}

// nolint gomnd
func TestDecode(t *testing.T) {
	a := assert.New(t)
	a.Equal(2, Decode("a", "a", 2).(int))
	a.Equal(3, Decode("a", "b", 2, 3).(int))
	a.Nil(Decode("a", "b", 2))
}

func TestSplitToMap(t *testing.T) {
	a := assert.New(t)
	a.Equal(SplitToMap("m:mean;s:sum", ":", ";"), map[string]string{"m": "mean", "s": "sum"})
	a.Equal(SplitToMap("", ":", ";"), map[string]string{})
	a.Equal(SplitToMap("aa", ":", ";"), map[string]string{"aa": ""})
	a.Equal(SplitToMap("aa;bb:1", ":", ";"), map[string]string{"aa": "", "bb": "1"})
}
