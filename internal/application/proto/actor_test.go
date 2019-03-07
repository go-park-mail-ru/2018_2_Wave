package proto

import (
	"sync"
	"testing"
)

func TestActorPanicRecovery(t *testing.T) {
	var (
		actor = MakeActor(10)
	)
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Unexpected panic %v", err)
		}
	}()
	actor.PanicRecovery(func() {
		panic("panic")
	})
}

func TestActorPanicRecoveryWithRepeat(t *testing.T) {
	var (
		actor = MakeActor(10)
		hits  = 0
	)
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Unexpected panic %v", err)
		}
	}()
	actor.OnPanic = func(interface{}) bool {
		if hits == 0 {
			hits++
			return true
		}
		if hits == 1 {
			hits++
			return false
		}
		t.Fatal("Unexpected OnPanic function call")
		return false
	}

	actor.PanicRecovery(func() {
		panic("panic")
	})
}

func TestActorAsyncPanicRecovery(t *testing.T) {
	var (
		actor = MakeActor(10)
		wg    = sync.WaitGroup{}
	)
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Unexpected panic %v", err)
		}
	}()
	wg.Add(1)
	actor.OnPanic = func(interface{}) bool {
		wg.Done()
		return false
	}
	actor.PanicRecoveryAsync(func() {
		panic("panic")
	})
	wg.Wait()
}

