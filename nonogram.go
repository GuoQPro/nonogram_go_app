package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"time"
)

var (
	STAGE_W = 400
	STAGE_H = 600
)

type Game struct {
	board     *Board
	puzzle    [][]int
	input     *Input
	row       int
	col       int
	startTime time.Time
	endTime   time.Time
	state     gameState
	op_mode   opMode
	options   []Option
}

// outline
type Bound struct {
	x float64
	y float64
	w float64
	h float64
}

type opMode int
type gameState int

const (
	opModeLeftClick opMode = iota
	opModeRightClick
	opModeReserved
)

var (
	color_white = color.RGBA{252, 245, 239, 255}
	color_blue  = color.RGBA{109, 181, 202, 255}
	color_red   = color.RGBA{255, 104, 53, 255}
	color_black = color.RGBA{0, 0, 0, 255}
)

func (gs gameState) String() string {
	if gs == gameStatePlaying {
		return "Playing!!!"
	} else if gs == gameStateSucc {
		return "Congrats!!!"
	}

	return ""
}

const (
	gameStatePlaying gameState = iota
	gameStateSucc
	gameStateSettle
)

type nonogramErr struct {
	desc string
}

func (e *nonogramErr) Error() string {
	return e.desc
}

var textFont font.Face

func StartGame() (*Game, error) {
	game := &Game{}
	game.op_mode = opModeLeftClick
	game.input = NewInput()
	InitFonts()
	game.SetOptions()

	default_row := 5
	default_col := 5
	err := game.InitGame(default_row, default_col)
	return game, err
}

func (g *Game) SetOptions() {
	// Draw Options
	availableSize := [][]int{
		{5, 5}, {9, 9}, {12, 12}, {15, 15},
	}

	g.options = make([]Option, 0)

	org_x := (float64(STAGE_W) / 3)
	org_y := float64(STAGE_H) * 1 / 4
	height := 20.0

	optionNumInRow := 2

	counter := 0

	x := org_x
	y := org_y
	for i := range availableSize {
		row := availableSize[i][0]
		col := availableSize[i][1]
		txt := fmt.Sprintf("%d x %d", row, col)
		o := NewOption(row, col, txt, Bound{
			x: x,
			y: y,
			w: float64(STAGE_W),
			h: height,
		})

		g.options = append(g.options, *o)

		counter += 1

		if counter >= optionNumInRow {
			y -= height
			x = org_x
			counter = 0
		} else {
			x += (float64(STAGE_W) / 3)
		}
	}
}

func (g *Game) SwitchOpMode() {
	g.op_mode++
	g.op_mode %= opModeReserved
}

func (g *Game) ShowHint() {
	// todo: if the current solution is already wrong.
	t := make([][]int, g.row)

	for r := 0; r < g.row; r++ {
		t[r] = make([]int, g.col)
		for c := 0; c < g.col; c++ {
			t[r][c] = int(g.board.grids[r][c].value)
		}
	}
	hint_row, hint_col, _ := GetHint(g.board.row_ind, g.board.col_ind, t)

	g.board.grids[hint_row][hint_col].Hint()
}

func InitFonts() {
	tt, _ := truetype.Parse(fonts.ArcadeN_ttf)
	textFont = truetype.NewFace(tt, &truetype.Options{
		Size:    6,
		DPI:     96,
		Hinting: font.HintingFull,
	})
}

func (g *Game) RestartGame(row int, col int) error {
	err := g.InitGame(row, col)
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	if err := g.board.DrawBoard(screen); err != nil {
		log.Println(err)
	}

	for i := range g.options {
		isCurrentOption := g.options[i].IsCurrentOption(g.row, g.col)
		g.options[i].DrawOption(screen, isCurrentOption)
	}

	op_indicator_x := g.board.start_x
	op_indicator_y := float64(g.board.start_y + g.board.height + 30)
	if g.op_mode == opModeLeftClick {
		ebitenutil.DrawRect(screen, op_indicator_x, op_indicator_y, grid_w, grid_h, color_black)
	} else {
		ebitenutil.DrawRect(screen, op_indicator_x, op_indicator_y, grid_w, grid_h, color_red)
	}

	var timeLapse string
	if g.state == gameStateSettle {
		timeLapse = "Congrats: " + g.GetSolvingTime()
	} else {
		timeLapse = g.GetLapse()
	}

	timer_pos_x := int(op_indicator_x + grid_w*2)
	timer_pos_y := int(op_indicator_y + 10)
	text.Draw(screen, timeLapse, textFont, timer_pos_x, timer_pos_y, color_black)
}

