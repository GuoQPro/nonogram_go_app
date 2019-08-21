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
