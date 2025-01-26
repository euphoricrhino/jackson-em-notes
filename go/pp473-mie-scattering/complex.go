package main

import "math/big"

type bigComplex struct {
	re *big.Float
	im *big.Float
}

// i^n.
func iPow(n int) *bigComplex {
	switch n % 4 {
	case 0:
		return newBigComplex(fromInt(1), blankFloat())
	case 1:
		return newBigComplex(blankFloat(), fromInt(1))
	case 2:
		return newBigComplex(fromInt(-1), blankFloat())
	case 3:
		return newBigComplex(blankFloat(), fromInt(-1))
	}
	return nil
}

func blankBigComplex() *bigComplex {
	return newBigComplex(blankFloat(), blankFloat())
}

func newBigComplex(re, im *big.Float) *bigComplex {
	return &bigComplex{re: re, im: im}
}

func bigComplexFromBigFloat(r *big.Float) *bigComplex {
	return &bigComplex{re: r, im: blankFloat()}
}

func bigComplexFromFloat64(r float64) *bigComplex {
	return newBigComplex(fromFloat64(r), blankFloat())
}

func bigComplexFromComplex128(r complex128) *bigComplex {
	return newBigComplex(fromFloat64(real(r)), fromFloat64(imag(r)))
}

func bigComplexFromInt(n int) *bigComplex {
	return newBigComplex(fromInt(n), blankFloat())
}

func (c *bigComplex) add(d *bigComplex) *bigComplex {
	v := blankBigComplex()
	v.re.Add(c.re, d.re)
	v.im.Add(c.im, d.im)
	return v
}

func (c *bigComplex) sub(d *bigComplex) *bigComplex {
	v := blankBigComplex()
	v.re.Sub(c.re, d.re)
	v.im.Sub(c.im, d.im)
	return v
}

func (c *bigComplex) mul(d *bigComplex) *bigComplex {
	v := blankBigComplex()
	v.re.Mul(c.re, d.re)
	tmp1 := blankFloat().Mul(c.im, d.im)
	v.re.Sub(v.re, tmp1)
	v.im.Mul(c.re, d.im)
	tmp2 := blankFloat().Mul(c.im, d.re)
	v.im.Add(v.im, tmp2)
	return v
}

func (c *bigComplex) quo(d *bigComplex) *bigComplex {
	v := blankBigComplex()
	v.re.Mul(c.re, d.re)
	tmp1 := blankFloat().Mul(c.im, d.im)
	v.re.Add(v.re, tmp1)
	v.im.Mul(c.im, d.re)
	tmp2 := blankFloat().Mul(c.re, d.im)
	v.im.Sub(v.im, tmp2)
	denom := blankFloat().Mul(d.re, d.re)
	tmp3 := blankFloat().Mul(d.im, d.im)
	denom.Add(denom, tmp3)
	v.re.Quo(v.re, denom)
	v.im.Quo(v.im, denom)
	return v
}

func (c *bigComplex) mod2() *big.Float {
	tmp1 := blankFloat().Mul(c.re, c.re)
	tmp2 := blankFloat().Mul(c.im, c.im)
	return tmp1.Add(tmp1, tmp2)
}
