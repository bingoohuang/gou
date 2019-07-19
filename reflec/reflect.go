package reflec

import (
	"reflect"
	"strings"
	"sync"
)

// nolint
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
