package gou

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapDefault(t *testing.T) {
	a := assert.New(t)

	m := map[int]int{1: 11}

	a.Equal(11, MapDefault(m, 1, 100))
	a.Equal(22, MapDefault(m, 2, 22))

	m2 := map[string]string{"1": "11"}

	a.Equal("11", MapDefault(m2, "1", "100"))
	a.Equal("22", MapDefault(m2, "2", "22"))
}

func TestIterateMapInt(t *testing.T) {
	a := assert.New(t)

	m := map[int]int{
		1: 11,
		2: 22,
		3: 33,
	}

	ks := ""
	vs := ""

	IterateMapSorted(m, func(k, v int) {
		ks += fmt.Sprintf("%d", k)
		vs += fmt.Sprintf("%d", v)
	})
	a.Equal("123", ks)
	a.Equal("112233", vs)
}

func TestIterateMapString(t *testing.T) {
	a := assert.New(t)

	m := map[string]int{
		"1": 11,
		"2": 22,
		"3": 33,
	}

	ks := ""
	vs := ""

	IterateMapSorted(m, func(k string, v int) {
		ks += k
		vs += fmt.Sprintf("%v", v)
	})
	a.Equal("123", ks)
	a.Equal("112233", vs)
}

func TestIterateMapFloat64(t *testing.T) {
	a := assert.New(t)

	m := map[float64]int{
		1.1: 11,
		2.2: 22,
		3.3: 33,
	}

	ks := ""
	vs := ""

	IterateMapSorted(m, func(k float64, v int) {
		ks += fmt.Sprintf("%.1f", k)
		vs += fmt.Sprintf("%v", v)
	})
	a.Equal("1.12.23.3", ks)
	a.Equal("112233", vs)
}

func TestIterateMapOther(t *testing.T) {
	a := assert.New(t)

	m := map[float32]int{
		1.1: 11,
		2.2: 22,
		3.3: 33,
	}

	ks := ""
	vs := ""

	IterateMapSorted(m, func(k float32, v int) {
		ks += fmt.Sprintf("%.1f", k)
		vs += fmt.Sprintf("%v", v)
	})
	a.Equal("1.12.23.3", ks)
	a.Equal("112233", vs)
}
