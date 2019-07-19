package math

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

// MaxInt 返回a和b中最大值
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt 返回a和b中最小值
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxFloat32 返回a和b中最大值
func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// MinFloat32 返回a和b中最小值
func MinFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Sum 计算slice中指定值的和
func Sum(arr interface{}, f func(v interface{}) int) int {
	value := redirectValue(reflect.ValueOf(arr))
	kind := value.Kind()
	s := 0
	if kind == reflect.Array || kind == reflect.Slice {
		length := value.Len()
		for i := 0; i < length; i++ {
			elem := redirectValue(value.Index(i)).Interface()
			s += f(elem)
		}

		return s
	}

	logrus.Warnf("Type %s is not supported", value.Type().String())
	return s
}

func redirectValue(value reflect.Value) reflect.Value {
	for {
		if !value.IsValid() || value.Kind() != reflect.Ptr {
			return value
		}

		res := reflect.Indirect(value)

		// Test for a circular type.
		if res.Kind() == reflect.Ptr && value.Pointer() == res.Pointer() {
			return value
		}

		value = res
	}
}

// Valuer 值计算器
type Valuer interface {
	// Reset 重置状态
	Reset()
	// Tag 打点数据
	Tap(float32)
	// Value 获取数据，以及数据是否有效
	Value() (float32, bool)
	// Value 获取数据
	PureValue() float32
}

// ValuerBase 表示Valuer的基础实现
type ValuerBase struct {
	started bool
	value   float32
}

// Maxer 求最大值
type Maxer struct {
	ValuerBase
}

// Reset 重置状态
func (m *ValuerBase) Reset() {
	m.started = false
	m.value = 0
}

// PureValue 获取数据
func (m ValuerBase) PureValue() float32 { return m.value }

// Value 获取数据，以及数据是否有效
func (m ValuerBase) Value() (float32, bool) { return m.value, m.started }

// Tap 打点数据
func (m *Maxer) Tap(v float32) {
	if !m.started {
		m.started = true
		m.value = v
	} else if v > m.value {
		m.value = v
	}
}

// Miner 求最小值
type Miner struct {
	ValuerBase
}

// Tap 打点数据
func (m *Miner) Tap(v float32) {
	if !m.started {
		m.started = true
		m.value = v
	} else if v < m.value {
		m.value = v
	}
}
