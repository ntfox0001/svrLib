package debug

import (
	"runtime"
)

// The maximum stack buffer size (64K). This is the same value used by
// net.Conn's serve() method.
const maxStackBufferSize = (64 << 10)

func RuntimeStacks() string {
	buf := make([]byte, maxStackBufferSize)
	return string(buf[:runtime.Stack(buf, true)])
}
