package main

import (
	"math"
	"math/big"
)

const (
	prec = 100
)

func blankFloat() *big.Float {
	return big.NewFloat(0.0).SetPrec(prec)
}

func fromFloat64(f float64) *big.Float {
	return big.NewFloat(f).SetPrec(prec)
}

func fromInt(n int) *big.Float {
	return blankFloat().SetInt64(int64(n)).SetPrec(prec)
}

// j_l(x), j_l'(x) for l=0..maxL
func sphericalBessel1(maxL int, x float64) ([]*big.Float, []*big.Float) {
	vals := make([]*big.Float, maxL+1)
	derivs := make([]*big.Float, maxL+1)
	// Check small-argument cutoff.
	leadingTerm := fromFloat64(1)
	cutoff := fromFloat64(1e-5)
	lCutoff := 0
	for leadingTerm.Cmp(cutoff) > 0 && lCutoff <= maxL {
		lCutoff++
		leadingTerm.Mul(leadingTerm, fromFloat64(x))
		leadingTerm.Quo(leadingTerm, fromInt(2*lCutoff+1))
	}

	// Use recursion to calculate j_l(x) for l=0..lCutoff-1
	if lCutoff > 0 {
		vals[0] = fromFloat64(math.Sin(x))
		vals[0].Quo(vals[0], fromFloat64(x))
	}
	if lCutoff > 1 {
		vals[1] = fromFloat64(math.Sin(x))
		vals[1].Quo(vals[1], fromFloat64(x*x))
		tmp := fromFloat64(math.Cos(x))
		tmp.Quo(tmp, fromFloat64(x))
		vals[1].Sub(vals[1], tmp)
		derivs[1] = fromFloat64(-2.0 / x)
		derivs[1].Mul(derivs[1], vals[1])
		derivs[1].Add(derivs[1], vals[0])
	}
	sphericalBesselRecurse(x, vals, derivs, lCutoff)

	sphericalBessel1Small(x, lCutoff, vals, derivs)
	return vals, derivs
}

// y_l(x), y_l'(x) for l=0..maxL.
func sphericalBessel2(maxL int, x float64) ([]*big.Float, []*big.Float) {
	vals := make([]*big.Float, maxL+1)
	derivs := make([]*big.Float, maxL+1)

	vals[0] = fromFloat64(-math.Cos(x))
	vals[0].Quo(vals[0], fromFloat64(x))
	if maxL > 0 {
		vals[1] = fromFloat64(-math.Cos(x))
		vals[1].Quo(vals[1], fromFloat64(x*x))
		tmp := fromFloat64(math.Sin(x))
		tmp.Quo(tmp, fromFloat64(x))
		vals[1].Sub(vals[1], tmp)
		derivs[1] = fromFloat64(-2.0 / x)
		derivs[1].Mul(derivs[1], vals[1])
		derivs[1].Add(derivs[1], vals[0])
	}
	sphericalBesselRecurse(x, vals, derivs, maxL+1)

	return vals, derivs
}

func sphericalBessel1Small(x float64, lCutoff int, vals, derivs []*big.Float) {
	if lCutoff == 0 {
		// We don't care about l=0 in Mie scattering.
		lCutoff = 1
	}
	c1 := fromFloat64(1)
	c2 := fromFloat64(2)
	b := fromFloat64(1) // x^{l}
	for l := 1; l <= lCutoff; l++ {
		c1.Mul(c1, fromInt(2*l+1))
		c2.Mul(c2, fromInt(2*l+1))
		b.Mul(b, fromFloat64(x))
	}
	c2.Mul(c2, fromInt(2*lCutoff+3))
	a := blankFloat().Quo(b, fromFloat64(x)) // x^{l-1}
	c := blankFloat().Mul(b, fromFloat64(x)) // x^{l+1}
	d := blankFloat().Mul(c, fromFloat64(x)) // x^{l+2}

	bigx := fromFloat64(x)
	for l := lCutoff; l < len(vals); l++ {
		vals[l] = blankFloat().Quo(b, c1)
		tmp := blankFloat().Quo(d, c2)
		vals[l].Sub(vals[l], tmp)
		derivs[l] = blankFloat().Mul(fromInt(l), a)
		derivs[l].Quo(derivs[l], c1)
		tmp = blankFloat().Mul(fromInt(l+2), c)
		tmp.Quo(tmp, c2)
		derivs[l].Sub(derivs[l], tmp)
		a.Mul(a, bigx)
		b.Mul(b, bigx)
		c.Mul(c, bigx)
		d.Mul(d, bigx)
		c1.Mul(c1, fromInt(2*l+3))
		c2.Mul(c2, fromInt(2*l+5))
	}
}

func sphericalBesselRecurse(x float64, vals, derivs []*big.Float, lCutoff int) {
	for l := 2; l < lCutoff; l++ {
		vals[l] = fromFloat64((2*float64(l) - 1) / x)
		vals[l].Mul(vals[l], vals[l-1])
		vals[l].Sub(vals[l], vals[l-2])

		derivs[l] = fromFloat64(-float64(l+1) / x)
		derivs[l].Mul(derivs[l], vals[l])
		derivs[l].Add(derivs[l], vals[l-1])
	}
}

// P_l(x), P_l'(x) for l=0..maxL
func legendre(maxL int, x float64) ([]*complex, []*complex) {
	vals := make([]*big.Float, maxL+1)
	derivs := make([]*big.Float, maxL+1)
	vals[0] = fromFloat64(1)
	derivs[0] = blankFloat()
	if maxL > 0 {
		vals[1] = fromFloat64(x)
		derivs[1] = fromFloat64(1)
	}
	for l := 2; l <= maxL; l++ {
		tmp1 := fromFloat64(float64(2*l-1) / float64(l) * x)
		tmp1.Mul(tmp1, vals[l-1])
		tmp2 := fromFloat64(float64(l-1) / float64(l))
		tmp2.Mul(tmp2, vals[l-2])
		vals[l] = tmp1.Sub(tmp1, tmp2)

		tmp1 = fromInt(l)
		tmp1.Mul(tmp1, vals[l-1])
		tmp2 = fromFloat64(x)
		tmp2.Mul(tmp2, derivs[l-1])
		derivs[l] = tmp1.Add(tmp1, tmp2)
	}
	cvals := make([]*complex, maxL+1)
	cderivs := make([]*complex, maxL+1)
	for i := 0; i <= maxL; i++ {
		cvals[i] = complexFromReal(vals[i])
		cderivs[i] = complexFromReal(derivs[i])
	}
	return cvals, cderivs
}
