package lo

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// Recover 在系统崩溃是，恢复系统
func Recover() {
	if err := recover(); err != nil {
		logrus.Warnln(err)
		debug.PrintStack()
	}
}

// Err logs err if it is not nil
func Err(err error) {
	if err != nil {
		logrus.Warnf("error %v", err)
	}
}
