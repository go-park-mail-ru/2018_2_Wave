/**
 *
 */
package proto

import (
	"sync"
	"sync/atomic"
)

// ----------------| NumerationType

// NumerationType - caunting type enumeration
type NumerationType int

const (
	// Counting - serial value
	Counting NumerationType = iota
	// FillGaps - the lowest empty value
	FillGaps
)

// ----------------| Counter

// Counter - type generates IDs for users
type Counter struct {
	userCounterMap  map[UserToken]int64
	UserCounterType NumerationType
	userCounter     int64
	mu              sync.Mutex
}

// MakeCounter - constructor
func MakeCounter(CounterType NumerationType) Counter {
	return Counter{
		userCounterMap:  map[UserToken]int64{},
		UserCounterType: CounterType,
	}
}

// Register a new user
func (r *Counter) Register(t UserToken) {
	counter := r.getNextCounter()
	r.userCounterMap[t] = counter
}

// Unregister the user
func (r *Counter) Unregister(t UserToken) {
	delete(r.userCounterMap, t)
}

// GetUserID if the user exists
func (r *Counter) GetUserID(u IUser) (counter int64, err error) {
	if u == nil {
		return 0, ErrorNil
	}
	return r.GetTokenID(u.GetToken())
}

// GetTokenID - get user id by token
func (r *Counter) GetTokenID(t UserToken) (counter int64, err error) {
	if counter, ok := r.userCounterMap[t]; ok {
		return counter, nil
	}
	return 0, ErrorNotExists
}

// get the next empty id
func (r *Counter) getNextCounter() int64 {
	// find the lowest gap
	if r.UserCounterType == FillGaps {
		r.mu.Lock()
		defer r.mu.Unlock()

		set := make([]bool, r.userCounter+1)
		for _, c := range r.userCounterMap {
			set[c] = true
		}
		for i, bUsed := range set {
			if bUsed {
				continue
			}
			return int64(i)
		}
	}
	// get a next value
	return atomic.AddInt64(&r.userCounter, 1)
}
