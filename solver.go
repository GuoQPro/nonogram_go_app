package main

import (
	"fmt"
)

const (
	PUZZLE_VALUE_NOT_EXIST = 2
	PUZZLE_VALUE_EXIST     = 1
	PUZZLE_VALUE_NULL      = 0
)

func IsSoluable(row_ind [][]int, col_ind [][]int) (bool, [][]int) {
	row_num := len(row_ind)
	col_num := len(col_ind)

	progressed := true

	t := InitTable(row_num, col_num)

	for progressed {
		progressed = false

		for r := 0; r < 1; r++ {
			progressed = (Analyze(&t[r], row_ind[r]) || progressed)
		}

		for c := 0; c < col_num; c++ {

		}
	}

	fmt.Println("result = ", t)

	return true, nil
}

func InitTable(row int, col int) [][]int {
	t := make([][]int, row)

	for r := 0; r < row; r++ {
		t[r] = make([]int, col)
		for c := 0; c < col; c++ {
			t[r][c] = PUZZLE_VALUE_NULL
		}
	}

	return t
}

type Candidate struct {
	data        []int
	isAvailable bool
	ind         []int
}

func CreateCandidate(v int, src *Candidate) Candidate {
	c := Candidate{
		isAvailable: true,
		data:        []int{},
		ind:         []int{},
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
	//fmt.Println("UpdateInd: ", ind)
	if v == PUZZLE_VALUE_EXIST {
		ind[len(ind)-1]++
	} else if v == PUZZLE_VALUE_NOT_EXIST {
		if ind[len(ind)-1] != 0 {
			ind = append(ind, 0)
		}
	}

	return ind
}

func Analyze(data *[]int, ind []int) bool {
	//fmt.Println(data, ind)
	dataLen := len(*data)

	pre_list := &[]Candidate{}

	for data_index := 0; data_index < dataLen; data_index++ {
		cur_value := (*data)[data_index]
		candidatesLen := len(*pre_list)

		newCadidates := []Candidate{}

		if cur_value == PUZZLE_VALUE_NULL {
			if candidatesLen == 0 {
				newCadidates = append(newCadidates, CreateCandidate(PUZZLE_VALUE_EXIST, nil))
				newCadidates = append(newCadidates, CreateCandidate(PUZZLE_VALUE_NOT_EXIST, nil))
			} else {
				for c_index := 0; c_index < candidatesLen; c_index++ {
					if (*pre_list)[c_index].isAvailable {
						newCadidates = append(newCadidates, CreateCandidate(PUZZLE_VALUE_EXIST, &(*pre_list)[c_index]))
						newCadidates = append(newCadidates, CreateCandidate(PUZZLE_VALUE_NOT_EXIST, &(*pre_list)[c_index]))
					}
				}
			}
		} else {
			if candidatesLen == 0 {
				newCadidates = append(newCadidates, CreateCandidate(cur_value, nil))
			} else {
				for c_index := 0; c_index < candidatesLen; c_index++ {
					if (*pre_list)[c_index].isAvailable {
						newCadidates = append(newCadidates, CreateCandidate(cur_value, &(*pre_list)[c_index]))
					}
				}
			}
		}

		pre_list = &newCadidates

	}

	for i := 0; i < len(*pre_list); i++ {
		fmt.Println((*pre_list)[i])
	}

	fmt.Println("final = ", len(*pre_list))

	return false
}

var test = [][]int{
	{1, 1, 1, 1, 1},
	{0, 1, 1, 0, 1},
	{1, 1, 0, 0, 1},
	{0, 1, 0, 0, 0},
	{1, 0, 1, 0, 0},
}

var row_ind = [][]int{
	{5}, {2, 1}, {2, 1}, {1}, {1, 1},
}

var col_ind = [][]int{
	{1, 1, 1}, {4}, {2, 1}, {1}, {3},
}

func main() {
	ok, answer := IsSoluable(row_ind, col_ind)

	if ok {
		fmt.Println("Well done: ", answer)
	} else {
		fmt.Println("What a pity")
	}
}
