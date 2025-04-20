package main

func permute(nums []int) [][]int {
	var results [][]int
	var currentPermutation []int
	used := make([]bool, len(nums))

	backtrack(nums, currentPermutation, used, &results)

	return results
}

func backtrack(nums []int, currentPermutation []int, used []bool, results *[][]int) {
	if len(currentPermutation) == len(nums) {
		temp := make([]int, len(currentPermutation))
		copy(temp, currentPermutation)
		*results = append(*results, temp)
		return
	}

	for i := 0; i < len(nums); i++ {
		if !used[i] {
			used[i] = true
			currentPermutation = append(currentPermutation, nums[i])

			backtrack(nums, currentPermutation, used, results)

			used[i] = false
			currentPermutation = currentPermutation[:len(currentPermutation)-1]
		}
	}
}