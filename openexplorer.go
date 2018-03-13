package go_utils

import (
	"github.com/skratchdot/open-golang/open"
	"runtime"
	"time"
)

func OpenExplorer(port string) {
	go func() {
		time.Sleep(100 * time.Millisecond)

		switch runtime.GOOS {
		case "windows":
			fallthrough
		case "darwin":
			open.Run("http://127.0.0.1:" + port + "/?" + RandString(10))
		}
	}()
}
