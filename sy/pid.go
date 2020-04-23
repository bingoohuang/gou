package sy

import (
	"os"
	"strconv"

	"github.com/bingoohuang/gou/file"
	"github.com/pkg/errors"
)

// UpdatePidFile update the pid to pidFile like var/pid (kill -USR2 {pid} 会执行重启)
// If pidFile is empty, it will try env PID_FILE' value, then to "var/pid" file.
func UpdatePidFile(pidFileVars ...string) error {
	pidFile := ""

	if len(pidFileVars) > 0 {
		pidFile = pidFileVars[0]
	}

	if pidFile == "" {
		pidFile = os.Getenv("PID_FILE")
	}

	if pidFile == "" {
		pidFile = "var/pid"
	}

	currentPID := strconv.Itoa(os.Getpid())
	oldPid, err := file.ReadValue(pidFile, currentPID)

	if err != nil {
		return errors.Wrapf(err, "read pid file %s", pidFile)
	}

	if oldPid != currentPID {
		if err := file.WriteValue(pidFile, currentPID); err != nil {
			return errors.Wrapf(err, "write pid file %s", pidFile)
		}
	}

	return nil
}
