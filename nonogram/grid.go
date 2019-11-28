package nonogram

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	//"image/color"
	//"log"
)

// GridState is determined by grid value
type GridState int

const (
	gridStateNull GridState = iota
	gridStateMarkExist
	gridStateMarkNotExist
)

// Grid is the data structure of grid
type Grid struct {
	value GridState
	posX  float64
	posY  float64
	hint  bool
}

// NewGrid creates a new grid instance.
func NewGrid(x float64, y float64) *Grid {
	return &Grid{
		value: gridStateNull,
		posX:  x,
		posY:  y,
		hint:  false,
	}
}

// SetValue is the value setter.
func (g *Grid) SetValue(v GridState) {
	g.value = v
}

// GetValue is the value getter.
func (g *Grid) GetValue() GridState {
	return g.value
}

// Hint make the grid display as hinting.
func (g *Grid) Hint() {
	g.hint = true
}

// Draw defines the method to draw a grid.
func (g *Grid) Draw(screen *ebiten.Image) error {

	ebitenutil.DrawRect(screen, g.posX, g.posY, gridWidth, gridHeight, colorWhite)

	innerOffsetW := (gridWidth - innerGridW) * 0.5
	innerOffsetH := (gridHeight - innerGridH) * 0.5

	if g.value == gridStateMarkExist {
		ebitenutil.DrawRect(screen, g.posX+innerOffsetW, g.posY+innerOffsetH, innerGridW, innerGridH, colorBlack)
	} else if g.value == gridStateMarkNotExist {
		ebitenutil.DrawRect(screen, g.posX+innerOffsetW, g.posY+innerOffsetH, innerGridW, innerGridH, colorRed)
	} else if g.value == gridStateNull {
		if g.hint {
			ebitenutil.DrawRect(screen, g.posX+innerOffsetW, g.posY+innerOffsetH, innerGridW, innerGridH, colorBlack)
		}
	}

	return nil
}

// OnLeftClick is the handler of left button click event.
func (g *Grid) OnLeftClick() error {
	if g.value == gridStateMarkExist {
		g.value = gridStateNull
	} else if g.value == gridStateMarkNotExist {
		g.value = gridStateNull
	} else if g.value == gridStateNull {
		g.value = gridStateMarkExist
	}
	return nil
}

// OnRightClick is the handler of right button click event.
func (g *Grid) OnRightClick() error {
	if g.value == gridStateMarkExist {
		g.value = gridStateMarkNotExist
	} else if g.value == gridStateMarkNotExist {
		g.value = gridStateNull
	} else if g.value == gridStateNull {
		g.value = gridStateMarkNotExist
	}
	return nil
}

// OnLeftDragOn is the handler of left button drag event.
func (g *Grid) OnLeftDragOn() error {
	if g.value == gridStateMarkExist {

	} else if g.value == gridStateMarkNotExist {
		g.value = gridStateMarkExist
	} else if g.value == gridStateNull {
		g.value = gridStateMarkExist
	}
	return nil
}

// OnRightDragOn is the handler of right button drag event.
func (g *Grid) OnRightDragOn() error {
	if g.value == gridStateMarkExist {
		g.value = gridStateMarkNotExist
	} else if g.value == gridStateMarkNotExist {

	} else if g.value == gridStateNull {
		g.value = gridStateMarkNotExist
	}
	return nil
}

// IsSameGrid determines if the givin grid is the same grid.
func (g *Grid) IsSameGrid(rh *Grid) bool {
	if g.posX == rh.posX && g.posY == rh.posY {
		return true
	}
	return false
}
