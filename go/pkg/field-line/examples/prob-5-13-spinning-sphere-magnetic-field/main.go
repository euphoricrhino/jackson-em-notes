package main

import (
	"flag"
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

	b0 := 1.0
	a := 0.2
	tangentIn := func(p fieldline.Vec3) fieldline.Vec3 {
		return fieldline.Vec3{0, 2 * b0 * a, 0}
	}
	tangentOut := func(p fieldline.Vec3) fieldline.Vec3 {
		r := p.Norm()
		st, ct := p[0]/r, p[1]/r
		v := fieldline.Vec3{
			3 * st * ct,
			2*ct*ct - st*st,
			0.0,
		}
		return v.Scale(a * math.Pow(a/r, 3))
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		r := p.Norm()
		if r <= a {
			return tangentIn(p)
		}
		return tangentOut(p)
	}

	makeAtEnd := func(start fieldline.Vec3) func(p, v fieldline.Vec3) bool {
		// We end the tracing when the field line curved back to the z=0 plane.
		return func(p, v fieldline.Vec3) bool {
			return v[1] < 0 && p[1] < 1e-3
		}
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.0,
		FadingGamma: 0.5,
	}

	rand.Seed(time.Now().UnixNano())
	var trajs []fieldline.Trajectory
	samples := 21
	for i := 1; i < samples-1; i++ {
		theta := math.Pi * float64(i) / float64(samples-1)
		start := fieldline.Vec3{a * math.Cos(theta), 0.0, 0.0}
		color := fieldline.RandColor()
		traj := fieldline.Trajectory{
			Start: start,
			AtEnd: makeAtEnd(start),
			Color: color,
		}
		// Lower half is reflected using symmetry.
		traj.AddSymmetry(
			func(v fieldline.Vec3) fieldline.Vec3 {
				v[1] = -v[1]
				return v
			},
			color,
		)
		trajs = append(trajs, traj)
	}

	fieldline.Run(opts, trajs)
}
