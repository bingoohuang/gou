package go_utils

import "unicode"

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
