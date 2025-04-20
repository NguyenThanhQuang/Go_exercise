package main

func maxCount(m int, n int, ops [][]int) int {
	min_row := m
	min_col := n

	for _, op := range ops {
		if op[0] < min_row {
			min_row = op[0]
		}
		if op[1] < min_col {
			min_col = op[1]
		}
	}

	return min_row * min_col
}
