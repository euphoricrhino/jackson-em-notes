package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/draw"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"

	fieldrenderer "github.com/euphoricrhino/jackson-em-notes/go/pkg/field-renderer"
)

// Example command:
// go run main.go --heatmap ../heatmaps/wikipedia.png --output diff --slit-width 6 --slit-height 4 --z=200 --gamma .4
var (
	width      = flag.Int("width", 800, "Width of the image")
	height     = flag.Int("height", 800, "Height of the image")
	heatmap    = flag.String("heatmap", "", "heatmap file")
	gamma      = flag.Float64("gamma", 1.0, "gamma correction")
	output     = flag.String("output", "", "output file")
	slitWidth  = flag.Float64("slit-width", 0.0, "slit width in units of lambda")
	slitHeight = flag.Float64("slit-height", 0.0, "slit height in units of lambda")
	z          = flag.Float64("z", 10.0, "observation point z in units of lambda")
)

const (
	screenWidthLambdas = 1000
	threshold          = 1e-8
)

func main() {
	flag.Parse()

	for f := 0; f <= 180; f++ {
		saveFrame(f)
	}
}

func saveFrame(f int) {
	beta := float64(f) * math.Pi / 180
	if err := fieldrenderer.Run(fieldrenderer.Options{
		HeatMapFile: *heatmap,
		OutputFile:  fmt.Sprintf("%v-%03d.png", *output, f),
		Gamma:       *gamma,
		Width:       *width,
		Height:      *height,
		Field:       renderField(beta),
		PostEdit:    postEdit(beta),
	}); err != nil {
		panic(err)
	}
}

func renderField(beta float64) func(int, int) float64 {
	return func(x, y int) float64 {
		sbeta, cbeta := math.Sin(beta), math.Cos(beta)
		fx := float64(x-*width/2) / float64(*width) * screenWidthLambdas
		fy := float64(*height/2-y) / float64(*height) * screenWidthLambdas
		rho := math.Sqrt(fx*fx + fy*fy)
		r := math.Sqrt(fx*fx + fy*fy + *z**z)

		ct := *z / r
		st := math.Sqrt(1 - ct*ct)
		var f1, f2, f3 float64

		ka := 2 * math.Pi * *slitWidth
		kb := 2 * math.Pi * *slitHeight

		// Deal with singularities by returning the limit.
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
}

func postEdit(beta float64) func(img draw.Image) {
	return func(img draw.Image) {
		yaxisColor := color.RGBA{0, 0xcc, 0xcc, 0xff}
		gc := draw2dimg.NewGraphicContext(img)
		gc.SetLineWidth(.8)
		gc.SetStrokeColor(yaxisColor)

		cx, cy := 100.0, 100.0
		d := 55.0
		gc.MoveTo(cx, cy-d)
		gc.LineTo(cx, cy+d)

		gc.Stroke()

		v := 45.0
		vectorColor := color.RGBA{0xcc, 0, 0, 0xff}
		gc.SetLineWidth(1)
		gc.SetStrokeColor(vectorColor)
		sbeta, cbeta := math.Sin(beta), math.Cos(beta)
		tip1x, tip1y := cx+v*sbeta, cy-v*cbeta
		tip2x, tip2y := cx-v*sbeta, cy+v*cbeta
		gc.MoveTo(tip1x, tip1y)
		gc.LineTo(tip2x, tip2y)
		gc.Stroke()

		s := 20.0
		ang := 15 * math.Pi / 180

		fx1, fy1 := tip1x-s*math.Sin(ang+beta), tip1y+s*math.Cos(ang+beta)
		fx2, fy2 := tip1x-s*math.Sin(beta-ang), tip1y+s*math.Cos(beta-ang)
		gc.SetFillColor(vectorColor)
		gc.MoveTo(tip1x, tip1y)
		gc.LineTo(fx1, fy1)
		gc.LineTo(fx2, fy2)
		gc.Close()
		gc.FillStroke()

		fx1, fy1 = tip2x+s*math.Sin(beta-ang), tip2y-s*math.Cos(beta-ang)
		fx2, fy2 = tip2x+s*math.Sin(beta+ang), tip2y-s*math.Cos(beta+ang)

		gc.SetFillColor(vectorColor)
		gc.MoveTo(tip2x, tip2y)
		gc.LineTo(fx1, fy1)
		gc.LineTo(fx2, fy2)
		gc.Close()
		gc.FillStroke()

		draw2d.SetFontFolder("/Users/xni/Library/Fonts")
		draw2d.SetFontNamer(func(_ draw2d.FontData) string { return "MonoLisaVariableNormal.ttf" })
		text := fmt.Sprintf("β=%.0f°", beta*180.0/math.Pi)
		gc.SetFillColor(vectorColor)
		gc.SetStrokeColor(vectorColor)
		gc.SetDPI(288)
		gc.SetFontSize(5.5)
		gc.FillStringAt(text, cx+d, cy-d)
	}
}
