package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	//"log"
)

const (
	grid_w = 15
	grid_h = 15
)

var (
	color_white = color.RGBA{255, 255, 255, 255}
	color_blue  = color.RGBA{0, 0, 255, 255}
	color_red   = color.RGBA{255, 0, 0, 255}
	color_black = color.RGBA{0, 0, 0, 255}
)

type grid_state int

const (
	GRID_NULL grid_state = iota
	GRID_MARK_EXIST
	GRID_MARK_NOTEXIST
)

type Grid struct {
	value grid_state
	pos_x float64
	pos_y float64
}

func NewGrid(x float64, y float64) *Grid {
	return &Grid{
		value: GRID_NULL,
		pos_x: x,
		pos_y: y,
	}
}

func (g *Grid) SetValue(v grid_state) {
	g.value = v
}

func (g *Grid) GetValue() grid_state {
	return g.value
}

func (g *Grid) GetPos() (float64, float64) {
	return g.pos_x, g.pos_y
}

func (g *Grid) Draw(screen *ebiten.Image) error {
	var gridColor color.Color
	if g.value == GRID_MARK_EXIST {
		gridColor = color_black
	} else if g.value == GRID_MARK_NOTEXIST {
		gridColor = color_red
	} else if g.value == GRID_NULL {
		gridColor = color_white
	}

	ebitenutil.DrawRect(screen, g.pos_x, g.pos_y, grid_w, grid_h, gridColor)

	return nil
}

func (g *Grid) OnLeftClick() error {
	if g.value == GRID_MARK_EXIST {
		g.value = GRID_NULL
	} else if g.value == GRID_MARK_NOTEXIST {
		g.value = GRID_NULL
	} else if g.value == GRID_NULL {
		g.value = GRID_MARK_EXIST
	}
	return nil
}

func (g *Grid) OnRightClick() error {
	if g.value == GRID_MARK_EXIST {
		g.value = GRID_MARK_NOTEXIST
	} else if g.value == GRID_MARK_NOTEXIST {
		g.value = GRID_NULL
	} else if g.value == GRID_NULL {
		g.value = GRID_MARK_NOTEXIST
	}
	return nil
}

func (g *Grid) OnLeftDragOn() error {
	if g.value == GRID_MARK_EXIST {

	} else if g.value == GRID_MARK_NOTEXIST {
		g.value = GRID_MARK_EXIST
	} else if g.value == GRID_NULL {
		g.value = GRID_MARK_EXIST
	}
	return nil
}

func (g *Grid) OnRightDragOn() error {
	if g.value == GRID_MARK_EXIST {
		g.value = GRID_MARK_NOTEXIST
	} else if g.value == GRID_MARK_NOTEXIST {

	} else if g.value == GRID_NULL {
		g.value = GRID_MARK_NOTEXIST
	}
	return nil
}
