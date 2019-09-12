package main

import (
//"fmt"
)

/*
The methodology is to iterate the table row by row and then column by column.
For each row(column), enumerate all possible combinations which match the indicator.
If any value of same position in all possible combinations are same(exist or not),
then the value could be considered fixed.

After one full iteration(all rows and columns), if at least one fixed value is added to the table, start a new
iteration with status quo, stop otherwise.

For the final table, is any value is NOT fixed, then the puzzle is unsoluable since there are more than 1 answers
to the given indicators.
*/

func IsSoluable(row_ind [][]int, col_ind [][]int) (bool, [][]int) {
	row_num := len(row_ind)
	col_num := len(col_ind)

	progressed := true

	t := InitTable(row_num, col_num)

	for progressed {
		progressed = false

		for r := 0; r < row_num; r++ {
			cur_row := GetRow(&t, r)
			p, _ := Analyze(cur_row, row_ind[r])
			progressed = (p || progressed)
		}

		for c := 0; c < col_num; c++ {
			cur_col := GetCol(&t, c)
			p, _ := Analyze(cur_col, col_ind[c])
			progressed = (p || progressed)
		}
	}

	for r := 0; r < row_num; r++ {
		for c := 0; c < col_num; c++ {
			if t[r][c] == PUZZLE_VALUE_NULL {
				return false, nil
			}
		}
	}

	return true, t
}

func GetHint(row_ind [][]int, col_ind [][]int, puzzle [][]int) (int, int, int) {
	row_num := len(row_ind)
	col_num := len(col_ind)

	t := make([][]int, row_num)

	for i := 0; i < row_num; i++ {
		t[i] = make([]int, col_num)
		copy(t[i], puzzle[i])
	}

	for r := 0; r < row_num; r++ {
		cur_row := GetRow(&t, r)
		progressed, new_value_index := Analyze(cur_row, row_ind[r])
		if progressed && len(new_value_index) > 0 {
			return r, new_value_index[0], PUZZLE_VALUE_EXIST
		}
	}

	for c := 0; c < col_num; c++ {
		cur_col := GetCol(&t, c)
		progressed, new_value_index := Analyze(cur_col, col_ind[c])
		if progressed && len(new_value_index) > 0 {
			return new_value_index[0], c, PUZZLE_VALUE_EXIST
		}
	}

	return 0, 0, 0
}

// Get references of a specific column for modification purpose.
func GetCol(t *[][]int, col int) []*int {
	col_data := []*int{}

	for row_index := range *t {
		col_data = append(col_data, &(*t)[row_index][col])
	}

	return col_data
}

// Get references of a specific row for modification purpose.
func GetRow(t *[][]int, row int) []*int {
	row_data := []*int{}
	col_num := len((*t)[0])

	for col_index := 0; col_index < col_num; col_index++ {
		row_data = append(row_data, &(*t)[row][col_index])
	}

	return row_data
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
	if v == PUZZLE_VALUE_EXIST {
		ind[len(ind)-1]++
	} else if v == PUZZLE_VALUE_NOT_EXIST {
		if ind[len(ind)-1] != 0 {
			ind = append(ind, 0)
		}
	}

	return ind
}

func ValidateCandidate(c *Candidate, ind []int, is_final bool) bool {
	cand_ind_len := len(c.ind)

	if is_final {
		if cand_ind_len < len(ind) {
			return false
		}
	}
	for i := 0; i < cand_ind_len; i++ {
		// is last one
		if i == cand_ind_len-1 {
			if i < len(ind) {
				if is_final {
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
	data_len := len(data)
	pre_list := &[]Candidate{}

	for data_index := 0; data_index < data_len; data_index++ {
		cur_value := *((data)[data_index])
		candidates_len := len(*pre_list)
		new_cadidates := []Candidate{}
		is_last_data := (data_index == data_len-1)

		if cur_value == PUZZLE_VALUE_NULL {
			if candidates_len == 0 {
				new_cadidates = append(new_cadidates, CreateCandidate(PUZZLE_VALUE_EXIST, nil))
				new_cadidates = append(new_cadidates, CreateCandidate(PUZZLE_VALUE_NOT_EXIST, nil))
			} else {
				for c_index := 0; c_index < candidates_len; c_index++ {

					new_c1 := CreateCandidate(PUZZLE_VALUE_EXIST, &(*pre_list)[c_index])

					// Early stop goes into effect here. Any candidate which violated the given indicator
					// will be filtered ASAP.
					if ValidateCandidate(&new_c1, ind, is_last_data) {
						new_cadidates = append(new_cadidates, new_c1)
					}

					new_c2 := CreateCandidate(PUZZLE_VALUE_NOT_EXIST, &(*pre_list)[c_index])
					if ValidateCandidate(&new_c2, ind, is_last_data) {
						new_cadidates = append(new_cadidates, new_c2)
					}
				}
			}
		} else {
			if candidates_len == 0 {
				new_cadidates = append(new_cadidates, CreateCandidate(cur_value, nil))
			} else {
				for c_index := 0; c_index < candidates_len; c_index++ {
					new_c3 := CreateCandidate(cur_value, &(*pre_list)[c_index])
					if ValidateCandidate(&new_c3, ind, is_last_data) {
						new_cadidates = append(new_cadidates, new_c3)
					}
				}
			}
		}
		pre_list = &new_cadidates
	}

	final_result := make([]int, data_len)
	if len(*pre_list) > 0 {
		// check if any col in all candidates are the same
		for i := 0; i < data_len; i++ {
			is_same := true
			for c := 1; c < len(*pre_list); c++ {
				if (*pre_list)[0].data[i] != (*pre_list)[c].data[i] {
					is_same = false
				}
			}

			if is_same {
				final_result[i] = (*pre_list)[0].data[i]
			}
		}
	} else {
		// already solved.
		return false, []int{}
	}

	// check if progressed
	has_changed := false
	new_exist_value_index := []int{}
	for i := 0; i < data_len; i++ {
		if final_result[i] != PUZZLE_VALUE_NULL {
			if final_result[i] != *(data)[i] {
				has_changed = true
				*(data)[i] = final_result[i]

				if final_result[i] == PUZZLE_VALUE_EXIST {
					new_exist_value_index = append(new_exist_value_index, i)
				}
			}
		}
	}

	return has_changed, new_exist_value_index
}

/*var test = [][]int{
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
}*/
