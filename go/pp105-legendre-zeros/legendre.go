package legendrezeros

import (
	"math/big"
	"sync"
)

var (
	// Precision of all big.Float computation.
	floatPrec uint = 2000
	once      sync.Once
)

// Will be set only once after flag is parsed.
func SetPrecOnce(prec uint) {
	once.Do(func() {
		floatPrec = prec
	})
}

func NewFloat(v float64) *big.Float {
	return big.NewFloat(v).SetPrec(floatPrec)
}

func BlankFloat() *big.Float { return NewFloat(0.0) }

func NewFloatFromInt(n int) *big.Float { return NewFloat(float64(n)) }

// Evaluates the Legendre function of first kind of order 𝜈, at x, with error bound eb.
func Legendre(x, nu, eb *big.Float) *big.Float {
	xi := NewFloat(1.0)
	xi.Sub(xi, x)
	xi.Mul(xi, NewFloat(0.5))
	negNu := BlankFloat().Neg(nu)
	term := NewFloat(1.0)
	sum := NewFloat(1.0)
	for l := 1; true; l++ {
		v := BlankFloat().Add(negNu, NewFloatFromInt(l-1))
		term.Mul(term, v)
		v.Add(nu, NewFloatFromInt(l))
		term.Mul(term, v)
		term.Mul(term, xi)
		term.Quo(term, NewFloatFromInt(l*l))
		abst := BlankFloat().Abs(term)
		sum.Add(sum, term)
		if abst.Cmp(eb) < 0 {
			break
		}
	}
	return sum
}
