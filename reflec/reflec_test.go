package reflec_test

import (
	"fmt"
	"reflect"
	"testing"
)

type Thing struct {
}

func (tp *Thing) Goodbye(greeting string) {
	fmt.Printf("adios %s\n", greeting)
}

func (tp *Thing) Hello(answer int) string {
	fmt.Printf("Hi There!\n")
	return "rob"
}

func (tp *Thing) Invoke(name string, args ...interface{}) []reflect.Value {
	return invoke(tp, name, args...)
}

func inspect(v []reflect.Value) {
	for i := range v {
		fmt.Printf("v[%d] => %q, ", i, v[i].Kind())
		switch v[i].Kind() {
		case reflect.Int:
			fmt.Printf("%d", v[i].Int())
		case reflect.Interface:
			fmt.Printf("%v", v[i].Interface())
		case reflect.Invalid, reflect.Bool,
			reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8,
			reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr, reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128, reflect.Array,
			reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr,
			reflect.Slice, reflect.String, reflect.Struct,
			reflect.UnsafePointer:
		}
		fmt.Println()
	}
}

func invoke(any interface{}, name string, args ...interface{}) []reflect.Value {
	if len(args) > 0 {
		fmt.Printf("\tthere are %d args\n", len(args))
		inputs := make([]reflect.Value, len(args))
		for i := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		return reflect.ValueOf(any).MethodByName(name).Call(inputs)
	}

	return reflect.ValueOf(any).MethodByName(name).Call([]reflect.Value{})
}

func TestInvoke(tt *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("failed to find method: %s\n", r)
		}
	}()

	var a Thing
	inspect(a.Invoke("Hello", 42))
	a.Invoke("Goodbye", "Rob", "DEAD")
}
