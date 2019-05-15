package gou

import "strings"

// Split1 将s按分隔符sep分成x份，取第1份
func Split1(s, sep string, trimSpace, ignoreEmpty bool) (s0 string) {
	s0, _, _, _ = Split4(s, sep, trimSpace, ignoreEmpty)

	return
}

// Split2 将s按分隔符sep分成x份，取第1、2份
func Split2(s, sep string, trimSpace, ignoreEmpty bool) (s0, s1 string) {
	s0, s1, _, _ = Split4(s, sep, trimSpace, ignoreEmpty)

	return
}

// Split3 将s按分隔符sep分成x份，取第1、2、3份
func Split3(s, sep string, trimSpace, ignoreEmpty bool) (s0, s1, s2 string) {
	s0, s1, s2, _ = Split4(s, sep, trimSpace, ignoreEmpty)
	return
}

// Split4 将s按分隔符sep分成x份，取第x份，取第1、2、3、4份
func Split4(s, sep string, trimSpace, ignoreEmpty bool) (s0, s1, s2, s3 string) {
	s0, s1, s2, s3, _ = Split5(s, sep, trimSpace, ignoreEmpty)
	return
}

// Split5 将s按分隔符sep分成x份，取第x份，取第1、2、3、4、5份
func Split5(s, sep string, trimSpace, ignoreEmpty bool) (s0, s1, s2, s3, s4 string) {
	parts := SplitN(s, sep, trimSpace, ignoreEmpty)
	l := len(parts)

	if l > 0 {
		s0 = strings.TrimSpace(parts[0])
	}
	if l > 1 {
		s1 = strings.TrimSpace(parts[1])
	}
	if l > 2 {
		s2 = strings.TrimSpace(parts[2])
	}
	if l > 3 {
		s3 = strings.TrimSpace(parts[3])
	}
	if l > 4 {
		s4 = strings.TrimSpace(parts[4])
	}
	return s0, s1, s2, s3, s4
}

func SplitN(s, sep string, trimSpace, ignoreEmpty bool) []string {
	parts := strings.SplitN(s, sep, -1)

	result := make([]string, 0)

	for _, p := range parts {
		if trimSpace {
			p = strings.TrimSpace(p)
		}

		if ignoreEmpty && p == "" {
			continue
		}

		result = append(result, p)
	}

	return result
}
