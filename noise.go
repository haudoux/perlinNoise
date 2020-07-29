package main

import (
	"fmt"

	noise "github.com/haudoux/perlinNoise/pkg"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeigth int = 800, 600

func lerp(b1, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}
func colorLerp(c1, c2 color, pct float32) color {
	return color{lerp(c1.red, c2.red, pct), lerp(c1.green, c2.green, pct), lerp(c1.blue, c2.blue, pct)}
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}
func getDualGradient(c1, c2, c3, c4 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)-float32(0.5))
		}
	}
	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}
func rescaleAndDraw(noises []float32, min, max float32, gradient []color) []byte {
	scale := 255.0 / (max - min)
	result := make([]byte, winWidth*winHeigth*4)
	offset := min * scale

	for i := range noises {
		noises[i] = noises[i]*scale - offset
		c := gradient[clamp(0, 255, int(noises[i]))]
		p := i * 4
		result[p] = c.red
		result[p+1] = c.green
		result[p+2] = c.blue
	}
	return result
}

type color struct {
	red, green, blue byte
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.red
		pixels[index+1] = c.green
		pixels[index+2] = c.blue
	}
}
func main() {
	window, err := sdl.CreateWindow("Noise", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeigth), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()
	if err != nil {
		fmt.Println(err)
		return
	}

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeigth))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer texture.Destroy()

	keyState := sdl.GetKeyboardState()
	octaves := 3
	lacunarity := float32(3.0)
	gain := float32(0.2)
	frequency := float32(0.01)

	/*for y := 0; y < winHeigth; y++ {
		for x := 0; x < winWidth; x++ {
			setPixel(x, y, color{255, 100, 0}, pixels)
		}
	}*/
	mult := 1
	//Mult : 1 frequency : 0.010000 lacunarity : 3.000000 gain 0.680000, octaves : 3
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if keyState[sdl.SCANCODE_0] != 0 {
			mult *= -1
			fmt.Printf("Mult : %d frequency : %f lacunarity : %f gain %f, octaves : %d \n", mult, frequency, lacunarity, gain, octaves)
		}
		if keyState[sdl.SCANCODE_1] != 0 {
			octaves = octaves + mult
			fmt.Printf("Mult : %d frequency : %f lacunarity : %f gain %f, octaves : %d \n", mult, frequency, lacunarity, gain, octaves)

		}
		if keyState[sdl.SCANCODE_2] != 0 {
			frequency = frequency + 0.001*float32(mult)
			fmt.Printf("Mult : %d frequency : %f lacunarity : %f gain %f, octaves : %d \n", mult, frequency, lacunarity, gain, octaves)

		}
		if keyState[sdl.SCANCODE_3] != 0 {
			gain = gain + 0.01*float32(mult)
			fmt.Printf("Mult : %d frequency : %f lacunarity : %f gain %f, octaves : %d \n", mult, frequency, lacunarity, gain, octaves)

		}
		if keyState[sdl.SCANCODE_4] != 0 {
			lacunarity = lacunarity + 0.01*float32(lacunarity)
			fmt.Printf("Mult : %d frequency : %f lacunarity : %f gain %f, octaves : %d \n", mult, frequency, lacunarity, gain, octaves)

		}
		noises, min, max := noise.MakeNoise(noise.FBM, frequency, lacunarity, gain, octaves, winWidth, winHeigth)
		gradient := getDualGradient(color{98, 68, 72}, color{98, 93, 68}, color{68, 98, 93}, color{68, 72, 98})
		pixels := rescaleAndDraw(noises, min, max, gradient)

		texture.Update(nil, pixels, winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}

}
