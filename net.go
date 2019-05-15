package gou

import (
	"net"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// IsIP 判断 host 字符串表达式是不是IP(v4/v6)的格式
func IsIP(host string) bool {
	ip := net.ParseIP(host)
	return ip != nil
}

// RspBase 代表返回的公共结构
type RspBase struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// RspBaseData 代表返回的公共结构
type RspBaseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrHandled 处理 err 错误，并且返回给 c
func ErrHandled(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}

	c.JSON(200, RspBase{Status: 400, Message: err.Error()})
	return true
}

// BuildURL 创建一个url
func BuildURL(base string, querys map[string]string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	q := u.Query()

	for k, v := range querys {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// RestGet 发起一次HTTP GET调用，并且反序列化JSON到v代表的指针中。
func RestGet(url string, v interface{}) error {
	logrus.Debugf("RestGet %v", url)
	req, err := UrlGet(url)
	if err != nil {
		logrus.Warnf(" UrlUrlGet failed %v", err)
		return err
	}

	if err := req.ToJson(v); err != nil {
		logrus.Warnf("Unmarshal failed %v", err)
		return err
	}
	logrus.Debugf("RestGet  Result %+v", v)

	return nil
}

// HTTPGet 表示一次HTTP的Get调用
func HTTPGet(url string) ([]byte, error) {
	logrus.Debugf("http get %s", url)

	resp, err := UrlGet(url)
	if err != nil {
		return nil, err
	}

	return resp.Bytes()
}

// RestPost 表示一次HTTP的POST调用
func RestPost(url string, requestBody interface{}, responseStruct interface{}) ([]byte, error) {
	logrus.Debugf("post body %#v", requestBody)

	resp, err := UrlPost(url)
	if err != nil {
		return nil, err
	}

	if err = resp.JsonBody(requestBody); err != nil {
		return nil, err
	}

	if responseStruct != nil {
		return nil, resp.ToJson(responseStruct)
	}

	return resp.Bytes()
}
