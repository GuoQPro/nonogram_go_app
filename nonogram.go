package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"time"
)

var (
	STAGE_W = 320
	STAGE_H = 320
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
}

type gameState int

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

func StartGame(row int, col int) (*Game, error) {
	game := &Game{}
	game.row = row
	game.col = col
	game.input = NewInput()
	InitFonts()
	err := game.InitGame(row, col)
	return game, err
}

func InitFonts() {
	tt, _ := truetype.Parse(fonts.MPlus1pRegular_ttf)
	textFont = truetype.NewFace(tt, &truetype.Options{
		Size:    8,
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

	timer_pos_x := int(g.board.start_x) + 50
	timer_pos_y := int(g.board.start_y + g.board.height + 50)
	text.Draw(screen, g.GetLapse(), textFont, timer_pos_x, timer_pos_y, color.RGBA{255, 255, 255, 255})
}

func (g *Game) InitGame(row int, col int) error {
	g.puzzle = GetPuzzle(row, col)
	g.GenerateBoard(g.puzzle)
	g.state = gameStatePlaying
	g.startTime = time.Now()
	return nil
}

func (g *Game) GetLapse() string {
	lapse := time.Since(g.startTime)
	return lapse.Truncate(time.Millisecond).String()
}

func (g *Game) GenerateBoard(puzzle [][]int) {
	g.board = NewBoard(puzzle)
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.input.Update()

	if g.state == gameStatePlaying {
		switch g.input.mouseState {
		case mouseStateLeftPress:
			g.board.OnLeftClick(g.input.mouseInitPosX, g.input.mouseInitPosY)
		case mouseStateRightPress:
			g.board.OnRightClick(g.input.mouseInitPosX, g.input.mouseInitPosY)
		case mouseStateLeftDrag:
			g.board.OnLeftDrag(g.input.mouseCurPosX, g.input.mouseCurPosY, g.input.mouseInitPosX, g.input.mouseInitPosY)
		case mouseStateRightDrag:
			g.board.OnRightDrag(g.input.mouseCurPosX, g.input.mouseCurPosY, g.input.mouseInitPosX, g.input.mouseInitPosY)
		}

		if g.IsCorrectAnswer() {
			g.state = gameStateSucc
		}
	} else if g.state == gameStateSucc {
		g.endTime = time.Now()
		g.state = gameStateSettle
	} else if g.state == gameStateSettle {

	}

	g.Draw(screen)

	return nil
}

func (g *Game) SubmitAnswer() bool {
	return g.IsCorrectAnswer()
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
