package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/graphix/zraster"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Example commands:
// `go run main.go --p=0 --m=0 --xmn=2.405 --mode="TM (mnp=010)" --out-dir=./frames/tm-010`
// `go run main.go --p=4 --m=3 --xmn=9.761 --mode="TM (mnp=324)" --out-dir=./frames/tm-324`
// `go run main.go --p=2 --m=5 --xmn=18.9801 --mode="TM (mnp=542)" --out-dir=./frames/tm-542`
// `go run main.go --p=1 --m=1 -xmn=1.841 --mode="TE (mnp=111)" --out-dir=./frames/te-111`
// `go run main.go --p=3 --m=1 -xmn=5.331 --mode="TE (mnp=123)" --out-dir=./frames/te-123`
// `go run main.go --p=6 --m=4 -xmn=19.196 --mode="TE (mnp=456)" --out-dir=./frames/te-456`

var (
	hotHeatmap  = flag.String("hot-heatmap", "../heatmaps/hot.png", "hot heatmap")
	coldHeatmap = flag.String("cold-heatmap", "../heatmaps/cold.png", "cold heatmap")

	p    = flag.Int("p", 0, "longitudinal mode number")
	m    = flag.Int("m", 0, "order of Bessel function")
	xmn  = flag.Float64("xmn", 0.0, "nth root of Jm(x) or J'm(x), depending on TM or TE")
	mode = flag.String("mode", "", "has to start with TM|TE")

	outDir = flag.String("out-dir", "", "output dir")
)

const (
	width           = 1280
	height          = 1280
	orbitPeriods    = 24
	framesPerPeriod = 30

	maxAmp     = .3
	rhoSamples = 50
	phiSamples = 120
	zSamples   = 50

	radius = 2.0
	d      = 7.0
	// Camera orbiting circle axis, which is in the x-y plane, forming this angle with the x-axis.
	cameraOrbitAxisAngle   = math.Pi * 17.0 / 180.0
	cameraOrbitAngleOffset = -math.Pi * 7.0 / 180.0
	tmMode                 = "TM"
	teMode                 = "TE"

	minFading   = 0.2
	maxFading   = 1
	fadingGamma = 0.5
)

