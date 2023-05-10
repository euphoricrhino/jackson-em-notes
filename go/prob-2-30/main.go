package main

import (
	"flag"
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

/*
Program to solve Poisson equation for Jackson prob 2.30 using FEA.
This is by no means an optimized program, the purpose is to showcase how to construct the linear system using FEA,
in particular how to use the integral relations (2.81) between neighboring grids.
*/

var (
	spacings = flag.Int("spacings", 0, "spacings")
)

func main() {
	flag.Parse()
	if *spacings%4 != 0 {
		panic("--spacings is not multiple of 4")
	}
	n := *spacings - 1
	// We have n^2 unknowns, and n equations, each with up to 9 unknowns in that equation.
	grid := make([]float64, n*n*n*n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fillMat(grid, i, j, n)
		}
	}

	b := make([]float64, n*n)
	// All n equation's RHS is 4πh^2.
	for i := range b {
		b[i] = 4 * math.Pi * (1.0 / float64(n+1) * 1.0 / float64(n+1))
	}
	// Solve the n^2xn^2 system using brute force, no optimization whatsoever.
	bmat := mat.NewDense(n*n, 1, b)
	a := mat.NewSymDense(n*n, grid)
	var solution mat.Dense
	if err := solution.Solve(a, bmat); err != nil {
		panic(err)
	}

	getSolution := func(i, j int) float64 {
		return solution.At(j*n+i, 0)
	}
	fmt.Printf("potential at (0.25, 0.25): %v\n", getSolution(*spacings/4, *spacings/4))
	fmt.Printf("potential at (0.5, 0.25): %v\n", getSolution(*spacings/2, *spacings/4))
	fmt.Printf("potential at (0.5, 0.5): %v\n", getSolution(*spacings/2, *spacings/2))
}

func fillMat(grid []float64, i, j, n int) {
	// Index for the equation for grid (i,j).
	idxij := j*n + i
	// Starting position of the unknowns in the linear system's n^2*n^2 matrix.
	idxStart := idxij * n * n
	inRange := func(x, y int) bool {
		return x >= 0 && x < n && y >= 0 && y < n
	}
	set := func(x, y int, v float64) {
		idx := y*n + x
		grid[idxStart+idx] = v
	}
	const (
		selfWeight     = 8.0 / 3.0
		neighborWeight = -1.0 / 3.0
	)
	// (i,j) on itself.
	set(i, j, selfWeight)

	// Potentially 8 neighboring grids.
	// Northwest.
	if inRange(i-1, j-1) {
		set(i-1, j-1, neighborWeight)
	}
	// North.
	if inRange(i, j-1) {
		set(i, j-1, neighborWeight)
	}
	// Northeast.
	if inRange(i+1, j-1) {
		set(i+1, j-1, neighborWeight)
	}
	// East.
	if inRange(i+1, j) {
		set(i+1, j, neighborWeight)
	}
	// Southeast.
	if inRange(i+1, j+1) {
		set(i+1, j+1, neighborWeight)
	}
	// South.
	if inRange(i, j+1) {
		set(i, j+1, neighborWeight)
	}
	// Southwest.
	if inRange(i-1, j+1) {
		set(i-1, j+1, neighborWeight)
	}
	// West.
	if inRange(i-1, j) {
		set(i-1, j, neighborWeight)
	}
}
