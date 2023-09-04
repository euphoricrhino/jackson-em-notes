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

	b := 0.6
	a := b * 0.9
	mu := 100.0
	den := (mu+1)*(mu+1)*b*b - (mu-1)*(mu-1)*a*a
	inFactor := (mu*mu - 1) * (b*b - a*a) / den / (a * a)
	ringFactor1 := 2 * (mu - 1) / den
	ringFactor2 := 2 * (mu + 1) * b * b / den
	outFactor := 4 * mu * b * b / den
	// Potential is lower in the center, so we are tracing backwards.
	scale := -0.1
	tangentIn := func(p fieldline.Vec3) fieldline.Vec3 {
		rho := p.Norm()
		cp, sp := p[0]/rho, p[1]/rho
		rho2 := rho * rho
		return fieldline.Vec3{
			-2 * cp * sp / rho2,
			(cp*cp-sp*sp)/rho2 - inFactor,
			0,
		}.Scale(scale)
	}

	tangentRing := func(p fieldline.Vec3) fieldline.Vec3 {
		rho := p.Norm()
		cp, sp := p[0]/rho, p[1]/rho
		rho2 := rho * rho
		return fieldline.Vec3{
			-2 * cp * sp / rho2 * ringFactor2,
			(cp*cp-sp*sp)/rho2*ringFactor2 + ringFactor1,
			0,
		}.Scale(scale)
	}
	tangentOut := func(p fieldline.Vec3) fieldline.Vec3 {
		rho := p.Norm()
		cp, sp := p[0]/rho, p[1]/rho
		rho2 := rho * rho
		return fieldline.Vec3{
			-2 * cp * sp / rho2 * outFactor,
			(cp*cp - sp*sp) / rho2 * outFactor,
			0,
		}.Scale(scale)
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		r := p.Norm()
		if r <= a {
			return tangentIn(p)
		}
		if r <= b {
			return tangentRing(p)
		}
		return tangentOut(p)
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return v[1] < 0 && p[1] < 1e-3 || v.Norm() < 1e-10
	}

	// Deviation from phi=𝜋/2.
	phis := []float64{0, 0.25, 0.5, 0.75, 1, 2, 3, 4, 5, 6, 6.5, 7, 7.5, 8, 8.5, 9, 9.5, 9.75, 10, 10.5, 11, 12, 13, 14, 15, 18, 20, 23, 25, 30, 35, 40, 50, 60}
	var trajs []fieldline.Trajectory
	sr := 0.05
	for i, phiDeg := range phis {
		phi := (90.0 - phiDeg) * math.Pi / 180.0
		color := fieldline.RandColor()
		traj := fieldline.Trajectory{
			Start: fieldline.Vec3{sr * math.Cos(phi), sr * math.Sin(phi), 0},
			AtEnd: atEnd,
			Color: color,
		}
		traj.AddSymmetry(
			func(v fieldline.Vec3) fieldline.Vec3 {
				v[1] = -v[1]
				return v
			}, color,
		)
		if i != 0 {
			leftColor := fieldline.RandColor()
			traj.AddSymmetry(
				func(v fieldline.Vec3) fieldline.Vec3 {
					v[0] = -v[0]
					return v
				}, leftColor,
			)
			traj.AddSymmetry(
				func(v fieldline.Vec3) fieldline.Vec3 {
					v[0] = -v[0]
					v[1] = -v[1]
					return v
				}, leftColor,
			)
		}
		trajs = append(trajs, traj)
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.0,
		FadingGamma: .25,
	}

	fieldline.Run(opts, trajs)
}
