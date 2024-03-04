package lib

import "fmt"

func LongestCommonSubsequence() {
	s := "ABCD"
	t := "ACDE"
	lcs := lcs(s, t)
	fmt.Println(lcs)
}

func lcs(s, t string) string {
	// DPテーブルの作成
	dp := make([][]int, len(s)+1)
	for i := range dp {
		dp[i] = make([]int, len(t)+1)
	}

	// DPテーブルの初期化
	for i := 0; i <= len(s); i++ {
		dp[i][0] = 0
	}
	for i := 0; i <= len(t); i++ {
		dp[0][i] = 0
	}

	fmt.Println(dp)

	// DPテーブルの計算
	for i := 1; i <= len(s); i++ {
		for j := 1; j <= len(t); j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	fmt.Println(dp)

	// LCSの復元
	i := len(s)
	j := len(t)
	lcs := ""
	for i > 0 && j > 0 {
		if s[i-1] == t[j-1] {
			fmt.Println("s", i-1, "t", j-1)
			lcs = string(s[i-1]) + lcs
			i--
			j--
		} else {
			if dp[i-1][j] > dp[i][j-1] {
				i--
			} else {
				j--
			}
		}
	}

	return lcs
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
