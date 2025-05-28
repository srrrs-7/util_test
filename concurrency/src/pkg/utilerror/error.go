package utilerror

import "fmt"

var (
	ErrInvalidUserID  = fmt.Errorf("invalid user ID, must be a positive integer greater than 0")
	ErrInvalidQueueID = fmt.Errorf("invalid queue ID, must be a non-empty string")
	ErrInvalidContent = fmt.Errorf("invalid content, must be a non-empty map with string keys and any values")
	ErrInvalidStatus  = fmt.Errorf("invalid status, must be one of the predefined statuses")
)
