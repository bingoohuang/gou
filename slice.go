package gou

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

func SliceContains(arr interface{}, elem interface{}) bool {
	arrValue := reflect.ValueOf(arr)
	arrType := arrValue.Type()
	kind := arrType.Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < arrValue.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrValue.Index(i).Interface() == elem {
				return true
			}
		}
		return false
	}

	panic(fmt.Sprintf("Type %s is not supported by Map", arrType.String()))
	return false
}

func RandomIterateSlice(arr interface{}, fn interface{}) (bool, interface{}) {
	if !funk.IsFunction(fn) {
		panic("Second argument must be function")
	}

	arrValue := reflect.ValueOf(arr)
	arrType := arrValue.Type()
	kind := arrType.Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		start := RandomIntN(uint64(arrValue.Len()))
		return iterateSlice(arrValue, start, reflect.ValueOf(fn))
	}

	panic(fmt.Sprintf("Type %s is not supported by Map", arrType.String()))
}

func IterateSlice(arr interface{}, start int, fn interface{}) (bool, interface{}) {
	if !funk.IsFunction(fn) {
		panic("Second argument must be function")
	}

	arrValue := reflect.ValueOf(arr)
	arrType := arrValue.Type()
	kind := arrType.Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		return iterateSlice(arrValue, start, reflect.ValueOf(fn))
	}

	panic(fmt.Sprintf("Type %s is not supported by Map", arrType.String()))
}

var ErrorInterface = reflect.TypeOf((*error)(nil)).Elem()

func iterateSlice(arrValue reflect.Value, start int, funcValue reflect.Value) (bool, interface{}) {
	funcType := funcValue.Type()
	numOut := funcType.NumOut()
	numIn := funcType.NumIn()
	if !(numIn == 1 || numIn == 2) || numOut > 2 {
		panic("Iterate function with an array must have 1/2 parameter and must return 0/1(bool)/2(bool,error) parameter")
	}

	if numOut >= 1 && funcType.Out(0).Kind() != reflect.Bool {
		panic("Iterate function must return bool when there is 1 parameters")
	}
	if numOut >= 2 && funcType.Out(1).Kind() != reflect.Interface {
		panic("Iterate function must return (bool, error) when there is 2 parameters")
	}

	arrElemType := arrValue.Type().Elem()

	// Checking whether element type is convertible to function's first argument's type.
	elemPos := 0
	if numIn == 2 {
		elemPos = 1
	}
	if !arrElemType.ConvertibleTo(funcType.In(elemPos)) {
		panic("Iterate function's argument is not compatible with type of array.")
	}

	if numIn == 2 && reflect.Int != funcType.In(0).Kind() {
		panic("Iterate function's 1st argument is not int.")
	}

	if numOut == 0 {
		internalIterateSlice0(start, arrValue.Len(), arrValue, numIn, funcValue)
		internalIterateSlice0(0, start, arrValue, numIn, funcValue)
		return false, nil
	}

	if over, inte := internalIterateSlice1(start, arrValue.Len(), arrValue, numIn, numOut, funcValue); over {
		return true, inte
	}
	return internalIterateSlice1(0, start, arrValue, numIn, numOut, funcValue)
}

func internalIterateSlice1(from, to int, arrValue reflect.Value, numIn, numOut int, funcValue reflect.Value) (bool, interface{}) {
	for i := from; i < to; i++ {
		var values []reflect.Value
		if numIn == 1 {
			values = []reflect.Value{arrValue.Index(i)}
		} else if numIn == 2 {
			values = []reflect.Value{reflect.ValueOf(i), arrValue.Index(i)}
		}

		if results := funcValue.Call(values); results[0].Bool() {
			if numOut >= 2 {
				return true, results[1].Interface()
			}
			return true, nil
		}
	}

	return false, nil
}

func internalIterateSlice0(from, to int, arrValue reflect.Value, numIn int, funcValue reflect.Value) {
	for i := from; i < to; i++ {
		var values []reflect.Value
		if numIn == 1 {
			values = []reflect.Value{arrValue.Index(i)}
		} else if numIn == 2 {
			values = []reflect.Value{reflect.ValueOf(i), arrValue.Index(i)}
		}
		_ = funcValue.Call(values)
	}
}
