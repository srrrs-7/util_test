package lib

func longestPalindrome(s string) string {
	l := len(s)
	if l < 2 {
		return s
	}

	str := ""
	for i := range s {
		// Odd length palindromes
		left, right := i, i
		for left >= 0 && right < l && s[left] == s[right] {
			if right-left+1 > len(str) {
				str = s[left : right+1]
			}
			left--
			right++
		}

		// Even length palindromes
		left, right = i, i+1
		for left >= 0 && right < l && s[left] == s[right] {
			if right-left+1 > len(str) {
				str = s[left : right+1]
			}
			left--
			right++
		}
	}

	return str
}
