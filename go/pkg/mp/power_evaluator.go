package mp

import "math/big"

// PowerEvaluator evaluates for x raised to some power, after construction, all power evaluations can be run in logarithmic time.
type PowerEvaluator struct {
	x *big.Float
	// Precomputed all powers of form x^(2^n).
	powers []*big.Float
}

// NewPowerEvaluator constructs power evaluator given the maximum possible power to be computed subsequently.
func NewPowerEvaluator(x *big.Float, maxPower int) *PowerEvaluator {
	bits := 0
	n := maxPower
	for n != 0 {
		n = n >> 1
		bits++
	}
	pEval := &PowerEvaluator{
		x:      x,
		powers: make([]*big.Float, bits+1),
	}
	// powers[0] was never explicitly used by pow() below, so let it remain nil.
	if bits == 0 {
		return pEval
	}
	pEval.powers[1] = BlankFloat().Set(x)
	for k := 2; k <= bits; k++ {
		pEval.powers[k] = BlankFloat().Mul(pEval.powers[k-1], pEval.powers[k-1])
	}

	return pEval
}

// Pow computes the n-th power of 'x' used to construct this PowerEvaluator. The calculation is logarithmic time.
func (pEval *PowerEvaluator) Pow(n int) *big.Float {
	ans := NewFromFloat64(1.0)
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
