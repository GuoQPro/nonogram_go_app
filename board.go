package main

import (
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
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
}

const (
	gap_w = 1.0
	gap_h = 1.0
)

type nonogrameErr struct {
	desc string
}

func (e *nonogrameErr) Error() string {
	return e.desc
}

var textFont font.Face

func (b *Board) InitBoard(puzzle [][]int) {

	InitFonts()
	b.CalcIndicator(puzzle)

	b.start_x = float64(STAGE_W) * 0.2
	b.start_y = float64(STAGE_H) * 0.2

	row := len(puzzle)
	col := len(puzzle[0])

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

func InitFonts() {
	tt, _ := truetype.Parse(fonts.MPlus1pRegular_ttf)
	textFont = truetype.NewFace(tt, &truetype.Options{
		Size:    18,
		DPI:     54,
		Hinting: font.HintingFull,
	})
}

func (b *Board) DrawIndicators(screen *ebiten.Image) {
	textColor := color.RGBA{255, 255, 255, 255}

	start_y := int(b.start_y)
	start_x := int(b.start_x)

	cur_y := start_y

	for row := range b.row_ind {
		ind_num := len(b.row_ind[row])
		for i := range b.row_ind[row] {
			str := fmt.Sprintf("%d", b.row_ind[row][i])
			bound, _ := font.BoundString(textFont, str)
			h := (bound.Max.Y - bound.Min.Y).Ceil()
			cur_x := start_x - (ind_num-i)*(grid_w+gap_w)
			text.Draw(screen, str, textFont, cur_x, cur_y+h/2+grid_h/2, textColor)
		}
		cur_y += (grid_h + gap_h)
	}

	cur_x := start_x

	for col := range b.col_ind {
		ind_num := len(b.col_ind[col])

		for i := range b.col_ind[col] {
			str := fmt.Sprintf("%d", b.col_ind[col][i])
			bound, _ := font.BoundString(textFont, str)
			w := (bound.Max.X - bound.Min.X).Ceil()
			cur_y := start_y - (ind_num-i)*(grid_h+gap_h)
			text.Draw(screen, str, textFont, cur_x+w/2, cur_y, textColor)
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

		b.col_ind = append(b.col_ind, cur_col_ind)
	}
}

func (b *Board) DrawBoard(screen *ebiten.Image) error {
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

	return nil, &nonogrameErr{desc: "click outside board"}
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

func (b *Board) OnLeftDrag(x int, y int) error {
	grid, err := b.GetGridByPos(x, y)

	if err == nil {
		grid.OnLeftDragOn()
	}

	return err
}

func (b *Board) OnRightDrag(x int, y int) error {
	grid, err := b.GetGridByPos(x, y)

	if err == nil {
		grid.OnRightDragOn()
	}

	return err
}
