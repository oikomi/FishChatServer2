package itime

import (
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	tr := NewTicker(globalTimer, 1*time.Second)
	now := time.Now().Unix()
	for i := 0; i < 3; i++ {
		if after := <-tr.C; after.Unix()-now != int64(i+1) {
			t.FailNow()
		}
	}
	tr.Stop()
}
