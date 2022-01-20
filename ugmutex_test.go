package mutex

import (
	"time"
	"testing"
)

func TestUgMutexNew(t *testing.T) {
	ugl := NewUgMutex()

	if ugl == nil {
		t.Fatal("can't create mutex.")
	}
}

func TestUgMutexLock(t *testing.T) {
	l := NewUgMutex()
	l.Lock()
	l.Unlock()
}

func TestUgMutexLockRace(t *testing.T) {
	l := NewUgMutex()

	flg := true

	l.Lock()
	go func() {
		defer l.Unlock()

		flg = false
	}()

	l.Lock()
	defer l.Unlock()

	if flg {
		t.Fatal("can't block when locked.")
	}
}

func TestUgMutexRLock(t *testing.T) {
	l := NewUgMutex()
	l.RLock()
	l.RUnlock()
}

func TestUgMutexRLockRace(t *testing.T) {
	l := NewUgMutex()

	flg := true

	l.RLock()
	go func() {
		defer l.RUnlock()

		time.Sleep(time.Second * 2)

		if !flg {
			t.Fatal("can't block when locked.")
		}
	}()

	l.Lock()
	go func() {
		defer l.Unlock()

		time.Sleep(time.Second)

		flg = false
	}()

	l.RLock()
	defer l.RUnlock()
	if flg {
		t.Fatal("can't block when locked.")
	}
}

func TestUgMutexUgLock(t *testing.T) {
	l := NewUgMutex()
	l.UgLock()
	l.UgUnlock()
}

func TestUgMutexUgLockWithUpgrade(t *testing.T) {
	l := NewUgMutex()
	l.UgLock()
	l.Upgrade()
	l.UgUnlock()
}

func TestUgMutexUgLockRace(t *testing.T) {
	l := NewUgMutex()

	flg := true

	l.UgLock()

	go func() {
		l.RLock()
		defer l.RUnlock()

		if !flg {
			t.Fatal("can't block when locked.")
		}
	}()

	time.Sleep(time.Second)
	l.Upgrade()

	go func() {
		l.RLock()
		defer l.RUnlock()

		if flg {
			t.Fatal("can't block when locked.")
		}
	}()
	go func() {
		l.Lock()
		defer l.Unlock()

		if flg {
			t.Fatal("can't block when locked.")
		}
	}()

	defer l.UgUnlock()
	flg = false
}
