package sy

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

// nolint gochecknoglobals
var goroutineSpace = []byte("goroutine ")

// CurGoroutineID return current goroutine ID.
func CurGoroutineIDString() string {
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)

	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)

	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}

	return string(b[:i])
}

// CurGoroutineID return current goroutine ID.
func CurGoroutineID() uint64 {
	return curGoroutineID()
}

func curGoroutineID() uint64 {
	gid := CurGoroutineIDString()
	n, err := strconv.ParseUint(gid, 10, 64)

	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %s: %v", gid, err))
	}

	return n
}

// nolint gochecknoglobals
var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}
