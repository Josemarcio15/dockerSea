package stacks

import (
	"sync"
	"testing"
	"time"
)

func TestStackLockManager_SerializesSameKey(t *testing.T) {
	mgr := NewStackLockManager()
	key := "prof_1:stk_1"

	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unlock := mgr.Lock(key)
			defer unlock()

			current := counter
			time.Sleep(2 * time.Millisecond)
			counter = current + 1
		}()
	}

	wg.Wait()

	if counter != 20 {
		t.Fatalf("esperado contador = 20, obtido %d (corrida de concorrência detectada)", counter)
	}

	// Verifica se o mapa foi limpo após término
	mgr.mu.Lock()
	mapLen := len(mgr.locks)
	mgr.mu.Unlock()

	if mapLen != 0 {
		t.Fatalf("esperado locks map limpo (0 entradas), mas tem %d", mapLen)
	}
}

func TestStackLockManager_ParallelDifferentKeys(t *testing.T) {
	mgr := NewStackLockManager()
	key1 := "prof_1:stk_1"
	key2 := "prof_1:stk_2"

	started1 := make(chan struct{})
	finished1 := make(chan struct{})

	go func() {
		unlock1 := mgr.Lock(key1)
		close(started1)
		time.Sleep(50 * time.Millisecond)
		unlock1()
		close(finished1)
	}()

	<-started1

	// Deve conseguir adquirir lock da key2 imediatamente mesmo com key1 ocupada
	start2 := time.Now()
	unlock2 := mgr.Lock(key2)
	elapsed := time.Since(start2)
	unlock2()

	if elapsed > 30*time.Millisecond {
		t.Fatalf("lock em key2 demorou muito (%v), deveria ser imediato", elapsed)
	}

	<-finished1
}
