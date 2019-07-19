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
}
