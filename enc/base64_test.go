package enc

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	const s = "我是黄进兵哈哈哈1234"
	a := base64.StdEncoding.EncodeToString([]byte(s))
	b, _ := UnBase64(a)
	assert.Equal(t, s, b)

	a1 := base64.RawStdEncoding.EncodeToString([]byte(s))
	b1, _ := UnBase64(a1)
	assert.Equal(t, s, b1)

	a2 := base64.URLEncoding.EncodeToString([]byte(s))
	b2, _ := UnBase64(a2)
	assert.Equal(t, s, b2)

	b3, _ := UnBase64(Base64(s))
	assert.Equal(t, s, b3)
}
