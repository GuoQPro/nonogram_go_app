package puzzle

//"fmt"

/*
The methodology is to iterate the table row by row and then column by column.
For each row(column), enumerate all possible combinations which match the indicator.
If any value of same position in all possible combinations are same(exist or not),
then the value could be considered fixed.

After one full iteration(all rows and columns), if at least one fixed value is added to the table, start a new
iteration with status quo, stop otherwise.

For the final table, if any value is NOT fixed, then the puzzle is unsoluable since there are more than 1 answer
to the given indicators.
*/

func IsSoluable(rowInd [][]int, colInd [][]int) (bool, Puzzle) {
	rowNum := len(rowInd)
	colNum := len(colInd)

	progressed := true

	t := InitTable(rowNum, colNum)

	for progressed {
		progressed = false

		for r := 0; r < rowNum; r++ {
			curRow := GetRow(&t, r)
			p, _ := Analyze(curRow, rowInd[r])
			progressed = (p || progressed)
		}

		for c := 0; c < colNum; c++ {
			curCol := GetCol(&t, c)
			p, _ := Analyze(curCol, colInd[c])
			progressed = (p || progressed)
		}
	}

	for r := 0; r < rowNum; r++ {
		for c := 0; c < colNum; c++ {
			if t[r][c] == PuzzleValueNull {
				return false, nil
			}
		}
	}

	return true, t
}

func GetHint(rowInd [][]int, colInd [][]int, puzzle Puzzle) (int, int, int) {
	rowNum := len(rowInd)
	colNum := len(colInd)

	t := make(Puzzle, rowNum)

	for i := 0; i < rowNum; i++ {
		t[i] = make([]int, colNum)
		copy(t[i], puzzle[i])
	}

	for r := 0; r < rowNum; r++ {
		curRow := GetRow(&t, r)
		progressed, newValueIndex := Analyze(curRow, rowInd[r])
		if progressed && len(newValueIndex) > 0 {
			return r, newValueIndex[0], PuzzleValueExist
		}
	}

	for c := 0; c < colNum; c++ {
		curCol := GetCol(&t, c)
		progressed, newValueIndex := Analyze(curCol, colInd[c])
		if progressed && len(newValueIndex) > 0 {
			return newValueIndex[0], c, PuzzleValueExist
		}
	}

	return 0, 0, 0
}

// Get references of a specific column for modification purpose.
func GetCol(t *Puzzle, col int) []*int {
	colData := []*int{}

	for rowIndex := range *t {
		colData = append(colData, &(*t)[rowIndex][col])
	}

	return colData
}

// Get references of a specific row for modification purpose.
func GetRow(t *Puzzle, row int) []*int {
	rowData := []*int{}
	colNum := len((*t)[0])

	for colIndex := 0; colIndex < colNum; colIndex++ {
		rowData = append(rowData, &(*t)[row][colIndex])
	}

	return rowData
}

func InitTable(row int, col int) Puzzle {
	t := make(Puzzle, row)

	for r := 0; r < row; r++ {
		t[r] = make([]int, col)
		for c := 0; c < col; c++ {
			t[r][c] = PuzzleValueNull
		}
	}
	return t
}

type Candidate struct {
	data []int
	ind  []int
}

func CreateCandidate(v int, src *Candidate) Candidate {
	c := Candidate{
		data: []int{},
		ind:  []int{},
	}

	if src != nil {
		c.ind = make([]int, len((*src).ind))
		c.data = make([]int, len((*src).data))
		copy(c.data, (*src).data)
		copy(c.ind, (*src).ind)
		c.data = append(c.data, v)
	} else {
		c.data = append(c.data, v)
		c.ind = append(c.ind, 0)
	}

	c.ind = UpdateInd(c.ind, v)

	return c
}

func UpdateInd(ind []int, v int) []int {
	if v == PuzzleValueExist {
		ind[len(ind)-1]++
	} else if v == PuzzleValueNotExist {
		if ind[len(ind)-1] != 0 {
			ind = append(ind, 0)
		}
	}

	return ind
}

func ValidateCandidate(c *Candidate, ind []int, isFinal bool) bool {
	candIndLen := len(c.ind)

	if isFinal {
		if candIndLen < len(ind) {
			return false
		}
	}
	for i := 0; i < candIndLen; i++ {
		// is last one
		if i == candIndLen-1 {
			if i < len(ind) {
				if isFinal {
					if c.ind[i] != ind[i] {
						return false
					}
				} else {
					if c.ind[i] > ind[i] {
						return false
					}
				}
			} else {
				if c.ind[i] != 0 {
					return false
				}
			}
		} else {
			if i < len(ind) {
				if c.ind[i] != ind[i] {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func Analyze(data []*int, ind []int) (bool, []int) {
	dataLen := len(data)
	preList := &[]Candidate{}

	for dataIndex := 0; dataIndex < dataLen; dataIndex++ {
		curValue := *((data)[dataIndex])
		candidatesLen := len(*preList)
		newCadidates := []Candidate{}
		isLastData := (dataIndex == dataLen-1)

		if curValue == PuzzleValueNull {
			if candidatesLen == 0 {
				newCadidates = append(newCadidates, CreateCandidate(PuzzleValueExist, nil))
				newCadidates = append(newCadidates, CreateCandidate(PuzzleValueNotExist, nil))
			} else {
				for cIndex := 0; cIndex < candidatesLen; cIndex++ {

					newC1 := CreateCandidate(PuzzleValueExist, &(*preList)[cIndex])

					// Early stop goes into effect here. Any candidate solution which does not comply with the given indicator
					// will be filtered ASAP.
					if ValidateCandidate(&newC1, ind, isLastData) {
						newCadidates = append(newCadidates, newC1)
					}

					newC2 := CreateCandidate(PuzzleValueNotExist, &(*preList)[cIndex])
					if ValidateCandidate(&newC2, ind, isLastData) {
						newCadidates = append(newCadidates, newC2)
					}
				}
			}
		} else {
			if candidatesLen == 0 {
				newCadidates = append(newCadidates, CreateCandidate(curValue, nil))
			} else {
				for cIndex := 0; cIndex < candidatesLen; cIndex++ {
					newC3 := CreateCandidate(curValue, &(*preList)[cIndex])
					if ValidateCandidate(&newC3, ind, isLastData) {
						newCadidates = append(newCadidates, newC3)
					}
				}
			}
		}
		preList = &newCadidates
	}

	finalResult := make([]int, dataLen)
	if len(*preList) > 0 {
		// check if any col in all candidates are the same
		for i := 0; i < dataLen; i++ {
			isSame := true
			for c := 1; c < len(*preList); c++ {
				if (*preList)[0].data[i] != (*preList)[c].data[i] {
					isSame = false
				}
			}

			if isSame {
				finalResult[i] = (*preList)[0].data[i]
			}
		}
	} else {
		// already solved.
		return false, []int{}
	}

	// check if progressed
	hasChanged := false
	newExistValueIndex := []int{}
	for i := 0; i < dataLen; i++ {
		if finalResult[i] != PuzzleValueNull {
			if finalResult[i] != *(data)[i] {
				hasChanged = true
				*(data)[i] = finalResult[i]

				if finalResult[i] == PuzzleValueExist {
					newExistValueIndex = append(newExistValueIndex, i)
				}
			}
		}
	}

	return hasChanged, newExistValueIndex
}
