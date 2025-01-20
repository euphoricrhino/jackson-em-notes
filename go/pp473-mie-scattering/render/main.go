package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/euphoricrhino/jackson-em-notes/go/pkg/heatmap"
)

var (
	heatmapFile = flag.String("heatmap-file", "", "heatmap file")
	output      = flag.String("output", "", "output file")
	dataFile    = flag.String("data-file", "", "data file")
	count       = flag.Int("count", 1, "number of frames")
	gamma       = flag.Float64("gamma", 1.0, "gamma correction")
	width       = flag.Int("width", 800, "Width of the image")
	height      = flag.Int("height", 800, "Height of the image")
)

// e.g.,
// go run main.go --heatmap-file ../../heatmaps/wikipedia.png --output ./mie-scattered --data-file=../mie-scattered --count 376 --gamma=.5 --width 800 --height 800
func main() {
	flag.Parse()
	hm, err := heatmap.Load(*heatmapFile, *gamma)
	if err != nil {
		panic(fmt.Sprintf("failed to load heatmap: %v", err))
	}

	max, min := math.NaN(), math.NaN()
	frames := make([][]float64, *count)
	for i := 0; i < *count; i++ {
		frames[i] = loadData(fmt.Sprintf("%v-%03v.data", *dataFile, i))
		for _, v := range frames[i] {
			if math.IsNaN(max) || max < v {
				max = v
			}
			if math.IsNaN(min) || min > v {
				min = v
			}
		}
	}

	spread := max - min
	for frame, data := range frames {
		for i := range data {
			data[i] = (data[i] - min) / spread
		}
		savePNG(data, hm, frame)
	}
}

func savePNG(data []float64, hm []color.Color, frame int) {
	img := image.NewRGBA(image.Rect(0, 0, *width, *height))
	for x := 0; x < *width; x++ {
		for y := 0; y < *height; y++ {
			pixel := data[y**width+x]
			pos := int(pixel * float64(len(hm)-1))
			r, g, b, a := hm[pos].RGBA()
			img.SetRGBA64(
				x,
				y,
				color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)},
			)
		}
	}
	filename := fmt.Sprintf("%v-%03v.png", *output, frame)
	out, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to create output file '%v': %v", filename, err))
	}
	defer out.Close()
	if err := png.Encode(out, img); err != nil {
		panic(fmt.Sprintf("failed to encode to PNG: %v", err))
	}
}

func loadData(filename string) []float64 {
	f, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %v", err))
	}
	defer f.Close()

	data := make([]float64, *width**height)
	for i := range data {
		if err := binary.Read(f, binary.LittleEndian, &data[i]); err != nil {
			panic(fmt.Sprintf("failed to read data: %v", err))
		}
	}
	return data
}
