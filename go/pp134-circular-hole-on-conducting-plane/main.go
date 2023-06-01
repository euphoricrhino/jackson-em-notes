package main

import (
	"flag"
	"math"

	fieldrenderer "github.com/euphoricrhino/jackson-em-notes/go/pkg/field-renderer"
)

var (
	heatmap = flag.String("heatmap", "", "heatmap file")
	output  = flag.String("output", "", "output file")
	gamma   = flag.Float64("gamma", 1.0, "gamma correction")
	width   = flag.Int("width", 640, "output width")
	height  = flag.Int("height", 640, "output height")
)

func main() {
	flag.Parse()

	field := func(x, y int) float64 {
		fx := float64(x) - float64(*width-1)/2
		fy := float64(*height-1-y) - float64(*height-1)/2
		a := float64(*width) / 8
		a2 := a * a
		l := (fy*fy + fx*fx - a2) / a2
		r := math.Sqrt(l*l + 4.0*fy*fy/a2)
		v1 := math.Sqrt((r - l) / 2)
		v2 := math.Abs(fy) / a * math.Atan(math.Sqrt(2/(r+l)))
		const e0 = 100.0
		const e1 = 0.0
		ret := (e0 - e1) * a / math.Pi * (v1 - v2)
		if fy > 0 {
			ret += e0 * fy
		} else {
			ret += e1 * fy
		}
		return ret
	}

	if err := fieldrenderer.Run(fieldrenderer.Options{
		HeatMapFile: *heatmap,
		OutputFile:  *output,
		Gamma:       *gamma,
		Width:       *width,
		Height:      *height,
		Field:       field,
	}); err != nil {
		panic(err)
	}
}
