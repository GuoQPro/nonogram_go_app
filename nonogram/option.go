package nonogram

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
)

// Option is used for users to select puzzle size.
type Option struct {
	bound   Bound
	content string
	row     int
	col     int
}

// NewOption is constructor of option.
func NewOption(row int, col int, txt string, bound Bound) *Option {
	o := &Option{}
	o.bound = bound
	o.content = txt
	o.row = row
	o.col = col
	return o
}

// DrawOption is drawer of option instance.
func (o *Option) DrawOption(screen *ebiten.Image, isSelected bool) error {
	squareW := 10.0
	squareH := 10.0
	textX := int(o.bound.x + squareW + 10.0)
	textY := int(o.bound.y + 10.0)

	bgColor := colorBlue
	if isSelected {
		bgColor = colorRed
	}
	ebitenutil.DrawRect(screen, o.bound.x-10, o.bound.y-2, squareW+100, squareH+4, bgColor)

	// a check-box like controller.
	lineWidth := 1.0
	ebitenutil.DrawRect(screen, o.bound.x, o.bound.y, squareW, squareH, colorBlack)

	if !isSelected {
		ebitenutil.DrawRect(screen, o.bound.x+lineWidth, o.bound.y+lineWidth, squareW-lineWidth*2, squareH-lineWidth*2, bgColor)
	}

	text.Draw(screen, o.content, textFont, textX, textY, colorBlack)

	return nil
}

// TestTouch tests if a click event occurs inside the option area.
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

/*
IsCurrentOption test if a given row/col is current option.
*/
func (o *Option) IsCurrentOption(row int, col int) bool {
	if o.row == row && o.col == col {
		return true
	}
	return false
}
