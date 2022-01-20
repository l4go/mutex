package mutex

import (
	"sync"
	"sync/atomic"
)

type UgMutex struct {
	l   *sync.Mutex

	c_l *sync.Cond
	c_u *sync.Cond

	wlf int32
	ulf int32

	cnt int32
}

func NewUgMutex() *UgMutex {
	l := new(sync.Mutex)
	return &UgMutex{l:l, c_l:sync.NewCond(l), c_u:sync.NewCond(l)}
}

func (m *UgMutex) Lock() {
	for {
		m.lockSlow()

		if atomic.CompareAndSwapInt32(&m.wlf, MODE_OFF, MODE_ON) {
			break
		}
	}
}

func (m *UgMutex) RLock() {
	m.l.Lock()
	if atomic.LoadInt32(&m.wlf) == MODE_ON {
		m.c_l.Wait()
	}

	n_cnt := atomic.AddInt32(&m.cnt, int32(1))
	if n_cnt > MAX_CNT || n_cnt < 0 {
		panic("counter out of range.")
	}
	if n_cnt > 1 {
		m.l.Unlock()
	}
}

func (m *UgMutex) UgLock() {
	m.lockSlow()
	m.ulf = MODE_ON

	m.l.Unlock()
}

func (m *UgMutex) Upgrade() {
	atomic.StoreInt32(&m.wlf, MODE_ON)

	m.l.Lock()
}

func (m *UgMutex) UgUnlock() {
	if !atomic.CompareAndSwapInt32(&m.wlf, MODE_ON, MODE_OFF) {
		m.l.Lock()
	}
	m.ulf = MODE_OFF
	m.l.Unlock()

	m.c_u.Signal()
	m.c_l.Broadcast()
}

func (m *UgMutex) RUnlock() {
	if n_cnt := atomic.AddInt32(&m.cnt, int32(-1)); n_cnt < 1 {
		m.l.Unlock()
	}
}

func (m *UgMutex) Unlock() {
	atomic.StoreInt32(&m.wlf, MODE_OFF)

	m.l.Unlock()
	m.c_u.Signal()
	m.c_l.Broadcast()
}

func (m *UgMutex) lockSlow() {
	m.l.Lock()

	for {
		if m.ulf != MODE_OFF {
			m.c_u.Wait()
			continue
		}

		break
	}
}
