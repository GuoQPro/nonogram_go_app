package nonogram_go_app

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var (
	stageWidth  = 400
	stageHeight = 600
)

type Game struct {
	board     *Board
	puzzle    Puzzle //[][]int
	input     *Input
	row       int
	col       int
	startTime time.Time
	endTime   time.Time
	state     gameState
	opMode    OpMode
	options   []Option
}

// outline
type Bound struct {
	x float64
	y float64
	w float64
	h float64
}

type OpMode int
type gameState int

const (
	opModeLeftClick OpMode = iota
	opModeRightClick
	opModeReserved
)

var (
	colorWhite = color.RGBA{252, 245, 239, 255}
	colorBlue  = color.RGBA{109, 181, 202, 255}
	colorRed   = color.RGBA{255, 104, 53, 255}
	colorBlack = color.RGBA{0, 0, 0, 255}
)

func (gs gameState) String() string {
	switch gs {
	case gameStatePlaying:
		return "Playing"
	case gameStateSucc:
		return "Congrats"
	case gameStateSettle:
		return "Congrats"
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
	game.opMode = opModeLeftClick
	game.input = NewInput()
	InitFonts()
	game.SetOptions()

	defaultRow := 5
	defaultCol := 5
	err := game.InitGame(defaultRow, defaultCol)
	return game, err
}

func (g *Game) SetOptions() {
	// Draw Options
	availableSize := [][]int{
		{5, 5}, {9, 9}, {12, 12}, {15, 15},
	}

	g.options = make([]Option, 0)

	orgX := (float64(stageWidth) / 3)
	orgY := float64(stageHeight) * 1 / 4
	height := 20.0

	optionNumInRow := 2

	counter := 0

	x := orgX
	y := orgY
	for i := range availableSize {
		row := availableSize[i][0]
		col := availableSize[i][1]
		txt := fmt.Sprintf("%d x %d", row, col)
		o := NewOption(row, col, txt, Bound{
			x: x,
			y: y,
			w: float64(stageWidth),
			h: height,
		})

		g.options = append(g.options, *o)

		counter++

		if counter >= optionNumInRow {
			y -= height
			x = orgX
			counter = 0
		} else {
			x += (float64(stageWidth) / 3)
		}
	}
}

func (g *Game) SwitchOpMode() {
	g.opMode++
	g.opMode %= opModeReserved
}

func (g *Game) ShowHint() {
	// todo: if the current solution is already wrong.
	t := make(Puzzle, g.row)

	for r := 0; r < g.row; r++ {
		t[r] = make([]int, g.col)
		for c := 0; c < g.col; c++ {
			t[r][c] = int(g.board.grids[r][c].value)
		}
	}
	hintRow, hintCol, _ := GetHint(g.board.rowInd, g.board.colInd, t)

	g.board.grids[hintRow][hintCol].Hint()
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

	opIndicatorX := g.board.startX
	opIndicatorY := float64(g.board.startY + g.board.height + 30)
	if g.opMode == opModeLeftClick {
		ebitenutil.DrawRect(screen, opIndicatorX, opIndicatorY, gridWidth, gridHeight, colorBlack)
	} else {
		ebitenutil.DrawRect(screen, opIndicatorX, opIndicatorY, gridWidth, gridHeight, colorRed)
	}

	var timeLapse string
	if g.state == gameStateSettle {
		timeLapse = fmt.Sprintf("[%s] %s", g.state, g.GetSolvingTime()) // "Congrats: " + g.GetSolvingTime()
	} else {
		timeLapse = fmt.Sprintf("[%s] %s", g.state, g.GetLapse()) //g.GetLapse()
	}

	timerPosX := int(opIndicatorX + gridWidth*2)
	timerPosY := int(opIndicatorY + 10)
	text.Draw(screen, timeLapse, textFont, timerPosX, timerPosY, colorBlack)
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

func (g *Game) GenerateBoard(p Puzzle) {
	bound := Bound{}
	boardMarginWidth := float64(stageWidth) * 0.1
	bound.w = float64(stageWidth)*2/3 - boardMarginWidth
	bound.h = float64(stageHeight)*2/3 - boardMarginWidth
	bound.x = float64(stageWidth) - (bound.w + boardMarginWidth)
	bound.y = float64(stageHeight) - (bound.h + boardMarginWidth)
	g.board = NewBoard(p, bound)
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
			if g.opMode == opModeLeftClick {
				g.board.OnLeftClick(g.input.touchInitPosX, g.input.touchInitPosY)
			} else {
				g.board.OnRightClick(g.input.touchInitPosX, g.input.touchInitPosY)
			}
		case touchStateDrag:
			if g.opMode == opModeLeftClick {
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
			qValue := g.puzzle[i][j]
			aValue := g.board.grids[i][j].GetValue()

			if aValue == gridStateMarkNotExist {
				aValue = gridStateNull
			}

			if qValue != int(aValue) {
				return false
			}
		}
	}

	return true
}
