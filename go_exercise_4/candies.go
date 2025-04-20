package main

func distributeCandies(candyType []int) int {
	n := len(candyType)
	maxAllowedToEat := n / 2

	uniqueCandyTypes := make(map[int]struct{})

	for _, candy := range candyType {
		uniqueCandyTypes[candy] = struct{}{}
	}

	numUniqueTypes := len(uniqueCandyTypes)

	if numUniqueTypes < maxAllowedToEat {
		return numUniqueTypes
	}

	return maxAllowedToEat
}
