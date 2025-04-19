package main

func isPalindrome(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	reversedHalf := 0
	for x > reversedHalf {
		digit := x % 10
		reversedHalf = reversedHalf*10 + digit
		x /= 10
	}
	return x == reversedHalf || x == reversedHalf/10
}

