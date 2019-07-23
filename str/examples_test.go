package str_test

import (
	"fmt"
	"testing"

	"github.com/bingoohuang/gou/str"
)

func TestExamples(t *testing.T) {
	fmt.Println(str.SingleLine("hello\nworld"))                             // hello world
	fmt.Println(str.FirstWord("hello world"))                               // hello
	fmt.Println(str.ParseMapString("k1=v1;k2=v2", ";", "="))                // map[k1:v1 k2:v2]
	fmt.Println(str.MapOf("k1", "v1", "k2", "v2"))                          // map[k1:v1 k2:v2]
	fmt.Println(str.MapToString(map[string]string{"k1": "v1", "k2": "v2"})) // map[k1:v1 k2:v2]
	fmt.Println(str.IndexOf("k1", "k0", "k1"))                              // 1
	fmt.Println(str.SplitTrim("k1,,k2", ","))                               // [k1 k2]
	fmt.Println(str.EmptyThen("", "default"))                               // default
	fmt.Println(str.ContainsIgnoreCase("ÑOÑO", "ñoño"))                     // true
	fmt.Println(str.HasPrefix("http://www.abc.com", "http://", "https://")) // true

	a := ".tar.gz"
	fmt.Println(str.AnyOf(a, ".tar", ".tar.gz")) // true
	fmt.Println(str.NoneOf(a, ".xls", ".xlsx"))  // true

	fmt.Println(str.ParseFloat32("1.1")) // 1.1
	fmt.Println(str.ParseFloat64("1.1")) // 1.1
	fmt.Println(str.ParseInt8("-11"))    // -11
	fmt.Println(str.ParseInt16("11"))    // 11
	fmt.Println(str.ParseInt32("11"))    // 11
	fmt.Println(str.ParseInt64("11"))    // 11

	fmt.Println(str.ParseUint("11"))   // 11
	fmt.Println(str.ParseUint8("11"))  // 11
	fmt.Println(str.ParseUint16("11")) // 11
	fmt.Println(str.ParseUint32("11")) // 11
	fmt.Println(str.ParseUint64("11")) // 11

	fmt.Println(str.ParseInt("11"))   // 11
	fmt.Println(str.ParseInt8("11"))  // 11
	fmt.Println(str.ParseInt16("11")) // 11
	fmt.Println(str.ParseInt32("11")) // 11
	fmt.Println(str.ParseInt64("11")) // 11
}
