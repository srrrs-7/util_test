package lib

import "fmt"

type Item struct {
	weight int
	value  int
}

func KonpSack() {

	itemMap := make(map[int]Item, 4)
	itemMap[1] = Item{weight: 2, value: 5}
	itemMap[2] = Item{weight: 3, value: 4}
	itemMap[3] = Item{weight: 9, value: 15}
	itemMap[4] = Item{weight: 1, value: 5}

	// dp table
	dp := make([][]int, len(itemMap)+1)
	for i := range dp {
		dp[i] = make([]int, 11)
	}

	for i := 1; i <= len(itemMap)+1; i++ {
		if item, ok := itemMap[i]; ok {
			for w := 1; w <= 10; w++ {
				if w < item.weight {
					dp[i][w] = dp[i-1][w]
				} else {
					dp[i][w] = max(dp[i-1][w], dp[i-1][w-item.weight]+item.value)
				}

			}
		}
	}

	fmt.Println(dp)
}
