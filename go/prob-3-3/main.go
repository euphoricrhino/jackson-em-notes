package main

import (
	"flag"
	"math"
	"math/big"

	fieldrenderer "github.com/euphoricrhino/jackson-em-notes/go/pkg/field-renderer"
	"github.com/euphoricrhino/jackson-em-notes/go/pkg/mp"
)

var (
	heatmap = flag.String("heatmap", "", "heatmap file")
	output  = flag.String("output", "", "output file")
	gamma   = flag.Float64("gamma", 1.0, "gamma correction")
	width   = flag.Int("width", 640, "output width")
	height  = flag.Int("height", 640, "output height")
	terms   = flag.Int("terms", 10, "number of terms to keep in the series sum")
	prec    = flag.Uint("prec", 100, "floating point precision")
)

func main() {
	flag.Parse()
	mp.SetPrecOnce(*prec)

	// Construct legendre polynomials up to the 2*terms+1 order.
	legs := constructLegendre(2*(*terms) + 1)

	field := func(x, y int) float64 {
		rad := float64(*width) / 8
		fx := float64(x) - float64(*width-1)/2
		fy := float64(*height-1-y) - float64(*height-1)/2

		r := math.Sqrt(fx*fx + fy*fy)

		// cosθ.
		ct := fy / r
		fct := mp.NewFromFloat64(ct)

		maxRPower := 2*(*terms) + 1
		if r >= rad {
			// Use the even formula.
			rr := mp.BlankFloat().Quo(mp.NewFromFloat64(rad), mp.NewFromFloat64(r))
			// Power evaluator for radial polynomial.
			rpe := mp.NewPowerEvaluator(rr, maxRPower)

			maxPPower := 2 * (*terms)
			// Power evaluator for angle legendre polynomials.
			ppe := mp.NewPowerEvaluator(fct, maxPPower)
			res := mp.BlankFloat()
			sgn := 1.0
			for l := 0; l <= *terms; l++ {
				v := mp.NewFromFloat64(sgn / float64(2*l+1))
				v.Mul(v, rpe.Pow(2*l+1))
				v.Mul(v, legs[2*l].eval(ppe))
				res.Add(res, v)
				sgn *= -1.0
			}
			fv, _ := res.Float64()
			return fv * 2.0 / math.Pi
		}

		sgn := 1.0
		if ct < 0 {
			sgn = -1.0
		}

		rr := mp.BlankFloat().Quo(mp.NewFromFloat64(r), mp.NewFromFloat64(rad))
		rpe := mp.NewPowerEvaluator(rr, maxRPower)
		maxPPower := 2*(*terms) + 1
		ppe := mp.NewPowerEvaluator(fct, maxPPower)
		res := mp.BlankFloat()
		// Flip the sign one more time to account for the (-1)^{k+1}.
		sgn *= -1.0
		for k := 0; k <= *terms; k++ {
			v := mp.NewFromFloat64(sgn / float64(2*k+1))
			v.Mul(v, rpe.Pow(2*k+1))
			v.Mul(v, legs[2*k+1].eval(ppe))
			res.Add(res, v)
			sgn *= -1.0
		}
		fv, _ := res.Float64()
		return fv*2.0/math.Pi + 1.0
	}

	if err := fieldrenderer.Run(fieldrenderer.Options{
		HeatMapFile: *heatmap,
		OutputFile:  *output,
		Gamma:       *gamma,
		Width:       *width,
		Height:      *height,
		Field:       field,
	}); err != nil {
		panic(err)
	}
}

// Represents a polynomial.
type polynomial struct {
	coeff []*big.Float
}

// Constructs Legendre polynomials up to degree maxL.
func constructLegendre(maxL int) []*polynomial {
	ret := make([]*polynomial, maxL+1)
	ret[0] = legendre(0, nil, nil)
	ret[1] = legendre(1, nil, nil)
	for l := 2; l <= maxL; l++ {
		ret[l] = legendre(l, ret[l-2], ret[l-1])
	}
	return ret
}

// Computes order-l legendre polynomial from previous two orders.
func legendre(l int, pprev, prev *polynomial) *polynomial {
	if l == 0 {
		return &polynomial{
			coeff: []*big.Float{mp.NewFromFloat64(1.0)},
		}
	}
	if l == 1 {
		return &polynomial{
			coeff: []*big.Float{nil, mp.NewFromFloat64(1.0)},
		}
	}

	ret := &polynomial{
		coeff: make([]*big.Float, l+1),
	}

	factor1 := mp.NewFromRat(2*l-1, l)
	factor2 := mp.NewFromRat(l-1, l)
	// P_l only has powers with the same parity as l.
	for k := l; k >= 0; k -= 2 {
		ret.coeff[k] = mp.BlankFloat()
		if k > 0 {
			ret.coeff[k].Mul(prev.coeff[k-1], factor1)
		}
		if k <= l-2 {
			ret.coeff[k].Sub(ret.coeff[k], mp.BlankFloat().Mul(pprev.coeff[k], factor2))
		}
	}
	return ret
}

// Evaluates the legendre polynomial at the value 'x' which was used to construct the given PowerEvaluator.
func (p *polynomial) eval(pe *mp.PowerEvaluator) *big.Float {
	res := mp.BlankFloat()
	for i, c := range p.coeff {
		if c == nil {
			continue
		}
		res.Add(res, mp.BlankFloat().Mul(c, pe.Pow(i)))
	}
	return res
}
