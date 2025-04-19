package main

import (
	"fmt"
	"strings"
	"unicode"
)

func CountChar(s string, char rune) int {
	count := 0
	for _, c := range s {
		if c == char {
			count++
		}
	}
	return count
}

func IsPalindrome(s string) bool {
	runes := []rune(s)
	length := len(runes)
	
	for i := 0; i < length/2; i++ {
		if runes[i] != runes[length-1-i] {
			return false
		}
	}
	return true
}

func RemoveWhitespace(s string) string {
	var result []rune
	for _, char := range s {
		if !unicode.IsSpace(char) {
			result = append(result, char)
		}
	}
	return string(result)
}

func ContainsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}

func main1() {
	var choice int
	fmt.Println("Chọn bài tập để chạy:")
	fmt.Println("1. Đếm số lần xuất hiện của ký tự")
	fmt.Println("2. Kiểm tra chuỗi đối xứng")
	fmt.Println("3. Loại bỏ khoảng trắng")
	fmt.Println("4. Kiểm tra chuỗi con")
	fmt.Print("Nhập lựa chọn (1-4): ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nBài tập 1: Đếm số lần xuất hiện của ký tự")
		fmt.Println(CountChar("hello", 'l')) 
	case 2:
		fmt.Println("\nBài tập 2: Kiểm tra chuỗi đối xứng")
		fmt.Println(IsPalindrome("madam"))   
		fmt.Println(IsPalindrome("hello"))   
		fmt.Println(IsPalindrome("racecar")) 
		fmt.Println(IsPalindrome("golang")) 
	case 3:
		fmt.Println("\nBài tập 3: Loại bỏ khoảng trắng")
		fmt.Println(RemoveWhitespace("hello world"))      
		fmt.Println(RemoveWhitespace("  go  lang  "))    
		fmt.Println(RemoveWhitespace("  spaces  \t\n"))   
	case 4:
		fmt.Println("\nBài tập 4: Kiểm tra chuỗi con")
		fmt.Println(ContainsSubstring("hello", "ell"))    
		fmt.Println(ContainsSubstring("hello", "world"))  
		fmt.Println(ContainsSubstring("golang", "go"))    
		fmt.Println(ContainsSubstring("golang", "java"))  
	default:
		fmt.Println("Lựa chọn không hợp lệ!")
	}
} 