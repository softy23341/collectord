package delayedjob

import (
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolBasic(t *testing.T) {
	var (
		pool = NewPool(1, 10*time.Millisecond)
		job  = NewJob(50*time.Millisecond, func() {})
	)

	if pool.Enqueue(job) == nil {
		t.Error(".Enqueue on non running pool should return error")
	}

	if err := pool.AsyncRun(); err != nil {
		t.Fatalf("No error expected during AsyncRun, but got: %v", err)
	}

	if pool.AsyncRun() == nil {
		t.Error("Error expected when invoke AsyncRun on running pool")
	}

	// add delayed job to pool to prevent fast pool stopping
	if err := pool.Enqueue(job); err != nil {
		t.Fatalf("No error expected during Enqueue, but got: %v", err)
	}

	doneCh := pool.AsyncStop()

	if pool.Enqueue(job) == nil {
		t.Error(".Enqueue on stopping pool should return error")
	}

	select {
	case <-doneCh:
		break
	case <-time.Tick(100 * time.Millisecond):
		t.Fatalf("Pool stop timeout")
	}

	if pool.Enqueue(job) == nil {
		t.Error(".Enqueue on stopped pool should return error")
	}

	if pool.AsyncStop() != doneCh {
		t.Error("Multiply AsyncStop calls should return same channel")
	}
}

func TestPoolUnderLoad(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		njobs       = 100000
		nworkers    = 10
		granularity = 100 * time.Millisecond
	)

	pool := NewPool(nworkers, granularity)
	if err := pool.AsyncRun(); err != nil {
		t.Fatalf("Can't run pool, error: %v", err)
	}

	var (
		cumalativeSum         uint64
		expectedCumalativeSum uint64
		nfinished             uint64
		rnd                   = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 1; i <= njobs; i++ {
		expectedCumalativeSum += uint64(i)
		delay := time.Duration(rnd.Intn(1000)) * time.Millisecond
		iClone := i // clone to use inside job
		job := NewJob(delay, func() {
			atomic.AddUint64(&cumalativeSum, uint64(iClone))
			atomic.AddUint64(&nfinished, 1)
		})
		if err := pool.Enqueue(job); err != nil {
			t.Fatalf("Can't enqueue job to pool, error: %v", err)
		}
	}

	select {
	case <-pool.AsyncStop():
		break
	case <-time.After(2 * time.Second):
		t.Fatalf("Pool execution timeout, number of finished jobs: %d", nfinished)
	}

	if nfinished != uint64(njobs) {
		t.Errorf("Not all jobs processed, expected: %d, actual: %d", njobs, nfinished)
	}

	if cumalativeSum != expectedCumalativeSum {
		t.Errorf(
			"Wrong jobs result, expected cumulative sum: %d, actual: %d",
			expectedCumalativeSum,
			cumalativeSum,
		)
	}
}
