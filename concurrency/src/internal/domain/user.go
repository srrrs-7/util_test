package domain

import "concurrency/pkg/utilerror"

type UserID int

type User struct {
	ID UserID `json:"user_id"`
}

func (u User) validate() error {
	if u.ID <= 0 {
		return utilerror.ErrInvalidUserID
	}

	return nil
}
