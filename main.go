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
	frames               = 0
	second               = time.Tick(time.Second)
	last                 = time.Now()
	pWidth, pHeight      = 60.0, 100.0 // Width and height of player character
	player_starting_pos  = pixel.V(0, 1500)
	camera_position      = player_starting_pos
	camera_zoom          = 1.0
	playerRec            = pixel.R(-pWidth/2, -pHeight/2, pWidth/2, pHeight/2).Moved(player_starting_pos) // Used to handle player physics
	imd                  = imdraw.New(nil)                                                                // Used to draw shapes (player rectangle, platforms)
	playerVel, playerAcc = pixel.ZV, pixel.ZV
	touching_ground      = false
	transition_speed     = 12.0
	debug                = false // Press D to enter debug mode
	line_completed       = true  // Used to track if a new platform is being drawn

	platforms = ReadPlatformData("world_data/level1_platforms.yaml")
)

// Setup ground lines
type Platform struct {
	Line pixel.Line
}

// Main Functions /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
This is where all the game code runs, it is basically
the new main function since main is used by pixel.
*/
func run() {

	//hardcode platforms
	// platforms = append(platforms,
	// 	Platform{Line: pixel.L(pixel.V(-400, 1400), pixel.V(400, 1400))},
	// )

	// load the picture and create the player sprite
	playerImg, err := loadPicture("assets/smile.png")
	if err != nil {
		panic(err)
	}
	playerSprite := pixel.NewSprite(playerImg, playerImg.Bounds())

	// Window configuration
	cfg := pixelgl.WindowConfig{
		Title:     "Caleb's Game",
		Bounds:    pixel.R(0, 0, 900, 700),
		VSync:     true,
		Resizable: true,
	}

	// Create game window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Set up foreground tiles
	// Want to have a list of sprites to draw, and a game location.

	type tile struct {
		sprite *pixel.Sprite
		x      float64
		y      float64
	}

	tileImg, err := loadPicture("assets/tile3.png")
	if err != nil {
		panic(err)
	}

	tile1 := tile{
		sprite: pixel.NewSprite(tileImg, pixel.R(0, 0, 1000, 1000)),
		x:      0,
		y:      0,
	}
	cam := pixel.IM
	point1 := cam.Unproject(win.MousePosition())

	tileImg2, err := loadPicture("assets/tile4.png")
	if err != nil {
		panic(err)
	}

	tile2 := tile{
		sprite: pixel.NewSprite(tileImg2, pixel.R(0, 0, 1000, 1000)),
		x:      1,
		y:      0,
	}

	// MAIN LOOP //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	for !win.Closed() {

		imd.Clear()
		win.Clear(color.Black)

		// TIME
		dt := time.Since(last).Seconds()
		last = time.Now()

		// CAMERA
		// Make camera follow the player
		camera_position = pixel.Lerp(camera_position, playerRec.Center(), 1-math.Pow(1.0/128, dt))
		camera_zoom *= math.Pow(1.1, win.MouseScroll().Y)
		// Create camera matrix to translate all sprites by
		cam = pixel.IM.Scaled(camera_position, camera_zoom).Moved(win.Bounds().Center().Sub(camera_position))
		// Translates screen space by the cam matrix.
		win.SetMatrix(cam)

		// DEBUG
		mouse_circle := pixel.C(cam.Unproject(win.MousePosition()), 15)
		if win.JustPressed(pixelgl.KeyD) {
			debug = !debug
		}

		// draw new lines
		if debug == true {
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				if line_completed {
					point1 = cam.Unproject(win.MousePosition())
					for _, p := range platforms {

						if mouse_circle.Contains(p.Line.A) {
							point1 = p.Line.A
							break
						} else if mouse_circle.Contains(p.Line.B) {
							point1 = p.Line.B
							break
						}
					}

					line_completed = false
				} else {
					point2 := cam.Unproject(win.MousePosition())
					for _, p := range platforms {

						if mouse_circle.Contains(p.Line.A) {
							point2 = p.Line.A
							break
						} else if mouse_circle.Contains(p.Line.B) {
							point2 = p.Line.B
							break
						}
					}
					line_completed = true
					CreateNewPlatform(point1, point2)
				}
			}
		}
		if !line_completed {

			imd.Color = pixel.RGB(0, 255, 0)
			imd.Push(point1, cam.Unproject(win.MousePosition()))
			imd.Line(5)
		}

		// CONTROLS
		switch {
		case win.Pressed(pixelgl.KeyRight):
			playerAcc.X = 5.0
		case win.Pressed(pixelgl.KeyLeft):
			playerAcc.X = -5.0
		default:
			playerAcc.X = 0.0
		}
		if win.Pressed(pixelgl.KeyUp) && touching_ground {
			playerVel.Y = 15
		}

		// If touching ground
		touching_any_platform := false
		for _, p := range platforms {
			line := p.Line
			//lineSlope = GetIntersectingLineSlope(line, playerRec)
			if playerRec.IntersectLine(line) == pixel.ZV {
				playerAcc.Y = -1.5
			} else if playerVel.Y <= 0 {
				touching_any_platform = true
				playerAcc.Y = 0.0
				playerVel.Y = 0.0
				playerRec = playerRec.Moved(MoveRectangleUp(line, playerRec))
				//break
			}
		}
		touching_ground = touching_any_platform

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

		// Draw tile
		tile1.sprite.Draw(win, pixel.IM.Moved(pixel.V(tile1.x*1000, tile1.y*1000)).Scaled(pixel.ZV, 4))
		tile2.sprite.Draw(win, pixel.IM.Moved(pixel.V(tile2.x*1000, tile2.y*1000)).Scaled(pixel.ZV, 4))

		if debug == true {
			imd.Color = pixel.RGB(255, 0, 0)
			for _, p := range platforms {
				line := p.Line
				imd.EndShape = imdraw.RoundEndShape
				imd.Push(line.A, line.B)
				imd.Line(5)
			}

			imd.Color = pixel.RGB(255, 0, 0)
			imd.Push(playerRec.Min, playerRec.Max)
			imd.Rectangle(1)

			imd.Push(mouse_circle.Center)
			imd.Circle(mouse_circle.Radius, 1)
		}

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

// Helper Functions ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

// Creates a new platform and adds it to the platforms array.
func CreateNewPlatform(p1 pixel.Vec, p2 pixel.Vec) {
	new_line := pixel.L(p1, p2)
	new_platform := Platform{Line: new_line}
	platforms = append(platforms, new_platform)
	WritePlatformData(new_platform, "world_data/level1_platforms.yaml")
}
