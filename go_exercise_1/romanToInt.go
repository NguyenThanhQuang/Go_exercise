package main

func romanToInt(s string) int {
	romanMap := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	if len(s) == 0 {
		return 0
	}

	total := romanMap[s[len(s)-1]]

	for i := len(s) - 2; i >= 0; i-- {
		currentValue := romanMap[s[i]]
		previousValue := romanMap[s[i+1]]

		if currentValue < previousValue {
			total -= currentValue
		} else {
			total += currentValue
		}
	}

	return total
}