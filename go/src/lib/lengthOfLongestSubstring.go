package lib

func LengthOfLongestSubstring(s string) int {
	m := make(map[rune]int)
	maxLength := 0
	start := 0

	for i, c := range s {
		// sliding window
		if lastIdx, exist := m[c]; exist && lastIdx >= start {
			start = lastIdx + 1
		}

		// update character index
		m[c] = i

		// update max length
		currentLength := i - start + 1
		if currentLength > maxLength {
			maxLength = currentLength
		}
	}

	return maxLength
}
