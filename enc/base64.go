package enc

import (
	"encoding/base64"
	"strings"
)

// Base64 安全编码, URL兼容，并且去除后面多余的等号。
func Base64(source string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(source))
}

// UnBase64 安全解码，兼容标准和URL编码，以及后面去除等号的情况。
func UnBase64(src string) (string, error) {
	s := src
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-')
	s = strings.Replace(s, "_", "/", -1)
	s = strings.Replace(s, "-", "+", -1)
	s = strings.TrimRight(s, "=")

	decoded, err := base64.RawStdEncoding.DecodeString(s)
	return string(decoded), err
}
