package io

import "io"

// Close ...
func Close(c io.Closer) {
	_ = c.Close()
}
