package lang

import (
	"github.com/bingoohuang/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

// nolint gochecknoinits
func init() {
	extra.SetNamingStrategy(strcase.ToCamelLower)
}

// nolint gochecknoglobals
var (
	JSONUnmarshal     = jsoniter.Unmarshal
	JSONMarshal       = jsoniter.Marshal
	JSONMarshalIndent = jsoniter.MarshalIndent
)
