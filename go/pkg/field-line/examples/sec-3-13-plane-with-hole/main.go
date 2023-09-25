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
	step   = flag.Float64("step", 0.01, "step")
)

func main() {
	flag.Parse()
	const a = 0.35
	const e0 = 1
	const e1 = .2

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		rho, z := p[0], p[1]
		lambda := (z*z + rho*rho - a*a) / (a * a)
		r := math.Sqrt(lambda*lambda + 4*z*z/(a*a))
		c1 := (e0 - e1) * a / math.Pi
		c2 := .25 / math.Sqrt((r-lambda)/2)
		c3 := math.Abs(z) / a / (1 + 2/(r+lambda)) / math.Sqrt(2/(r+lambda)) / ((r + lambda) * (r + lambda))
		sgn := func(v float64) float64 {
			if v >= 0.0 {
				return 1
			}
			return -1
		}
		c4 := sgn(z) / a * math.Atan(math.Sqrt(2/(r+lambda)))
		dlambdadrho := 2 * rho / (a * a)
		dlambdadz := 2 * z / (a * a)
		drdrho := .5 / r * 2 * lambda * dlambdadrho
		drdz := .5 / r * (2*lambda*dlambdadz + 8*z/(a*a))

		erho := -c1 * (c2*(drdrho-dlambdadrho) + c3*(drdrho+dlambdadrho))
		ez := -c1 * (c2*(drdz-dlambdadz) - c4 + c3*(drdz+dlambdadz))
		if z >= 0 {
			ez -= e0
		} else {
			ez -= e1
		}
		return fieldline.Vec3{erho, ez, 0}
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return v.Norm() < 1e-2
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.5,
		FadingGamma: 1,
	}

	minrho := -.95 * a
	maxrho := .95 * a
	samples := 21
	gap := (maxrho - minrho) / float64(samples-1)
	var trajs []fieldline.Trajectory
	for i := 0; i < samples; i++ {
		rho := minrho + gap*float64(i)
		trajs = append(trajs, fieldline.Trajectory{
			Start: fieldline.Vec3{rho, 0.5, 0},
			AtEnd: atEnd,
			Color: fieldline.RandColor(),
		})
	}

	minout := 1.25 * a
	maxout := 0.95
	outsamples := 21
	outgap := (maxout - minout) / float64(outsamples)
	for i := 0; i < outsamples; i++ {
		rho := minout + outgap*float64(i)
		trajs = append(trajs, fieldline.Trajectory{
			Start: fieldline.Vec3{rho, 0.5, 0},
			AtEnd: atEnd,
			Color: fieldline.RandColor(),
		})
		trajs = append(trajs, fieldline.Trajectory{
			Start: fieldline.Vec3{-rho, 0.5, 0},
			AtEnd: atEnd,
			Color: fieldline.RandColor(),
		})
	}
	fieldline.Run(opts, trajs)
}
