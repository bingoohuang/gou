package gou

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestTpl(t *testing.T) {
	a := assert.New(t)
	a.Equal("bingoohuang", Tpl("{name}", map[string]interface{}{"name": "bingoohuang"}))
}

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

func TestRand(t *testing.T) {
	a := assert.New(t)
	num := RandomNum(10)
	fmt.Println(num)
	a.NotEmpty(num)
}

var c = make(chan bool, 3)

func TestSlice(t *testing.T) {
	a := make(map[string][]string)
	a["k1"] = make([]string, 1)
	a["k2"] = make([]string, 1)
	a["k3"] = make([]string, 1)

	for k, v := range a {
		go f(k, v)
	}

	<-c
	<-c
	<-c
}

func f(k string, v []string) {
	if k == "k1" {
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("2:%p\n", v)
	c <- true
}

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
