package lib

func TwoSum(nums []int, target int) []int {
	m := make(map[int]int, 0)
	for idx, num := range nums {
		if requireIdx, exists := m[target-num]; exists {
			return []int{requireIdx, idx}
		}
		m[num] = idx
	}
	return []int{}
}
