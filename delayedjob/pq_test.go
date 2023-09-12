package delayedjob

import (
	"testing"
	"time"
)

func createDummyJobs() (j100 Job, j200 Job, j300 Job) {
	j100 = NewJob(time.Duration(100), func() {})
	j200 = NewJob(time.Duration(200), func() {})
	j300 = NewJob(time.Duration(200), func() {})
	return
}

func TestPqLen(t *testing.T) {
	j100, j200, j300 := createDummyJobs()
	pq := NewPriorityQueue(j300, j200, j100)

	if pq.Len() != 3 {
		t.Error("Init pq len must be 3")
	}
}

// using Len for testing
func TestPqPop(t *testing.T) {
	j100, j200, j300 := createDummyJobs()
	pq := NewPriorityQueue(j300, j200, j100)

	if pq.Pop() != j100 {
		t.Error("Dirst pop must return j100")
	}
	if pq.Pop() != j200 {
		t.Error("Second pop must return j200")
	}
	if pq.Pop() != j300 {
		t.Error("Third pop must return j300")
	}

	if pq.Len() != 0 {
		t.Error("After all pops pq must be empty")
	}
}

// using Pop for testing
func TestPqPeek(t *testing.T) {
	j100, j200, j300 := createDummyJobs()
	pq := NewPriorityQueue(j300, j200, j100)

	if pq.Peek() != j100 {
		t.Error("j100 must be first in queue")
	}
	pq.Pop()

	if pq.Peek() != j200 {
		t.Error("j200 must be second in queue")
	}
	pq.Pop()

	if pq.Peek() != j300 {
		t.Error("j300 must be third in queue")
	}
	pq.Pop()
}

// using Peek for testing
func TestPqPush(t *testing.T) {
	j100, j200, j300 := createDummyJobs()
	pq := NewPriorityQueue()

	pq.Push(j300)
	if pq.Peek() != j300 {
		t.Error("After push j300 peek should return j300")
	}

	pq.Push(j200)
	if pq.Peek() != j200 {
		t.Error("After push j200 peek should return j200")
	}

	pq.Push(j100)
	if pq.Peek() != j100 {
		t.Error("After push j100 peek should return j100")
	}

	if pq.Len() != 3 {
		t.Error("After all pushes pq len must be 3")
	}
}
