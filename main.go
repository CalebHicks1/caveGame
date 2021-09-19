package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
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
		playerRec = pixel.R(-56, -56, 56, 56).Moved(win.Bounds().Center())
	)

	// create IMDraw object
	imd := imdraw.New(nil)

	// // Create line
	// groundLine := line{pixel.V(24, 360), pixel.V(1000, 360)}
	// imd.EndShape = imdraw.RoundEndShape
	// imd.Push(groundLine.p1, groundLine.p2)
	// imd.Line(15)

	rectangle := pixel.R(0, 0, 1024, 40)
	// imd.Push(rectangle.Min, rectangle.Max)
	// imd.Rectangle(0)

	// MAIN LOOP ////////////////////////////////////
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		fmt.Printf("Player position: (%.2f, %.2f)", playerMat.Project(pixel.ZV).X, playerMat.Project(pixel.ZV).Y)

		// Clear the screen and redraw the sprite
		imd.Clear()
		win.Clear(color.Black)
		imd.Push(playerRec.Min, playerRec.Max)
		imd.Rectangle(0)
		imd.Push(rectangle.Min, rectangle.Max)
		imd.Rectangle(0)
		imd.Draw(win)
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

		if !playerRec.Intersects(rectangle) {
			playerAcc.Y = -5.0
			fmt.Printf(" not touching ground  \r")
		} else {
			playerAcc.Y = 0.0
			fmt.Printf(" touching ground      \r")
		}

		// Apply velocity
		playerVel.X = playerVel.X*(1-dt*transition_speed) + playerAcc.X*(dt*transition_speed)
		playerVel.Y = 0.5 * playerAcc.Y
		playerRec = playerRec.Moved(playerVel)
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
