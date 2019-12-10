package enc

import (
	"encoding/base64"
	"strings"
)

// Base64 安全编码, URL兼容，并且去除后面多余的等号。
func Base64(source string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(source))
}

// Base64Decode 安全解码，兼容标准和URL编码，以及后面去除等号的情况。
func Base64Decode(src string) (string, error) {
	s := src
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-')
	s = strings.Replace(s, "_", "/", -1)
	s = strings.Replace(s, "-", "+", -1)
	s = strings.TrimRight(s, "=")

	decoded, err := base64.RawStdEncoding.DecodeString(s)

	return string(decoded), err
}

// Base64SafeEncode 安全URL编码，去除后面多余的等号
func Base64SafeEncode(source []byte) string {
	dest := base64.URLEncoding.EncodeToString(source)

	return strings.TrimRight(dest, "=")
}

// Base64SafeDecode 安全解码，兼容标准和URL编码，以及后面等号是否多余
func Base64SafeDecode(source string) ([]byte, error) {
	src := source
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+'
	// (replaced by '_' and '-') and trailing '=' are removed.
	src = strings.Replace(src, "_", "/", -1)
	src = strings.Replace(src, "-", "+", -1)

	if i := len(src) % 4; i != 0 {
		src += strings.Repeat("=", 4-i)
	}

	return base64.StdEncoding.DecodeString(src)
}
