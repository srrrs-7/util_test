package request

import "fmt"

type CreateReq struct {
	UserID string `json:"user_id"`
}

func (c *CreateReq) Validate() error {
	if c.UserID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	return nil
}

type CheckReq struct {
	ID string `json:"id"`
}

func (c *CheckReq) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	return nil
}
