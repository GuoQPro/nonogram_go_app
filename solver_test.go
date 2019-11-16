package main

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

var row_ind = [][]int{
	{5}, {2, 1}, {2, 1}, {1}, {1, 1},
}

var col_ind = [][]int{
	{1, 1, 1}, {4}, {2, 1}, {1}, {3},
}

func TestIsSoluable(t *testing.T) {
	ok, answer := IsSoluable(row_ind, col_ind)

	if ok {
		t.Log("Well done: ", answer)
	} else {
		//fmt.Println("What a pity")
		t.Errorf("What a pity")
	}
}
