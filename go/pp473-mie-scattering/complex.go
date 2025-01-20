package main

import "math/big"

type complex struct {
	re *big.Float
	im *big.Float
}

// i^n.
func iPow(n int) *complex {
	switch n % 4 {
	case 0:
		return newComplex(fromInt(1), blankFloat())
	case 1:
		return newComplex(blankFloat(), fromInt(1))
	case 2:
		return newComplex(fromInt(-1), blankFloat())
	case 3:
		return newComplex(blankFloat(), fromInt(-1))
	}
	return nil
}

func blankComplex() *complex {
	return newComplex(blankFloat(), blankFloat())
}

func newComplex(re, im *big.Float) *complex {
	return &complex{re: re, im: im}
}

func complexFromReal(r *big.Float) *complex {
	return &complex{re: r, im: blankFloat()}
}

func complexFromFloat64(r float64) *complex {
	return newComplex(fromFloat64(r), blankFloat())
}

func complexFromInt(n int) *complex {
	return newComplex(fromInt(n), blankFloat())
}

func (c *complex) add(d *complex) *complex {
	v := blankComplex()
	v.re.Add(c.re, d.re)
	v.im.Add(c.im, d.im)
	return v
}

func (c *complex) sub(d *complex) *complex {
	v := blankComplex()
	v.re.Sub(c.re, d.re)
	v.im.Sub(c.im, d.im)
	return v
}

func (c *complex) mul(d *complex) *complex {
	v := blankComplex()
	v.re.Mul(c.re, d.re)
	tmp1 := blankFloat().Mul(c.im, d.im)
	v.re.Sub(v.re, tmp1)
	v.im.Mul(c.re, d.im)
	tmp2 := blankFloat().Mul(c.im, d.re)
	v.im.Add(v.im, tmp2)
	return v
}

func (c *complex) quo(d *complex) *complex {
	v := blankComplex()
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

func (c *complex) mod2() *big.Float {
	tmp1 := blankFloat().Mul(c.re, c.re)
	tmp2 := blankFloat().Mul(c.im, c.im)
	return tmp1.Add(tmp1, tmp2)
}
