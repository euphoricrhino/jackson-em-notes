package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	width  = flag.Int("width", 640, "Width of the image")
	height = flag.Int("height", 640, "Height of the image")
	maxL   = flag.Int("max-l", 20, "max l value")
	outDir = flag.String("out-dir", "./", "output directory")

	n    = flag.Float64("n", 1.2, "Refractive index")
	nImg = flag.Float64("n-img", 0.0, "Imaginary part of refractive index")
	mu   = flag.Float64("mu", 1.0, "Permeability")
)

const (
	lambda = 0.2
	k      = 2.0 * math.Pi / lambda
	minRad = lambda * 0.25
	maxRad = lambda * 4.005

	incRad = 0.1 * lambda
)

func main() {
	flag.Parse()

	frame := 0
	for rad := minRad; rad <= maxRad; rad += incRad {
		totalField, scatteredField := computeOneFrame(rad)
		saveFrame(totalField, fmt.Sprintf("%v/mie-total-%03v.data", *outDir, frame))
		saveFrame(scatteredField, fmt.Sprintf("%v/mie-scattered-%03v.data", *outDir, frame))
		frame++
	}
}

func saveFrame(data []float64, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to create output file '%v': %v", filename, err))
	}
	defer f.Close()
	for _, v := range data {
		binary.Write(f, binary.LittleEndian, v)
	}
}

func computeOneFrame(rad float64) ([]float64, []float64) {
	workers := 2 * runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workers)
	totalField := make([]float64, *width**height)
	scatteredField := make([]float64, *width**height)
	cnt := int32(0)
	for w := 0; w < workers; w++ {
		go func(i int) {
			defer wg.Done()
			for x := *width / 2; x < *width; x++ {
				if x%workers != i {
					continue
				}
				for y := 0; y < *height; y++ {
					totalVal, scatteredVal := mieField(x, y, rad)
					totalField[y**width+x] = totalVal
					scatteredField[y**width+x] = scatteredVal
					// Symmetrically fill the other half.
					totalField[y**width+(*width-1-x)] = totalVal
					scatteredField[y**width+(*width-1-x)] = scatteredVal
					atomic.AddInt32(&cnt, 2)
				}
			}
		}(w)
	}
	total := *width * *height
	// Progress counter.
	counterDone := make(chan struct{})
	go func() {
		erase := strings.Repeat(" ", 80)
		nextMark := 1.0
		for {
			doneCnt := int(atomic.LoadInt32(&cnt))
			if doneCnt == total {
				fmt.Printf("\r%v\rR=%.2fλ rendering complete\n", erase, rad/lambda)
				close(counterDone)
				return
			}
			progress := float64(doneCnt) / float64(total) * 100.0
			if progress >= nextMark {
				fmt.Printf(
					"\r%v\rrendering for R=%.2fλ... %.2f%% done",
					erase,
					rad/lambda,
					progress,
				)
				nextMark = math.Ceil(progress)
			}
			runtime.Gosched()
		}
	}()
	wg.Wait()
	<-counterDone

	return totalField, scatteredField
}

