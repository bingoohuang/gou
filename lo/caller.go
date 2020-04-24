package lo

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/bingoohuang/gou/str"
)

// nolint gochecknoglobals
var (
	// qualified package name, cached at first use
	loPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once

	// UserHome is the directory of user home directory.
	UserHome string
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		loPackage = getPackageName(runtime.FuncForPC(pcs[1]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLogrusFrames

		UserHome, _ = os.UserHomeDir()
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		// If the caller isn't part of this package, we're done
		if pkg := getPackageName(f.Function); str.NoneOf(pkg,
			"github.com/sirupsen/logrus",
			"github.com/rifflock/lfshook",
			loPackage) {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")

		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
