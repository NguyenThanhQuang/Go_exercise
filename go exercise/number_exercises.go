package main

import (
	"fmt"
	"math"
)

func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func Factorial(n int) int {
	if n < 0 {
		return -1
	}
	if n == 0 || n == 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func Fibonacci(n int) int {
	if n <= 0 {
		return -1
	}
	if n == 1 || n == 2 {
		return 1
	}
	a, b := 1, 1
	for i := 3; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func IsPerfect(n int) bool {
	if n <= 1 {
		return false
	}
	sum := 1
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			sum += i
			if i != n/i {
				sum += n / i
			}
		}
	}
	return sum == n
}

func SumOfDigits(n int) int {
	sum := 0
	for n != 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

func main() {
	var choice int
	fmt.Println("Chọn bài tập để chạy:")
	fmt.Println("1. Kiểm tra số nguyên tố")
	fmt.Println("2. Tính giai thừa")
	fmt.Println("3. Tìm số Fibonacci thứ n")
	fmt.Println("4. Kiểm tra số hoàn hảo")
	fmt.Println("5. Tính tổng các chữ số")
	fmt.Print("Nhập lựa chọn (1-5): ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nBài tập 1: Kiểm tra số nguyên tố")
		fmt.Println(IsPrime(7))
		fmt.Println(IsPrime(10))
		fmt.Println(IsPrime(13))
		fmt.Println(IsPrime(1))
	case 2:
		fmt.Println("\nBài tập 2: Tính giai thừa")
		fmt.Println(Factorial(5))
		fmt.Println(Factorial(3))
		fmt.Println(Factorial(0))
		fmt.Println(Factorial(-1))
	case 3:
		fmt.Println("\nBài tập 3: Tìm số Fibonacci thứ n")
		fmt.Println(Fibonacci(7))
		fmt.Println(Fibonacci(10))
		fmt.Println(Fibonacci(1))
		fmt.Println(Fibonacci(2))
	case 4:
		fmt.Println("\nBài tập 4: Kiểm tra số hoàn hảo")
		fmt.Println(IsPerfect(6))
		fmt.Println(IsPerfect(28))
		fmt.Println(IsPerfect(10))
		fmt.Println(IsPerfect(496))
	case 5:
		fmt.Println("\nBài tập 5: Tính tổng các chữ số")
		fmt.Println(SumOfDigits(123))
		fmt.Println(SumOfDigits(456))
		fmt.Println(SumOfDigits(1000))
		fmt.Println(SumOfDigits(0))
	default:
		fmt.Println("Lựa chọn không hợp lệ!")
	}
} 