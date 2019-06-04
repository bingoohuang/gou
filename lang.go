package gou

import (
	"bytes"
	"encoding/json"
	"regexp"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// Recover 在系统崩溃是，恢复系统
func Recover() {
	if err := recover(); err != nil {
		logrus.Warnln(err)
		debug.PrintStack()
		recover()
	}
}

// If 代替三元表达式 c ? a : b
func If(c bool, a, b string) string {
	if c {
		return a
	}
	return b
}

func ContainsWord(s, sep, word string) bool {
	parts := SplitN(s, sep, true, true)
	for _, p := range parts {
		if p == word {
			return true
		}
	}

	return false
}

// AnyOf 给定x是否在any中
func AnyOf(x interface{}, any ...interface{}) bool {
	for _, a := range any {
		if x == a {
			return true
		}
	}

	return false
}

// EmptyTo 在s为空时，返回def，否则返回s
func EmptyTo(s, def string) string {
	if s == "" {
		return def
	}

	return s
}

// JoinNonEmpty 组合x
func JoinNonEmpty(sep string, x ...string) string {
	s := ""
	for _, i := range x {
		if i != "" {
			s += sep + i
		}
	}

	if s != "" {
		s = s[len(sep):]
	}

	return s
}

// Join 组合x
func Join(x ...string) string {
	s := ""
	for _, i := range x {
		s += i
	}
	return s
}

// PrependIf 在cond为真时，在s前添加前缀
func PrependIf(cond bool, s string, subs ...string) string {
	if cond {
		return Join(subs...) + s
	}
	return s
}

// JoinIf 在cond为真时，组合后续的子串
func JoinIf(cond bool, subs ...string) string {
	if cond {
		return Join(subs...)
	}

	return ""
}

// Decode acts like oracle decode.
func Decode(target interface{}, decodeVars ...interface{}) interface{} {
	length := len(decodeVars)
	i := 0
	for ; i+1 < length; i += 2 {
		if target == decodeVars[i] {
			return decodeVars[i+1]
		}
	}

	if i < length {
		return decodeVars[i]
	}

	return nil
}

// LogErr logs err if it is not nil
func LogErr(err error) {
	if err != nil {
		logrus.Warnf("error %v", err)
	}
}

var reCrLn = regexp.MustCompile(`\r?\n`)
var reBlanks = regexp.MustCompile(`\s+`)

// SingleLine 替换s中的换行以及连续空白，为单个空白
func SingleLine(s string) string {
	return reBlanks.ReplaceAllString(reCrLn.ReplaceAllString(s, " "), " ")
}

// JSON 序列化v到JSON串
func JSON(v interface{}) []byte {
	s, err := json.Marshal(v)
	if err != nil {
		logrus.Warnf("JSON error %v for value %+v", err, v)
		return nil
	}

	return s
}

// PrettyJsontify 以格式化后的形式输出JSON
func PrettyJsontify(v interface{}) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	enc.Encode(v)
	return buf.String()
}
