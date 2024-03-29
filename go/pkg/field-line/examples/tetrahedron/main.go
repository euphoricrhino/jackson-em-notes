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
			{Start: charge.Add(lz), AtEnd: atEnd, Color: fieldline.RandColor()},
			{Start: charge.Subtract(lz), AtEnd: atEnd, Color: fieldline.RandColor()},
		}
		for _, theta := range thetaDeg {
			thetaRad := theta * math.Pi / 180.0
			for _, phi := range phiDeg {
				phiRad := phi * math.Pi / 180.0
				disp := lx.Scale(math.Sin(thetaRad) * math.Cos(phiRad))
				disp = disp.Add(ly.Scale(math.Sin(thetaRad) * math.Sin(phiRad)))
				disp = disp.Add(lz.Scale(math.Cos(thetaRad)))
				ret = append(ret, fieldline.Trajectory{Start: charge.Add(disp), AtEnd: atEnd, Color: fieldline.RandColor()})
			}
		}
		return ret
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.5,
		FadingGamma: .25,
		CameraOrbit: fieldline.NewCameraOrbit(30, 180),
	}

	trajs := generateTraj(positives[0], fieldline.Vec3{1, 1, 1}, fieldline.Vec3{0, -1, -1})
	trajs = append(trajs, generateTraj(positives[1], fieldline.Vec3{-1, 1, -1}, fieldline.Vec3{1, -1, 0})...)
	fieldline.Run(opts, trajs)
}
