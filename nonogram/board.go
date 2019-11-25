package nonogram

import (
	"fmt"

	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"

	//"log"
	"math"
)

type Board struct {
	grids  [][]*Grid
	rowInd [][]int
	colInd [][]int
	startX float64
	startY float64
	width  float64
	height float64
	bound  Bound
}

var (
	gapWidth   float64
	gapHeight  float64
	gridWidth  float64
	gridHeight float64
	innerGridW float64
	innerGridH float64
)

const minGridSize = 20

func NewBoard(puzzle Puzzle, bound Bound) *Board {
	board := &Board{}
	board.InitBoard(puzzle, bound)
	return board
}

func (b *Board) InitBoard(puzzle Puzzle, bound Bound) {
	b.bound = bound
	b.CalcIndicator(puzzle)

	row := len(puzzle)
	col := len(puzzle[0])

	b.CentralizeBoard(row, col)

	b.grids = make([][]*Grid, row)

	curY := b.startY
	for r := 0; r < row; r++ {
		curX := b.startX
		b.grids[r] = make([]*Grid, col)
		for c := 0; c < col; c++ {
			b.grids[r][c] = NewGrid(float64(curX), float64(curY))
			curX += (gridWidth + gapWidth)
		}
		curY += (gridHeight + gapHeight)
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

	gridTotalW = math.Min(minGridSize, gridTotalW)
	gridTotalH = math.Min(minGridSize, gridTotalH)

	gridWidth = math.Ceil(gridTotalW * 0.9)
	gridHeight = math.Ceil(gridTotalH * 0.9)
	innerGridW = math.Ceil(gridWidth * 0.8)
	innerGridH = math.Ceil(gridHeight * 0.8)
	gapWidth = math.Ceil(gridTotalW * 0.1)
	gapHeight = math.Ceil(gridTotalH * 0.1)

	b.width = float64(col)*gridWidth + float64(col-1)*gapWidth
	b.height = float64(row)*gridHeight + float64(row-1)*gapHeight

	b.startX = ((b.bound.w - b.width) / 2) + b.bound.x
	b.startY = ((b.bound.h - b.height) / 2) + b.bound.y
}

func (b *Board) DrawIndicators(screen *ebiten.Image) {
	textColor := colorBlack

	startY := b.startY
	startX := b.startX

	curY := startY
	curX := startX

	rowIndGap := 5.0
	for row := range b.rowInd {
		curX = startX
		len := len(b.rowInd[row])
		for i := range b.rowInd[row] {
			str := fmt.Sprintf("%d|", b.rowInd[row][len-i-1])
			bound, _ := font.BoundString(textFont, str)
			w := float64((bound.Max.X - bound.Min.X).Ceil())
			h := float64((bound.Max.Y - bound.Min.Y).Ceil())
			curX = curX - (w + rowIndGap)
			textY := int(curY + h/2 + gridHeight/2)
			text.Draw(screen, str, textFont, int(curX), textY, textColor)
		}
		curY += (gridHeight + gapHeight)
	}

	curX = startX

	for col := range b.colInd {
		indNum := len(b.colInd[col])

		for i := range b.colInd[col] {
			str := fmt.Sprintf(" %d", b.colInd[col][i])
			if b.colInd[col][i] >= 10 {
				str = fmt.Sprintf("%d", b.colInd[col][i])
			}
			bound, _ := font.BoundString(textFont, str)
			//w := float64((bound.Max.X - bound.Min.X).Ceil())
			h := float64((bound.Max.Y - bound.Min.Y).Ceil())
			curY := startY - float64(indNum-i)*(gridHeight+gapHeight)*0.75
			textX := int(curX - gapWidth*0.5) // + w/2)
			textY := int(curY + h)
			text.Draw(screen, str, textFont, textX, textY, textColor)
		}

		curX += (gridWidth + gapWidth)
	}
}

func (b *Board) CalcIndicator(puzzle Puzzle) {
	rowNum := len(puzzle)
	colNum := len(puzzle[0])

	row, col := 0, 0

	b.rowInd = [][]int{}
	for row = 0; row < rowNum; row++ {
		curRowInd := []int{}
		curInd := 0
		for col = 0; col < colNum; col++ {
			value := puzzle[row][col]
			if value == puzzleValueExist {
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

		b.rowInd = append(b.rowInd, curRowInd)
	}

	b.colInd = [][]int{}
	for col = 0; col < colNum; col++ {
		curColInd := []int{}
		curInd := 0

		for row = 0; row < rowNum; row++ {
			value := puzzle[row][col]
			if value == puzzleValueExist {
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

		b.colInd = append(b.colInd, curColInd)
	}
}

func (b *Board) DrawBoard(screen *ebiten.Image) error {
	ebitenutil.DrawRect(screen, 0, 0, float64(stageWidth), float64(stageHeight), color.RGBA{255, 255, 255, 255})
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
	rowNum := len(b.grids)
	colNum := len(b.grids[0])

	colIndex := int(math.Floor((float64(x) - b.startX) / (gridWidth + gapWidth)))
	rowIndex := int(math.Floor((float64(y) - b.startY) / (gridHeight + gapHeight)))

	if rowIndex >= 0 && rowIndex < rowNum && colIndex >= 0 && colIndex < colNum {
		return b.grids[rowIndex][colIndex], nil
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

func (b *Board) OnLeftDrag(curX int, curY int, startX int, startY int) error {
	curGrid, curErr := b.GetGridByPos(curX, curY)

	if curErr != nil {
		return curErr
	}

	startGrid, startErr := b.GetGridByPos(startX, startY)

	if startErr != nil {
		return startErr
	}

	if curGrid.IsSameGrid(startGrid) {
		return nil
	}

	curGrid.OnLeftDragOn()

	return nil
}

func (b *Board) OnRightDrag(curX int, curY int, startX int, startY int) error {
	curGrid, curErr := b.GetGridByPos(curX, curY)

	if curErr != nil {
		return curErr
	}

	startGrid, startErr := b.GetGridByPos(startX, startY)

	if startErr != nil {
		return startErr
	}

	if curGrid.IsSameGrid(startGrid) {
		return nil
	}

	curGrid.OnRightDragOn()

	return nil
}
