package count

import "time"

func getMsecPos(t time.Time) int {
	msec := t.Nanosecond() / int(time.Millisecond)
	return msec % 1000
}
