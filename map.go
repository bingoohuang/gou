package gou

import (
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
		p := strings.TrimSpace(pair)
		if p == "" {
			continue
		}

		k, v := Split2(p, kvSep, true, false)
		m[k] = v
	}

	return m
}
