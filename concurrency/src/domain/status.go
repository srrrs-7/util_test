package domain

const (
	Completed = "completed"
	Failed    = "failed"
	Running   = "running"
	Pending   = "pending"

	StatusCompleted Status = iota + 1
	StatusFailed
	StatusRunning
	StatusPending
)

type Status int

func (s Status) String() string {
	switch s {
	case StatusCompleted:
		return Completed
	case StatusFailed:
		return Failed
	case StatusRunning:
		return Running
	case StatusPending:
		return Pending
	default:
		return "unknown"
	}
}
