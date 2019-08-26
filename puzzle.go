package main

import "math/rand"
import "time"

const (
	PUZZLE_VALUE_EXIST = 1
	PUZZLE_VALUE_NULL  = 0
)

func GetPuzzle(row int, col int) [][]int {
	/*return [][]int{
		{1, 1, 1, 1, 0},
		{1, 1, 0, 1, 1},
		{1, 1, 1, 1, 0},
		{0, 1, 0, 1, 1},
		{0, 1, 1, 0, 1},
	}*/

	return GeneratePuzzle(row, col)
}

func GeneratePuzzle(row int, col int) [][]int {
	seed := rand.NewSource(time.Now().UnixNano())
	new_rand := rand.New(seed)
	p := make([][]int, row)

	for r := 0; r < row; r++ {
		p[r] = make([]int, col)
		for c := 0; c < col; c++ {
			p[r][c] = new_rand.Intn(2)
		}
	}

	return p
}

//func ValidatePuzzle(p [][]int) bool {
//	row_ind, col_ind = CalcIndicator(p)
//	return true
//}

func CalcIndicator(puzzle [][]int) ([][]int, [][]int) {
	row_num := len(puzzle)
	col_num := len(puzzle[0])

	row, col := 0, 0

	row_ind := [][]int{}
	for row = 0; row < row_num; row++ {
		cur_row_ind := []int{}
		cur_ind := 0
		for col = 0; col < col_num; col++ {
			value := puzzle[row][col]
			if value == PUZZLE_VALUE_EXIST {
				cur_ind += 1
			} else {
				if cur_ind != 0 {
					cur_row_ind = append(cur_row_ind, cur_ind)
					cur_ind = 0
				}
			}
		}

		if cur_ind > 0 {
			cur_row_ind = append(cur_row_ind, cur_ind)
		}

		if len(cur_row_ind) == 0 {
			cur_row_ind = append(cur_row_ind, 0)
		}

		row_ind = append(row_ind, cur_row_ind)
	}

	col_ind := [][]int{}
	for col = 0; col < col_num; col++ {
		cur_col_ind := []int{}
		cur_ind := 0

		for row = 0; row < row_num; row++ {
			value := puzzle[row][col]
			if value == PUZZLE_VALUE_EXIST {
				cur_ind += 1
			} else {
				if cur_ind != 0 {
					cur_col_ind = append(cur_col_ind, cur_ind)
					cur_ind = 0
				}
			}
		}

		if cur_ind > 0 {
			cur_col_ind = append(cur_col_ind, cur_ind)
		}

		if len(cur_col_ind) == 0 {
			cur_col_ind = append(cur_col_ind, 0)
		}

		col_ind = append(col_ind, cur_col_ind)
	}

	return row_ind, col_ind
}
