// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package urllib is a httplib for golang
package gou

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var defaultUrlSetting = UrlHttpSettings{
	ShowDebug:        false,
	UserAgent:        "GiterLab",
	ConnectTimeout:   10 * time.Second,
	ReadWriteTimeout: 10 * time.Second,
	TlsClientConfig:  nil,
	Proxy:            nil,
	Transport:        nil,
	EnableCookie:     false,
	Gzip:             true,
	DumpBody:         true,
}
var defaultCookieJar http.CookieJar
var settingMutex sync.Mutex

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// Overwrite default settings
func SetDefaultSetting(setting UrlHttpSettings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultUrlSetting = setting
	if defaultUrlSetting.ConnectTimeout == 0 {
		defaultUrlSetting.ConnectTimeout = 10 * time.Second
	}
	if defaultUrlSetting.ReadWriteTimeout == 0 {
		defaultUrlSetting.ReadWriteTimeout = 10 * time.Second
	}
}

// UrlGet current default settings
func GetDefaultSetting() *UrlHttpSettings {
	settingMutex.Lock()
	defer settingMutex.Unlock()

	return &defaultUrlSetting
}

// return *UrlHttpRequest with specific method
func newUrlRequest(rawURL, method string) (*UrlHttpRequest, error) {
	var resp http.Response
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	req := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	return &UrlHttpRequest{
		url:     rawURL,
		req:     &req,
		params:  map[string]string{},
		files:   map[string]string{},
		setting: defaultUrlSetting,
		resp:    &resp,
		body:    nil,
	}, nil
}

// UrlGet returns *UrlHttpRequest with GET method.
func UrlGet(url string) (*UrlHttpRequest, error) {
	return newUrlRequest(url, "GET")
}

// MustUrlGet  returns *UrlHttpRequest with GET method.
func MustUrlGet(url string) *UrlHttpRequest {
	req, err := UrlGet(url)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

// UrlPost returns *UrlHttpRequest with POST method.
func UrlPost(url string) (*UrlHttpRequest, error) {
	return newUrlRequest(url, "POST")
}

// MustUrlPost returns *UrlHttpRequest with POST method.
func MustUrlPost(url string) *UrlHttpRequest {
	req, err := UrlPost(url)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

// UrlPut returns *UrlHttpRequest with PUT method.
func UrlPut(url string) (*UrlHttpRequest, error) {
	return newUrlRequest(url, "PUT")
}

// MustUrlPut returns *UrlHttpRequest with PUT method.
func MustUrlPut(url string) *UrlHttpRequest {
	req, err := UrlPut(url)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

// UrlDelete returns *UrlHttpRequest with DELETE method.
func UrlDelete(url string) (*UrlHttpRequest, error) {
	return newUrlRequest(url, "DELETE")
}

// MustDelete returns *UrlHttpRequest with DELETE method.
func MustDelete(url string) *UrlHttpRequest {
	req, err := UrlDelete(url)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

// UrlHead returns *UrlHttpRequest with HEAD method.
func UrlHead(url string) (*UrlHttpRequest, error) {
	return newUrlRequest(url, "HEAD")
}

// MustUrlHead returns *UrlHttpRequest with UrlHead method.
func MustUrlHead(url string) *UrlHttpRequest {
	req, err := UrlHead(url)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

// UrlHead returns *UrlHttpRequest with PATCH method.
func UrlPatch(url string) (*UrlHttpRequest, error) {
	return newUrlRequest(url, "PATCH")
}

// MustUrlPatch returns *UrlHttpRequest with UrlPatch method.
func MustUrlPatch(url string) *UrlHttpRequest {
	req, err := UrlPatch(url)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

// UrlHttpSettings
type UrlHttpSettings struct {
	ShowDebug        bool
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TlsClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	EnableCookie     bool
	Gzip             bool
	DumpBody         bool
}

// UrlHttpRequest provides more useful methods for requesting one url than http.Request.
type UrlHttpRequest struct {
	url     string
	req     *http.Request
	params  map[string]string
	files   map[string]string
	setting UrlHttpSettings
	resp    *http.Response
	body    []byte
	dump    []byte
}

// Change request settings
func (b *UrlHttpRequest) Setting(setting UrlHttpSettings) *UrlHttpRequest {
	b.setting = setting
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *UrlHttpRequest) SetBasicAuth(username, password string) *UrlHttpRequest {
	b.req.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *UrlHttpRequest) SetEnableCookie(enable bool) *UrlHttpRequest {
	b.setting.EnableCookie = enable
	return b
}

// SetUserAgent sets User-Agent header field
func (b *UrlHttpRequest) SetUserAgent(useragent string) *UrlHttpRequest {
	b.setting.UserAgent = useragent
	return b
}

// Debug sets show debug or not when executing request.
func (b *UrlHttpRequest) Debug(isdebug bool) *UrlHttpRequest {
	b.setting.ShowDebug = isdebug
	return b
}

// Dump Body.
func (b *UrlHttpRequest) DumpBody(isdump bool) *UrlHttpRequest {
	b.setting.DumpBody = isdump
	return b
}

// return the DumpRequest
func (b *UrlHttpRequest) DumpRequest() []byte {
	return b.dump
}

// return the DumpRequest string
func (b *UrlHttpRequest) DumpRequestString() string {
	return string(b.DumpRequest())
}

// SetTimeout sets connect time out and read-write time out for Request.
func (b *UrlHttpRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *UrlHttpRequest {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *UrlHttpRequest) SetTLSClientConfig(config *tls.Config) *UrlHttpRequest {
	b.setting.TlsClientConfig = config
	return b
}

// Header add header item string in request.
func (b *UrlHttpRequest) Header(key, value string) *UrlHttpRequest {
	b.req.Header.Set(key, value)
	return b
}

// Set HOST
func (b *UrlHttpRequest) SetHost(host string) *UrlHttpRequest {
	b.req.Host = host
	return b
}

// Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (b *UrlHttpRequest) SetProtocolVersion(vers string) *UrlHttpRequest {
	if len(vers) == 0 {
		vers = "HTTP/1.1"
	}

	major, minor, ok := http.ParseHTTPVersion(vers)
	if ok {
		b.req.Proto = vers
		b.req.ProtoMajor = major
		b.req.ProtoMinor = minor
	}

	return b
}

// SetCookie add cookie into request.
func (b *UrlHttpRequest) SetCookie(cookie *http.Cookie) *UrlHttpRequest {
	b.req.Header.Add("Cookie", cookie.String())
	return b
}

// UrlGet default CookieJar
func GetDefaultCookieJar() http.CookieJar {
	return defaultCookieJar
}

// Set transport to
func (b *UrlHttpRequest) SetTransport(transport http.RoundTripper) *UrlHttpRequest {
	b.setting.Transport = transport
	return b
}

// Set http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (b *UrlHttpRequest) SetProxy(proxy func(*http.Request) (*url.URL, error)) *UrlHttpRequest {
	b.setting.Proxy = proxy
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *UrlHttpRequest) Param(key, value string) *UrlHttpRequest {
	b.params[key] = value
	return b
}

func (b *UrlHttpRequest) PostFile(formname, filename string) *UrlHttpRequest {
	b.files[formname] = filename
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *UrlHttpRequest) Body(data interface{}) *UrlHttpRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	}
	return b
}

// JsonBody adds request raw body encoding by JSON.
func (b *UrlHttpRequest) JsonBody(obj interface{}) error {
	if b.req.Body != nil || obj == nil {
		return errors.New("body should not be nil")
	}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return err
	}
	b.req.Body = ioutil.NopCloser(buf)
	b.req.ContentLength = int64(buf.Len())
	b.req.Header.Set("Content-Type", "application/json;charset=utf-8")
	return nil
}

func (b *UrlHttpRequest) buildUrl(paramBody string) {
	// build GET url with query string
	m := b.req.Method
	if m == "GET" && len(paramBody) > 0 {
		if strings.Index(b.url, "?") != -1 {
			b.url += "&" + paramBody
		} else {
			b.url += "?" + paramBody
		}
		return
	}

	// build POST/PUT/PATCH url and body
	if (m == "POST" || m == "PUT" || m == "PATCH") && b.req.Body == nil {
		// with files
		if len(b.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for formname, filename := range b.files {
					fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
					if err != nil {
						log.Fatal(err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Fatal(err)
					}
					//iocopy
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Fatal(err)
					}
				}
				for k, v := range b.params {
					bodyWriter.WriteField(k, v)
				}
				bodyWriter.Close()
				pw.Close()
			}()
			b.Header("Content-Type", bodyWriter.FormDataContentType())
			b.req.Body = ioutil.NopCloser(pr)
			return
		}

		// with params
		if len(paramBody) > 0 {
			b.Header("Content-Type", "application/x-www-form-urlencoded")
			b.Body(paramBody)
		}
	}
}

func (b *UrlHttpRequest) getResponse() (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
	}
	resp, err := b.SendOut()
	if err != nil {
		return nil, err
	}
	b.resp = resp
	return resp, nil
}

