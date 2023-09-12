package delayedjob

import "time"

// Immideate TBD
var Immideate = time.Duration(0)

// Job TBD
type Job interface {
	Run()
	RunAt() time.Time
}

// NewJob TBD
func NewJob(delay time.Duration, fn func()) Job {
	return &simpleJob{
		fn:    fn,
		runAt: time.Now().Add(delay),
	}
}

type simpleJob struct {
	fn    func()
	runAt time.Time
}

func (s simpleJob) Run() { s.fn() }

func (s simpleJob) RunAt() time.Time { return s.runAt }
