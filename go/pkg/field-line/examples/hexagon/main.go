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
	// Uncomment to use alternate sign of charges.
	//posAngles := []float64{0, 120, 240}
	//negAngles := []float64{60, 180, 300}
	posAngles := []float64{0, 60, 120, 180, 240, 300}
	negAngles := []float64{}
	var positives []fieldline.Vec3
	for _, ang := range posAngles {
		arad := math.Pi * ang / 180.0
		positives = append(positives, fieldline.Vec3{
			a * math.Cos(arad), a * math.Sin(arad), 0.0,
		})
	}
	var negatives []fieldline.Vec3
	for _, ang := range negAngles {
		arad := math.Pi * ang / 180.0
		negatives = append(negatives, fieldline.Vec3{
			a * math.Cos(arad), a * math.Sin(arad), 0.0,
		})
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

	var trajs []fieldline.Trajectory
	samples := 30
	for _, charge := range positives {
		for i := 0; i < samples; i++ {
			theta := math.Pi * 2.0 * float64(i) / float64(samples)
			trajs = append(trajs, fieldline.Trajectory{
				Start: charge.Add(fieldline.Vec3{sr * math.Cos(theta), sr * math.Sin(theta), 0.0}),
				AtEnd: atEnd, Color: fieldline.RandColor(),
			})
		}
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.5,
		FadingGamma: .5,
	}

	fieldline.Run(opts, trajs)
}
