package itime

import (
	"errors"
	itime "time"
)

// A Ticker holds a channel that delivers `ticks' of a clock
// at intervals.
type Ticker struct {
	C  <-chan itime.Time // The channel on which the ticks are delivered.
	c  chan itime.Time
	td *TimerData
	t  *Timer
}

// NewTicker returns a new Ticker containing a channel that will send the
// time with a period specified by the duration argument.
// It adjusts the intervals or drops ticks to make up for slow receivers.
// The duration d must be greater than zero; if not, NewTicker will panic.
// Stop the ticker to release associated resources.
func NewTicker(t *Timer, d itime.Duration) *Ticker {
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	var (
		c  = timePool.Get().(chan itime.Time)
		td = t.StartPeriod(d, func() {
			// Non-blocking send of time on c.
			// Used in NewTimer, it cannot block anyway (buffer).
			// Used in NewTicker, dropping sends on the floor is
			// the desired behavior when the reader gets behind,
			// because the sends are periodic.
			select {
			case c <- itime.Now():
			default:
			}
		})
	)
	return &Ticker{
		C:  c,
		c:  c,
		t:  t,
		td: td,
	}
}

// Stop turns off a ticker.  After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t *Ticker) Stop() {
	if t.td.Stop() {
		timePool.Put(t.c)
	}
}
