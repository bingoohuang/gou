package enc

import (
	"crypto/hmac"
	"crypto/sha1" // nolint
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Sha256 generate a sha256 hash of src in HEX format
func Sha256(src string) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(src)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// HmacSha1 generate HmacSha1 of src in BASE64 format
func HmacSha1(src string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	_, _ = h.Write([]byte(src))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
