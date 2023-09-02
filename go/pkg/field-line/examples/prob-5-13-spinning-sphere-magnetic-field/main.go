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
			return v[1] < 0 && p[1] < 1e-3 || v.Norm() < 1e-4
		}
	}

	thetas := []float64{0, 15, 30, 45, 60, 75}
	phis := []float64{45, 90, 135, 180, 225, 270, 315}
	colors := make([][3]float64, 1+(len(thetas)-1)*len(phis))
	for c := range colors {
		colors[c] = fieldline.RandColor()
	}
	rand.Seed(time.Now().UnixNano())
	c := 0
	makeTrajs := func() []fieldline.Trajectory {
		var trajs []fieldline.Trajectory
		for i, thetaDeg := range thetas {
			theta := thetaDeg * math.Pi / 180.0
			start := fieldline.Vec3{a * math.Sin(theta), 0, 0.0}
			traj := fieldline.Trajectory{
				Start: start,
				AtEnd: makeAtEnd(start),
				Color: colors[c],
			}
			// Lower half is reflected using symmetry.
			traj.AddSymmetry(
				func(v fieldline.Vec3) fieldline.Vec3 {
					v[1] = -v[1]
					return v
				},
				colors[c],
			)
			if i != 0 {
				for j, phiDeg := range phis {
					phi := phiDeg * math.Pi / 180.0
					color := colors[1+(i-1)*len(phis)+j]
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
		return trajs
	}

	camCircleAngle := math.Pi * .17
	camCircleZ := fieldline.Vec3{-math.Sin(camCircleAngle), math.Cos(camCircleAngle), 0}
	camCircleX := fieldline.Vec3{0, 0, 1}
	camCircleY := camCircleZ.Cross(camCircleX)

	frames := 180
	dtheta := math.Pi * 2.0 / float64(frames)
	for f := 0; f < frames; f++ {
		theta := dtheta * float64(f)
		camPos := camCircleX.Scale(math.Cos(theta))
		camPos = camPos.Add(camCircleY.Scale(math.Sin(theta)))
		camRight := camCircleX.Scale(-math.Sin(theta))
		camRight = camRight.Add(camCircleY.Scale(math.Cos(theta)))
		opts := fieldline.Options{
			OutputFile:  fmt.Sprintf("%v-%03d.png", *output, f),
			Width:       *width,
			Height:      *height,
			Step:        *step,
			TangentAt:   tangentAt,
			Camera:      fieldline.NewCamera(camPos, camRight),
			LineWidth:   1.0,
			FadingGamma: .5,
		}

		trajs := makeTrajs()
		fieldline.Run(opts, trajs)
		fmt.Println(opts.OutputFile)
	}
}
