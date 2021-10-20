package sdp

import (
	"fmt"
	"time"
)

type Time struct {
	StartTime int
	EndTime   int
	Repeat    int
	Active    int
	Offsets   []time.Duration
}

func (t *Time) Println() string {
	var s string
	s += fmt.Sprintf("t=%d %d\r\n", t.StartTime, t.EndTime)
	s += fmt.Sprintf("r=7d 1h 0 25h\r\n")
	return s
}
