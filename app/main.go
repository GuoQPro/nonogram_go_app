package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/mrsep18th/nonogram_go_app/nonogram"
)

var game *nonogram.Game

// update is called every frame (1/60 [s]).
func update(screen *ebiten.Image) error {

	// Write your game's logical update.

	if ebiten.IsDrawingSkipped() {
		// When the game is running slowly, the rendering result
		// will not be adopted.
		return nil
	}

	// Write your game's rendering.
	return game.Update(screen)
}

func main() {
	// Call ebiten.Run to start your game loop.
	game, _ = nonogram.StartGame()
	if err := ebiten.Run(update, 400, 600, 2, "Nonogram Game"); err != nil {
		log.Fatal(err)
	}
}