func (g *Game) InitGame(row int, col int) error {
	g.puzzle = GetPuzzle(row, col)
	g.GenerateBoard(g.puzzle)
	g.state = gameStatePlaying
	g.startTime = time.Now()
	g.row = row
	g.col = col
	return nil
}

func (g *Game) GetLapse() string {
	lapse := time.Since(g.startTime)
	return lapse.Truncate(time.Millisecond).String()
}

func (g *Game) GetSolvingTime() string {
	lapse := g.endTime.Sub(g.startTime)
	return lapse.Truncate(time.Millisecond).String()
}

func (g *Game) GenerateBoard(puzzle [][]int) {
	bound := Bound{}
	BOARD_MARGIN := float64(STAGE_W) * 0.1

	bound.w = float64(STAGE_W)*2/3 - BOARD_MARGIN
	bound.h = float64(STAGE_H)*2/3 - BOARD_MARGIN
	bound.x = float64(STAGE_W) - (bound.w + BOARD_MARGIN)
	bound.y = float64(STAGE_H) - (bound.h + BOARD_MARGIN)
	g.board = NewBoard(puzzle, bound)
}

func (g *Game) UpdateStateMachine() error {
	switch g.state {
	case gameStatePlaying:
		switch g.input.mouseState {
		case mouseStateLeftPress:
			for i := range g.options {
				if g.options[i].TestTouch(g.input.mouseInitPosX, g.input.mouseInitPosY) {
					g.RestartGame(g.options[i].row, g.options[i].col)
				}
			}
			g.board.OnLeftClick(g.input.mouseInitPosX, g.input.mouseInitPosY)
		case mouseStateRightPress:
			g.board.OnRightClick(g.input.mouseInitPosX, g.input.mouseInitPosY)
		case mouseStateLeftDrag:
			g.board.OnLeftDrag(g.input.mouseCurPosX, g.input.mouseCurPosY, g.input.mouseInitPosX, g.input.mouseInitPosY)
		case mouseStateRightDrag:
			g.board.OnRightDrag(g.input.mouseCurPosX, g.input.mouseCurPosY, g.input.mouseInitPosX, g.input.mouseInitPosY)
		}

		switch g.input.touchState {
		case touchStateMultiTouch:
			if g.input.touchPointNum == 2 {
				g.SwitchOpMode()
			} else if g.input.touchPointNum == 3 {
				g.ShowHint()
			}

		case touchStatePress:
			if g.op_mode == opModeLeftClick {
				g.board.OnLeftClick(g.input.touchInitPosX, g.input.touchInitPosY)
			} else {
				g.board.OnRightClick(g.input.touchInitPosX, g.input.touchInitPosY)
			}
		case touchStateDrag:
			if g.op_mode == opModeLeftClick {
				g.board.OnLeftDrag(g.input.touchCurPosX, g.input.touchCurPosY, g.input.touchInitPosX, g.input.touchInitPosY)
			} else {
				g.board.OnRightDrag(g.input.touchCurPosX, g.input.touchCurPosY, g.input.touchInitPosX, g.input.touchInitPosY)
			}
		}

		if g.IsCorrectAnswer() {
			g.state = gameStateSucc
		}

	case gameStateSucc:
		g.endTime = time.Now()
		g.state = gameStateSettle

	case gameStateSettle:
		if g.input.mouseState == mouseStateLeftPress ||
			g.input.touchState == touchStatePress {
			g.RestartGame(g.row, g.col)
		}
	}
	return nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.input.Update()
	g.UpdateStateMachine()
	g.Draw(screen)

	return nil
}

func (g *Game) IsCorrectAnswer() bool {
	for i := range g.puzzle {
		for j := range g.puzzle[i] {
			q_value := g.puzzle[i][j]
			a_value := g.board.grids[i][j].GetValue()

			if a_value == GRID_MARK_NOTEXIST {
				a_value = GRID_NULL
			}

			if q_value != int(a_value) {
				return false
			}
		}
	}

	return true
}
