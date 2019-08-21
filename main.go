package main

import (
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"github.com/hajimehoshi/ebiten/inpututil"
	//"image/color"
	"log"
)

var game *Game

// update is called every frame (1/60 [s]).
func update(screen *ebiten.Image) error {

	// Write your game's logical update.

	if ebiten.IsDrawingSkipped() {
		// When the game is running slowly, the rendering result
		// will not be adopted.
		return nil
	}

	// Write your game's rendering.
	//draw_table(screen)
	//update_input()
	game.Update(screen)

	return nil
}

func main() {
	// Call ebiten.Run to start your game loop.
	game = StartGame()
	if err := ebiten.Run(update, STAGE_W, STAGE_H, 2, "Nonogram Game"); err != nil {
		log.Fatal(err)
	}
}
