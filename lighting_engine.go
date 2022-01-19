package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var fragmentShader = `
#version 330 core
// The first line in glsl source code must always start with a version directive as seen above.

// vTexCoords are the texture coordinates, provided by Pixel
in vec2  vTexCoords;

// fragColor is what the final result is and will be rendered to your screen.
out vec4 fragColor;

// uTexBounds is the texture's boundries, provided by Pixel.
uniform vec4 uTexBounds;

// uTexture is the actualy texture we are sampling from, also provided by Pixel.
uniform sampler2D uTexture;
uniform sampler2D lightMap;

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	if(texture(uTexture, t).a	< 0.001) {
		discard;
	}
	
	fragColor = texture(uTexture, t);
}
`

/*
TODO:
1. set up canvas to draw light map to
2. draw shapes to canvas
3. combine with scene
4. find way to ambiently draw rest of background

*/

// Make sure lightmap is off screen
func GenerateLightMap(canvas pixelgl.Canvas, imd imdraw.IMDraw, lightSprite pixel.Sprite) {
	//rec := pixel.R(-100, -100, 100, 100)
	imd.Color = pixel.RGB(255, 0, 0)
	imd.Push(canvas.Bounds().Center(), canvas.Bounds().Center().ScaledXY(pixel.V(2, 2)))
	imd.Rectangle(0)
	imd.Push(canvas.Bounds().Min, canvas.Bounds().Center())
	imd.Rectangle(0)
	canvas.SetComposeMethod(pixel.ComposePlus)
	imd.Draw(&canvas)
	canvas.SetComposeMethod(pixel.ComposeOut)
	lightSprite.Draw(&canvas, pixel.IM.Moved(canvas.Bounds().Center()).Scaled(canvas.Bounds().Center(), 1.2))
}
