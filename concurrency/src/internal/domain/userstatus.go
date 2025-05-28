package domain

type UserStatus struct {
	QueueID QueueID `json:"queue_id,omitempty"`
	UserID  UserID  `json:"id"`
	Status  Status  `json:"status"`
	Context string  `json:"context,omitempty"`
}
