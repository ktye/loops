// Package plot provides a plot block for loops
package plot

import (
	"fmt"
	"image"
	"image/draw"
	"image/color"
	"image/png"
	"os"
)

// Colors are the individual line colors, cycled through per channel.
var Colors []color.Color = []color.Color{
	color.RGBA{0, 0, 255, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{255, 0, 0, 255},
	color.RGBA{255, 0, 255, 255},
	color.RGBA{0, 255, 255, 255},
	color.RGBA{255, 255, 0, 255},
}

// Plot is a terminal block which writes a png image.
// The plot is very primitive: a pixel per value.
// The x-axis is stretched horizontally at the center of the image
// with one pixel per DT.
type Plot struct {
	NumChannels int         // Number of input channels.
	Scale       float64     // Y-axis contains data from [-Scale,+Scale]
	Size        image.Point // Image dimensions.
	img         *image.RGBA // Image structure.
	x           int         // current x pixel position
}

func (p *Plot) Inputs() int {
	return p.NumChannels
}
func (p *Plot) Outputs() int {
	return 0
}
func (p *Plot) Step(in, out []float64) bool {
	// At the first step, initialize the image.
	if p.img == nil {
		width, height := p.Size.X, p.Size.Y
		if width <= 0 || height <= 0 {
			p.Size = image.Point{512, 512}
		}
		p.img = image.NewRGBA(image.Rect(0, 0, p.Size.X, p.Size.Y))
		draw.Draw(p.img, p.img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

		// Set default y-scale to [+1,-1].
		if p.Scale == 0 {
			p.Scale = 1
		}

		// Draw the x-axis.
		for i := 0; i < p.Size.X; i++ {
			p.img.Set(i, p.Size.Y/2, color.Black)
		}
	}

	// Draw one pixel per input channel, in it's own color.
	for i, v := range in {
		y := p.Size.Y/2 - int(v*float64(p.Size.Y)/(2*p.Scale))
		fmt.Println("Plot:", v, y, p.Size)
		p.img.Set(p.x, y, Colors[i%len(Colors)])
	}
	p.x++
	return true
}

// Write stores the image as a png file.
// It must be called manually at the end of the simulation.
func (p *Plot) Write(filename string) error {
	if f, err := os.Create(filename); err != nil {
		return err
	} else {
		defer f.Close()
		return png.Encode(f, p.img)
	}
}
