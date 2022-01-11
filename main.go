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
	touchingGround   = false
	transition_speed = 12.0
)

// Returns the slope of the line that is intersecting r
func GetIntersectingLineSlope(l pixel.Line, r pixel.Rect) float64 {
	points := r.IntersectionPoints(l)
	if len(points) == 0 {
		return -1
	} else {
		if len(points) == 2 {
			return math.Abs((points[0].Y - points[1].Y) / (points[0].X - points[1].X))
		}
		return 0
	}
}

// Returns a vec that moves r up so that it is no longer intersecting l
func MoveRectangleUp(l pixel.Line, r pixel.Rect) pixel.Vec {
	testRec := r
	retVec := pixel.ZV
	for len(testRec.IntersectionPoints(l)) > 0 {
		retVec.Y++
		testRec = testRec.Moved(retVec)
		//fmt.Print(retVec)
	}
	retVec.Y-- // subtract one from the y component in order to keep sprite on the ground and keep it from bouncing.
	return retVec
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
		Bounds: pixel.R(0, 0, 900, 700),
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
		playerHeight   = 100.0
		playerRec      = pixel.R(-playerWidth/2, -playerHeight/2, playerWidth/2, playerHeight/2).Moved(pixel.ZV)
	)

	// create IMDraw object
	imd := imdraw.New(nil)

	// // Create line
	// groundLine := line{pixel.V(24, 360), pixel.V(1000, 360)}
	// imd.EndShape = imdraw.RoundEndShape
	// imd.Push(groundLine.p1, groundLine.p2)
	// imd.Line(15)

	type platform struct {
		line pixel.Line
	}
	platforms := []platform{
		{line: pixel.L(pixel.V(-400, -70), pixel.V(400, -50))},
		{line: pixel.L(pixel.V(450, -20), pixel.V(900, 100))},
		{line: pixel.L(pixel.V(900, 100), pixel.V(1500, 100))},
	}

	// MAIN LOOP ////////////////////////////////////
	for !win.Closed() {

		// TIME
		dt := time.Since(last).Seconds()
		last = time.Now()

		// CAMERA
		// Make camera follow the player
		cameraPosition = pixel.Lerp(cameraPosition, playerRec.Center(), 1-math.Pow(1.0/128, dt))
		// Create camera matrix to translate all sprites by
		cam := pixel.IM.Moved(win.Bounds().Center().Sub(cameraPosition))
		// Translates screen space by the cam matrix.
		win.SetMatrix(cam)

		// CONTROLS
		switch {
		case win.Pressed(pixelgl.KeyRight):
			playerAcc.X = 5.0
		case win.Pressed(pixelgl.KeyLeft):
			playerAcc.X = -5.0
		default:
			playerAcc.X = 0.0
		}
		if win.Pressed(pixelgl.KeyUp) && touchingGround {
			playerVel.Y = 15
		}

		// If touching ground
		for _, p := range platforms {
			line := p.line
			//lineSlope = GetIntersectingLineSlope(line, playerRec)
			if playerRec.IntersectLine(line) == pixel.ZV {
				playerAcc.Y = -1.5
				touchingGround = false
			} else if playerVel.Y <= 0 {
				touchingGround = true
				playerAcc.Y = 0.0
				playerVel.Y = 0.0
				//playerRec = playerRec.Moved(pixel.Lerp(playerRec.Center(), MoveRectangleUp(line, playerRec), 1-math.Pow(1.0/128, dt)))
				playerRec = playerRec.Moved(MoveRectangleUp(line, playerRec))
				break
			}
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
		fmt.Printf(" Player position: (%.2f, %.2f)\r", playerRec.Center().X, playerRec.Center().Y)

		// DRAW SPRITES
		// Clear the screen and redraw the sprite
		imd.Clear()
		win.Clear(color.Black)
		imd.Color = pixel.RGB(255, 255, 255)
		for _, p := range platforms {
			line := p.line
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(line.A, line.B)
			imd.Line(5)
		}
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
