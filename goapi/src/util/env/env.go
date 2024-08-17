package env

import "fmt"

type Env struct{}

func NewEnv() Env {
	return Env{}
}

func (e Env) Validate() bool {
	return true
}

func (e Env) OutPut() string {
	return fmt.Sprintf("%v", e)
}
