package main

func lengthOfLongestSubstring(s string) int {

	charMap := make(map[byte]int)

	maxLength := 0
	left := 0

	for right := 0; right < len(s); right++ {
		char := s[right]

		lastSeenIndex, exists := charMap[char]

		if exists && lastSeenIndex >= left {
			left = lastSeenIndex + 1
		}

		charMap[char] = right

		currentLength := right - left + 1
		if currentLength > maxLength {
			maxLength = currentLength
		}
	}

	return maxLength
}