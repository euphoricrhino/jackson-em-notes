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

	// We end the tracing when the field line curved back to the z=0 plane.
	atEnd := func(p, v fieldline.Vec3) bool {
		return v[1] < 0 && p[1] < 1e-3 || v.Norm() < 1e-3
	}

	thetas := []float64{0, 15, 30, 45, 60, 75}
	phis := []float64{45, 90, 135, 180, 225, 270, 315}
	var trajs []fieldline.Trajectory
	for i, thetaDeg := range thetas {
		theta := thetaDeg * math.Pi / 180.0
		start := fieldline.Vec3{a * math.Sin(theta), 0, 0.0}
		color := fieldline.RandColor()
		traj := fieldline.Trajectory{
			Start: start,
			AtEnd: atEnd,
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
		if i != 0 {
			for _, phiDeg := range phis {
				phi := phiDeg * math.Pi / 180.0
				traj.AddSymmetry(
					func(v fieldline.Vec3) fieldline.Vec3 {
						rho := v[0]
						return fieldline.Vec3{
							rho * math.Cos(phi),
							v[1],
							rho * math.Sin(phi),
						}
					},
					color,
				)
				traj.AddSymmetry(
					func(v fieldline.Vec3) fieldline.Vec3 {
						rho := v[0]
						return fieldline.Vec3{
							rho * math.Cos(phi),
							-v[1],
							rho * math.Sin(phi),
						}
					},
					color,
				)
			}
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
		FadingGamma: .5,
		CameraOrbit: fieldline.NewCameraOrbit(30, 180),
	}

	fieldline.Run(opts, trajs)
}
