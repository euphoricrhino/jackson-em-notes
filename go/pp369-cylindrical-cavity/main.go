package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/euphoricrhino/jackson-em-notes/go/pkg/heatmap"
)

// Example commands:
// `go run main.go --p=0 --m=0 --xmn=2.405 --mode="TM (mnp=010)" --out-dir=./frames/tm-010`
// `go run main.go --p=4 --m=3 --xmn=9.761 --mode="TM (mnp=324)" --out-dir=./frames/tm-324`
// `go run main.go --p=2 --m=5 --xmn=18.9801 --mode="TM (mnp=542)" --out-dir=./frames/tm-542`
// `go run main.go --p=1 --m=1 -xmn=1.841 --mode="TE (mnp=111)" --out-dir=./frames/te-111`
// `go run main.go --p=3 --m=1 -xmn=5.331 --mode="TE (mnp=123)" --out-dir=./frames/te-123`
// `go run main.go --p=6 --m=4 -xmn=19.196 --mode="TE (mnp=456)" --out-dir=./frames/te-456`

var (
	width        = flag.Int("width", 1280, "image width")
	height       = flag.Int("height", 1280, "image height")
	unitInPixels = flag.Int("unit-in-pixels", 150, "unit length in pixels")
	hotHeatmap   = flag.String("hot-heatmap", "../heatmaps/hot.png", "hot heatmap")
	coldHeatmap  = flag.String("cold-heatmap", "../heatmaps/cold.png", "cold heatmap")

	orbitPeriods    = flag.Int("orbit-periods", 24, "one camera orbit in wave periods")
	framesPerPeriod = flag.Int("frames-per-period", 30, "frames per period")

	maxAmp     = flag.Float64("max-amp", .3, "maximum amplitude of field vectors")
	rhoSamples = flag.Int("rho-samples", 50, "samples in radial direction")
	phiSamples = flag.Int("phi-samples", 120, "samples in angular direction")
	zSamples   = flag.Int("z-samples", 50, "samples in longitudinal direction")
	p          = flag.Int("p", 0, "longitudinal mode number")
	m          = flag.Int("m", 0, "order of Bessel function")
	xmn        = flag.Float64("xmn", 0.0, "nth root of Jm(x) or J'm(x), depending on TM or TE")
	mode       = flag.String("mode", "", "has to start with TM|TE")

	outDir = flag.String("out-dir", "", "output dir")
)

const (
	radius       = 2.0
	d            = 7.0
	initRotation = -math.Pi * 7.0 / 180.0
	axisTheta    = math.Pi * 17.0 / 180.0
	tmMode       = "TM"
	teMode       = "TE"
)

type Vec [3]float64

