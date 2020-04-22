package lang

import (
	"io"

	"github.com/hashicorp/go-multierror"
	errs "github.com/pkg/errors"
)

// Closef runs function and on error return error by argument including the given error (usually
// from caller function).
func Closef(err *error, closer io.Closer, format string, a ...interface{}) {
	*err = multierror.Append(*err, errs.Wrapf(closer.Close(), format, a...))
}
