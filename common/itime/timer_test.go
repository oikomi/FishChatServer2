package itime

import (
	"testing"
	"time"
)

func TestTimerFreeList(t *testing.T) {
	// get, put, grow
	timer := new(Timer)
	timer.size = 1
	td := timer.get()
	if td == nil {
		t.FailNow()
	}
	if timer.free != nil {
		t.FailNow()
	}
	td1 := timer.get()
	if td1 == nil {
		t.FailNow()
	}
	timer.put(td)
	if timer.free != td {
		t.FailNow()
	}
	timer.put(td1)
	if timer.free != td1 {
		t.FailNow()
	}
	if td1.next != td {
		t.FailNow()
	}
}

func TestTimer(t *testing.T) {
	timer := NewTimer(100)
	tds := make([]*TimerData, 100)
	for i := 0; i < 100; i++ {
		tds[i] = timer.Start(time.Duration(i)*time.Second+5*time.Minute, nil)
	}
	for i := 0; i < 100; i++ {
		tds[i].Stop()
	}
	for i := 0; i < 100; i++ {
		tds[i] = timer.Start(time.Duration(i)*time.Second+5*time.Minute, nil)
	}
	for i := 0; i < 100; i++ {
		tds[i].Stop()
	}
	// Start
	timer.Start(time.Second, nil)
	time.Sleep(time.Second * 2)
	if len(timer.timers) != 0 {
		t.FailNow()
	}
	i := 0
	timer.Start(time.Second, func() {
		i++
	})
	timer.Start(time.Millisecond*500, func() {
		i++
	})
	time.Sleep(time.Millisecond * 510)
	if i != 1 {
		t.Errorf("i: %d", i)
		t.FailNow()
	}
	time.Sleep(time.Millisecond * 510)
	if i != 2 {
		t.Errorf("i: %d", i)
		t.FailNow()
	}
	// StartPeriod
}

func TestAfter(t *testing.T) {
	now := time.Now().Unix()
	after := After(time.Second * 1)
	if after.Unix()-now != 1 {
		t.FailNow()
	}
}

func TestAfterFunc(t *testing.T) {
	i := 0
	td := AfterFunc(time.Second*1, func() {
		i++
	})
	time.Sleep(time.Second * 2)
	td.Stop()
	if i != 1 {
		t.Errorf("i: %d", i)
		t.FailNow()
	}
}
