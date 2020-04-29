package str

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseFloat32E ...
func ParseFloat32E(s string) (float32, error) {
	f, e := strconv.ParseFloat(s, 32)
	return float32(f), e
}

// ParseFloat64E ...
func ParseFloat64E(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

// ParseIntE ...
func ParseIntE(s string) (int, error) { p, e := ParseInt64E(s); return int(p), e }

// ParseInt8E ...
func ParseInt8E(s string) (int8, error) { p, e := strconv.ParseInt(s, 0, 8); return int8(p), e }

// ParseInt16E ...
func ParseInt16E(s string) (int16, error) { p, e := strconv.ParseInt(s, 0, 16); return int16(p), e }

// ParseInt32E ...
func ParseInt32E(s string) (int32, error) { p, e := strconv.ParseInt(s, 0, 32); return int32(p), e }

// ParseInt64E ...
func ParseInt64E(s string) (int64, error) { return strconv.ParseInt(s, 0, 64) }

// ParseUintE ...
func ParseUintE(s string) (uint, error) { p, e := ParseUint64E(s); return uint(p), e }

// ParseUint8E ...
func ParseUint8E(s string) (uint8, error) { p, e := strconv.ParseUint(s, 0, 8); return uint8(p), e }

// ParseUint16E ...
func ParseUint16E(s string) (uint16, error) { p, e := strconv.ParseUint(s, 0, 16); return uint16(p), e }

// ParseUint32E ...
func ParseUint32E(s string) (uint32, error) { p, e := strconv.ParseUint(s, 0, 32); return uint32(p), e }

// ParseUint64E ...
func ParseUint64E(s string) (uint64, error) { return strconv.ParseUint(s, 0, 64) }

// ParseInt8 ...
func ParseInt8(s string) int8 { i, _ := ParseInt8E(s); return i }

// ParseInt16 ...
func ParseInt16(s string) int16 { i, _ := ParseInt16E(s); return i }

// ParseInt32 ...
func ParseInt32(s string) int32 { i, _ := ParseInt32E(s); return i }

// ParseInt64 ...
func ParseInt64(s string) int64 { i, _ := ParseInt64E(s); return i }

// ParseUint8 ...
func ParseUint8(s string) uint8 { i, _ := ParseUint8E(s); return i }

// ParseUint16 ...
func ParseUint16(s string) uint16 { i, _ := ParseUint16E(s); return i }

// ParseUint32 ...
func ParseUint32(s string) uint32 { i, _ := ParseUint32E(s); return i }

// ParseUint64 ...
func ParseUint64(s string) uint64 { i, _ := ParseUint64E(s); return i }

// ParseFloat32 ...
func ParseFloat32(s string) float32 { f, _ := ParseFloat32E(s); return f }

// ParseFloat64 ...
func ParseFloat64(s string) float64 { f, _ := ParseFloat64E(s); return f }

// ParseInt ...
func ParseInt(s string) int { f, _ := ParseIntE(s); return f }

// ParseUint ...
func ParseUint(s string) uint { f, _ := ParseUintE(s); return f }

// ParseBool returns the boolean value represented by the string.
// It accepts 1, t, true, y, yes, on with camel case incentive.
func ParseBool(s string) bool {
	b, _ := ParseBoolE(s)
	return b
}

// ParseBoolE returns the boolean value represented by the string.
// It accepts 1, t, true, y, yes, on as true with camel case incentive
// and accepts 0, f false, n, no, off as false with camel case incentive
// Any other value returns an error.
func ParseBoolE(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "1", "t", "true", "y", "yes", "on":
		return true, nil
	case "0", "f", "false", "n", "no", "off":
		return false, nil
	}

	return false, &strconv.NumError{Func: "ParseBoolE", Num: s, Err: strconv.ErrSyntax}
}

// ParseDuration ...
func ParseDuration(s string) time.Duration { f, _ := ParseDurationE(s); return f }

// ParseDurationE ...
func ParseDurationE(s string) (time.Duration, error) {
	return time.ParseDuration(StripSpaces(s))
}

// Strip strips runes that predicates returns true.
func Strip(str string, predicates ...func(rune) bool) string {
	var b strings.Builder

	b.Grow(len(str))

	for _, ch := range str {
		if !matchesAny(predicates, ch) {
			b.WriteRune(ch)
		}
	}

	return b.String()
}

func matchesAny(predicates []func(rune) bool, ch rune) bool {
	for _, p := range predicates {
		if p(ch) {
			return true
		}
	}

	return false
}

// StripSpaces strips all spaces from the string.
func StripSpaces(str string) string {
	return Strip(str, unicode.IsSpace)
}

// HasSpaces test if any spaces in the string.
func HasSpaces(str string) bool {
	return Has(str, unicode.IsSpace)
}

// Has test if any rune which predicates tells it it true.
func Has(str string, predicates ...func(rune) bool) bool {
	for _, ch := range str {
		if matchesAny(predicates, ch) {
			return true
		}
	}

	return false
}

func Not(p func(rune) bool) func(rune) bool {
	return func(r rune) bool { return !p(r) }
}