func (b *UrlHttpRequest) SendOut() (*http.Response, error) {
	var paramBody string
	if len(b.params) > 0 {
		var buf bytes.Buffer
		for k, v := range b.params {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}

	b.buildUrl(paramBody)
	u, err := url.Parse(b.url)
	if err != nil {
		return nil, err
	}

	b.req.URL = u

	trans := b.setting.Transport

	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig: b.setting.TlsClientConfig,
			Proxy:           b.setting.Proxy,
			DialContext:     UrlTimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout),
		}
	} else {
		// if b.transport is *http.Transport then set the settings.
		if t, ok := trans.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = b.setting.TlsClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = b.setting.Proxy
			}
			if t.DialContext == nil {
				t.DialContext = UrlTimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout)
			}
		}
	}

	var jar http.CookieJar = nil
	if b.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}
		jar = defaultCookieJar
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.req.Header.Get("User-Agent") == "" {
		b.req.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.ShowDebug {
		dump, err := httputil.DumpRequest(b.req, b.setting.DumpBody)
		if err != nil {
			log.Println(err.Error())
		}
		b.dump = dump
	}
	return client.Do(b.req)
}

// String returns the body string in response.
// it calls Response inner.
func (b *UrlHttpRequest) String() (string, error) {
	data, err := b.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *UrlHttpRequest) Bytes() ([]byte, error) {
	if b.body != nil {
		return b.body, nil
	}
	resp, err := b.getResponse()
	if err != nil || resp.Body == nil {
		return nil, err
	}

	return b.ReadResponseBody(resp)
}

func (b *UrlHttpRequest) ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	if b.setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(reader)
	}

	return ioutil.ReadAll(resp.Body)
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *UrlHttpRequest) ToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := b.getResponse()
	if err != nil || resp.Body == nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// ToJson returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *UrlHttpRequest) ToJson(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ToXml returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *UrlHttpRequest) ToXml(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// Response executes request client gets response mannually.
func (b *UrlHttpRequest) Response() (*http.Response, error) {
	return b.getResponse()
}

// UrlTimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func UrlTimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(ctx context.Context, net, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, err
	}
}
