package main

func isValid(s string) bool {
	if len(s)%2 != 0 {
		return false
	}

	var stack []byte

	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	for i := 0; i < len(s); i++ {
		char := s[i]

		if opener, isCloser := pairs[char]; isCloser {
			if len(stack) == 0 || stack[len(stack)-1] != opener {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, char)
		}
	}
	return len(stack) == 0
}