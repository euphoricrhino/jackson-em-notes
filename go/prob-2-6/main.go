package main

import (
	"flag"
	"fmt"
)

/*
This program calculates Jackson Prob 2.6 (c). But I have a different answer than 0.6189, what's wrong?
Example run:
go run main.go --iters 20000
*/
var (
	iters = flag.Int("iters", 100, "iterations")
)

func main() {
	flag.Parse()

	// Image charges on sphere a and b.
	cha := make([]charge, *iters)
	chb := make([]charge, *iters)
	// Initialize, will rescale total sphere charge to unity later.
	cha[0].q = 1.0
	cha[0].x = 0.0
	chb[0].q = 1.0
	chb[0].x = 0.0

	// Iteratively calculate image charges.
	qa := 1.0
	qb := 1.0
	for n := 1; n < *iters; n++ {
		cha[n].q = -chb[n-1].q / (2.0 - chb[n-1].x)
		qa += cha[n].q
		cha[n].x = 1.0 / (2.0 - chb[n-1].x)
		chb[n].q = -cha[n-1].q / (2.0 - cha[n-1].x)
		qb += chb[n].q
		chb[n].x = 1.0 / (2.0 - cha[n-1].x)
	}

	// Rescale each sphere's charge so each is charged with unit total charge.
	for _, ch := range cha {
		ch.q /= qa
	}

	for _, ch := range chb {
		ch.q /= qb
	}

	// Calculate force.
	// If all charges are located at each sphere's center, the force would have been 1/4 (ignoring 4πε_0 factor).
	f0 := 0.25
	f := 0.0
	for n := 0; n < *iters; n++ {
		for m := 0; m < *iters; m++ {
			d := 2.0 - cha[n].x - chb[m].x
			f += cha[n].q * chb[m].q / (d * d)
		}
	}
	fmt.Printf("f/f0 = %v\n", f/f0)
}

type charge struct {
	q float64
	x float64
}
