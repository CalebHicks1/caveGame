package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Definitions ///////////////////////////////////////////////////////

var (
	gravity          float64
	walkSpeed        = 50.0
	runSpeed         float64
	rect             pixel.Rect
	vel              pixel.Vec
	onGround         bool
	desired_velocity = 0.2
)

type playerSprite struct {
}

var (
	frames    = 0
	second    = time.Tick(time.Second)
	playerVel = pixel.ZV
	playerAcc = pixel.ZV
	camPos    = pixel.ZV
	playerPos = pixel.ZV
)

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

// Main Functions ///////////////////////////////////////////////////

/*
This is where all the game code runs, it is basically
the new main function since main is used by pixel.
*/
func run() {
	// variables
	//playerPhys := &playerPhysics{}
	playerImg, err := loadPicture("assets/smile.png")
	if err != nil {
		panic(err)
	}
	playerSprite := pixel.NewSprite(playerImg, playerImg.Bounds())

	cfg := pixelgl.WindowConfig{
		Title:  "Caleb's Game",
		Bounds: pixel.R(0, 0, 1024, 720),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	playerPos = win.Bounds().Center()
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		playerPos.Y = playerVel.Y * 20
		win.Clear(color.Black)
		playerMat := pixel.IM
		playerMat = playerMat.Scaled(pixel.ZV, 8)
		playerMat = playerMat.Moved(playerPos)
		playerSprite.Draw(win, playerMat)
		win.Update()

		// Controls
		desired_velocity = 0.0
		if win.Pressed(pixelgl.KeyLeft) {
			desired_velocity = -200.0
		}
		if win.Pressed(pixelgl.KeyRight) {
			desired_velocity = 200.0
		}
		playerVel.X = playerVel.X*(1-dt*20) + desired_velocity*(dt*20)
		playerPos.X += playerVel.X * dt

		// Draw fps on window title
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

/*
This is the main function. It starts the
pixel engine with the run function.
*/
func main() {
	pixelgl.Run(run)
}
