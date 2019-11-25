package puzzle

import (
	//"fmt"
	"testing"
)

var test = [][]int{
	{1, 1, 1, 1, 1},
	{0, 1, 1, 0, 1},
	{1, 1, 0, 0, 1},
	{0, 1, 0, 0, 0},
	{1, 0, 1, 0, 0},
}

var rowInd = [][]int{
	{5}, {2, 1}, {2, 1}, {1}, {1, 1},
}

var colInd = [][]int{
	{1, 1, 1}, {4}, {2, 1}, {1}, {3},
}

func TestIsSoluable(t *testing.T) {
	ok, answer := IsSoluable(rowInd, colInd)

	if ok {
		t.Log("Well done: ", answer)
	} else {
		//fmt.Println("What a pity")
		t.Errorf("What a pity")
	}
}
