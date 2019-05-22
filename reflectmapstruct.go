package gou

import (
	"fmt"
	"reflect"
	"time"

	"github.com/araddon/dateparse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// Map2Struct map m to struct by its type
func Map2Struct(m map[string]interface{}, result interface{}) error {
	structTypePtr := reflect.TypeOf(result)
	if structTypePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("non struct ptr %v", structTypePtr)
	}
	v := reflect.ValueOf(result).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("non struct ptr %v", structTypePtr)
	}

	structFields := CachedStructFields(v.Type(), "iql")
	for i, sf := range structFields {
		fillField(m, sf, v.Field(i))
	}

	return nil
}

func fillField(m map[string]interface{}, sf StructField, f reflect.Value) {
	for _, tag := range sf.Tag {
		if v, ok := m[tag]; ok {
			setFieldValue(sf, f, v)
			return
		}
	}

	name := SnakeCase(sf.Name)
	if v, ok := m[name]; ok {
		setFieldValue(sf, f, v)
		return
	}

	if sf.Kind == reflect.Struct && sf.Anonymous {
		fv := reflect.New(f.Type()).Interface()
		err := Map2Struct(m, fv)
		LogErr(err)

		v := reflect.ValueOf(fv).Elem()
		f.Set(v)
	}
}

var timeType = reflect.TypeOf(time.Time{})

func setFieldValue(sf StructField, f reflect.Value, v interface{}) {
	vt := reflect.TypeOf(v)
	if vt == sf.Type {
		f.Set(reflect.ValueOf(v))
		return
	}

	if timeType == f.Type() {
		t, _ := dateparse.ParseAny(cast.ToString(v))
		f.Set(reflect.ValueOf(t))
		return
	}

	switch fk := sf.Kind; fk {
	case reflect.Float32, reflect.Float64:
		f.SetFloat(cast.ToFloat64(v))
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		f.SetUint(cast.ToUint64(v))
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		f.SetInt(cast.ToInt64(v))
	case reflect.String:
		f.SetString(cast.ToString(v))
	default:
		logrus.Warnf("unable to convert from %+v to %v", v, fk)
	}
}
