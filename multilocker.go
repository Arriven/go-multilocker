package multilocker

import (
	"runtime"
	"unsafe"
	"sort"
)


//Lockable defines minimal interface for types that are supported
type Lockable interface {
	Lock()
	Unlock()
}

//TryLockable extended Lockable interface that allows better functionality
type TryLockable interface {
	Lockable
	TryLock() bool
}

//Locker is a struct used for locking/unlocking
//It's possible to implement with free functions but the code will become more complex.
//Feel free to make a PR if you need such functionality
type Locker struct {
	locks []Lockable
}

//Lock is a function that locks all resources provided.
//It's guaranteed to avoid deadlock for TryLockable types but might potentially fail into livelock.
//It has high chance of avoiding deadlock on regular locks as well, but deadlock avoidance is not guaranteed.
//If a panic is thrown during locking of one of the resouces it'll unlock all the acquired resources.
func (l *Locker) Lock(locks ...Lockable) {
	if tryableLocks, ok := getTryLockable(locks...); ok {
		for !l.TryLock(tryableLocks...) {
			runtime.Gosched()
		}
	} else {
		defer l.unlockOnPanic()
		sort.Slice(locks, func(i, j int) bool { return uintptr(unsafe.Pointer(&locks[i])) < uintptr(unsafe.Pointer(&locks[j])) })
		for _, lock := range locks {
			lock.Lock()
			l.locks = append(l.locks, lock)
		}
	}
}

//TryLock tries to acquire all provided resources.
//If a panic is thrown during locking of one of the resouces it'll unlock all the acquired resources.
func (l *Locker) TryLock(locks ...TryLockable) bool {
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

func getTryLockable(locks ...Lockable) ([]TryLockable, bool) {
	result := []TryLockable{}
	for _, lock := range locks {
		if tryableLock, ok := lock.(TryLockable); ok {
			result = append(result, tryableLock)
		} else {
			return nil, false
		}
	}
	return result, true
}

func (l *Locker) unlockOnPanic() {
	if r := recover(); r != nil {
		l.Unlock()
		panic(r)
	}
}
