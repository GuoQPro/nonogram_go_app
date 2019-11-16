package main

import "math/rand"
import "time"
import "fmt"

const (
	PUZZLE_VALUE_EXIST     = 1
	PUZZLE_VALUE_NULL      = 0
	PUZZLE_VALUE_NOT_EXIST = 2
)

type Row []int
type Puzzle []Row

func (r *Row) String() string {
	result := "|"
	for i := range *r {
		if (*r)[i] == PUZZLE_VALUE_EXIST {
			result += " x"
		} else {
			result += "  "
		}
	}
	result += " |"

	return result
}

func (p *Puzzle) String() string {
	result := ""
	topLine := "---"
	bottomLine := "---"

	for r := range *p {
		topLine += "--"
		bottomLine += "--"
		result += fmt.Sprintf("\n%s", &((*p)[r]))
	}
	result = topLine + result + "\n" + bottomLine
	return result
}

func GetPuzzle(row int, col int) Puzzle {
	for {
		p := GeneratePuzzle(row, col)
		row_ind, col_ind := CalcIndicator(p)
		ok, result := IsSoluable(row_ind, col_ind)

		if ok {
			fmt.Println(&result)
			return p
		} else {
			//fmt.Println("What a pity")
		}
	}
}

func GeneratePuzzle(row int, col int) Puzzle {
	seed := rand.NewSource(time.Now().UnixNano())
	new_rand := rand.New(seed)
	p := make(Puzzle, row)

	for r := 0; r < row; r++ {
		p[r] = make([]int, col)
		for c := 0; c < col; c++ {
			p[r][c] = new_rand.Intn(2)
		}
	}

	return p
}

func CalcIndicator(puzzle Puzzle) ([][]int, [][]int) {
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
