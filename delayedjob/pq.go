package delayedjob

import "container/heap"

// NewPriorityQueue TBD
func NewPriorityQueue(jobs ...Job) (pq PriorityQueue) {
	pq = PriorityQueue{internal: (*internalPriorityQueue)(&jobs)}
	heap.Init(pq.internal)
	return
}

// PriorityQueue TBD
type PriorityQueue struct {
	internal *internalPriorityQueue
}

// Len TBD
func (pq PriorityQueue) Len() int {
	return pq.internal.Len()
}

// Peek TBD
func (pq PriorityQueue) Peek() Job {
	return (*pq.internal)[0]
}

// Push TBD
func (pq PriorityQueue) Push(job Job) {
	heap.Push(pq.internal, job)
}

// Pop TBD
func (pq PriorityQueue) Pop() Job {
	return heap.Pop(pq.internal).(Job)
}

// internalPriorityQueue implements heap.Interface
type internalPriorityQueue []Job

func (ipq internalPriorityQueue) Len() int { return len(ipq) }

func (ipq internalPriorityQueue) Less(i, j int) bool {
	return ipq[i].RunAt().Before(ipq[j].RunAt())
}

func (ipq internalPriorityQueue) Swap(i, j int) {
	ipq[i], ipq[j] = ipq[j], ipq[i]
}

func (ipq *internalPriorityQueue) Push(x interface{}) {
	*ipq = append(*ipq, x.(Job))
}

func (ipq *internalPriorityQueue) Pop() interface{} {
	old := *ipq
	n := len(old)
	item := old[n-1]
	*ipq = old[0 : n-1]
	return item
}
