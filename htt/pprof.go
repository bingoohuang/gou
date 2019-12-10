package htt

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// PprofAddrPflag declares pprof plags.
func PprofAddrPflag() *string {
	return pflag.StringP("pprof-addr", "", "",
		"pprof address to listen on, not activate pprof if empty, eg. --pprof-addr localhost:6060")
}

// StartPprof starts pprof.
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
