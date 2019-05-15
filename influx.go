package gou

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Query execute influxQl (refer to https://docs.influxdata.com/influxdb/v1.7/query_language)
// influxDBAddr  InfluxDB的连接地址， 例如http://localhost:8086, 注意：1. 右边没有/ 2. 右边不带其它path，例如/query等。
func InfluxQuery(influxDBAddr, influxQl string) (string, error) {
	req, err := UrlGet(influxDBAddr + `/query`)
	if err != nil {
		return "", err
	}

	req.Param("q", influxQl)

	return req.String()
}

// InfluxWrite 写入打点值
// refer https://github.com/DCSO/fluxline/blob/master/encoder.go
func InfluxWrite(influxDBWriteAddr, line string) error {
	req, err := UrlPost(influxDBWriteAddr)
	if err != nil {
		return err
	}
	req.Body([]byte(line))

	rsp, err := req.SendOut()
	if err != nil {
		return err
	}

	rspBody, err := req.ReadResponseBody(rsp)
	logrus.Debugf("influx write %s returned status %s msg %s", line, rsp.Status, string(rspBody))

	return nil
}

// LineProtocol format inputs to line protocol
// https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_tutorial/
func LineProtocol(name string, tags map[string]string, fields map[string]interface{}, t time.Time) (string, error) {
	if len(fields) == 0 {
		return "", errors.New("fields are empty")
	}

	tagstr := ""
	IterateMapSorted(tags, func(k, v string) {
		tagstr += fmt.Sprintf(",%s=%s", escapeSpecialChars(k), escapeSpecialChars(v))
	})

	out := ""
	// serialize fields
	var err error
	IterateMapSorted(fields, func(k string, v interface{}) {
		var repr string
		repr, err = toInfluxRepr(v)
		if err != nil {
			return
		}
		out += fmt.Sprintf(",%s=%s", escapeSpecialChars(k), repr)
	})

	if err != nil {
		return "", err
	}

	if out != "" {
		out = out[1:]
	}

	// construct line protocol string
	return fmt.Sprintf("%s%s %s %d", name, tagstr, out, uint64(t.UnixNano())), nil
}

func escapeSpecialChars(in string) string {
	str := strings.Replace(in, ",", `\,`, -1)
	str = strings.Replace(str, "=", `\=`, -1)
	str = strings.Replace(str, " ", `\ `, -1)
	return str
}

// toInfluxRepr 将val转换为Influx表示形式
func toInfluxRepr(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return stringToInfluxRepr(v)
	case []byte:
		return stringToInfluxRepr(string(v))
	case int32, int64, int16, int8, int, uint32, uint64, uint16, uint8, uint:
		return fmt.Sprintf("%d", v), nil
	case float64, float32:
		return fmt.Sprintf("%g", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case time.Time:
		return fmt.Sprintf("%d", uint64(v.UnixNano())), nil
	default:
		return "", fmt.Errorf("%+v: unsupported type for Influx Line Protocol", val)
	}
}

func stringToInfluxRepr(v string) (string, error) {
	if len(v) > 64000 {
		return "", fmt.Errorf("string too long (%d characters, max. 64K)", len(v))
	}
	return fmt.Sprintf("%q", v), nil
}