package go_utils

import (
	"unicode"
	"strings"
)

func FirstWord(value string) string {
	started := -1
	ended := 0
	// Loop over all indexes in the string.
	for i, c := range value {
		// If we encounter a space, reduce the count.
		if started == -1 && !unicode.IsSpace(c) {
			started = i
		}
		if started >= 0 && unicode.IsSpace(c) {
			ended = i
			break
		}
	}
	// Return the entire string.
	return value[started:ended]
}

func IfElse(ifCondition bool, ifValue, elseValue string) string {
	if ifCondition {
		return ifValue
	}
	return elseValue
}

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