func main() {
	flag.Parse()

	gridCnt := zSamples * (phiSamples*(rhoSamples-1) + 1)

	grids := make([]*graphix.Vec3, gridCnt)
	dz := d / float64(zSamples-1)
	dphi := 2.0 * math.Pi / float64(phiSamples)
	drho := radius / float64(rhoSamples-1)

	efields := make([]*graphix.Vec3, gridCnt*framesPerPeriod)
	hfields := make([]*graphix.Vec3, gridCnt*framesPerPeriod)
	idx := 0
	// First record the polar coordinates into grids [rho,phi,z].
	for nz := 0; nz < zSamples; nz++ {
		z := dz * float64(nz)
		// Center of slice nz.
		grids[idx] = graphix.NewVec3(0, 0, z)
		idx++
		for nrho := 1; nrho < rhoSamples; nrho++ {
			rho := drho * float64(nrho)
			for nphi := 0; nphi < phiSamples; nphi++ {
				phi := dphi * float64(nphi)
				grids[idx] = graphix.NewVec3(rho, phi, z)
				idx++
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(framesPerPeriod)

	// For each frame, get its max-E and max-H separately for concurrency safety.
	maxes := make([]float64, framesPerPeriod)
	maxhs := make([]float64, framesPerPeriod)
	for f := 0; f < framesPerPeriod; f++ {
		go func(fr int) {
			idxStart := fr * gridCnt
			ef, hf := efields[idxStart:], hfields[idxStart:]
			omegat := 2 * math.Pi * float64(fr) / float64(framesPerPeriod)
			for i, grid := range grids {
				ef[i], hf[i] = field(grid, omegat)
				maxes[fr] = math.Max(maxes[fr], ef[i].Norm())
				maxhs[fr] = math.Max(maxhs[fr], hf[i].Norm())
			}
			wg.Done()
		}(f)
	}
	wg.Wait()

	// Normalize by the maximum value of all frames.
	maxe, maxh := math.Inf(-1), math.Inf(-1)
	for f := 0; f < framesPerPeriod; f++ {
		maxe = math.Max(maxe, maxes[f])
		maxh = math.Max(maxh, maxhs[f])
	}

	escale := maxAmp / maxe
	for _, e := range efields {
		e.Scale(e, escale)
	}
	hscale := maxAmp / maxh
	for _, h := range hfields {
		h.Scale(h, hscale)
	}

	// Convert grids to cartesian, and shift the cylinder to center at origin.
	for _, grid := range grids {
		rho, phi := grid[0], grid[1]
		grid[0], grid[1] = rho*math.Cos(phi), rho*math.Sin(phi)
		grid[2] -= d / 2
	}

	ehm, err := graphix.LoadHeatmap(*hotHeatmap, 1.0)
	if err != nil {
		panic(fmt.Sprintf("failed to load hot heatmap: %v", err))
	}
	hhm, err := graphix.LoadHeatmap(*coldHeatmap, 1.0)
	if err != nil {
		panic(fmt.Sprintf("failed to load cold heatmap: %v", err))
	}

	frameCnt := orbitPeriods * framesPerPeriod

	// Configuration for the camera orbit.
	n := graphix.NewVec3(math.Cos(cameraOrbitAxisAngle), math.Sin(cameraOrbitAxisAngle), 0)
	pos := graphix.NewVec3(0, 0, -d)
	forward := graphix.NewVec3(0, 0, 1)
	up := graphix.NewVec3(0, 1, 0)
	pr := graphix.NewPerspective(5 * d / 7)
	sc := graphix.NewScreen(width, height, -2*radius, -2*radius, 2*radius, 2*radius)
	cir := graphix.NewCircularCameraOrbit(n, pos, forward, up, frameCnt, cameraOrbitAngleOffset, pr, sc)

	for f := 0; f < frameCnt; f++ {
		paths := make([]*zraster.SpacePath, 0, 2*gridCnt)
		idxStart := (f % framesPerPeriod) * gridCnt
		ef, hf := efields[idxStart:], hfields[idxStart:]
		for i := 0; i < gridCnt; i++ {
			// Look up the color of e,h from heatmap based on original z.
			t := (grids[i][2] + d/2) / d
			// Showing E field only for the first 1/3 of frames, H field only for the second 1/3, and E+H for the last.
			if f < frameCnt/3 || f >= frameCnt*2/3 {
				epos := int(t * float64(len(ehm)-1))
				paths = append(paths, &zraster.SpacePath{
					Segments: []*zraster.SpaceVertex{{
						Pos:   grids[i],
						Color: makeTransparency(ehm[epos], ef[i].Norm(), maxe),
					}},
					End:       graphix.BlankVec3().Add(grids[i], ef[i]),
					LineWidth: 1,
				})
			}
			if f >= frameCnt/3 {
				hpos := int(t * float64(len(hhm)-1))
				paths = append(paths, &zraster.SpacePath{
					Segments: []*zraster.SpaceVertex{{
						Pos:   grids[i],
						Color: makeTransparency(hhm[hpos], hf[i].Norm(), maxh),
					}},
					End:       graphix.BlankVec3().Add(grids[i], hf[i]),
					LineWidth: 1,
				})
			}
		}

		img := zraster.Run(zraster.Options{
			Camera:  cir.GetCamera(f),
			Paths:   paths,
			Workers: runtime.NumCPU(),
		})
		renderCaption(img, f, frameCnt)
		fn := filepath.Join(*outDir, fmt.Sprintf("frame-%04v.png", f))
		file, err := os.Create(fn)
		if err != nil {
			panic(fmt.Sprintf("failed to create output file '%v': %v", fn, err))
		}
		if err := png.Encode(file, img); err != nil {
			panic(fmt.Sprintf("failed to encode to PNG: %v", err))
		}
		file.Close()
		fmt.Fprintf(os.Stdout, "generated %v\n", fn)
	}
}

// Derivative of Jm at x.
func besselDer(x float64) float64 {
	if *m == 0 {
		return -math.J1(x)
	}
	return (math.Jn(*m-1, x) - math.Jn(*m+1, x)) / 2
}

// Returns E,H field at (rho,phi,z,omegat).
func field(grid *graphix.Vec3, omegat float64) (*graphix.Vec3, *graphix.Vec3) {
	rho, phi, z := grid[0], grid[1], grid[2]
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

	var e, h *graphix.Vec3
	if strings.HasPrefix(*mode, tmMode) {
		et := v * szarg
		eta, etb := et*a, et*b
		e = graphix.NewVec3(-eta*cc-etb*ss, -eta*cs+etb*sc, jm*czarg*cphase)

		ht := czarg
		hta, htb := ht*a, ht*b
		h = graphix.NewVec3(htb*cc+hta*ss, htb*cs-hta*sc, 0)
	} else if strings.HasPrefix(*mode, teMode) {
		ht := v * czarg
		hta, htb := ht*a, ht*b
		h = graphix.NewVec3(hta*cc+htb*ss, hta*cs-htb*sc, jm*szarg*cphase)

		et := szarg
		eta, etb := et*a, et*b
		e = graphix.NewVec3(-etb*cc-eta*ss, -etb*cs+eta*sc, 0)
	}
	return e, h
}

func renderCaption(img draw.Image, f, frameCnt int) {
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
}

// Make the transparency for a color based on the field strength compared to the maximum field strength.
func makeTransparency(c color.Color, field, maxField float64) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(255 * (minFading + (maxFading-minFading)*math.Pow(field/maxField, fadingGamma))),
	}
}
