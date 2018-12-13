package room

import "sync/atomic"

type NumerationType int

const (
	Counting NumerationType = iota
	FillGaps
)

type Counter struct {
	userCounterMap  map[UserID]int64
	UserCounterType NumerationType
	userCounter     int64
}

func NewCounter(CounterType NumerationType) *Counter {
	return &Counter{
		userCounterMap:  map[UserID]int64{},
		UserCounterType: CounterType,
	}
}

func (r *Counter) Add(t UserID) {
	counter := r.getNextCounter()
	r.userCounterMap[t] = counter
}

func (r *Counter) Delete(t UserID) {
	delete(r.userCounterMap, t)
}

func (r *Counter) GetUserCounter(u IUser) (counter int64, err error) {
	if u == nil {
		return 0, ErrorNil
	}
	if counter, ok := r.userCounterMap[u.GetID()]; ok {
		return counter, nil
	}
	return 0, ErrorNotExists
}

func (r *Counter) GetTokenCounter(t UserID) (counter int64, err error) {
	if counter, ok := r.userCounterMap[t]; ok {
		return counter, nil
	}
	return 0, ErrorNotExists
}

func (r *Counter) getNextCounter() int64 {
	if r.UserCounterType == FillGaps {
		set := make([]bool, r.userCounter+1)
		for _, c := range r.userCounterMap {
			set[c] = true
		}
		for i, ok := range set {
			if !ok {
				return int64(i)
			}
		}
	}
	return atomic.AddInt64(&r.userCounter, 1)
}
