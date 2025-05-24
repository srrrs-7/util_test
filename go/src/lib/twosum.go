package lib

func TwoSum(nums []int, target int) []int {
	m := make(map[int]int, len(nums))

	for i, n := range nums {
		diff := target - n
		if v, ok := m[diff]; ok && v != i {
			return []int{v, i}
		}
		m[n] = i
	}
	return nil
}
