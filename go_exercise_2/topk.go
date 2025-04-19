package main

func topKFrequent(nums []int, k int) []int {
	freqMap := make(map[int]int)
	for _, num := range nums {
		freqMap[num]++
	}

	buckets := make([][]int, len(nums)+1)
	for num, freq := range freqMap {
		buckets[freq] = append(buckets[freq], num)
	}

	result := make([]int, 0, k)

	for i := len(buckets) - 1; i >= 1 && len(result) < k; i-- {
		if len(buckets[i]) > 0 {
			for _, num := range buckets[i] {
				if len(result) < k {
					result = append(result, num)
				} else {
					break
				}
			}
		}
	}

	return result
}
