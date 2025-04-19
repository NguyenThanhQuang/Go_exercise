package main

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	for i := 0; i < len(strs[0]); i++ {
		char := strs[0][i]

		for j := 1; j < len(strs); j++ {
			currentStr := strs[j]
			if i >= len(currentStr) || currentStr[i] != char {
				return strs[0][:i]
			}
		}
	}

	return strs[0]
}