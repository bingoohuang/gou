package coll

// MultiMap 多值map
type MultiMap struct {
	M map[interface{}][]interface{}
}

// MakeMultiMap 创建一个MultiMap
func MakeMultiMap() *MultiMap {
	return &MultiMap{M: make(map[interface{}][]interface{})}
}

// Get 根据k或者v,ok
func (m *MultiMap) Get(k interface{}) (v []interface{}, ok bool) {
	v, ok = m.M[k]
	return
}

// Put 加入k,v
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

// Delete 删除k
func (m *MultiMap) Delete(k interface{}) {
	delete(m.M, k)
}
