package mp

import (
	"math/big"
	"sync"
)

var (
	// Precision of all big.Float computation.
	floatPrec uint = 2000
	once      sync.Once
)

// SetPrecOnce will be set only once after config is loaded.
func SetPrecOnce(prec uint) {
	once.Do(func() {
		floatPrec = prec
	})
}

// BlankFloat creates a zero-value float at the global precision.
func BlankFloat() *big.Float { return big.NewFloat(0).SetPrec(floatPrec) }

// NewFromFloat64 creates a float with the given initial value at the global precision.
func NewFromFloat64(val float64) *big.Float {
	return big.NewFloat(val).SetPrec(floatPrec)
}

// NewFromInt creates a float with the given initial value at the global precision.
func NewFromInt(val int) *big.Float {
	return BlankFloat().SetInt64(int64(val)).SetPrec(floatPrec)
}

// NewFromRat creates a float with the given initial value at the global precision.
func NewFromRat(n, d int) *big.Float {
	return BlankFloat().SetRat(big.NewRat(int64(n), int64(d))).SetPrec(floatPrec)
}