func main() {
	flag.Parse()

	gridCnt := *zSamples * ((*phiSamples)*(*rhoSamples-1) + 1)

	grids := make([]Vec, gridCnt)
	dz := d / float64(*zSamples-1)
	dphi := 2.0 * math.Pi / float64(*phiSamples)
	drho := radius / float64(*rhoSamples-1)

	efields := make([]Vec, gridCnt*(*framesPerPeriod))
	hfields := make([]Vec, gridCnt*(*framesPerPeriod))
	idx := 0
	// First record the polar coordinates into grids [rho,phi,z].
	for nz := 0; nz < *zSamples; nz++ {
		z := dz * float64(nz)
		// Center of slice nz.
		grids[idx][2] = dz * float64(nz)
		idx++
		for nrho := 1; nrho < *rhoSamples; nrho++ {
			rho := drho * float64(nrho)
			for nphi := 0; nphi < *phiSamples; nphi++ {
				phi := dphi * float64(nphi)
				grids[idx][0], grids[idx][1], grids[idx][2] = rho, phi, z
				idx++
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(*framesPerPeriod)

	maxes := make([]float64, *framesPerPeriod)
	maxhs := make([]float64, *framesPerPeriod)
	for f := 0; f < *framesPerPeriod; f++ {
		go func(fr int) {
			idxStart := fr * gridCnt
			ef, hf := efields[idxStart:], hfields[idxStart:]
			omegat := 2 * math.Pi * float64(fr) / float64(*framesPerPeriod)
			for i := range grids {
				field(grids[i][0], grids[i][1], grids[i][2], omegat, &ef[i], &hf[i], &maxes[fr], &maxhs[fr])
			}
			wg.Done()
		}(f)
	}

	wg.Wait()

	// Normalize by the maximum value of all frames.
	maxe, maxh := maxes[0], maxhs[0]
	for f := 1; f < *framesPerPeriod; f++ {
		if maxes[f] > maxe {
			maxe = maxes[f]
		}
		if maxhs[f] > maxh {
			maxh = maxhs[f]
		}
	}

	escale := *maxAmp / maxe
	hscale := *maxAmp / maxh
	for i := 0; i < gridCnt*(*framesPerPeriod); i++ {
		efields[i][0] *= escale
		efields[i][1] *= escale
		efields[i][2] *= escale
		hfields[i][0] *= hscale
		hfields[i][1] *= hscale
		hfields[i][2] *= hscale
	}

	// Convert grids to cartesian.
	for i := range grids {
		rho, phi := grids[i][0], grids[i][1]
		grids[i][0], grids[i][1] = rho*math.Cos(phi), rho*math.Sin(phi)
		grids[i][2] -= d / 2
	}

	workers := runtime.NumCPU()
	wg.Add(workers)

	stheta, ctheta := math.Sin(axisTheta), math.Cos(axisTheta)
	frameCnt := *orbitPeriods * (*framesPerPeriod)
	drot := 2 * math.Pi / float64(frameCnt)

	ehm, err := heatmap.Load(*hotHeatmap, 1.0)
	if err != nil {
		panic(fmt.Sprintf("failed to load hot heatmap: %v", err))
	}
	hhm, err := heatmap.Load(*coldHeatmap, 1.0)
	if err != nil {
		panic(fmt.Sprintf("failed to load cold heatmap: %v", err))
	}
	for w := 0; w < workers; w++ {
		go func(wk int) {
			// Rotated grid and field vectors.
			rg := make([]Vec, gridCnt)
			re := make([]Vec, gridCnt)
			rh := make([]Vec, gridCnt)
			for f := 0; f < frameCnt; f++ {
				if f%workers != wk {
					continue
				}

				rot := initRotation + float64(f)*drot
				srot, crot := math.Sin(rot), math.Cos(rot)
				rotate := func(v Vec) Vec {
					u := Vec{stheta*v[0] - ctheta*v[1], ctheta*v[0] + stheta*v[1], v[2]}
					u = Vec{crot*u[0] - srot*u[2], u[1], srot*u[0] + crot*u[2]}
					return Vec{stheta*u[0] + ctheta*u[1], -ctheta*u[0] + stheta*u[1], u[2]}
				}
				// Rotate the grid points.
				for i := range rg {
					rg[i] = rotate(grids[i])
				}
				idxStart := (f % *framesPerPeriod) * gridCnt
				// Rotate the field vectors.
				ef, hf := efields[idxStart:], hfields[idxStart:]
				for i := range grids {
					re[i] = rotate(ef[i])
					rh[i] = rotate(hf[i])
				}
				if err := render(grids, rg, re, rh, ehm, hhm, f, frameCnt); err != nil {
					panic(err)
				}
			}
			wg.Done()
		}(w)
	}
	wg.Wait()
}

// Returns E,H field at (rho,phi,z,omegat)
func field(rho, phi, z, omegat float64, e, h *Vec, maxe, maxh *float64) {
	gamma := *xmn / radius
	gr := gamma * rho
	jm := math.Jn(*m, gr)
	jmder := besselDer(gr)

	zarg := float64(*p) * math.Pi * z / d
	szarg, czarg := math.Sin(zarg), math.Cos(zarg)
	sphi, cphi := math.Sin(phi), math.Cos(phi)
	phase := float64(*m)*phi - omegat
	sphase, cphase := math.Sin(phase), math.Cos(phase)
	cc, cs, sc, ss := cphase*cphi, cphase*sphi, sphase*cphi, sphase*sphi

	v := float64(*p) * math.Pi / (float64(d) * gamma * gamma)
	var a, b float64
	if rho == 0.0 {
		if *m == 1 {
			a, b = gamma/2, gamma/2
		}
	} else {
		a, b = gamma*jmder, float64(*m)/rho*jm
	}

	if strings.HasPrefix(*mode, tmMode) {
		e[2] = jm * czarg * cphase

		et := v * szarg
		eta, etb := et*a, et*b
		e[0] = -eta*cc - etb*ss
		e[1] = -eta*cs + etb*sc

		ht := czarg
		hta, htb := ht*a, ht*b
		h[0] = htb*cc + hta*ss
		h[1] = htb*cs - hta*sc
	} else if strings.HasPrefix(*mode, teMode) {
		h[2] = jm * szarg * cphase
		ht := v * czarg
		hta, htb := ht*a, ht*b
		h[0] = hta*cc + htb*ss
		h[1] = hta*cs - htb*sc

		et := szarg
		eta, etb := et*a, et*b
		e[0] = -etb*cc - eta*ss
		e[1] = -etb*cs + eta*sc
	}

	elen := math.Sqrt(e[0]*e[0] + e[1]*e[1] + e[2]*e[2])
	if elen > *maxe {
		*maxe = elen
	}
	hlen := math.Sqrt(h[0]*h[0] + h[1]*h[1] + h[2]*h[2])
	if hlen > *maxh {
		*maxh = hlen
	}
}

func toPixels(x, y float64) (float64, float64) {
	return float64(*width/2) - x*float64(*unitInPixels), float64(*height/2) - y*float64(*unitInPixels)
}

func render(grids, rg, re, rh []Vec, ehm, hhm []color.Color, f, frameCnt int) error {
	img := image.NewRGBA(image.Rect(0, 0, *width, *height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)

	gc := draw2dimg.NewGraphicContext(img)
	gc.SetLineWidth(1)
	draw2d.SetFontFolder("/Users/xni/Library/Fonts")
	draw2d.SetFontNamer(func(_ draw2d.FontData) string { return "MonoLisaVariableNormal.ttf" })
	textColor := color.RGBA{0, 0xcc, 0xcc, 0xff}
	gc.SetFillColor(textColor)
	gc.SetStrokeColor(textColor)
	gc.SetDPI(288)
	gc.SetFontSize(5)
	gc.FillStringAt(*mode, 40.0, 40.0)

	text := ""
	if f < frameCnt/3 {
		text += "E only"
		textColor = color.RGBA{0xff, 0, 0, 0xff}
		gc.SetFillColor(textColor)
		gc.SetStrokeColor(textColor)
	} else if f < frameCnt*2/3 {
		textColor = color.RGBA{0, 0xff, 0, 0xff}
		gc.SetFillColor(textColor)
		gc.SetStrokeColor(textColor)
		text += "H only"
	} else {
		textColor = color.RGBA{0, 0xcc, 0xcc, 0xff}
		gc.SetFillColor(textColor)
		gc.SetStrokeColor(textColor)
		text += "both E and H"
	}
	gc.FillStringAt(text, 40.0, 70.0)

	for i := range rg {
		// Look up the color of e,h from heatmap based on original z.
		t := (grids[i][2] + d/2) / d

		// Draw E/H line in orthographic projection.
		if f < frameCnt/3 || f >= frameCnt*2/3 {
			epos := int(t * float64(len(ehm)-1))
			r, g, b, a := ehm[epos].RGBA()
			gc.SetStrokeColor(color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
			gc.MoveTo(toPixels(rg[i][0], rg[i][1]))
			gc.LineTo(toPixels(rg[i][0]+re[i][0], rg[i][1]+re[i][1]))
			gc.Stroke()
		}

		if f >= frameCnt/3 {
			hpos := int(t * float64(len(hhm)-1))
			r, g, b, a := hhm[hpos].RGBA()
			gc.SetStrokeColor(color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
			gc.MoveTo(toPixels(rg[i][0], rg[i][1]))
			gc.LineTo(toPixels(rg[i][0]+rh[i][0], rg[i][1]+rh[i][1]))
			gc.Stroke()
		}
	}
	fn := filepath.Join(*outDir, fmt.Sprintf("frame-%04v.png", f))
	file, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("failed to create output file '%v': %v", fn, err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode to PNG: %v", err)
	}
	fmt.Fprintf(os.Stdout, "generated %v\n", fn)
	return nil
}

// Derivative of Jm at x.
func besselDer(x float64) float64 {
	if *m == 0 {
		return -math.J1(x)
	}
	return (math.Jn(*m-1, x) - math.Jn(*m+1, x)) / 2
}
