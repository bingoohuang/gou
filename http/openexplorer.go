package http

import (
	"runtime"
	"time"

	"github.com/bingoohuang/gou/rand"

	"github.com/skratchdot/open-golang/open"
)

func OpenExplorerWithContext(contextPath, port string) {
	go func() {
		time.Sleep(100 * time.Millisecond)

		switch runtime.GOOS {
		case "windows":
			fallthrough
		case "darwin":
			_ = open.Run("http://127.0.0.1:" + port + contextPath + "/?" + rand.String(10))
		}
	}()
}

func OpenExplorer(port string) {
	OpenExplorerWithContext("", port)
}
