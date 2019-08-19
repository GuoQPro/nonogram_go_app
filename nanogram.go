package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
)

var game_state int
var cur_table [][]int
var game_table [][]int

var (
	STAGE_W = 320
	STAGE_H = 240
)

type Game struct {
	board  *Board
	puzzle [][]int
	input  *Input
}

func StartGame() *Game {
	game := &Game{}
	game.InitGame()
	return game
}

func (g *Game) get_table() ([][]int, int, int) {
	return cur_table, 0, 0
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.board.DrawBoard(screen) != nil {
		log.Println("Fatal Error")
	}
}

func (g *Game) InitGame() {
	g.puzzle = GetPuzzle(5, 5)
	g.GenerateBoard(g.puzzle)
	g.input = NewInput()
}

func (g *Game) GenerateBoard(puzzle [][]int) {
	g.board = &Board{}
	g.board.InitBoard(puzzle)
}

func (g *Game) Update(screen *ebiten.Image) {
	g.input.Update()

	switch g.input.mouseState {

	case mouseStateLeftSettled:
		g.board.OnLeftClick(g.input.mouseRelPosX, g.input.mouseRelPosY)
	case mouseStateRightSettled:
		g.board.OnRightClick(g.input.mouseRelPosX, g.input.mouseRelPosY)
	case mouseStateLeftDrag:
		g.board.OnLeftDrag(g.input.mouseCurPosX, g.input.mouseCurPosY)
	case mouseStateRightDrag:
		g.board.OnRightDrag(g.input.mouseCurPosX, g.input.mouseCurPosY)
	}

	g.Draw(screen)
}
