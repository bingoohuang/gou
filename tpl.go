package gou

import (
	"fmt"
	"strings"
)

// Tpl 模板替换
func Tpl(tpl string, vars map[string]interface{}) string {
	s := tpl

	for k, v := range vars {
		old := fmt.Sprintf("{%s}", k)
		nld := fmt.Sprintf("%v", v)
		s = strings.ReplaceAll(s, old, nld)
	}

	return s
}
