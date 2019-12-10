package htt

import (
	"runtime"
	"time"

	"github.com/bingoohuang/gou/ran"

	"github.com/skratchdot/open-golang/open"
)

// OpenExplorerWithContext ...
func OpenExplorerWithContext(contextPath, port string) {
	go func() {
		time.Sleep(100 * time.Millisecond)

		switch runtime.GOOS {
		case "windows":
			fallthrough
		case "darwin":
			_ = open.Run("http://127.0.0.1:" + port + contextPath + "/?" + ran.String(10))
		}
	}()
}

// OpenExplorer ...
func OpenExplorer(port string) {
	OpenExplorerWithContext("", port)
}
