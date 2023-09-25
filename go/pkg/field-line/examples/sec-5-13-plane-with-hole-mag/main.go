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
	const a2 = a * a
	const h0 = 1
	const hh = h0 * 2 * a / math.Pi

	// See https://github.com/euphoricrhino/jackson-em-notes/blob/main/pdf/ch-5/pp203-circular-hole-conducting-plane-magnetic.pdf
	// for the formulas.
	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		y, z, x := p[0], p[1], p[2]
		rho := math.Sqrt(x*x + y*y)
		sphi, cphi := y/rho, x/rho
		neg := false
		if z < 0 {
			z = -z
			neg = true
		}
		lambda := (z*z + rho*rho - a2) / a2
		r := math.Sqrt(lambda*lambda + 4*z*z/a2)
		v1 := math.Sqrt((r - lambda) / 2)
		v2 := math.Sqrt((r + lambda) / 2)
		c1 := z / (8 * rho) / v1
		c2 := -rho / (8 * a) / (1 + 1/(v2*v2)) / (v2 * v2 * v2)
		c3 := -a / (8 * rho) / v2
		dlambdadz := 2 * z / a2
		dlambdadrho := 2 * rho / a2
		drdz := lambda/r*dlambdadz + 4*z/a2/r
		drdrho := lambda / r * dlambdadrho
		hz := v1/(2*rho) + c1*(drdz-dlambdadz)
		hz += c2 * (drdz + dlambdadz)
		hz += c3 * (drdz + dlambdadz)
		hz *= -sphi

		hrho := -z/(2*rho*rho)*v1 + c1*(drdrho-dlambdadrho)
		hrho += 1/(2*a)*math.Atan(1/v2) + c2*(drdrho+dlambdadrho)
		hrho += a/(2*rho*rho)*v2 + c2*(drdrho+dlambdadrho)
		hrho *= -sphi

		hphi := z / (2 * rho * rho) * v1
		hphi += 1 / (2 * a) * math.Atan(1/v2)
		hphi -= a / (2 * rho * rho) * v2
		hphi *= -cphi

		hx := hrho*cphi - hphi*sphi
		hy := hrho*sphi + hphi*cphi
		if !neg {
			return fieldline.Vec3{
				h0 + hh*hy,
				hh * hz,
				hh * hx,
			}
		}
		return fieldline.Vec3{
			-hh * hy,
			hh * hz,
			-hh * hx,
		}
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return false
	}

	zs := []float64{0.03, 0.05, 0.1, 0.2}
	var zcolors [][3]float64
	for i := 0; i < len(zs); i++ {
		zcolors = append(zcolors, fieldline.RandColor())
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.0,
		FadingGamma: .3,
		CameraOrbit: fieldline.NewCameraOrbit(45, 180),
	}

	samples := 30
	gap := 1.0 / float64(samples)
	var trajs []fieldline.Trajectory
	for i, z := range zs {
		for j := 0; j <= samples; j++ {
			xstart := gap * float64(j)
			traj := fieldline.Trajectory{
				Start: fieldline.Vec3{-.99, z, xstart},
				AtEnd: atEnd,
				Color: zcolors[i],
			}
			if j != 0 {
				traj.AddSymmetry(
					func(v fieldline.Vec3) fieldline.Vec3 {
						return fieldline.Vec3{v[0], v[1], -v[2]}
					}, zcolors[i],
				)
			}
			trajs = append(trajs, traj)
		}
	}

	fieldline.Run(opts, trajs)
}
