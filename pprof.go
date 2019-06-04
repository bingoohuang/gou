package gou

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartPprof(pprofAddr string) {
	if pprofAddr == "" {
		return
	}

	logrus.Infof("Starting pprof HTTP server at: http://%s/debug/pprof", pprofAddr)
	go func() {
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			logrus.Fatalf("error %v", err)
		}
	}()
}
