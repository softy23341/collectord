package delayedjob

import (
	"errors"
	"sync"
	"time"

	"github.com/eapache/channels"
)

type poolState int

const (
	stateInit poolState = iota
	stateRunning
	stateStopping
)

// Pool TBD
type Pool struct {
	workersCnt   int
	granularity  time.Duration
	state        poolState
	mu           *sync.Mutex
	jobCh        *channels.InfiniteChannel
	delayedJobCh chan Job
	doneCh       chan struct{}
}

// NewPool TBD
func NewPool(workersCnt int, granularity time.Duration) *Pool {
	return &Pool{
		workersCnt:   workersCnt,
		granularity:  granularity,
		state:        stateInit,
		mu:           new(sync.Mutex),
		jobCh:        channels.NewInfiniteChannel(),
		delayedJobCh: make(chan Job),
		doneCh:       make(chan struct{}),
	}
}

// AsyncRun TBD
func (p *Pool) AsyncRun() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch p.state {
	case stateRunning:
		return errors.New("Can't run pool: already running")
	case stateStopping:
		return errors.New("Can't run pool: pool is stopping")
	}
	p.state = stateRunning
	go p.startWorkers()
	go p.startDelayedExecutor()

	return nil
}

// Enqueue TBD
func (p *Pool) Enqueue(jobs ...Job) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.state {
	case stateInit:
		return errors.New("Can't add job: pool is not running")
	case stateStopping:
		return errors.New("Can't add job: pool is stopping")
	}

	now := time.Now()
	for _, job := range jobs {
		if job.RunAt().Before(now) {
			p.jobCh.In() <- job
		} else {
			p.delayedJobCh <- job
		}
	}

	return nil
}

func (p *Pool) startWorkers() {
	wg := new(sync.WaitGroup)

	for i := 0; i < p.workersCnt; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range p.jobCh.Out() {
				job.(Job).Run()
			}
		}()
	}

	wg.Wait()
	close(p.doneCh)
}

func (p *Pool) startDelayedExecutor() {
	var (
		pq           = NewPriorityQueue()
		ticker       = time.NewTicker(p.granularity)
		delayedJobCh = p.delayedJobCh
		stopping     = false
	)

	for !(stopping && pq.Len() == 0) {
		select {
		case job, open := <-delayedJobCh:
			if open {
				pq.Push(job)
			} else {
				delayedJobCh = nil
				stopping = true
			}
		case now := <-ticker.C:
			for pq.Len() > 0 {
				if pq.Peek().RunAt().Before(now) {
					p.jobCh.In() <- pq.Pop()
				} else {
					break // for
				}
			}
		}
	}

	ticker.Stop()
	p.jobCh.Close()
}

// AsyncStop TBD
func (p *Pool) AsyncStop() (doneCh <-chan struct{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch p.state {
	case stateInit:
		close(p.doneCh)
	case stateRunning:
		close(p.delayedJobCh)
	}
	p.state = stateStopping
	return p.doneCh
}
