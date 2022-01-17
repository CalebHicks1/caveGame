package main

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
