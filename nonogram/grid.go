package nonogram

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	//"image/color"
	//"log"
)

type GridState int

const (
	gridStateNull GridState = iota
	gridStateMarkExist
	gridStateMarkNotExist
)

type Grid struct {
	value GridState
	posX  float64
	posY  float64
	hint  bool
}

func NewGrid(x float64, y float64) *Grid {

	return &Grid{
		value: gridStateNull,
		posX:  x,
		posY:  y,
		hint:  false,
	}
}

func (g *Grid) SetValue(v GridState) {
	g.value = v
}

func (g *Grid) GetValue() GridState {
	return g.value
}

func (g *Grid) GetPos() (float64, float64) {
	return g.posX, g.posY
}

func (g *Grid) Hint() {
	g.hint = true
}

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

func (g *Grid) OnLeftDragOn() error {
	if g.value == gridStateMarkExist {

	} else if g.value == gridStateMarkNotExist {
		g.value = gridStateMarkExist
	} else if g.value == gridStateNull {
		g.value = gridStateMarkExist
	}
	return nil
}

func (g *Grid) OnRightDragOn() error {
	if g.value == gridStateMarkExist {
		g.value = gridStateMarkNotExist
	} else if g.value == gridStateMarkNotExist {

	} else if g.value == gridStateNull {
		g.value = gridStateMarkNotExist
	}
	return nil
}

func (g *Grid) IsSameGrid(rh *Grid) bool {
	if g.posX == rh.posX && g.posY == rh.posY {
		return true
	}
	return false
}
