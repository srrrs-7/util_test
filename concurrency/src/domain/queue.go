package domain

type QueueID string

func (q QueueID) String() string {
	if q == "" {
		return ""
	}
	return string(q)
}
