package main

func findLHS(nums []int) int {
	freqMap := make(map[int]int)
	for _, num := range nums {
		freqMap[num]++
	}

	maxLength := 0

	for num, count := range freqMap {
		countNext, exists := freqMap[num+1]
		if exists {
			currentLength := count + countNext
			if currentLength > maxLength {
				maxLength = currentLength
			}
		}
	}

	return maxLength
}