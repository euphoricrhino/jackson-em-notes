package main

import (
	"flag"
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
	l := 2.5
	a := .25
	b := a * 2
	denom := (l+1.0)*(l+1.0)*b*b - (l-1.0)*(l-1.0)*a*a
	b1 := -4.0 * l * e * b * b / denom
	d1 := e * b * b * (l*l - 1.0) * (b*b - a*a) / denom
	f1 := -2.0 * (l + 1.0) * e * b * b / denom
	h1 := -2.0 * (l - 1.0) * e * a * a * b * b / denom

	tangentA := func(p fieldline.Vec3) fieldline.Vec3 {
		return fieldline.Vec3{-b1, 0.0, 0.0}
	}
	tangentB := func(p fieldline.Vec3) fieldline.Vec3 {
		x, y := p[0], p[1]
		r2 := x*x + y*y
		r4 := r2 * r2
		return fieldline.Vec3{
			-f1 - h1/r2 + 2*h1*x*x/r4,
			2.0 * h1 * x * y / r4,
			0.0,
		}
	}
	tangentC := func(p fieldline.Vec3) fieldline.Vec3 {
		x, y := p[0], p[1]
		r2 := x*x + y*y
		r4 := r2 * r2
		return fieldline.Vec3{
			e - d1/r2 + 2*d1*x*x/r4,
			2.0 * d1 * x * y / r4,
			0.0,
		}
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		x, y := p[0], p[1]
		r2 := x*x + y*y
		if r2 < a*a {
			return tangentA(p)
		}
		if r2 < b*b {
			return tangentB(p)
		}
		return tangentC(p)
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return false
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.5,
		FadingGamma: 0.8,
	}

	rand.Seed(time.Now().UnixNano())
	var trajs []fieldline.Trajectory
	maxy, miny := 0.98, -0.98
	startx := -0.98
	lines := 61
	gap := (maxy - miny) / float64(lines-1)
	for i := 0; i < lines; i++ {
		starty := miny + float64(i)*gap
		trajs = append(trajs, fieldline.Trajectory{
			Start: fieldline.Vec3{startx, starty, 0.0},
			AtEnd: atEnd,
			Color: fieldline.RandColor(),
		})
	}

	fieldline.Run(opts, trajs)
}
