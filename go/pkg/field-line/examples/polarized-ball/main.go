package main

import (
	"flag"
	"math"

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

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.5,
		FadingGamma: 1.2,
		CameraOrbit: fieldline.NewCameraOrbit(30, 180),
	}
	var trajs []fieldline.Trajectory
	for _, ctr := range centers {
		for j := 0; j < samples; j++ {
			theta := dtheta * float64(j)
			trajs = append(trajs, fieldline.Trajectory{
				Start: ctr.Add(fieldline.Vec3{0, r * math.Sin(theta), r * math.Cos(theta)}),
				AtEnd: atEnd,
				Color: fieldline.RandColor(),
			})
		}
	}

	fieldline.Run(opts, trajs)
}
