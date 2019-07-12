package gou

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/bingoohuang/gonet"
)

func ServeFavicon(path string, mustAsset func(name string) []byte,
	assetInfo func(name string) (os.FileInfo, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fi, _ := assetInfo(path)
		buffer := bytes.NewReader(mustAsset(path))
		w.Header().Set("Content-Type", gonet.DetectContentType(fi.Name()))
		w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, buffer)
	}
}
