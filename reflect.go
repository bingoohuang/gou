package gou

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// GetSliceByPtr 检查一个值v是否是Slice的指针，返回slice切片的reflect值
func GetSliceByPtr(v interface{}) (reflect.Value, error) {
	iv := reflect.ValueOf(v)
	nilValue := reflect.ValueOf(nil)
	if iv.Kind() != reflect.Ptr {
		return nilValue, fmt.Errorf("non-pointer %v", iv.Type())
	}

	// get the value that the pointer v points to.
	ve := iv.Elem()
	if ve.Kind() != reflect.Slice {
		return nilValue, fmt.Errorf("can't fill non-slice value")
	}

	return ve, nil
}

// EnsureSliceLen grows the slice capability
func EnsureSliceLen(v reflect.Value, len int) {
	// Grow slice if necessary
	if len >= v.Cap() {
		cap2 := v.Cap() + v.Cap()/2
		if cap2 < 4 {
			cap2 = 4
		}
		if cap2 < len {
			cap2 = len
		}

		v2 := reflect.MakeSlice(v.Type(), v.Len(), cap2)
		reflect.Copy(v2, v)
		v.Set(v2)
	}
	if len >= v.Len() {
		v.SetLen(len + 1)
	}
}

var fieldCache sync.Map // map[reflect.Type][]field

// StructField 表示一个struct的字段属性
type StructField struct {
	Index     int
	Name      string
	Tag       []string
	Kind      reflect.Kind
	Type      reflect.Type
	Anonymous bool
}

// CachedStructFields caches fields of struct type
func CachedStructFields(t reflect.Type, tagKey string) []StructField {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]StructField)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t, tagKey))
	return f.([]StructField)
}

func typeFields(t reflect.Type, tagKey string) []StructField {
	ff := t.NumField()
	fields := make([]StructField, ff)

	for fi := 0; fi < ff; fi++ {
		f := t.Field(fi)

		fields[fi] = StructField{
			Index:     fi,
			Name:      f.Name,
			Tag:       readTag(f, tagKey),
			Kind:      f.Type.Kind(),
			Type:      f.Type,
			Anonymous: f.Anonymous,
		}
	}

	return fields
}

// CanAssign accept strings, or empty interfaces.
func CanAssign(t reflect.Type, kinds ...reflect.Kind) bool {
	vk := t.Kind()
	if vk == reflect.Interface && t.NumMethod() == 0 {
		return true
	}

	for _, kind := range kinds {
		if kind == vk {
			return true
		}
	}

	return false

}

func readTag(f reflect.StructField, tagKey string) []string {
	val, ok := f.Tag.Lookup(tagKey)
	if !ok {
		return []string{}
	}
	return strings.Split(val, "/")
}
