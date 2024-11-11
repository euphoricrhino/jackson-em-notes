package main

import (
	"math"
	"math/big"

	"github.com/euphoricrhino/go-common/graphix"
)

var floatPrec uint = 2000

// Represents a polynomial.
type polynomial struct {
	coeff []*big.Float
}

// Evaluates the polynomial at x, all powers are computed in logarithmic time.
func (poly *polynomial) eval(x *big.Float) *big.Float {
	pEval := newPowerEvaluator(x, len(poly.coeff)-1)
	ans := blankFloat()
	for k := 0; k < len(poly.coeff); k++ {
		if poly.coeff[k] != nil {
			ans.Add(ans, blankFloat().Mul(poly.coeff[k], pEval.pow(k)))
		}
	}
	return ans
}

// Evaluator for x raised to some power, after construction, all power evaluations can be run in logarithmic time.
type powerEvaluator struct {
	x *big.Float
	// Precomputed all powers of form x^(2^n).
	powers []*big.Float
}

func angular(l, m int) (*polynomial, *polynomial) {
	pl := legendre(l)
	poly := &polynomial{
		coeff: make([]*big.Float, len(pl.coeff)-m),
	}
	for k := len(pl.coeff) - 1; k >= m; k -= 2 {
		poly.coeff[k-m] = blankFloat().Set(pl.coeff[k])
		for j := 0; j < m; j++ {
			poly.coeff[k-m].Mul(poly.coeff[k-m], newFromInt(k-j))
		}
	}

	// Derivative of poly.
	der := &polynomial{
		coeff: make([]*big.Float, len(pl.coeff)-1),
	}
	for k := 1; k < len(poly.coeff); k++ {
		if poly.coeff[k] != nil {
			der.coeff[k-1] = blankFloat().Mul(poly.coeff[k], newFromInt(k))
		}
	}
	return poly, der
}

// Constructs Legendre polynomial P_l recursively.
func legendre(l int) *polynomial {
	if l == 0 {
		return &polynomial{
			coeff: []*big.Float{newFromFloat64(1.0)},
		}
	}
	if l == 1 {
		return &polynomial{
			coeff: []*big.Float{nil, newFromFloat64(1.0)},
		}
	}

	pprev := legendre(0)
	prev := legendre(1)
	for ll := 2; ll <= l; ll++ {
		cur := &polynomial{
			coeff: make([]*big.Float, ll+1),
		}

		factor1 := newFromRat(2*ll-1, ll)
		factor2 := newFromRat(ll-1, ll)
		// P_l only has powers with the same parity as l.
		for k := ll; k >= 0; k -= 2 {
			cur.coeff[k] = blankFloat()
			if k > 0 {
				cur.coeff[k].Mul(prev.coeff[k-1], factor1)
			}
			if k <= ll-2 {
				cur.coeff[k].Sub(cur.coeff[k], blankFloat().Mul(pprev.coeff[k], factor2))
			}
		}
		pprev = prev
		prev = cur
	}
	return prev
}

// Constructs power evaluator given the maximum possible power to be computed subsequently.
func newPowerEvaluator(x *big.Float, maxPower int) *powerEvaluator {
	bits := 0
	n := maxPower
	for n != 0 {
		n = n >> 1
		bits++
	}
	pEval := &powerEvaluator{
		x:      x,
		powers: make([]*big.Float, bits+1),
	}
	// powers[0] was never explicitly used by pow() below, so let it remain nil.
	if bits == 0 {
		return pEval
	}
	pEval.powers[1] = blankFloat().Set(x)
	for k := 2; k <= bits; k++ {
		pEval.powers[k] = blankFloat().Mul(pEval.powers[k-1], pEval.powers[k-1])
	}

	return pEval
}

func (pEval *powerEvaluator) pow(n int) *big.Float {
	ans := newFromFloat64(1.0)
	shift := 1
	for n != 0 {
		if n&1 == 1 {
			ans.Mul(ans, pEval.powers[shift])
		}
		shift++
		n = n >> 1
	}
	return ans
}

func blankFloat() *big.Float { return big.NewFloat(0).SetPrec(floatPrec) }

func newFromFloat64(val float64) *big.Float {
	return big.NewFloat(val).SetPrec(floatPrec)
}

func newFromInt(val int) *big.Float {
	return blankFloat().SetInt64(int64(val)).SetPrec(floatPrec)
}

func newFromRat(n, d int) *big.Float {
	return blankFloat().SetRat(big.NewRat(int64(n), int64(d))).SetPrec(floatPrec)
}

type sph struct {
	l, m    int
	c       *big.Float
	poly    *polynomial
	polyDer *polynomial
}

func newsph(l, m int) *sph {
	c := newFromInt(2*l + 1)
	c.Quo(c, newFromFloat64(4*math.Pi))
	for k := m; k >= -m+1; k-- {
		c.Quo(c, newFromInt(l+k))
	}
	c.Sqrt(c)

	poly, der := angular(l, m)
	return &sph{l: l, m: m, c: c, poly: poly, polyDer: der}
}

