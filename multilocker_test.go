package multilocker

import (
	"github.com/LK4D4/trylock"
	"testing"
	"time"
)

func lockTwoMutexes(l1 Lock, l2 Lock) {
	var l Locker
	l.Lock(l1, l2)
	defer l.Unlock()
	time.Sleep(100 * time.Millisecond)
}

func TestTwoResourcesSameOrder(t *testing.T) {
	var m1, m2 trylock.Mutex
	go lockTwoMutexes(&m1, &m2)
	go lockTwoMutexes(&m1, &m2)
}

func TestTwoResourcesReverseOrder(t *testing.T) {
	var m1, m2 trylock.Mutex
	go lockTwoMutexes(&m1, &m2)
	go lockTwoMutexes(&m2, &m1)
}

type panickingMutex struct {
}

func (m *panickingMutex) TryLock() bool {
	panic("wtf")
	return true
}

func (m *panickingMutex) Unlock() {}

func TestPanicBehavior(t *testing.T) {
	var m trylock.Mutex
	var p panickingMutex
	var l Locker
	defer func() {
		if r := recover(); r != nil {
			t.Log("recovered from panic")
		}
	}()
	l.TryLock(&m, &p)
	if !l.TryLock(&m) {
		t.Fail()
	}
	l.Unlock()
}
