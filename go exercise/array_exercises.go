package main

import (
	"fmt"
	"sort"
)

func FindMinMax(arr []int) (int, int) {
	if len(arr) == 0 {
		return 0, 0
	}
	min, max := arr[0], arr[0]
	for _, num := range arr {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}
	return min, max
}

func SortArray(arr []int, ascending bool) []int {
	sorted := make([]int, len(arr))
	copy(sorted, arr)
	if ascending {
		sort.Ints(sorted)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	}
	return sorted
}

func CountOccurrences(arr []int, value int) int {
	count := 0
	for _, num := range arr {
		if num == value {
			count++
		}
	}
	return count
}

func RemoveElement(arr []int, value int) []int {
	result := []int{}
	for _, num := range arr {
		if num != value {
			result = append(result, num)
		}
	}
	return result
}

func MergeArrays(arr1, arr2 []int) []int {
	result := make([]int, len(arr1)+len(arr2))
	copy(result, arr1)
	copy(result[len(arr1):], arr2)
	return result
}

func main2() {
	var choice int
	fmt.Println("Chọn bài tập để chạy:")
	fmt.Println("1. Tìm giá trị lớn nhất và nhỏ nhất trong mảng")
	fmt.Println("2. Sắp xếp mảng")
	fmt.Println("3. Tìm số lần xuất hiện của một phần tử")
	fmt.Println("4. Xóa phần tử khỏi mảng")
	fmt.Println("5. Gộp hai mảng thành một")
	fmt.Print("Nhập lựa chọn (1-5): ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nBài tập 1: Tìm giá trị lớn nhất và nhỏ nhất trong mảng")
		arr := []int{3, 5, 7, 2, 8}
		min, max := FindMinMax(arr)
		fmt.Printf("Mảng: %v\n", arr)
		fmt.Printf("Giá trị nhỏ nhất: %d\n", min)
		fmt.Printf("Giá trị lớn nhất: %d\n", max)
	case 2:
		fmt.Println("\nBài tập 2: Sắp xếp mảng")
		arr := []int{5, 3, 8, 4, 2}
		fmt.Printf("Mảng ban đầu: %v\n", arr)
		fmt.Printf("Sắp xếp tăng dần: %v\n", SortArray(arr, true))
		fmt.Printf("Sắp xếp giảm dần: %v\n", SortArray(arr, false))
	case 3:
		fmt.Println("\nBài tập 3: Tìm số lần xuất hiện của một phần tử")
		arr := []int{1, 2, 3, 2, 2, 4}
		value := 2
		fmt.Printf("Mảng: %v\n", arr)
		fmt.Printf("Số lần xuất hiện của %d: %d\n", value, CountOccurrences(arr, value))
	case 4:
		fmt.Println("\nBài tập 4: Xóa phần tử khỏi mảng")
		arr := []int{1, 2, 3, 4}
		value := 2
		fmt.Printf("Mảng ban đầu: %v\n", arr)
		fmt.Printf("Mảng sau khi xóa %d: %v\n", value, RemoveElement(arr, value))
	case 5:
		fmt.Println("\nBài tập 5: Gộp hai mảng thành một")
		arr1 := []int{1, 2}
		arr2 := []int{3, 4}
		fmt.Printf("Mảng 1: %v\n", arr1)
		fmt.Printf("Mảng 2: %v\n", arr2)
		fmt.Printf("Mảng sau khi gộp: %v\n", MergeArrays(arr1, arr2))
	default:
		fmt.Println("Lựa chọn không hợp lệ!")
	}
} 