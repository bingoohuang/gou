package gou

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeSliceGroup0(t *testing.T) {
	type Bean struct {
		name string
		age  int
	}
	var slice []Bean

	a := assert.New(t)

	_, er := MakeSliceGroup(nil, func(i interface{}) interface{} { return i.(Bean).name })
	a.NotNil(er)

	g, e := MakeSliceGroup(slice, func(i interface{}) interface{} { return i.(Bean).name })
	a.Nil(e)

	_, _, ok := g.NextGroup()
	a.False(ok)

}

func TestMakeSliceGroup1(t *testing.T) {
	type Bean struct {
		name string
		age  int
	}
	slice := []Bean{
		{name: "bingoo", age: 1},
		{name: "bingoo", age: 2},
	}

	a := assert.New(t)
	g, e := MakeSliceGroup(slice, func(i interface{}) interface{} { return i.(Bean).name })
	a.Nil(e)

	gv, gs, ok := g.NextGroup()

	a.True(ok)
	a.Equal("bingoo", gv)
	a.Equal(slice, gs)

	_, _, ok = g.NextGroup()
	a.False(ok)

}

func TestMakeSliceGroup2(t *testing.T) {
	type Bean struct {
		name string
		age  int
	}
	slice := []Bean{
		{name: "bingoo", age: 1},
		{name: "dingoo", age: 2},
	}

	a := assert.New(t)
	g, e := MakeSliceGroup(slice, func(i interface{}) interface{} { return i.(Bean).name })
	a.Nil(e)

	gv, gs, ok := g.NextGroup()

	a.True(ok)
	a.Equal("bingoo", gv)
	a.Equal([]Bean{{name: "bingoo", age: 1}}, gs)

	gv, gs, ok = g.NextGroup()

	a.True(ok)
	a.Equal("dingoo", gv)
	a.Equal([]Bean{{name: "dingoo", age: 2}}, gs)

	gv, gs, ok = g.NextGroup()

	a.False(ok)
	a.Equal("", gv)
	a.Nil(gs)
}
