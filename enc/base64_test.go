package enc

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	const s = "不忘初心牢记使命!"
	a := base64.StdEncoding.EncodeToString([]byte(s))
	b, _ := Base64Decode(a)
	assert.Equal(t, s, b)

	a1 := base64.RawStdEncoding.EncodeToString([]byte(s))
	b1, _ := Base64Decode(a1)
	assert.Equal(t, s, b1)

	a2 := base64.URLEncoding.EncodeToString([]byte(s))
	b2, _ := Base64Decode(a2)
	assert.Equal(t, s, b2)

	b3, _ := Base64Decode(Base64(s))
	assert.Equal(t, s, b3)
}
