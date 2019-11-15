package main

import (
	//"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
)

type Option struct {
	bound   Bound
	content string
	row     int
	col     int
}

func NewOption(row int, col int, txt string, bound Bound) *Option {
	o := &Option{}
	o.bound = bound
	o.content = txt
	o.row = row
	o.col = col
	return o
}

func (o *Option) DrawOption(screen *ebiten.Image, isSelected bool) error {
	squareW := 10.0
	squareH := 10.0
	text_x := int(o.bound.x + squareW + 10.0)
	text_y := int(o.bound.y + 10.0)

	bg_color := color_blue
	if isSelected {
		bg_color = color_red
	}
	ebitenutil.DrawRect(screen, o.bound.x-10, o.bound.y-2, squareW+100, squareH+4, bg_color)

	// a check-box like controller.
	line_width := 1.0
	ebitenutil.DrawRect(screen, o.bound.x, o.bound.y, squareW, squareH, color_black)

	if !isSelected {
		ebitenutil.DrawRect(screen, o.bound.x+line_width, o.bound.y+line_width, squareW-line_width*2, squareH-line_width*2, bg_color)
	}

	text.Draw(screen, o.content, textFont, text_x, text_y, color_black)

	return nil
}

func (o *Option) TestTouch(x int, y int) bool {
	left := int(o.bound.x)
	right := left + int(o.bound.w)
	top := int(o.bound.y)
	bottom := top + int(o.bound.h)

	if x >= left && x <= right && y >= top && y <= bottom {
		return true
	}

	return false
}

func (o *Option) IsCurrentOption(row int, col int) bool {
	if o.row == row && o.col == col {
		return true
	} else {
		return false
	}
}
