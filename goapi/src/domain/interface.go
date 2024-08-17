package domain

import "api/driver/model"

type Queuer interface {
	Send(msg string) (string, error)
	Receive() (string, error)
	Delete(id string) error
}

type Cacher interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Querier interface {
	Select(id string) ([]model.User, error)
	Insert(user model.User) error
	Update(user model.User) error
	Delete(user model.User) error
}
