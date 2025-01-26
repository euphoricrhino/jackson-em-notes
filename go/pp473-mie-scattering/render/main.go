package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/euphoricrhino/jackson-em-notes/go/pkg/heatmap"
)

var (
	heatmapFile = flag.String("heatmap-file", "", "heatmap file")
	n           = flag.String("n", "1.2+0.2i", "refractive index")
	rStart      = flag.Float64("r-start", 0.25, "start radius")
	rInc        = flag.Float64("r-inc", 0.01, "increment radius")
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

func postEdit(img draw.Image, frame int) {
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetLineWidth(1)

	draw2d.SetFontFolder("/Users/xni/Library/Fonts")
	draw2d.SetFontNamer(func(_ draw2d.FontData) string { return "MonoLisaVariableNormal.ttf" })

	text := fmt.Sprintf("n=%v", *n)
	text += fmt.Sprintf(", R=%.2fλ", *rStart+float64(frame)*(*rInc))
	textColor := color.RGBA{0, 0xcc, 0xcc, 0xff}
	gc.SetFillColor(textColor)
	gc.SetStrokeColor(textColor)
	gc.SetDPI(288)
	gc.SetFontSize(3.5)
	gc.FillStringAt(text, 20.0, 20.0)
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

	postEdit(img, frame)
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
