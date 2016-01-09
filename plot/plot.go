// Package plot provides a plot block for loops
package plot

import (
	"image"
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
func (p *Plot) Step(in, out []float64) {
	if p.img == nil {
		width, height := p.Size.X, p.Size.Y
		if width <= 0 || height <= 0 {
			p.Size = image.Point{512, 512}
		}
		p.img = image.NewRGBA(image.Rect(0, 0, p.Size.X, p.Size.Y))
	}
	for i, v := range in {
		y := p.Size.Y/2 - int(v*p.Scale/float64(2*p.Size.Y))
		p.img.Set(p.x, y, Colors[i%len(Colors)])
	}
	p.x++
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
