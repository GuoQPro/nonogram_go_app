package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
)

var (
	STAGE_W = 320
	STAGE_H = 240
)

type Game struct {
	board  *Board
	puzzle [][]int
	input  *Input
}

type nonogramErr struct {
	desc string
}

func (e *nonogramErr) Error() string {
	return e.desc
}

func StartGame(row int, col int) (*Game, error) {
	game := &Game{}
	err := game.InitGame(row, col)
	return game, err
}

func (g *Game) RestartGame(row int, col int) error {
	err := g.InitGame(row, col)
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.board.DrawBoard(screen) != nil {
		log.Println("Fatal Error")
	}
}

func (g *Game) InitGame(row int, col int) error {
	g.puzzle = GetPuzzle(row, col)
	g.GenerateBoard(g.puzzle)
	g.input = NewInput()
	return nil
}

func (g *Game) GenerateBoard(puzzle [][]int) {
	g.board = &Board{}
	g.board.InitBoard(puzzle)
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.input.Update()

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

	g.Draw(screen)

	if g.IsCorrectAnswer() {
		log.Println("Congrats!!!!")
	}

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
