package multilocker

import (
	"runtime"
)

type Lock interface {
	TryLock() bool
	Unlock()
}

type Locker struct {
	locks []Lock
}

func (l *Locker) Lock(locks ...Lock) {
	for !l.TryLock(locks...) {
		runtime.Gosched()
	}
}

func (l *Locker) TryLock(locks ...Lock) bool {
	defer l.unlockOnPanic()
	for _, lock := range locks {
		if lock.TryLock() {
			l.locks = append(l.locks, lock)
		} else {
			l.Unlock()
			return false
		}
	}
	return true
}

func (l *Locker) Unlock() {
	for _, lock := range l.locks {
		lock.Unlock()
	}
	l.locks = nil
}

func (l *Locker) unlockOnPanic() {
	if r := recover(); r != nil {
		l.Unlock()
		panic(r)
	}
}
