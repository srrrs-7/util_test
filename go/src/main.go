package main

import (
	"fmt"
	"test/lib"
)

func main() {
	ch := make(chan string, 10)

	lib.Multi(ch)
	for s := range ch {
		fmt.Println(s)
	}

	value := lib.NewValue(1)
	for range 100 {
		value.Increment()
	}
	fmt.Println(value)

	lib.LongestCommonSubsequence()

	lib.KonpSack()
}
