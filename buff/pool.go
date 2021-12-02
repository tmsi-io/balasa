package buff

import (
	"bytes"
	"sync"
)

var BufferPool = sync.Pool{
	New: func() interface{} {
		var b = bytes.NewBuffer(nil)
		return b
	},
}
