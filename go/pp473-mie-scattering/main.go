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
	width  = flag.Int("width", 800, "Width of the image")
	height = flag.Int("height", 800, "Height of the image")
	maxL   = flag.Int("max-l", 20, "max l value")

	n  = flag.Float64("n", 1.2, "Refractive index")
	mu = flag.Float64("mu", 1.0, "Permeability")
)

const (
	lambda = 0.2
	k      = 2.0 * math.Pi / lambda
	minRad = lambda * 0.25
	maxRad = lambda * 4.005

	incRad = 0.01 * lambda
)

func main() {
	flag.Parse()

	frame := 0
	for rad := minRad; rad <= maxRad; rad += incRad {
		totalField, scatteredField := computeOneFrame(rad)
		saveFrame(totalField, fmt.Sprintf("mie-total-%03v.data", frame))
		saveFrame(scatteredField, fmt.Sprintf("mie-scattered-%03v.data", frame))
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
	workers := runtime.NumCPU()
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
	z1 := func(t float64) ([]*complex, []*complex) {
		jval, jder := sphericalBessel1(*maxL, t)
		zval := make([]*complex, *maxL+1)
		zder := make([]*complex, *maxL+1)
		for i := 0; i <= *maxL; i++ {
			zval[i] = complexFromReal(jval[i])
			zder[i] = complexFromReal(jder[i])
		}
		return zval, zder
	}

	// h_l^1(x), h_l^1'(x).
	z3 := func(t float64) ([]*complex, []*complex) {
		jval, jder := sphericalBessel1(*maxL, t)
		yval, yder := sphericalBessel2(*maxL, t)
		zval := make([]*complex, *maxL+1)
		zder := make([]*complex, *maxL+1)
		for i := 0; i <= *maxL; i++ {
			zval[i] = newComplex(jval[i], yval[i])
			zder[i] = newComplex(jder[i], yder[i])
		}
		return zval, zder
	}

	ka := k * rad
	nka := *n * ka
	cka := complexFromFloat64(ka)
	cnka := complexFromFloat64(nka)
	intJ, intJder := z1(ka)
	intH, intHder := z3(ka)
	intN, intNder := z1(nka)
	for i := 1; i <= *maxL; i++ {
		intJder[i] = intJder[i].mul(cka)
		intJder[i] = intJder[i].add(intJ[i])

		intHder[i] = intHder[i].mul(cka)
		intHder[i] = intHder[i].add(intH[i])

		intNder[i] = intNder[i].mul(cnka)
		intNder[i] = intNder[i].add(intN[i])
	}

	cmu := complexFromFloat64(*mu)
	cn := complexFromFloat64(*n)

	if r < rad {
		// Internal field.
		alpha := func(l int) *complex {
			v := intJ[l].mul(intHder[l])
			v = v.sub(intJder[l].mul(intH[l]))
			v = v.mul(cmu)
			u := cmu.mul(intHder[l]).mul(intN[l])
			u = u.sub(intH[l].mul(intNder[l]))
			return v.quo(u)
		}
		beta := func(l int) *complex {
			v := intJder[l].mul(intH[l])
			v = v.sub(intJ[l].mul(intHder[l]))
			v = v.mul(cmu).mul(cn)
			u := cmu.mul(intH[l]).mul(intNder[l])
			u = u.sub(cn.mul(cn).mul(intHder[l]).mul(intN[l]))
			return v.quo(u)
		}
		intAmp := multipoleExpansion(*n*k*r, st, ct, alpha, beta, z1)
		return intAmp * intAmp, 0
	}
	// Scattered field coefficients.
	alpha := func(l int) *complex {
		v := intJ[l].mul(intNder[l])
		v = v.sub(cmu.mul(intJder[l]).mul(intN[l]))
		u := cmu.mul(intHder[l]).mul(intN[l])
		u = u.sub(intH[l].mul(intNder[l]))
		return v.quo(u)
	}
	beta := func(l int) *complex {
		v := cn.mul(cn).mul(intJder[l]).mul(intN[l])
		v = v.sub(cmu.mul(intJ[l]).mul(intNder[l]))
		u := cmu.mul(intH[l]).mul(intNder[l])
		u = u.sub(cn.mul(cn).mul(intHder[l]).mul(intN[l]))
		return v.quo(u)
	}
	// Scattered field plus incident field.
	scAmp := multipoleExpansion(k*r, st, ct, alpha, beta, z3)
	incAmp := multipoleExpansion(
		k*r,
		st,
		ct,
		func(int) *complex { return complexFromInt(1) },
		func(int) *complex { return complexFromInt(1) },
		z1,
	)
	amp := scAmp + incAmp
	return amp * amp, scAmp * scAmp
}

func multipoleExpansion(
	kr, st, ct float64,
	alpha, beta func(int) *complex,
	z func(float64) ([]*complex, []*complex),
) float64 {
	cst, cct := complexFromFloat64(st), complexFromFloat64(ct)

	pval, pder := legendre(*maxL, ct)
	zval, zder := z(kr)

	ckr := complexFromFloat64(kr)
	// r-component.
	rcom := blankComplex()
	// theta-component.
	tcom := blankComplex()

	for l := 1; l <= *maxL; l++ {
		alphal, betal := alpha(l), beta(l)
		coeff := iPow(l - 1).mul(complexFromFloat64(float64(2*l+1) / float64(l*(l+1))))
		zlkrkr := zval[l].quo(ckr)
		ll1 := complexFromInt(l * (l + 1))
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
