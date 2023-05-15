package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"

	legendrezeros "github.com/euphoricrhino/jackson-em-notes/go/pp105-legendre-zeros"
)

// Binary search for zero of a Legendre function of order 𝜈 (not necessarily integer).
// One of left or right is assumed to be a zero of a Legendre function of a different order.
// We start with knowing that the zero of 𝜈=1 is x=0, and then we take two branches:
// 1. 𝜈 is increased gradually to search for zero of higher order, in this branch, we set left to the previous iteration's zero, and set rightPivot=true.
// 2. 𝜈 is decreased gradually to search for zero of lower order, in this branch, we set right to the previous iteration's zero, and set leftPivot=true.
// The leftPivot and rightPivot save the calculation of Legendre function at x=-1 or x=1 (since when 𝜈<1, x=-1 is divergent).
func searchZero(left, right, nu, leb, reb *big.Float, leftPivot, rightPivot bool) *big.Float {
	var leftVal, rightVal *big.Float
	for i := 0; true; i++ {
		if !leftPivot && leftVal == nil {
			leftVal = legendrezeros.Legendre(left, nu, leb)
			if leftVal.Sign() > 0 {
				panic("invalid left value sign")
			}
		}
		if !rightPivot && rightVal == nil {
			if rightVal == nil {
				rightVal = legendrezeros.Legendre(right, nu, leb)
			}
			if rightVal.Sign() < 0 {
				panic("invalid right value sign")
			}
		}
		mid := legendrezeros.BlankFloat().Add(left, right)
		mid.Mul(mid, legendrezeros.NewFloat(0.5))
		r := legendrezeros.BlankFloat().Sub(right, left)
		// We terminate on one of two conditions: 1) the [left, right] range is small enough.
		if r.Abs(r).Cmp(reb) < 0 {
			fmt.Printf("𝜈=%.02f, %v iterations (converged by range)\n", nu, i+1)
			return mid
		}
		midVal := legendrezeros.Legendre(mid, nu, leb)
		abs := legendrezeros.BlankFloat().Abs(midVal)
		// Or 2), the mid point value is close enough to zero.
		if abs.Cmp(reb) < 0 {
			fmt.Printf("𝜈=%.02f, %v iterations (converged by value)\n", nu, i+1)
			return mid
		}
		if midVal.Sign() < 0 {
			left = mid
			leftPivot = false
			leftVal = midVal
		} else {
			right = mid
			rightPivot = false
			rightVal = midVal
		}
	}
	return nil
}

var (
	prec             = flag.Uint("prec", 1000, "precision")
	legendreErrBound = flag.Float64("legendre-err-bound", 1e-10, "error bound for computing legendre polynomials")
	rootErrBound     = flag.Float64("root-err-bound", 1e-6, "error bound for computing roots")
)

func main() {
	flag.Parse()
	legendrezeros.SetPrecOnce(*prec)
	filename := filepath.Join(os.TempDir(), "genleg.m")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	leb := legendrezeros.NewFloat(*legendreErrBound)
	reb := legendrezeros.NewFloat(*rootErrBound)
	var nuRight, betaRight []float64
	one := legendrezeros.NewFloat(1.0)
	left := legendrezeros.NewFloat(0.0)
	for i := 1; i <= 300; i++ {
		nu := 1.0 + float64(i)/100.0
		nuRight = append(nuRight, nu)
		root := searchZero(left, one, legendrezeros.NewFloat(nu), leb, reb, false, true)
		rootv, _ := root.Float64()
		betaRight = append(betaRight, math.Acos(rootv)*180.0/math.Pi)
		left = root
	}

	var nuLeft, betaLeft []float64
	negOne := legendrezeros.NewFloat(-1.0)
	right := legendrezeros.NewFloat(0.0)
	for i := 1; i <= 95; i++ {
		nu := 1.0 - float64(i)/100.0
		nuLeft = append([]float64{nu}, nuLeft...)
		root := searchZero(negOne, right, legendrezeros.NewFloat(nu), leb, reb, true, false)
		rootv, _ := root.Float64()
		betaLeft = append([]float64{math.Acos(rootv) * 180.0 / math.Pi}, betaLeft...)
		right = root
	}
	x := append(nuLeft, 1.0)
	x = append(x, nuRight...)
	y := append(betaLeft, 90.0)
	y = append(y, betaRight...)

	// Asymptotic forms.
	var y1, y2 []float64
	// 1. equation (3.48a), for large 𝜈.
	for _, nu := range x {
		beta := 2.405 / (nu + 0.5)
		beta *= 180.0 / math.Pi
		y1 = append(y1, beta)
	}
	y2End := 0
	// 2. equation (3.48b), for small 𝜈.
	for i := range x {
		nu := x[i]
		// y2's start corresponds to 𝜈=1.5 (to match Jackson Figure 3.6).
		if nu > 1.5 {
			y2End = i
			break
		}
		beta := math.Pi - 2.0/math.Exp(0.5/nu)
		beta *= 180.0 / math.Pi
		y2 = append(y2, beta)
	}
	fmt.Fprintf(f, "x=%v;\n", x)
	fmt.Fprintf(f, "y=%v;\n", y)
	fmt.Fprintf(f, "y1=%v;\n", y1)
	fmt.Fprintf(f, "x2=%v;\n", x[:y2End])
	fmt.Fprintf(f, "y2=%v;\n", y2)
	fmt.Fprintln(f, "plot(y,x,'LineWidth',2,y1,x,'--','LineWidth',2,y2,x2,':','LineWidth',2);")
	fmt.Fprintln(f, "xlabel('\\beta (deg)');")
	fmt.Fprintln(f, "ylabel('𝜈');")
	fmt.Fprintln(f, "xlim([0 200]);")
	fmt.Fprintln(f, "set(gca,'FontSize',18);")
	fmt.Fprintln(f, "set(gca,'XTick', 0:30:180);")

	fmt.Println(filename)
}
