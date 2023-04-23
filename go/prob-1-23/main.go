package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

/*
Program to calculate Jackson prob 1.23 using relaxation method.
Example plots:
octave --persist /var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/jackson_prob_1_23.m
*/
var (
	errBound = flag.Float64("err-bound", 1e-5, "error bound")
	spacings = flag.Int("spacings", 0, "spacings")
)

func main() {
	flag.Parse()
	n := *spacings
	if n%2 == 1 {
		panic("need even --spacings")
	}
	buf := make([]float64, (2*n+1)*(n+1))

	// Initialize inner y=0 line to high potential value.
	for x := 0; x <= n; x++ {
		buf[x] = 1.0
	}
	// Initialize y=n line to zero potential value.
	for x := 0; x <= 2*n; x++ {
		buf[n*(2*n+1)+x] = 0.0
	}
	// Initialize interior potential with interpolated y.
	for y := 1; y < n; y++ {
		v := float64(n-y) / float64(n)
		for x := 0; x <= n+y; x++ {
			buf[y*(2*n+1)+x] = v
		}
	}
	var e float64
	for i := 0; true; i++ {
		buf, e = oneIteration(buf, n)
		if i%100 == 0 {
			fmt.Printf("iteration %v, max-error: %v\n", i, e)
		}
		if e < *errBound {
			fmt.Printf("converged after %v iterations\n", i)
			break
		}
	}

	getVal := func(x, y int) float64 {
		return get(buf, n, x, y) * 100
	}

	fmt.Printf("Φ1=%v, Φ2=%v, Φ3=%v, Φ4=%v\n", getVal(0, n/2), getVal(n/2, n/2), getVal(n, n/2), getVal(n*3/2, n/2))

	mesh := make([]float64, (4*n+1)*(4*n+1))
	tx := make([]float64, 4*n+1)
	ty := make([]float64, 4*n+1)
	step := 1 / float64(n)
	abs := func(x int) int {
		if x >= 0 {
			return x
		}
		return -x
	}

	for x := -(2 * n); x <= 2*n; x++ {
		tx[x+2*n] = float64(x) * step
		ty[x+2*n] = float64(x) * step
		absx := abs(x)
		for y := -(2 * n); y <= 2*n; y++ {
			absy := abs(y)
			idx := (y+2*n)*(4*n+1) + x + 2*n
			if absx <= n && absy <= n {
				mesh[idx] = 100
			}
			if absx <= n && absy > n {
				mesh[idx] = getVal(absx, absy-n)
			}
			if absx > n && absy <= n {
				mesh[idx] = getVal(absy, absx-n)
			}
			if absx > n && absy > n {
				dx, dy := absx-n, absy-n
				if dx <= dy {
					mesh[idx] = getVal(n+dx, dy)
				} else {
					mesh[idx] = getVal(n+dy, dx)
				}
			}
		}
	}

	// Write an octave script for visualization of the potential.
	filename := filepath.Join(os.TempDir(), "jackson_prob_1_23.m")
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "tx=%v;\n", tx)
	fmt.Fprintf(file, "ty=%v;\n", ty)
	fmt.Fprintf(file, "phi=[\n")
	for y := 0; y <= 4*n; y++ {
		for x := 0; x <= 4*n; x++ {
			fmt.Fprintf(file, "%v ", mesh[y*(4*n+1)+x])
		}
		fmt.Fprintf(file, "\n")
	}
	fmt.Fprintf(file, "];\n")
	fmt.Fprintf(file, "mesh(tx, ty, phi);\n")
	fmt.Println(filename)
}

// "Improvied" averaging scheme, mixing cross and square scheme with 4:1 weighting.
func improvedAvg(n, ne, e, se, s, sw, w, nw float64) float64 {
	sc := n + e + s + w
	ss := ne + se + sw + nw

	return 0.2*sc + 0.05*ss
}

// Calculates the potential value for the grid for one iteration.
func oneIteration(buf []float64, n int) ([]float64, float64) {
	ret := make([]float64, len(buf))
	copy(ret, buf)

	const maxWorkers = 20
	var wg sync.WaitGroup
	workers := maxWorkers
	if n-1 < maxWorkers {
		workers = n - 1
	}
	wg.Add(workers)
	maxErr := make([]float64, workers)
	for w := 0; w < workers; w++ {
		go func(worker int) {
			maxErr[worker] = -1.0
			defer wg.Done()
			for y := 1; y < n; y++ {
				if y%workers != worker {
					continue
				}
				for x := 0; x <= n+y; x++ {
					newVal := improvedAvg(
						get(buf, n, x, y-1),
						get(buf, n, x+1, y-1),
						get(buf, n, x+1, y),
						get(buf, n, x+1, y+1),
						get(buf, n, x, y+1),
						get(buf, n, x-1, y+1),
						get(buf, n, x-1, y),
						get(buf, n, x-1, y-1),
					)
					ret[y*(2*n+1)+x] = newVal
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
		if maxE < maxErr[w] {
			maxE = maxErr[w]
		}
	}
	return ret, maxE
}

// Reads the grid potential values.
func get(buf []float64, n, x, y int) float64 {
	if y < 0 || y > n {
		panic("out of bound y")
	}
	if x > 2*n {
		panic("out of bound x")
	}
	if x-n > y+2 {
		panic(fmt.Sprintf("too far from diagonal: (%v, %v)\n", x, y))
	}

	if x < 0 {
		x = -x
	}
	if x-n > y {
		// Flip across diagonal.
		x, y = n+y, x-n
	}
	return buf[y*(2*n+1)+x]
}