func (s *sph) evalY(theta float64, phi []float64) ([]*graphix.Vec3, []*graphix.Vec3) {
	x := math.Cos(theta)
	sx := math.Sqrt(1 - x*x)
	sxm := newPowerEvaluator(newFromFloat64(sx), s.m).pow(s.m)

	v := s.poly.eval(newFromFloat64(x))
	v.Mul(v, sxm)
	v.Mul(v, s.c)
	vf, _ := v.Float64()
	re := make([]*graphix.Vec3, len(phi))
	im := make([]*graphix.Vec3, len(phi))
	for i, ph := range phi {
		cp, sp := math.Cos(float64(s.m)*ph), math.Sin(float64(s.m)*ph)
		re[i] = graphix.NewVec3(vf*cp, 0, 0)
		im[i] = graphix.NewVec3(vf*sp, 0, 0)
	}
	return re, im
}

func (s *sph) evalPsi(theta float64, phi []float64) ([]*graphix.Vec3, []*graphix.Vec3) {
	x := math.Cos(theta)
	sx := math.Sqrt(1 - x*x)

	bigsx := newFromFloat64(sx)
	sxEval := newPowerEvaluator(bigsx, s.m+1)
	tmp1 := sxEval.pow(s.m + 1)
	tmp1.Neg(tmp1)
	bigx := newFromFloat64(x)
	bigm := newFromInt(s.m)
	vt := blankFloat().Mul(tmp1, s.polyDer.eval(bigx))
	if s.m != 0 {
		tmp1.Copy(bigm)
		tmp1.Mul(tmp1, bigx)
		tmp1.Mul(tmp1, sxEval.pow(s.m-1))
		tmp1.Mul(tmp1, s.poly.eval(bigx))
		vt.Add(vt, tmp1)
	}
	vt.Mul(vt, s.c)
	vp := blankFloat()
	if sx > 1e-5 {
		tmp1.Copy(sxEval.pow(s.m))
		tmp1.Quo(tmp1, bigsx)
		tmp1.Mul(tmp1, s.poly.eval(bigx))
		tmp1.Mul(tmp1, bigm)
		vp.Mul(tmp1, s.c)
	}

	re := make([]*graphix.Vec3, len(phi))
	im := make([]*graphix.Vec3, len(phi))
	vtf, _ := vt.Float64()
	vpf, _ := vp.Float64()
	for i, ph := range phi {
		cp, sp := math.Cos(float64(s.m)*ph), math.Sin(float64(s.m)*ph)
		re[i] = graphix.NewVec3(0, vtf*cp, -vpf*sp)
		im[i] = graphix.NewVec3(0, vtf*sp, vpf*cp)
	}
	return re, im
}

func (s *sph) evalPhi(theta float64, phi []float64) ([]*graphix.Vec3, []*graphix.Vec3) {
	re, im := s.evalPsi(theta, phi)
	for i := range phi {
		re[i][1], re[i][2] = re[i][2], -re[i][1]
		im[i][1], im[i][2] = im[i][2], -im[i][1]
	}
	return re, im
}

// Calculate spherical bessel function jl(x) and derivative d(xjl(x))/dx.
func sphericalBessel(l int, x float64) (float64, float64) {
	// k-1=0
	preva, prevb := newFromFloat64(1.0), newFromFloat64(0.0)
	// k=1
	bigx := newFromFloat64(x)
	cura := newFromFloat64(1)
	cura.Quo(cura, bigx)
	curb := newFromFloat64(-1)
	curb.Quo(curb, bigx)

	s, c := newFromFloat64(math.Sin(x)/x), newFromFloat64(math.Cos(x))

	getVal := func(a, b *big.Float) *big.Float {
		v1 := blankFloat().Mul(a, s)
		v2 := blankFloat().Mul(b, c)
		return v1.Add(v1, v2)
	}

	// j_{k-1}
	prev := getVal(preva, prevb)
	// j_k
	cur := getVal(cura, curb)

	getParam := func(f, cur, prev *big.Float) *big.Float {
		v1 := blankFloat().Mul(f, cur)
		return v1.Sub(v1, prev)
	}

	for k := 1; k < l; k++ {
		f := newFromFloat64(float64(2*k+1) / x)
		// j_{k+1}
		nexta := getParam(f, cura, preva)
		nextb := getParam(f, curb, prevb)
		next := getVal(nexta, nextb)
		preva, prevb, prev = cura, curb, cur
		cura, curb, cur = nexta, nextb, next
	}
	ret1, _ := cur.Float64()
	v1 := blankFloat().Mul(bigx, prev)
	v2 := blankFloat().Mul(newFromInt(l), cur)
	ret2, _ := v1.Sub(v1, v2).Float64()
	return ret1, ret2
}
