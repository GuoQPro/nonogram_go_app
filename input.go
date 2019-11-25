package nonogram_go_app

import (
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/inpututil"
	"math"
)

type mouseState int

const (
	mouseStateNone mouseState = iota
	mouseStateLeftPress
	mouseStateLeftPressing
	mouseStateLeftDrag
	mouseStateLeftSettled
	mouseStateRightPress
	mouseStateRightPressing
	mouseStateRightDrag
	mouseStateRightSettled
)

type touchState int

const (
	touchStateNone touchState = iota
	touchStateInvalid
	touchStateMultiTouching
	touchStateMultiTouch
	touchStatePress
	touchStatePressing
	touchStateDrag
)

// Input represents the current key states.
type Input struct {
	mouseState    mouseState
	mouseInitPosX int
	mouseInitPosY int
	mouseRelPosX  int
	mouseRelPosY  int
	mouseCurPosX  int
	mouseCurPosY  int

	touchState    touchState
	touchID       int
	touchInitPosX int
	touchInitPosY int
	touchCurPosX  int
	touchCurPosY  int
	touchPointNum int
}

const (
	dragThreshold = 5.0
)

func NewInput() *Input {
	return &Input{}
}

// Update updates the current input states.
func (i *Input) Update() {
	switch i.mouseState {
	case mouseStateNone:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.mouseInitPosX = x
			i.mouseInitPosY = y
			i.mouseState = mouseStateLeftPress
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			i.mouseInitPosX = x
			i.mouseInitPosY = y
			i.mouseState = mouseStateRightPress
		}
	case mouseStateLeftPress:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.mouseRelPosX = x
			i.mouseRelPosY = y
			i.mouseState = mouseStateLeftSettled
		} else {
			x, y := ebiten.CursorPosition()
			i.mouseCurPosX = x
			i.mouseCurPosY = y
			i.mouseState = mouseStateLeftPressing
		}

	case mouseStateLeftPressing:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.mouseRelPosX = x
			i.mouseRelPosY = y
			i.mouseState = mouseStateLeftSettled
		} else {
			x, y := ebiten.CursorPosition()
			i.mouseCurPosX = x
			i.mouseCurPosY = y

			if math.Abs(float64(x-i.mouseInitPosX)) > dragThreshold || math.Abs(float64(y-i.mouseInitPosY)) > dragThreshold {
				i.mouseState = mouseStateLeftDrag
			}
		}
	case mouseStateLeftDrag:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			i.mouseState = mouseStateNone
		} else {
			x, y := ebiten.CursorPosition()
			i.mouseCurPosX = x
			i.mouseCurPosY = y
		}

	case mouseStateLeftSettled:
		i.mouseState = mouseStateNone

	case mouseStateRightPress:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			i.mouseRelPosX = x
			i.mouseRelPosY = y
			i.mouseState = mouseStateRightSettled
		} else {
			x, y := ebiten.CursorPosition()
			i.mouseCurPosX = x
			i.mouseCurPosY = y
			i.mouseState = mouseStateRightPressing
		}

	case mouseStateRightPressing:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			i.mouseRelPosX = x
			i.mouseRelPosY = y
			i.mouseState = mouseStateRightSettled
		} else {
			x, y := ebiten.CursorPosition()
			i.mouseCurPosX = x
			i.mouseCurPosY = y

			if math.Abs(float64(x-i.mouseInitPosX)) > dragThreshold || math.Abs(float64(y-i.mouseInitPosY)) > dragThreshold {
				i.mouseState = mouseStateRightDrag
			}
		}
	case mouseStateRightDrag:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			i.mouseState = mouseStateNone
		} else {
			x, y := ebiten.CursorPosition()
			i.mouseCurPosX = x
			i.mouseCurPosY = y
		}
	case mouseStateRightSettled:
		i.mouseState = mouseStateNone
	}

	switch i.touchState {
	case touchStateNone:
		ts := ebiten.TouchIDs()
		if len(ts) == 1 {
			i.touchID = ts[0]
			x, y := ebiten.TouchPosition(ts[0])
			i.touchInitPosX = x
			i.touchInitPosY = y
			i.touchState = touchStatePress
		} else if len(ts) >= 2 {
			i.touchState = touchStateMultiTouch
		}
	case touchStatePress:
		i.touchState = touchStatePressing

	case touchStatePressing:
		ts := ebiten.TouchIDs()
		if len(ts) >= 2 {
			i.touchState = touchStateMultiTouch
		} else if len(ts) == 1 {
			if ts[0] != i.touchID {
				i.touchState = touchStateInvalid
			} else {
				x, y := ebiten.TouchPosition(ts[0])
				i.touchCurPosX = x
				i.touchCurPosY = y

				if math.Abs(float64(x-i.touchInitPosX)) > dragThreshold || math.Abs(float64(y-i.touchInitPosY)) > dragThreshold {
					i.touchState = touchStateDrag
				}
			}
			break
		} else if len(ts) == 0 {
			i.touchState = touchStateNone
		}
	case touchStateInvalid:
		if len(ebiten.TouchIDs()) == 0 {
			i.touchState = touchStateNone
		}
	case touchStateMultiTouch:
		ts := ebiten.TouchIDs()
		i.touchPointNum = len(ts)
		i.touchState = touchStateMultiTouching
	case touchStateMultiTouching:
		ts := ebiten.TouchIDs()
		if len(ts) == 0 {
			i.touchState = touchStateNone
		}
	case touchStateDrag:
		ts := ebiten.TouchIDs()
		if len(ts) == 0 {
			i.touchState = touchStateNone
		} else if len(ts) == 2 {
			i.touchState = touchStateMultiTouch
		} else if len(ts) == 1 {
			x, y := ebiten.TouchPosition(ts[0])
			i.touchCurPosX = x
			i.touchCurPosY = y
		}
	}
}
