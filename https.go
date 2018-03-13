package go_utils

import (
	"bytes"
	"compress/gzip"
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
