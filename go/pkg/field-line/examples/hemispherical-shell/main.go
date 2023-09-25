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

// Integral of P_{2n+1}(x)dx over [0,1];
func legInt(n int) float64 {
	res := 1.0
	if n%2 == 1 {
		res = -1.0
	}
	for nn := 0; nn <= n; nn += 1 {
		num := 2*nn - 1
		if num < 1 {
			num = 1
		}
		den := 2 * (nn + 1)
		res *= float64(num) / float64(den)
	}
	return res
}

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
	const a = 0.4

	const maxn = 20
	const scale = 1 / 200.0

	tangentIn := func(p fieldline.Vec3) fieldline.Vec3 {
		r := p.Norm()
		ct := p[1] / r
		st := math.Sqrt(1 - ct*ct)
		val, der := legend(2*maxn+1, ct)
		f := r / a
		a2 := a * a
		er := 0.0
		et := 0.0
		for n := 0; n <= maxn; n++ {
			l := 2*n + 1
			ff := math.Pow(f, float64(l)-1.0) / a2 * legInt(n)
			er -= float64(l) * ff * val[l]
			et -= ff * der[l] * (-st)
		}
		return fieldline.Vec3{
			er*st + et*ct,
			er*ct - et*st,
			0.0,
		}.Scale(scale)
	}

	tangentOut := func(p fieldline.Vec3) fieldline.Vec3 {
		r := p.Norm()
		ct := p[1] / r
		st := math.Sqrt(1 - ct*ct)
		val, der := legend(2*maxn+1, ct)
		f := a / r
		r2 := r * r
		// Er has l=0 contribution
		er := 1.0 / r2
		et := 0.0
		for n := 0; n <= maxn; n++ {
			l := 2*n + 1
			ff := math.Pow(f, float64(l)) / r2 * legInt(n)
			er -= float64(-l-1) * ff * val[l]
			et -= ff * der[l] * (-st)
		}
		return fieldline.Vec3{
			er*st + et*ct,
			er*ct - et*st,
			0.0,
		}.Scale(scale)
	}

	tangentAt := func(p fieldline.Vec3) fieldline.Vec3 {
		if p.Norm() < a {
			return tangentIn(p)
		}
		return tangentOut(p)
	}

	atEnd := func(p, v fieldline.Vec3) bool {
		return false
	}

	opts := fieldline.Options{
		OutputFile:  *output,
		Width:       *width,
		Height:      *height,
		Step:        *step,
		TangentAt:   tangentAt,
		LineWidth:   1.0,
		FadingGamma: 1,
	}

	sr := 0.01
	trajs := []fieldline.Trajectory{
		{Start: fieldline.Vec3{0, a + sr, 0}, AtEnd: atEnd, Color: fieldline.RandColor()},
		{Start: fieldline.Vec3{0, a - sr, 0}, AtEnd: atEnd, Color: fieldline.RandColor()},
	}
	dtheta := math.Pi * 5 / 180
	for theta := math.Pi / 2.0; theta > 0.0; theta -= dtheta {
		color := fieldline.RandColor()
		for _, ar := range []float64{a - sr, a + sr} {
			traj := fieldline.Trajectory{
				Start: fieldline.Vec3{ar * math.Sin(theta), ar * math.Cos(theta), 0},
				AtEnd: atEnd,
				Color: color,
			}
			traj.AddSymmetry(
				func(p fieldline.Vec3) fieldline.Vec3 {
					return fieldline.Vec3{-p[0], p[1], 0.0}
				},
				color,
			)
			trajs = append(trajs, traj)
		}
	}

	fieldline.Run(opts, trajs)
}
