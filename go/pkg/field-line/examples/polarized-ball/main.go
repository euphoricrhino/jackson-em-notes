package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	fieldline "github.com/euphoricrhino/jackson-em-notes/go/pkg/field-line"
)

var (
	output = flag.String("output", "", "output file")
	width  = flag.Int("width", 800, "output width")
	height = flag.Int("height", 800, "output height")
	step   = flag.Float64("step", 0.005, "step")
)

func main() {
	flag.Parse()
	e := 1.0
	l := 1.2
	a := .5
	b := 3.0 * e / (l + 2.0)
	c := (l - 1.0) / (l + 2.0) * e

	tangentIn := func(p fieldline.Vec3) fieldline.Vec3 {
		return fieldline.Vec3{b, 0.0, 0.0}
	}
	tangentOut := func(p fieldline.Vec3) fieldline.Vec3 {
		x, y, z := p[0], p[1], p[2]
		r := math.Sqrt(x*x + y*y + z*z)
		r5 := math.Pow(r, 5)
		return fieldline.Vec3{
			e - c*(math.Pow(r, -3.0)-3*x*x/r5),
			c * 3.0 * x * y / r5,
			c * 3.0 * x * z / r5,
		}
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		x, y, z := p[0], p[1], p[2]
		r2 := x*x + y*y + z*z
		if r2 < a*a {
			return tangentIn(p)
		}
		return tangentOut(p)
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return false
	}
	rand.Seed(time.Now().UnixNano())
	var trajs []fieldline.Trajectory
	samples := 60
	r := a / 2.0
	dtheta := math.Pi * 2.0 / float64(samples)
	startx := -0.6
	centers := []fieldline.Vec3{
		{startx, a, 0},
		{startx, 0, a},
		{startx, -a, 0},
		{startx, 0, -a},
	}
	for _, ctr := range centers {
		for i := 0; i < samples; i++ {
			theta := dtheta * float64(i)
			trajs = append(trajs, fieldline.Trajectory{
				Start: ctr.Add(fieldline.Vec3{0, r * math.Sin(theta), r * math.Cos(theta)}),
				AtEnd: atEnd,
				Color: fieldline.RandColor(),
			})
		}
	}

	camCircleAngle := math.Pi * .17
	camCircleZ := fieldline.Vec3{-math.Sin(camCircleAngle), math.Cos(camCircleAngle), 0}
	camCircleX := fieldline.Vec3{0, 0, 1}
	camCircleY := camCircleZ.Cross(camCircleX)

	frames := 180
	dcamtheta := math.Pi * 2.0 / float64(frames)
	for f := 0; f < frames; f++ {
		camtheta := dcamtheta * float64(f)
		camPos := camCircleX.Scale(math.Cos(camtheta))
		camPos = camPos.Add(camCircleY.Scale(math.Sin(camtheta)))
		camRight := camCircleX.Scale(-math.Sin(camtheta))
		camRight = camRight.Add(camCircleY.Scale(math.Cos(camtheta)))
		opts := fieldline.Options{
			OutputFile:  fmt.Sprintf("%v-%03d.png", *output, f),
			Width:       *width,
			Height:      *height,
			Step:        *step,
			TangentAt:   tangentAt,
			Camera:      fieldline.NewCamera(camPos, camRight),
			LineWidth:   1.5,
			FadingGamma: 1.2,
		}
		fieldline.Run(opts, trajs)
		fmt.Println(opts.OutputFile)
	}
}
