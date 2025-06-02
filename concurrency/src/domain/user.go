package domain

import (
	"fmt"
)

type UserID string

func (u UserID) String() string {
	return string(u)
}

type User struct {
	ID UserID `json:"user_id"`
}

func (u User) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("invalid user ID")
	}

	return nil
}
