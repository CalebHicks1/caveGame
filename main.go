package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

// struct to hold line properties
type line struct {
	p1 pixel.Vec
	p2 pixel.Vec
}

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
		Bounds: pixel.R(0, 0, 1900, 700),
		VSync:  true,
	}

	// Create game window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Declare vars
	var (
		cameraPosition = pixel.ZV
		last           = time.Now()
		playerWidth    = 60.0
		playerHeight   = 60.0
		playerRec      = pixel.R(-playerWidth/2, -playerHeight/2, playerWidth/2, playerHeight/2).Moved(pixel.ZV)
	)

	// create IMDraw object
	imd := imdraw.New(nil)

	// // Create line
	// groundLine := line{pixel.V(24, 360), pixel.V(1000, 360)}
	// imd.EndShape = imdraw.RoundEndShape
	// imd.Push(groundLine.p1, groundLine.p2)
	// imd.Line(15)

	rectangle := pixel.R(-800, -100, 800, -80) // ground rectangle

	// MAIN LOOP ////////////////////////////////////
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		cameraPosition = pixel.Lerp(cameraPosition, playerRec.Center(), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(win.Bounds().Center().Sub(cameraPosition))
		win.SetMatrix(cam)
		last = time.Now()
		fmt.Printf("Player position: (%.2f, %.2f)\r", playerRec.Center().X, playerRec.Center().Y)

		// Clear the screen and redraw the sprite
		// Controls
		switch {
		case win.Pressed(pixelgl.KeyRight):
			playerAcc.X = 5.0
		case win.Pressed(pixelgl.KeyLeft):
			playerAcc.X = -5.0
		default:
			playerAcc.X = 0.0
		}
		if win.Pressed(pixelgl.KeyUp) && playerRec.Intersects(rectangle) {
			playerVel.Y = 15
		}

		// If touching ground
		if !playerRec.Intersects(rectangle) {
			playerAcc.Y = -1.5
		} else if playerVel.Y <= 0 {
			playerAcc.Y = 0.0
			playerVel.Y = 0.0
			playerRec = playerRec.Moved(pixel.V(0, rectangle.Max.Y-playerRec.Min.Y))
		}

		// Apply velocity
		playerVel.X = playerVel.X*(1-dt*transition_speed) + playerAcc.X*(dt*transition_speed)
		playerVel.Y += playerAcc.Y
		playerRec = playerRec.Moved(playerVel)

		// Draw fps on window title
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		imd.Clear()
		win.Clear(color.Black)
		imd.Color = pixel.RGB(255, 255, 255)
		imd.Push(rectangle.Min, rectangle.Max)
		imd.Rectangle(0)
		imd.Color = pixel.RGB(255, 0, 0)
		imd.Push(playerRec.Min, playerRec.Max)
		imd.Rectangle(1)
		imd.Draw(win)
		playerSprite.Draw(win, pixel.IM.ScaledXY(
			pixel.ZV, pixel.V(playerRec.W()/16, playerRec.H()/16)).Moved(
			playerRec.Center()))

		win.Update()

	}
	fmt.Printf("\n")
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
