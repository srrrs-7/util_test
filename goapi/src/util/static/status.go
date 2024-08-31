package static

type Status string

const (
	DONE    Status = "done"
	ERROR   Status = "error"
	RUNNING Status = "runnning"
	PENDING Status = "pending"
)
