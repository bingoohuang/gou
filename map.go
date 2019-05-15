package gou

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// MultiMap 多值map
type MultiMap struct {
	M map[interface{}][]interface{}
}

// MakeMultiMap 创建一个MultiMap
func MakeMultiMap() *MultiMap {
	return &MultiMap{M: make(map[interface{}][]interface{})}
}

// UrlGet 根据k或者v,ok
func (m *MultiMap) Get(k interface{}) (v []interface{}, ok bool) {
	v, ok = m.M[k]
	return
}

// UrlPut 加入k,v
func (m *MultiMap) Put(k, v interface{}) {
	l, ok := m.M[k]
	if ok {
		m.M[k] = append(l, v)
		return
	}

	l = make([]interface{}, 1)
	l[0] = v
	m.M[k] = l
}

// UrlDelete 删除k
func (m *MultiMap) Delete(k interface{}) {
	delete(m.M, k)
}

// SplitToMap 将字符串s分割成map,其中key和value之间的间隔符是kvSep, kv和kv之间的分隔符是kkSep
func SplitToMap(s string, kvSep, kkSep string) map[string]string {
	var m map[string]string

	ss := strings.Split(s, kkSep)
	m = make(map[string]string)
	for _, pair := range ss {
		k, v := Split2(pair, kvSep, true, false)
		m[k] = v
	}

	return m
}

// MapKeys 返回Map的key切片
func MapKeys(m interface{}) interface{} {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return nil
	}

	keyType := v.Type().Key()
	ks := reflect.MakeSlice(reflect.SliceOf(keyType), v.Len(), v.Len())
	i := 0
	for _, key := range v.MapKeys() {
		ks.Index(i).Set(key)
		i++
	}

	return ks.Interface()
}

// MapValues 返回Map的value切片
func MapValues(m interface{}) interface{} {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return nil
	}

	typ := v.Type().Elem()
	sl := reflect.MakeSlice(reflect.SliceOf(typ), v.Len(), v.Len())
	i := 0
	for _, key := range v.MapKeys() {
		sl.Index(i).Set(v.MapIndex(key))
		i++
	}

	return sl.Interface()
}

func MapDefault(m, k, defaultValue interface{}) interface{} {
	mv := reflect.ValueOf(m)
	if mv.Kind() != reflect.Map {
		return nil
	}

	v := mv.MapIndex(reflect.ValueOf(k))
	if v.IsValid() {
		return v.Interface()
	}

	return defaultValue
}

func IterateMapSorted(m interface{}, iterFunc interface{}) {
	mv := reflect.ValueOf(m)
	if mv.Kind() != reflect.Map {
		return
	}

	mapLen := mv.Len()
	keyType := mv.Type().Key()
	keyKind := keyType.Kind()
	var keyMap map[interface{}]string
	switch keyKind {
	case reflect.String:
	case reflect.Int:
	case reflect.Float64:
	default:
		keyMap = make(map[interface{}]string, mapLen)
	}

	ks := reflect.MakeSlice(reflect.SliceOf(keyType), mapLen, mapLen)
	i := 0
	for _, k := range mv.MapKeys() {
		if keyMap != nil {
			keyMap[k.Interface()] = fmt.Sprintf("%v", k.Interface())
		}

		ks.Index(i).Set(k)
		i++
	}

	ksi := ks.Interface()

	if keyMap != nil {
		sort.Slice(ksi, func(i, j int) bool {
			ki := keyMap[ks.Index(i).Interface()]
			kj := keyMap[ks.Index(j).Interface()]
			return ki < kj
		})
	} else {
		switch keyKind {
		case reflect.String:
			sort.Strings(ksi.([]string))
		case reflect.Int:
			sort.Ints(ksi.([]int))
		case reflect.Float64:
			sort.Float64s(ksi.([]float64))
		}
	}

	funcValue := reflect.ValueOf(iterFunc)

	for j := 0; j < mapLen; j++ {
		k := ks.Index(j)
		v := mv.MapIndex(k)
		funcValue.Call([]reflect.Value{k, v})
	}
}
