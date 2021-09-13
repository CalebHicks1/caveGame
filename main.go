package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Definitions ///////////////////////////////////////////////////////

var (
	frames           = 0
	second           = time.Tick(time.Second)
	playerPos        = pixel.ZV
	playerVel        = pixel.ZV
	playerAcc        = pixel.ZV
	transition_speed = 12.0
)

// Main Functions ///////////////////////////////////////////////////

/*
This is where all the game code runs, it is basically
the new main function since main is used by pixel.
*/
func run() {

	// load the picture and create the player sprite
	playerImg, err := loadPicture("assets/smile.png")
	if err != nil {
		panic(err)
	}
	playerSprite := pixel.NewSprite(playerImg, playerImg.Bounds())

	// Window configuration
	cfg := pixelgl.WindowConfig{
		Title:  "Caleb's Game",
		Bounds: pixel.R(0, 0, 1024, 720),
		VSync:  true,
	}

	// Create game window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Declare vars
	var (
		playerMat = pixel.IM.Scaled(pixel.ZV, 8).Moved(win.Bounds().Center())
		last      = time.Now()
	)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Clear the screen and redraw the sprite
		win.Clear(color.Black)
		playerSprite.Draw(win, playerMat)

		// Controls
		switch {
		case win.Pressed(pixelgl.KeyRight):
			playerAcc.X = 5.0
		case win.Pressed(pixelgl.KeyLeft):
			playerAcc.X = -5.0
		default:
			playerAcc.X = 0.0
		}

		// Apply velocity
		playerVel.X = playerVel.X*(1-dt*transition_speed) + playerAcc.X*(dt*transition_speed)
		playerMat = playerMat.Moved(playerVel)

		// Draw fps on window title
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}

		win.Update()

	}
}

/*
This is the main function. It starts the
pixel engine with the run function.
*/
func main() {
	pixelgl.Run(run)
}

// Helper Functions /////////////////////////////////////////////////

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
