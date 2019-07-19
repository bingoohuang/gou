package enc

import (
	"bytes"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// JSON 序列化v到JSON串
func JSON(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		logrus.Warnf("JSON error %v for value %+v", err, v)
		return ""
	}

	return string(s)
}

// JSONPretty 以格式化后的形式输出JSON
func JSONPretty(v interface{}) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	_ = enc.Encode(v)
	return buf.String()
}