func mieField(x, y int, rad float64) (float64, float64) {
	stepx, stepy := 2.0/float64(*width), 2.0/float64(*height)
	cx, cy := float64(*width-1)/2, float64(*height-1)/2

	fx, fy := stepx*(float64(x)-cx), stepy*(cy-float64(y))
	r := math.Sqrt(fx*fx + fy*fy)
	st, ct := math.Abs(fx/r), fy/r

	// j_l(x), j_l'(x).
	z1 := func(t float64) ([]*bigComplex, []*bigComplex) {
		jval, jder := sphericalBessel1(*maxL, t)
		zval := make([]*bigComplex, *maxL+1)
		zder := make([]*bigComplex, *maxL+1)
		for i := 0; i <= *maxL; i++ {
			zval[i] = bigComplexFromBigFloat(jval[i])
			zder[i] = bigComplexFromBigFloat(jder[i])
		}
		return zval, zder
	}

	z1c := func(z complex128) ([]*bigComplex, []*bigComplex) {
		return sphericalBessel1C(*maxL, z)
	}

	// h_l^1(x), h_l^1'(x).
	z3 := func(t float64) ([]*bigComplex, []*bigComplex) {
		jval, jder := sphericalBessel1(*maxL, t)
		yval, yder := sphericalBessel2(*maxL, t)
		zval := make([]*bigComplex, *maxL+1)
		zder := make([]*bigComplex, *maxL+1)
		for i := 0; i <= *maxL; i++ {
			zval[i] = newBigComplex(jval[i], yval[i])
			zder[i] = newBigComplex(jder[i], yder[i])
		}
		return zval, zder
	}

	useComplex := *nImg != 0

	ka := k * rad

	var (
		cn      *bigComplex
		cnka    *bigComplex
		intN    []*bigComplex
		intNder []*bigComplex
	)
	if useComplex {
		cn = bigComplexFromComplex128(complex(*n, *nImg))
		nka := complex(*n, *nImg) * complex(ka, 0.0)
		intN, intNder = z1c(nka)
		cnka = bigComplexFromComplex128(nka)
	} else {
		cn = bigComplexFromFloat64(*n)
		intN, intNder = z1(*n * ka)
		cnka = bigComplexFromFloat64(*n * ka)
	}
	cka := bigComplexFromFloat64(ka)
	intJ, intJder := z1(ka)
	intH, intHder := z3(ka)
	for i := 1; i <= *maxL; i++ {
		intJder[i] = intJder[i].mul(cka)
		intJder[i] = intJder[i].add(intJ[i])

		intHder[i] = intHder[i].mul(cka)
		intHder[i] = intHder[i].add(intH[i])

		intNder[i] = intNder[i].mul(cnka)
		intNder[i] = intNder[i].add(intN[i])
	}

	cmu := bigComplexFromFloat64(*mu)

	if r < rad {
		// Internal field.
		alpha := func(l int) *bigComplex {
			v := intJ[l].mul(intHder[l])
			v = v.sub(intJder[l].mul(intH[l]))
			v = v.mul(cmu)
			u := cmu.mul(intHder[l]).mul(intN[l])
			u = u.sub(intH[l].mul(intNder[l]))
			return v.quo(u)
		}
		beta := func(l int) *bigComplex {
			v := intJder[l].mul(intH[l])
			v = v.sub(intJ[l].mul(intHder[l]))
			v = v.mul(cmu).mul(cn)
			u := cmu.mul(intH[l]).mul(intNder[l])
			u = u.sub(cn.mul(cn).mul(intHder[l]).mul(intN[l]))
			return v.quo(u)
		}
		cnkr := complex(*n, *nImg) * complex(k*r, 0.0)
		intAmp := multipoleExpansion(
			st,
			ct,
			cnkr,
			alpha,
			beta,
			func(arg complex128) ([]*bigComplex, []*bigComplex) {
				if useComplex {
					return z1c(arg)
				}
				return z1(real(arg))
			},
		)
		return intAmp * intAmp, 0
	}
	// Scattered field coefficients.
	alpha := func(l int) *bigComplex {
		v := intJ[l].mul(intNder[l])
		v = v.sub(cmu.mul(intJder[l]).mul(intN[l]))
		u := cmu.mul(intHder[l]).mul(intN[l])
		u = u.sub(intH[l].mul(intNder[l]))
		return v.quo(u)
	}
	beta := func(l int) *bigComplex {
		v := cn.mul(cn).mul(intJder[l]).mul(intN[l])
		v = v.sub(cmu.mul(intJ[l]).mul(intNder[l]))
		u := cmu.mul(intH[l]).mul(intNder[l])
		u = u.sub(cn.mul(cn).mul(intHder[l]).mul(intN[l]))
		return v.quo(u)
	}
	// Scattered field plus incident field.
	scAmp := multipoleExpansion(
		st,
		ct,
		complex(k*r, 0),
		alpha,
		beta,
		func(arg complex128) ([]*bigComplex, []*bigComplex) {
			return z3(real(arg))
		},
	)
	incAmp := multipoleExpansion(
		st,
		ct,
		complex(k*r, 0),
		func(int) *bigComplex { return bigComplexFromInt(1) },
		func(int) *bigComplex { return bigComplexFromInt(1) },
		func(arg complex128) ([]*bigComplex, []*bigComplex) {
			return z1(real(arg))
		},
	)
	amp := scAmp + incAmp
	return amp * amp, scAmp * scAmp
}

func multipoleExpansion(
	st, ct float64,
	ckr complex128,
	alpha, beta func(int) *bigComplex,
	z func(complex128) ([]*bigComplex, []*bigComplex),
) float64 {
	cst, cct := bigComplexFromFloat64(st), bigComplexFromFloat64(ct)

	pval, pder := legendre(*maxL, ct)
	zval, zder := z(ckr)

	bigckr := bigComplexFromComplex128(ckr)
	// r-component.
	rcom := blankBigComplex()
	// theta-component.
	tcom := blankBigComplex()

	for l := 1; l <= *maxL; l++ {
		alphal, betal := alpha(l), beta(l)
		coeff := iPow(l - 1).mul(bigComplexFromFloat64(float64(2*l+1) / float64(l*(l+1))))
		zlkrkr := zval[l].quo(bigckr)
		ll1 := bigComplexFromInt(l * (l + 1))
		rr := ll1.mul(betal).mul(zlkrkr)
		rr = rr.mul(cst)
		rr = rr.mul(pder[l])

		g := zlkrkr.add(zder[l])
		h := cct.mul(pder[l])
		h = h.sub(ll1.mul(pval[l]))

		tt := iPow(1).mul(alphal).mul(zval[l]).mul(pder[l])
		tt = tt.sub(betal.mul(g).mul(h))

		rcom = rcom.add(coeff.mul(rr))
		tcom = tcom.add(coeff.mul(tt))
	}

	// Extract the x-polarized field.
	rre, tre := rcom.re, tcom.re
	rre.Mul(rre, fromFloat64(st))
	tre.Mul(tre, fromFloat64(ct))

	rre.Add(rre, tre)
	f64, _ := rre.Float64()
	return f64
}
