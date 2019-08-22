package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	//"log"
)

const (
	grid_w = 16
	grid_h = 16

	inner_grid_w = grid_w - 2
	inner_grid_h = grid_h - 2
)

var (
	color_white = color.RGBA{255, 255, 255, 255}
	color_blue  = color.RGBA{0, 150, 214, 255}
	color_red   = color.RGBA{250, 128, 114, 255}
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

	ebitenutil.DrawRect(screen, g.pos_x, g.pos_y, grid_w, grid_h, color_white)

	if g.value == GRID_MARK_EXIST {
		ebitenutil.DrawRect(screen, g.pos_x+1, g.pos_y+1, inner_grid_w, inner_grid_h, color_blue)
	} else if g.value == GRID_MARK_NOTEXIST {
		ebitenutil.DrawRect(screen, g.pos_x+1, g.pos_y+1, inner_grid_w, inner_grid_h, color_red)
	} else if g.value == GRID_NULL {
		// do nothing
	}

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

func (g *Grid) IsSameGrid(rh *Grid) bool {
	if g.pos_x == rh.pos_x && g.pos_y == rh.pos_y {
		return true
	} else {
		return false
	}
}
