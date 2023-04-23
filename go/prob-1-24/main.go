package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
)

/*
Program to solve the Poisson equation for Jackson prob 1.24 using relaxation method.
*/
var (
	errBound = flag.Float64("err-bound", 1e-5, "error bound")
	spacings = flag.Int("spacings", 0, "spacings")
)

func main() {
	flag.Parse()
	n := *spacings
	if n%4 != 0 {
		panic("--spacings is not multiple of 4")
	}
	n /= 2
	buf := make([]float64, (n+1)*(n+1))
	// Initial guess for interiors.
	for i := range buf {
		buf[i] = 1.0
	}

	// Boundary conditions.
	for i := 0; i <= n; i++ {
		buf[n*(n+1)+i] = 0.0
		buf[i*(n+1)+n] = 0.0
	}

	var e float64
	h := 0.5 / float64(n)
	for i := 0; true; i++ {
		buf, e = oneIteration(buf, n, h)
		if i%100 == 0 {
			fmt.Printf("iteration %v, max-error: %v\n", i, e)
		}
		if e < *errBound {
			fmt.Printf("converged after %v iterations\n", i)
			break
		}
	}

	fmt.Printf("potential at (0.25, 0.25): %v\n", get(buf, n, n/2, n/2))
	fmt.Printf("potential at (0.5, 0.25): %v\n", get(buf, n, 0, n/2))
	fmt.Printf("potential at (0.5, 0.5): %v\n", get(buf, n, 0, 0))
}

func get(buf []float64, n, x, y int) float64 {
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return buf[y*(n+1)+x]
}

func improvedAvg(h, n, ne, e, se, s, sw, w, nw float64) float64 {
	sc := n + e + s + w
	ss := ne + se + sw + nw

	// Recall that we are actually calculating 4πε_0 times the potential, where g=ε_0 uniformally on all grids.
	return 0.2*sc + 0.05*ss + 0.3*h*h*(math.Pi*4)
}

func oneIteration(buf []float64, n int, h float64) ([]float64, float64) {
	ret := make([]float64, len(buf))
	copy(ret, buf)

	const maxWorkers = 20
	workers := maxWorkers
	if workers > n {
		workers = n
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	maxErr := make([]float64, workers)

	for w := 0; w < workers; w++ {
		go func(worker int) {
			defer wg.Done()
			maxErr[worker] = -1.0
			for y := 0; y < n; y++ {
				if y%workers != worker {
					continue
				}
				for x := 0; x < n; x++ {
					newVal := improvedAvg(
						h,
						get(buf, n, x, y-1),
						get(buf, n, x+1, y-1),
						get(buf, n, x+1, y),
						get(buf, n, x+1, y+1),
						get(buf, n, x, y+1),
						get(buf, n, x-1, y+1),
						get(buf, n, x-1, y),
						get(buf, n, x-1, y-1),
					)
					ret[y*(n+1)+x] = newVal
					e := newVal - get(buf, n, x, y)
					if e < 0 {
						e = -e
					}
					if maxErr[worker] < 0 || maxErr[worker] < e {
						maxErr[worker] = e
					}
				}
			}
		}(w)
	}

	wg.Wait()
	maxE := maxErr[0]
	for w := 1; w < workers; w++ {
		if maxErr[w] < maxE {
			maxE = maxErr[w]
		}
	}
	return ret, maxE
}
