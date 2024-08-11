package main

import (
	"fmt"
	"test/lib"
)

func main() {
	lib.LongestCommonSubsequence()
	lib.KonpSack()
	fmt.Println(lib.TwoSum([]int{1, 2, 3}, 4))

	cus := lib.Custom{
		Name: []int{1, 2},
	}
	fmt.Println("count: ", cus.Count())

	cus.TypeCheck()
	arr, err := lib.TypeCast[int](cus.Name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(arr)
}
