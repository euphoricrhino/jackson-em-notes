package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	legendrezeros "github.com/euphoricrhino/jackson-em-notes/go/pp105-legendre-zeros"
)

var (
	prec             = flag.Uint("prec", 1000, "precision")
	legendreErrBound = flag.Float64("legendre-err-bound", 1e-10, "error bound for computing legendre polynomials")
)

func main() {
	flag.Parse()
	legendrezeros.SetPrecOnce(*prec)
	filename := filepath.Join(os.TempDir(), "plot-legendre.m")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	leb := legendrezeros.NewFloat(*legendreErrBound)

	var x []float64
	for i := -100; i <= 100; i++ {
		x = append(x, float64(i)/100)
	}
	smallNus := []float64{0, 0.05, 0.25, 0.75, 1}
	for i, nu := range smallNus {
		xstart := 0
		if math.Floor(nu) != nu {
			// Exclude x=-1.0 for non-integer 𝜈.
			xstart = 1
		}
		nuf := legendrezeros.NewFloat(nu)
		var y []float64
		for j := xstart; j < len(x); j++ {
			f64, _ := legendrezeros.Legendre(legendrezeros.NewFloat(x[j]), nuf, leb).Float64()
			y = append(y, f64)
		}
		fmt.Fprintf(f, "x%v=%v;\n", i, x[xstart:])
		fmt.Fprintf(f, "y%v=%v;\n", i, y)
		fmt.Fprintln(f, "subplot(2,1,1);")
		fmt.Fprintf(f, "plot(x%v,y%v,'LineWidth',2);hold on;\n", i, i)
	}
	fmt.Fprintln(f, "plot(x0,zeros(length(x0)),'k--');")
	fmt.Fprintf(f, "legend(")
	var parts []string
	for _, nu := range smallNus {
		parts = append(parts, fmt.Sprintf("'𝜈=%v'", nu))
	}
	fmt.Fprintf(f, strings.Join(parts, ","))

	fmt.Fprintf(f, ",'location','southeast');")
	fmt.Fprintln(f, "title('𝜈 ≤ 1');")
	fmt.Fprintln(f, "xlabel('x');")
	fmt.Fprintln(f, "ylabel('P_{𝜈}(x)');")
	fmt.Fprintln(f, "set(gca,'FontSize',18);")

	bigNus := []float64{1, 1.25, 3, 3.75, 4}
	for i, nu := range bigNus {
		xstart := 0
		if math.Floor(nu) != nu {
			// Exclude x=-1.0 for non-integer 𝜈.
			xstart = 1
		}
		nuf := legendrezeros.NewFloat(nu)
		var y []float64
		for j := xstart; j < len(x); j++ {
			f64, _ := legendrezeros.Legendre(legendrezeros.NewFloat(x[j]), nuf, leb).Float64()
			y = append(y, f64)
		}
		fmt.Fprintf(f, "x%v=%v;\n", i, x[xstart:])
		fmt.Fprintf(f, "y%v=%v;\n", i, y)
		fmt.Fprintln(f, "subplot(2,1,2);")
		fmt.Fprintf(f, "plot(x%v,y%v,'LineWidth',2);hold on;\n", i, i)
	}
	fmt.Fprintln(f, "plot(x0,zeros(length(x0)),'k--');")
	fmt.Fprintf(f, "legend(")
	parts = nil
	for _, nu := range bigNus {
		parts = append(parts, fmt.Sprintf("'𝜈=%v'", nu))
	}
	fmt.Fprintf(f, strings.Join(parts, ","))

	fmt.Fprintf(f, ",'location','southeast');")
	fmt.Fprintln(f, "title('𝜈 ≥ 1');")
	fmt.Fprintln(f, "xlabel('x');")
	fmt.Fprintln(f, "ylabel('P_{𝜈}(x)');")
	fmt.Fprintln(f, "set(gca,'FontSize',18);")

	fmt.Println(filename)
}
