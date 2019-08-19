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

/*func draw_table(screen *ebiten.Image) {
	t, _, _ := get_table()

	c1 := color.RGBA{255, 0, 0, 255}
	c0 := color.RGBA{255, 255, 255, 255}

	start_x := 10.0
	start_y := 10.0

	cur_x := start_x
	cur_y := start_y

	gap_x := 1.0
	gap_y := 1.0

	w := 20.0
	h := 20.0

	for row := 0; row < len(t); row++ {
		cur_x = start_x
		for col := 0; col < len(t[row]); col++ {
			v := t[row][col]
			if v == 0 {
				ebitenutil.DrawRect(screen, cur_x, cur_y, w, h, c0)
			} else if v == 1 {
				ebitenutil.DrawRect(screen, cur_x, cur_y, w, h, c1)
			}
			cur_x += (w + gap_x)
		}

		cur_y += (h + gap_y)
	}
}

func update_input() {
	leftMouseBtn := ebiten.MouseButtonLeft
	rightMouseBtn := ebiten.MouseButtonRight
	if inpututil.IsMouseButtonJustPressed(leftMouseBtn) {
		log.Println(ebiten.CursorPosition())
	} else if inpututil.IsMouseButtonJustPressed(rightMouseBtn) {
		log.Println(ebiten.CursorPosition())
	}
}*/
