package main

func matrixReshape(mat [][]int, r int, c int) [][]int {
	m := len(mat)

	if m == 0 {
		if r == 0 && c == 0 {
			return mat
		}
		return mat
	}
	n := len(mat[0])

	if m*n != r*c {
		return mat
	}

	reshapedMat := make([][]int, r)
	for i := range reshapedMat {
		reshapedMat[i] = make([]int, c)
	}

	currentRow := 0
	currentCol := 0

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			reshapedMat[currentRow][currentCol] = mat[i][j]

			currentCol++
			if currentCol == c {
				currentCol = 0
				currentRow++
			}
		}
	}

	return reshapedMat
}
