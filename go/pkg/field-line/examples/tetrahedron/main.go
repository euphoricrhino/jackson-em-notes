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
	rand.Seed(time.Now().UnixNano())

	a := 0.3
	// Location of the positive and negative charges.
	positives := []fieldline.Vec3{
		{a, a, a},
		{-a, a, -a},
	}
	negatives := []fieldline.Vec3{
		{a, -a, -a},
		{-a, -a, a},
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		var v fieldline.Vec3
		accum := func(sgn float64, charge fieldline.Vec3) {
			d := p.Subtract(charge)
			d3 := math.Pow(d.Dot(d), 1.5)
			v = v.Add(d.Scale(sgn / d3))
		}
		for _, charge := range positives {
			accum(1.0, charge)
		}
		for _, charge := range negatives {
			accum(-1.0, charge)
		}
		return v.Scale(1.0 / 2000.0)
	}

	sr := 0.02
	atEnd := func(p, v fieldline.Vec3) bool {
		// Stop if the field is too weak.
		if v.Dot(v) < 1e-8 {
			return true
		}
		// Stop if we are close the negative charges.
		for _, charge := range negatives {
			d := p.Subtract(charge)
			if d.Dot(d) < sr*sr {
				return true
			}
		}
		return false
	}

	thetaDeg := []float64{30.0, 60.0, 90.0, 120.0, 150.0}
	phiDeg := []float64{0.0, 60.0, 120.0, 180.0, 240.0, 300.0}
	colors := make([][3]float64, 2+len(thetaDeg)*len(phiDeg))
	rand.Seed(time.Now().UnixNano())
	for i := range colors {
		colors[i] = fieldline.RandColor()
	}
	generateTraj := func(charge, localz, localx fieldline.Vec3) []fieldline.Trajectory {
		lz := localz.Normalize()
		xonz := lz.Scale(localx.Dot(lz))
		lx := localx.Subtract(xonz).Normalize()
		ly := lz.Cross(lx)

		lx = lx.Scale(sr)
		ly = ly.Scale(sr)
		lz = lz.Scale(sr)
		// North and south poles.
		ret := []fieldline.Trajectory{
			{Start: charge.Add(lz), AtEnd: atEnd, Color: colors[0]},
			{Start: charge.Subtract(lz), AtEnd: atEnd, Color: colors[1]},
		}
		for i, theta := range thetaDeg {
			thetaRad := theta * math.Pi / 180.0
			for j, phi := range phiDeg {
				phiRad := phi * math.Pi / 180.0
				disp := lx.Scale(math.Sin(thetaRad) * math.Cos(phiRad))
				disp = disp.Add(ly.Scale(math.Sin(thetaRad) * math.Sin(phiRad)))
				disp = disp.Add(lz.Scale(math.Cos(thetaRad)))
				ret = append(ret, fieldline.Trajectory{Start: charge.Add(disp), AtEnd: atEnd, Color: colors[2+i*len(phiDeg)+j]})
			}
		}
		return ret
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
			LineWidth:   1.5,
			FadingGamma: .25,
		}

		trajs := generateTraj(positives[0], fieldline.Vec3{1, 1, 1}, fieldline.Vec3{0, -1, -1})
		trajs = append(trajs, generateTraj(positives[1], fieldline.Vec3{-1, 1, -1}, fieldline.Vec3{1, -1, 0})...)
		fieldline.Run(opts, trajs)
		fmt.Println(opts.OutputFile)
	}
}
