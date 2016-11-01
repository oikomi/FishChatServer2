package itime

import (
	"fmt"
	"log"
	"sync"
	itime "time"
)

const (
	maxInt64         = (1<<63 - 1)
	infiniteDuration = itime.Duration(maxInt64)
)

var (
	globalTimer = NewTimer(4096)
	timePool    = &sync.Pool{
		New: func() interface{} {
			return make(chan itime.Time, 1)
		},
	}
)

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
func After(d itime.Duration) (now itime.Time) {
	var (
		c  chan itime.Time
		td *TimerData
	)
	c = timePool.Get().(chan itime.Time)
	td = globalTimer.Start(d, func() {
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
	now = <-c
	td.Stop()
	timePool.Put(c)
	return
}

// AfterFunc waits for the duration to elapse and then calls f.
// MUST call Stop after expired or terminated.
// f MUST NOT BLOCK!!!!!!!!
func AfterFunc(d itime.Duration, f func()) (td *TimerData) {
	td = globalTimer.Start(d, f)
	return
}

type TimerData struct {
	fn     func() // must nonblock!!!
	index  int
	expire int64
	period int64
	next   *TimerData
	timer  *Timer
}

// Reset reset the timer data with a new duration.
func (td *TimerData) Reset(d itime.Duration) bool {
	return td.timer.reset(td, d)
}

// Stop stop the timer data, returned it stoped or expired.
func (td *TimerData) Stop() bool {
	return td.timer.stop(td)
}

func (td *TimerData) String() string {
	return fmt.Sprintf(`
-------------
index:  %d
expire: %d
fn:     %p
next:   %p
timer:  %p
-------------
`, td.index, td.expire, td.fn, td.next, td.timer)
}

type Timer struct {
	lock   sync.Mutex
	signal *itime.Timer
	free   *TimerData
	timers []*TimerData
	size   int
}

// NewTimer new a timer.
func NewTimer(size int) (t *Timer) {
	t = new(Timer)
	t.init(size)
	return t
}

// Init init the timer.
func (t *Timer) Init(size int) {
	t.init(size)
}

// init init the timer.
func (t *Timer) init(size int) {
	t.signal = itime.NewTimer(infiniteDuration) // never send signal when init
	t.timers = make([]*TimerData, 0, size)
	t.size = size
	t.grow()
	go t.timerproc()
}

// grow grow the freelist timerData.
// free-> []timerData -> []timerData -> []timerData
func (t *Timer) grow() {
	var (
		i   int
		td  *TimerData
		tds = make([]TimerData, t.size) // only one object, optimize GC
	)
	// use free list reuse object
	t.free = &(tds[0])
	td = t.free
	for i = 1; i < t.size; i++ {
		td.next = &(tds[i])
		td.timer = t
		td = td.next
	}
	td.timer = t
	td.next = nil
	return
}

// get get a free timer data, if no free call the grow.
func (t *Timer) get() (td *TimerData) {
	if td = t.free; td == nil {
		t.grow()
		td = t.free
	}
	t.free = td.next
	return
}

// put put back a timer data to free list.
func (t *Timer) put(td *TimerData) {
	td.next = t.free
	t.free = td
}

// when is a helper function for setting the 'when' field of a timer.
// It returns what the time will be, in nanoseconds, Duration d in the future.
// If d is negative, it is ignored.  If the returned value would be less than
// zero because of an overflow, MaxInt64 is returned.
func when(d itime.Duration) int64 {
	if d <= 0 {
		return itime.Now().UnixNano()
	}
	t := itime.Now().UnixNano() + int64(d)
	if t < 0 {
		t = maxInt64
	}
	return t
}

// Start start the timer, if expired then call fn, the returned TimerData must
// Stop after expired or terminated.
// fn MUST NOT BLOCK!!!!!!!!!!!!!!!
func (t *Timer) Start(d itime.Duration, fn func()) (td *TimerData) {
	t.lock.Lock()
	td = t.get()
	td.period = 0
	td.expire = when(d)
	td.fn = fn
	t.add(td)
	t.lock.Unlock()
	return
}

// StartPeriod start the timer, if expired then call fn, the returned TimerData
// must Stop after expired or terminated.
// fn MUST NOT BLOCK!!!!!!!!!!!!!!!
func (t *Timer) StartPeriod(d itime.Duration, fn func()) (td *TimerData) {
	t.lock.Lock()
	td = t.get()
	td.expire = when(d)
	td.period = int64(d)
	td.fn = fn
	t.add(td)
	t.lock.Unlock()
	return
}

// add add a timer data into timer.
func (t *Timer) add(td *TimerData) {
	var d itime.Duration
	td.index = len(t.timers)
	t.timers = append(t.timers, td) // add to the minheap last node
	t.up(td.index)
	if td.index == 0 {
		// if first node, signal start goroutine
		d = itime.Duration(td.expire - itime.Now().UnixNano())
		t.signal.Reset(d)
		if debug {
			log.Printf("timer: reset signal %d\n", d)
		}
	}
	if debug {
		log.Printf("timer: add %s\n", td)
	}
	return
}

// stop stop the timer data, returned the timer stoped or expired.
func (t *Timer) stop(td *TimerData) (ok bool) {
	t.lock.Lock()
	ok = t.del(td)
	t.put(td)
	t.lock.Unlock()
	return
}

// del del a timer data from timer.
func (t *Timer) del(td *TimerData) bool {
	var (
		i    = td.index
		last = len(t.timers) - 1
	)
	if i < 0 || i > last || t.timers[i] != td {
		// already remove, usually by expire
		if debug {
			log.Printf("timer: already del %s\n", td)
		}
		return false
	}
	if i != last {
		t.swap(i, last)
		t.down(i, last)
		t.up(i)
	}
	// remove item is the last node
	t.timers[last].index = -1 // for safety
	t.timers = t.timers[:last]
	if debug {
		log.Printf("timer: del %s\n", td)
	}
	return true
}

// reset reset the timer data with a new expire duration.
func (t *Timer) reset(td *TimerData, d itime.Duration) (ok bool) {
	t.lock.Lock()
	ok = t.del(td)
	td.expire = when(d)
	t.add(td)
	t.lock.Unlock()
	return
}

// timerproc runs the time-driven events.
// It sleeps until the next event in the timers heap.
// If addtimer inserts a new earlier event, addtimer1 wakes timerproc early.
func (t *Timer) timerproc() {
	var (
		now   int64
		delta int64
		fn    func()
		td    *TimerData
	)
	for {
		t.lock.Lock()
		now = itime.Now().UnixNano()
		for {
			if len(t.timers) == 0 {
				delta = maxInt64
				if debug {
					log.Printf("timer: no other instance\n")
				}
				break
			}
			td = t.timers[0]
			if delta = td.expire - now; delta > 0 {
				break
			}
			if td.period > 0 {
				td.expire += td.period * (1 + -delta/td.period)
				t.down(0, len(t.timers)-1)
			} else {
				t.del(td) // let caller put back
			}
			fn = td.fn
			t.lock.Unlock() // fn maybe blocking, release lock
			if fn != nil {
				if debug {
					log.Printf("timer: expire %s\n", td)
				}
				fn()
			} else {
				if debug {
					log.Printf("timer: expire timer no fn\n")
				}
			}
			t.lock.Lock()
		}
		t.signal.Reset(itime.Duration(delta))
		t.lock.Unlock()
		if debug {
			log.Printf("timer: reset signal %d\n", delta)
		}
		<-t.signal.C
	}
	return
}

// minheap

func (t *Timer) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !t.less(j, i) {
			break
		}
		t.swap(i, j)
		j = i
	}
}

func (t *Timer) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !t.less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}
		if !t.less(j, i) {
			break
		}
		t.swap(i, j)
		i = j
	}
}

func (t *Timer) less(i, j int) bool {
	return t.timers[i].expire < t.timers[j].expire
}

func (t *Timer) swap(i, j int) {
	t.timers[i], t.timers[j] = t.timers[j], t.timers[i]
	t.timers[i].index = i
	t.timers[j].index = j
}
