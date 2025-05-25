package lib

import "math"

// func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
// 	sumLen := len(nums1) + len(nums2)
// 	nums := make([]int, 0, sumLen)

// 	nums = append(nums, nums1...)
// 	nums = append(nums, nums2...)
// 	slices.Sort(nums)

// 	if sumLen%2 == 0 {
// 		// case of even length
// 		firstMid := nums[sumLen/2-1]
// 		secondMid := nums[sumLen/2]
// 		return float64(firstMid+secondMid) / 2.0
// 	}

// 	// case of odd length
// 	mid := math.Floor(float64(sumLen) / 2.0)
// 	return float64(nums[int(mid)])
// }

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	sumLen := len(nums1) + len(nums2)
	cap := int(math.Floor(float64(sumLen)/2.0)) + 1
	nums := make([]int, 0, cap)

	i, j := 0, 0
	for i+j < cap {
		if i < len(nums1) && (j >= len(nums2) || nums1[i] <= nums2[j]) {
			nums = append(nums, nums1[i])
			i++
		} else {
			nums = append(nums, nums2[j])
			j++
		}
	}

	l := len(nums)
	if sumLen%2 == 0 {
		// case of even length
		firstMid := nums[l-2]
		secondMid := nums[l-1]
		return float64(firstMid+secondMid) / 2.0
	}
	// case of odd length
	mid := l - 1
	return float64(nums[mid])
}
