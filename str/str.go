package str

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/bingoohuang/goreflect"
)

// HasPrefix tells s has any prefix of p...
func HasPrefix(s string, p ...string) bool {
	for _, i := range p {
		if strings.HasPrefix(s, i) {
			return true
		}
	}

	return false
}

var reCrLn = regexp.MustCompile(`\r?\n`) // nolint
var reBlanks = regexp.MustCompile(`\s+`) // nolint

// SingleLine 替换s中的换行以及连续空白，为单个空白
func SingleLine(s string) string {
	return reBlanks.ReplaceAllString(reCrLn.ReplaceAllString(s, " "), " ")
}

// FirstWord find the first word of value.
func FirstWord(value string) string {
	started := -1

	// Loop over all indexes in the string.
	for i, c := range value {
		// If we encounter a space, reduce the count.
		if started == -1 && !unicode.IsSpace(c) {
			started = i
		}

		if started >= 0 && unicode.IsSpace(c) {
			return value[started:i]
		}
	}

	return value[started:]
}

// ParseMapString parses str to a map.
func ParseMapString(str string, separator, keyValueSeparator string) map[string]string {
	parts := strings.Split(str, separator)

	m := make(map[string]string)

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}

		index := strings.Index(p, keyValueSeparator)

		if index > 0 {
			key := p[0:index]
			val := p[index+1:]
			k := strings.TrimSpace(key)
			v := strings.TrimSpace(val)

			if k != "" {
				m[k] = v
			}
		} else if index < 0 {
			m[p] = ""
		}
	}

	return m
}

// IndexOf find the index of word in data slice.
func IndexOf(word string, data ...string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}

	return -1
}

// MapOf composites arr to a map.
func MapOf(arr ...string) map[string]string {
	result := make(map[string]string)
	i := 0

	for ; i+1 < len(arr); i += 2 {
		result[arr[i]] = arr[i+1]
	}

	if i < len(arr) {
		result[arr[i]] = ""
	}

	return result
}

// MapToString joins map to a string.
func MapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	_, _ = fmt.Fprintf(b, "%v", m)

	return b.String()
}

// SplitTrim splits and trims.
func SplitTrim(str, sep string) []string {
	subs := strings.Split(str, sep)
	ret := make([]string, 0)

	for i, v := range subs {
		v := strings.TrimSpace(v)

		if len(subs[i]) > 0 {
			ret = append(ret, v)
		}
	}

	return ret
}

// EmptyThen ...
func EmptyThen(s, then string) string {
	return If(s != "", s, then)
}

// ContainsIgnoreCase ...
func ContainsIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

// StringContains detects an item contained in container in whole word mode separated by sep
// but if the absolute is specified and the container ==  absolute, true should always returned
func StringContains(container, item, sep, absolute string) bool {
	if absolute != "" && absolute == container {
		return true
	}

	items := strings.Split(container, sep)

	return goreflect.SliceContains(items, item)
}

// If 代替三元表达式 c ? a : b
func If(c bool, a, b string) string {
	if c {
		return a
	}

	return b
}

// ContainsWord contains word or not.
func ContainsWord(s, sep, word string) bool {
	parts := SplitN(s, sep, true, true)

	for _, p := range parts {
		if p == word {
			return true
		}
	}

	return false
}

// NoneOf 给定x是否不在所有的any中
func NoneOf(x interface{}, any ...interface{}) bool {
	return !AnyOf(x, any...)
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

// Repeat repeat s with times with separator.
func Repeat(s, sep string, times int) string {
	str := ""

	for i := 0; i < times; i++ {
		if i > 0 {
			str += sep
		}

		str += s
	}

	return str
}
