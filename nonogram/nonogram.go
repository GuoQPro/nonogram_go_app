package nonogram

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
	"github.com/mrsep18th/nonogram_go_app/nonogram/puzzle"
	"golang.org/x/image/font"
)

// define size of game canvas
var (
	StageWidth  = 400
	StageHeight = 600
)

// Game is the data structure of app.
type Game struct {
	board     *Board
	puzzle    puzzle.Puzzle
	input     *Input
	row       int
	col       int
	startTime time.Time
	endTime   time.Time
	state     gameState
	opMode    opMode
	options   []Option
}

// Bound is the outline of game board
type Bound struct {
	x float64
	y float64
	w float64
	h float64
}

type opMode int

const (
	opModeLeftClick opMode = iota
	opModeRightClick
	opModeReserved
)

type gameState int

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

// StartGame means literally.
func StartGame() (*Game, error) {
	game := &Game{}
	game.opMode = opModeLeftClick
	game.input = NewInput()
	initFonts()
	game.setOptions()

	defaultRow := 5
	defaultCol := 5
	err := game.initGame(defaultRow, defaultCol)
	return game, err
}

func (g *Game) setOptions() {
	// Draw Options
	availableSize := [][]int{
		{5, 5}, {9, 9}, {12, 12}, {15, 15},
	}

	g.options = make([]Option, 0)

	orgX := (float64(StageWidth) / 3)
	orgY := float64(StageHeight) * 1 / 4
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
			w: float64(StageWidth),
			h: height,
		})

		g.options = append(g.options, *o)

		counter++

		if counter >= optionNumInRow {
			y -= height
			x = orgX
			counter = 0
		} else {
			x += (float64(StageWidth) / 3)
		}
	}
}

func (g *Game) switchOpMode() {
	g.opMode++
	g.opMode %= opModeReserved
}

func (g *Game) showHint() {
	// todo: if the current solution is already wrong.
	t := make(puzzle.Puzzle, g.row)

	for r := 0; r < g.row; r++ {
		t[r] = make([]int, g.col)
		for c := 0; c < g.col; c++ {
			t[r][c] = int(g.board.grids[r][c].value)
		}
	}
	hintRow, hintCol, _ := puzzle.GetHint(g.board.rowInd, g.board.colInd, t)

	g.board.grids[hintRow][hintCol].Hint()
}

func initFonts() {
	tt, _ := truetype.Parse(fonts.ArcadeN_ttf)
	textFont = truetype.NewFace(tt, &truetype.Options{
		Size:    6,
		DPI:     96,
		Hinting: font.HintingFull,
	})
}

func (g *Game) restartGame(row int, col int) error {
	err := g.initGame(row, col)
	return err
}

func (g *Game) draw(screen *ebiten.Image) {
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
		timeLapse = fmt.Sprintf("[%s] %s", g.state, g.getSolvingTime())
	} else {
		timeLapse = fmt.Sprintf("[%s] %s", g.state, g.getLapse())
	}

	timerPosX := int(opIndicatorX + gridWidth*2)
	timerPosY := int(opIndicatorY + 10)
	text.Draw(screen, timeLapse, textFont, timerPosX, timerPosY, colorBlack)
}

func (g *Game) initGame(row int, col int) error {
	g.puzzle = puzzle.GetPuzzle(row, col)
	g.generateBoard(g.puzzle)
	g.state = gameStatePlaying
	g.startTime = time.Now()
	g.row = row
	g.col = col
	return nil
}

func (g *Game) getLapse() string {
	lapse := time.Since(g.startTime)
	return lapse.Truncate(time.Millisecond).String()
}

func (g *Game) getSolvingTime() string {
	lapse := g.endTime.Sub(g.startTime)
	return lapse.Truncate(time.Millisecond).String()
}

func (g *Game) generateBoard(p puzzle.Puzzle) {
	bound := Bound{}
	boardMarginWidth := float64(StageWidth) * 0.1
	bound.w = float64(StageWidth)*2/3 - boardMarginWidth
	bound.h = float64(StageHeight)*2/3 - boardMarginWidth
	bound.x = float64(StageWidth) - (bound.w + boardMarginWidth)
	bound.y = float64(StageHeight) - (bound.h + boardMarginWidth)
	g.board = NewBoard(p, bound)
}

func (g *Game) updateStateMachine() error {
	switch g.state {
	case gameStatePlaying:
		switch g.input.mouseState {
		case mouseStateLeftPress:
			for i := range g.options {
				if g.options[i].TestTouch(g.input.mouseInitPosX, g.input.mouseInitPosY) {
					g.restartGame(g.options[i].row, g.options[i].col)
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
				g.switchOpMode()
			} else if g.input.touchPointNum == 3 {
				g.showHint()
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

		if g.isCorrectAnswer() {
			g.state = gameStateSucc
		}

	case gameStateSucc:
		g.endTime = time.Now()
		g.state = gameStateSettle

	case gameStateSettle:
		if g.input.mouseState == mouseStateLeftPress ||
			g.input.touchState == touchStatePress {
			g.restartGame(g.row, g.col)
		}
	}
	return nil
}

/*
Update is the mainloop of the game.
*/
func (g *Game) Update(screen *ebiten.Image) error {
	g.input.Update()
	g.updateStateMachine()
	g.draw(screen)

	return nil
}

func (g *Game) isCorrectAnswer() bool {
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
