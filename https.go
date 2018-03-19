package go_utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
)

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := GzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func DumpRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))
		fn(w, r)
	}
}

func DetectContentType(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); t == "" {
		t = "application/octet-stream"
	}
	return
}

func ServeImage(imageBytes []byte, fi os.FileInfo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer := bytes.NewReader(imageBytes)
		w.Header().Set("Content-Type", DetectContentType(fi.Name()))
		w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, buffer)
	}
}

func ReadObjectString(object io.ReadCloser) string {
	defer object.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	return buf.String()
}

func ReadObjectBytes(object io.ReadCloser) []byte {
	defer object.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	return buf.Bytes()
}

func HttpPost(url string, requestBody interface{}) ([]byte, error) {
	b, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("json err:", err)
		return nil, err
	}

	body := bytes.NewBuffer([]byte(b))
	log.Println("url:", url)
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		return nil, err
	}

	respBody := ReadObjectBytes(resp.Body)
	return respBody, nil
}

func HttpGet(url string) ([]byte, error) {
	log.Println("url:", url)
	resp, err := http.Get(url)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		return nil, err
	}

	respBody := ReadObjectBytes(resp.Body)
	return respBody, nil
}
