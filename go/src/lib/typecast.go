package lib

import (
	"errors"
	"fmt"
	"reflect"
)

type Custom struct {
	Name interface{}
}

func (c *Custom) Count() int {
	arr, ok := c.Name.([]int)
	if !ok {
		return -1
	}

	return len(arr)
}

func (c *Custom) TypeCheck() {
	t := reflect.TypeOf(c.Name).Kind()
	switch t {
	case reflect.Chan:
		fmt.Printf("type is: %s \n", t)
	case reflect.String:
		fmt.Printf("type is: %s \n", t)
	case reflect.Bool:
		fmt.Printf("type is: %s \n", t)
	case reflect.Array:
		fmt.Printf("type is: %s \n", t)
	case reflect.Slice:
		fmt.Printf("type is: %s \n", t)
	default:
		fmt.Printf("default: %s \n", t)
	}
}

func TypeCast[T any](v any) ([]T, error) {
	t := reflect.TypeOf(v).Kind()
	fmt.Println(t)
	switch t {
	case reflect.Slice:
		tmp, ok := v.([]T)
		if !ok {
			return nil, errors.New("value type cast error")
		}
		return tmp, nil
	default:
		return nil, errors.New("type cast error")
	}
}
