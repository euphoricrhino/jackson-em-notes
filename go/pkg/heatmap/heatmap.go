package heatmap

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Load loads the heatmap from the given file and returns the color spectrum.
func Load(file string, gamma float64) ([]color.Color, error) {
	// Load heatmap file.
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open heatmap file: %v", err)
	}
	defer f.Close()
	hm, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode heatmap as PNG: %v", err)
	}
	rect := hm.Bounds()
	width := rect.Max.X - rect.Min.X
	heatmap := make([]color.Color, width)
	for i := 0; i < width; i++ {
		r, g, b, _ := hm.At(i+rect.Min.X, rect.Min.Y).RGBA()
		max := float64(math.MaxUint16)
		r16 := uint16(math.Pow(float64(r)/max, gamma) * max)
		g16 := uint16(math.Pow(float64(g)/max, gamma) * max)
		b16 := uint16(math.Pow(float64(b)/max, gamma) * max)
		heatmap[i] = color.RGBA64{R: r16, G: g16, B: b16, A: math.MaxUint16}
	}
	return heatmap, nil
}
