package multilocker

import (
	"runtime"
	"unsafe"
	"sort"
)

//Lock defines minimal interface for types that are supported
type Lockable interface {
	TryLock() bool
	Unlock()
}

//Locker is a struct used for locking/unlocking
//It's possible to implement with free functions but the code will become more complex.
//Feel free to make a PR if you need such functionality
type Locker struct {
	locks []Lockable
}

//Lock is a function that locks all resources provided.
//It's guaranteed to avoid deadlock but might potentially fail into livelock.
//If a panic is thrown during locking of one of the resouces it'll unlock all the acquired resources.
func (l *Locker) Lock(locks ...Lockable) {
	for !l.TryLock(locks...) {
		runtime.Gosched()
	}
}

//TryLock tries to acquire all provided resources.
//If a panic is thrown during locking of one of the resouces it'll unlock all the acquired resources.
func (l *Locker) TryLock(locks ...Lockable) bool {
	defer l.unlockOnPanic()
	sort.Slice(locks, func(i, j int) bool { return uintptr(unsafe.Pointer(&locks[i])) < uintptr(unsafe.Pointer(&locks[j])) })
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

//Unlock just releases all acquired resources
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
