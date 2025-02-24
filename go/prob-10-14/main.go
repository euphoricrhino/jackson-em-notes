package main

import (
	"flag"
	"math"

	fieldrenderer "github.com/euphoricrhino/jackson-em-notes/go/pkg/field-renderer"
)

// Example command:
// go run main.go --heatmap ../heatmaps/wikipedia.png --output diff-6-4-135.png --slit-width 6 --slit-height 4 --z=200 --gamma .45 --beta-deg 135
var (
	width      = flag.Int("width", 800, "Width of the image")
	height     = flag.Int("height", 800, "Height of the image")
	heatmap    = flag.String("heatmap", "", "heatmap file")
	gamma      = flag.Float64("gamma", 1.0, "gamma correction")
	output     = flag.String("output", "", "output file")
	slitWidth  = flag.Float64("slit-width", 0.0, "slit width in units of lambda")
	slitHeight = flag.Float64("slit-height", 0.0, "slit height in units of lambda")
	z          = flag.Float64("z", 10.0, "observation point z in units of lambda")
	betaDeg    = flag.Float64("beta-deg", 0.0, "angle of the incident beam in degrees")
)

const (
	screenWidthLambdas = 1000
	threshold          = 1e-8
)

func main() {
	flag.Parse()

	if err := fieldrenderer.Run(fieldrenderer.Options{
		HeatMapFile: *heatmap,
		OutputFile:  *output,
		Gamma:       *gamma,
		Width:       *width,
		Height:      *height,
		Field:       field,
		PostEdit:    nil,
	}); err != nil {
		panic(err)
	}
}

func field(x, y int) float64 {
	sbeta, cbeta := math.Sin(*betaDeg*math.Pi/180), math.Cos(*betaDeg*math.Pi/180)
	fx := float64(x-*width/2) / float64(*width) * screenWidthLambdas
	fy := float64(*height/2-y) / float64(*height) * screenWidthLambdas
	rho := math.Sqrt(fx*fx + fy*fy)
	r := math.Sqrt(fx*fx + fy*fy + *z**z)

	ct := *z / r
	st := math.Sqrt(1 - ct*ct)
	var f1, f2, f3 float64

	ka := 2 * math.Pi * *slitWidth
	kb := 2 * math.Pi * *slitHeight

	if rho < threshold {
		f1 = 1.0
		f2 = 1.0
		f3 = ka / 2 * kb / 2
	} else {
		cp := fx / rho
		sp := fy / rho

		f1 = sbeta*cp + cbeta*sp
		f1 *= f1
		f1 *= st * st
		f1 += ct * ct
		f2 = math.Pow(st, 4)
		if math.Abs(cp) < threshold {
			f2 *= sp * sp
			f3 = ka / 2 * st * math.Sin(kb/2*st*sp)
		} else if math.Abs(sp) < threshold {
			f2 *= cp * cp
			f3 = kb / 2 * st * math.Sin(ka/2*st*cp)
		} else {
			f2 *= sp * sp * cp * cp
			f3 = math.Sin(ka/2*st*cp) * math.Sin(kb/2*st*sp)
		}
	}

	return f1 / f2 * f3 * f3
}
