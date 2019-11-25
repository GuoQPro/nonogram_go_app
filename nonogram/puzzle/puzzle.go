package puzzle

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	PuzzleValueNull     = 0
	PuzzleValueExist    = 1
	PuzzleValueNotExist = 2
)

type Row []int

//Puzzle data structure of a Puzzle
type Puzzle []Row

func (r *Row) String() string {
	result := "|"
	for i := range *r {
		if (*r)[i] == PuzzleValueExist {
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
		rowInd, colInd := CalcIndicator(p)
		ok, result := IsSoluable(rowInd, colInd)

		if ok {
			fmt.Println(&result)
			return p
		}
	}
}

func GeneratePuzzle(row int, col int) Puzzle {
	seed := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(seed)
	p := make(Puzzle, row)

	for r := 0; r < row; r++ {
		p[r] = make([]int, col)
		for c := 0; c < col; c++ {
			p[r][c] = newRand.Intn(2)
		}
	}

	return p
}

func CalcIndicator(puzzle Puzzle) ([][]int, [][]int) {
	rowNum := len(puzzle)
	colNum := len(puzzle[0])

	row, col := 0, 0

	rowInd := [][]int{}
	for row = 0; row < rowNum; row++ {
		curRowInd := []int{}
		curInd := 0
		for col = 0; col < colNum; col++ {
			value := puzzle[row][col]
			if value == PuzzleValueExist {
				curInd++
			} else {
				if curInd != 0 {
					curRowInd = append(curRowInd, curInd)
					curInd = 0
				}
			}
		}

		if curInd > 0 {
			curRowInd = append(curRowInd, curInd)
		}

		if len(curRowInd) == 0 {
			curRowInd = append(curRowInd, 0)
		}

		rowInd = append(rowInd, curRowInd)
	}

	colInd := [][]int{}
	for col = 0; col < colNum; col++ {
		curColInd := []int{}
		curInd := 0

		for row = 0; row < rowNum; row++ {
			value := puzzle[row][col]
			if value == PuzzleValueExist {
				curInd++
			} else {
				if curInd != 0 {
					curColInd = append(curColInd, curInd)
					curInd = 0
				}
			}
		}

		if curInd > 0 {
			curColInd = append(curColInd, curInd)
		}

		if len(curColInd) == 0 {
			curColInd = append(curColInd, 0)
		}

		colInd = append(colInd, curColInd)
	}

	return rowInd, colInd
}
