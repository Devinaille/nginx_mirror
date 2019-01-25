package count

import (
	"testing"
	"time"
)

func TestGetMsecPos(t *testing.T) {
	testTime := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.Local)
	result := getMsecPos(testTime)
	if result != 0 {
		t.FailNow()
	}
}
