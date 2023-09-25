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

func legend(l int, x float64) ([]float64, []float64) {
	val := make([]float64, l+1)
	der := make([]float64, l+1)
	val[0] = 1
	val[1] = x
	der[0] = 0
	der[1] = 1
	for ll := 2; ll <= l; ll++ {
		val[ll] = float64(2*ll-1)/float64(ll)*x*val[ll-1] - float64(ll-1)/float64(ll)*val[ll-2]
		der[ll] = val[ll-1] + 2*x*der[ll-1] - der[ll-2]
	}
	return val, der
}

func main() {
	flag.Parse()
	const eps = 3.5
	const a = 0.4
	const d = 0.6

	const maxl = 20
	const scale = 1 / 2000.0
	tangentIn := func(p fieldline.Vec3) fieldline.Vec3 {
		r := math.Sqrt(p.Dot(p))
		ct := p[0] / r
		st := math.Sqrt(1 - ct*ct)
		val, der := legend(maxl, ct)
		er, et := 0.0, 0.0
		rd := r / d
		d2 := d * d
		for l := 1; l <= maxl; l++ {
			cc := float64(2*l+1) / ((eps+1.0)*float64(l) + 1) * math.Pow(rd, float64(l-1))
			er -= cc * float64(l) * val[l] / d2
			et += cc * der[l] * st / d2
		}
		ret := fieldline.Vec3{
			er*ct - et*st,
			er*st + et*ct,
			0.0,
		}.Scale(scale)
		return ret
	}

	tangentOut := func(p fieldline.Vec3) fieldline.Vec3 {
		// From point charge.
		disp := p.Subtract(fieldline.Vec3{d, 0, 0})
		ret1 := disp.Scale(1.0 / (math.Pow(disp.Dot(disp), 1.5))).Scale(scale)

		r := math.Sqrt(p.Dot(p))
		ct := p[0] / r
		st := math.Sqrt(1 - ct*ct)
		val, der := legend(maxl, ct)
		er, et := 0.0, 0.0
		aadr := a * a / (d * r)
		adr2 := a / (d * r * r)
		for l := 1; l <= maxl; l++ {
			cc := (eps - 1) * float64(l) / ((eps+1)*float64(l) + 1) * math.Pow(aadr, float64(l)) * adr2
			er -= float64(l+1) * cc * val[l]
			et -= cc * der[l] * st
		}
		ret2 := fieldline.Vec3{
			er*ct - et*st,
			er*st + et*ct,
			0.0,
		}.Scale(scale)
		return ret1.Add(ret2)
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		if p.Dot(p) < a*a {
			return tangentIn(p)
		}
		return tangentOut(p)
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return v.Dot(v) < 1e-8
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.5,
		FadingGamma: .3,
	}

	ctr := fieldline.Vec3{d, 0, 0}
	sr := 0.02
	trajs := []fieldline.Trajectory{
		{Start: ctr.Add(fieldline.Vec3{sr, 0, 0}), AtEnd: atEnd, Color: fieldline.RandColor()},
		{Start: ctr.Add(fieldline.Vec3{-sr, 0, 0}), AtEnd: atEnd, Color: fieldline.RandColor()},
	}
	thetas := []float64{30, 60, 90, 95, 100, 105, 110, 115, 120, 125, 130, 135, 140, 145, 150, 155, 160, 165, 170, 175}
	for _, theta := range thetas {
		rad := theta * math.Pi / 180
		traj := fieldline.Trajectory{
			Start: ctr.Add(fieldline.Vec3{sr * math.Cos(rad), sr * math.Sin(rad), 0}),
			AtEnd: atEnd,
			Color: fieldline.RandColor(),
		}
		samples := 2
		for j := 1; j < samples; j++ {
			phi := math.Pi * 2 * float64(j) / float64(samples)
			traj.AddSymmetry(
				func(p fieldline.Vec3) fieldline.Vec3 {
					rho := math.Sqrt(p[1]*p[1] + p[2]*p[2])
					return fieldline.Vec3{
						p[0],
						rho * math.Cos(phi),
						rho * math.Sin(phi),
					}
				},
				fieldline.RandColor(),
			)
		}
		trajs = append(trajs, traj)
	}

	fieldline.Run(opts, trajs)
}
