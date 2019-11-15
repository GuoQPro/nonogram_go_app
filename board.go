package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	//"log"
	"math"
)

type Board struct {
	grids   [][]*Grid
	row_ind [][]int
	col_ind [][]int
	start_x float64
	start_y float64
	width   float64
	height  float64
	bound   Bound
}

var (
	gap_w        float64
	gap_h        float64
	grid_w       float64
	grid_h       float64
	inner_grid_w float64
	inner_grid_h float64
)

const MIN_GRID_SIZE = 20

func NewBoard(puzzle [][]int, bound Bound) *Board {
	board := &Board{}
	board.InitBoard(puzzle, bound)
	return board
}

func (b *Board) InitBoard(puzzle [][]int, bound Bound) {
	b.bound = bound
	b.CalcIndicator(puzzle)

	row := len(puzzle)
	col := len(puzzle[0])

	b.CentralizeBoard(row, col)

	b.grids = make([][]*Grid, row)

	cur_y := b.start_y
	for r := 0; r < row; r++ {
		cur_x := b.start_x
		b.grids[r] = make([]*Grid, col)
		for c := 0; c < col; c++ {
			b.grids[r][c] = NewGrid(float64(cur_x), float64(cur_y))
			cur_x += (grid_w + gap_w)
		}
		cur_y += (grid_h + gap_h)
	}
}

func (b *Board) CentralizeBoard(row int, col int) {
	gridTotalW := math.Ceil(b.bound.w / float64(row))
	gridTotalH := math.Ceil(b.bound.h / float64(col))

	if gridTotalW > gridTotalH {
		gridTotalW = gridTotalH
	} else {
		gridTotalH = gridTotalW
	}

	gridTotalW = math.Min(MIN_GRID_SIZE, gridTotalW)
	gridTotalH = math.Min(MIN_GRID_SIZE, gridTotalH)

	grid_w = math.Ceil(gridTotalW * 0.9)
	grid_h = math.Ceil(gridTotalH * 0.9)
	inner_grid_w = math.Ceil(grid_w * 0.8)
	inner_grid_h = math.Ceil(grid_h * 0.8)
	gap_w = math.Ceil(gridTotalW * 0.1)
	gap_h = math.Ceil(gridTotalH * 0.1)

	b.width = float64(col)*grid_w + float64(col-1)*gap_w
	b.height = float64(row)*grid_h + float64(row-1)*gap_h

	b.start_x = ((b.bound.w - b.width) / 2) + b.bound.x
	b.start_y = ((b.bound.h - b.height) / 2) + b.bound.y
}

func (b *Board) DrawIndicators(screen *ebiten.Image) {
	textColor := color_black

	start_y := b.start_y
	start_x := b.start_x

	cur_y := start_y
	cur_x := start_x

	row_ind_gap := 10.0
	for row := range b.row_ind {
		cur_x = start_x
		len := len(b.row_ind[row])
		for i := range b.row_ind[row] {
			str := fmt.Sprintf("%d", b.row_ind[row][len-i-1])
			bound, _ := font.BoundString(textFont, str)
			w := float64((bound.Max.X - bound.Min.X).Ceil())
			h := float64((bound.Max.Y - bound.Min.Y).Ceil())
			cur_x = cur_x - (w + row_ind_gap)
			text_y := int(cur_y + h/2 + grid_h/2)
			text.Draw(screen, str, textFont, int(cur_x), text_y, textColor)
		}
		cur_y += (grid_h + gap_h)
	}

	cur_x = start_x

	for col := range b.col_ind {
		ind_num := len(b.col_ind[col])

		for i := range b.col_ind[col] {
			str := fmt.Sprintf(" %d ", b.col_ind[col][i])
			if b.col_ind[col][i] >= 10 {
				str = fmt.Sprintf("%d", b.col_ind[col][i])
			}
			bound, _ := font.BoundString(textFont, str)
			//w := float64((bound.Max.X - bound.Min.X).Ceil())
			h := float64((bound.Max.Y - bound.Min.Y).Ceil())
			cur_y := start_y - float64(ind_num-i)*(grid_h+gap_h)*0.75
			text_x := int(cur_x) // + w/2)
			text_y := int(cur_y + h)
			text.Draw(screen, str, textFont, text_x, text_y, textColor)
		}

		cur_x += (grid_w + gap_w)
	}
}

func (b *Board) CalcIndicator(puzzle [][]int) {
	row_num := len(puzzle)
	col_num := len(puzzle[0])

	row, col := 0, 0

	b.row_ind = [][]int{}
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

		b.row_ind = append(b.row_ind, cur_row_ind)
	}

	b.col_ind = [][]int{}
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

		b.col_ind = append(b.col_ind, cur_col_ind)
	}
}

func (b *Board) DrawBoard(screen *ebiten.Image) error {
	ebitenutil.DrawRect(screen, 0, 0, float64(STAGE_W), float64(STAGE_H), color.RGBA{255, 255, 255, 255})
	b.DrawIndicators(screen)
	for row := range b.grids {
		for col := range b.grids[row] {
			if err := b.grids[row][col].Draw(screen); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Board) GetGridByPos(x int, y int) (*Grid, error) {
	row_num := len(b.grids)
	col_num := len(b.grids[0])

	col_index := int(math.Floor((float64(x) - b.start_x) / (grid_w + gap_w)))
	row_index := int(math.Floor((float64(y) - b.start_y) / (grid_h + gap_h)))

	if row_index >= 0 && row_index < row_num && col_index >= 0 && col_index < col_num {
		return b.grids[row_index][col_index], nil
	}

	return nil, &nonogramErr{desc: "click outside board"}
}

func (b *Board) OnLeftClick(x int, y int) error {
	grid, err := b.GetGridByPos(x, y)

	if err == nil {
		grid.OnLeftClick()
	}

	return err
}

func (b *Board) OnRightClick(x int, y int) error {
	grid, err := b.GetGridByPos(x, y)

	if err == nil {
		grid.OnRightClick()
	}

	return err
}

func (b *Board) OnLeftDrag(x_cur int, y_cur int, x_init int, y_init int) error {
	cur_grid, cur_err := b.GetGridByPos(x_cur, y_cur)

	if cur_err != nil {
		return cur_err
	}

	init_grid, init_err := b.GetGridByPos(x_init, y_init)

	if init_err != nil {
		return init_err
	}

	if cur_grid.IsSameGrid(init_grid) {
		return nil
	}

	cur_grid.OnLeftDragOn()

	return nil
}

func (b *Board) OnRightDrag(x_cur int, y_cur int, x_init int, y_init int) error {
	cur_grid, cur_err := b.GetGridByPos(x_cur, y_cur)

	if cur_err != nil {
		return cur_err
	}

	init_grid, init_err := b.GetGridByPos(x_init, y_init)

	if init_err != nil {
		return init_err
	}

	if cur_grid.IsSameGrid(init_grid) {
		return nil
	}

	cur_grid.OnRightDragOn()

	return nil
}
