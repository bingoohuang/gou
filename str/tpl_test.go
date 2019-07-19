package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTpl(t *testing.T) {
	assert.Equal(t, "bingoohuang", Tpl("{name}", map[string]interface{}{"name": "bingoohuang"}))
}
