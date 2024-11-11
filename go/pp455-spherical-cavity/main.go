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
// `go run *.go --l=2 --m=0 --xln=9.09501 --mode="TE (lmn=202)" --out-dir=./frames/te-202`
// `go run *.go --l=3 --m=1 --xln=20.12181 --mode="TE (lmn=315)" --out-dir=./frames/te-315`
// `go run *.go --l=4 --m=3 --xln=15.03966 --mode="TE (lmn=433)" --out-dir=./frames/te-433`
// `go run *.go --l=2 --m=2 --xln=17.10274 --mode="TM (lmn=225)" --out-dir=./frames/tm-225`
// `go run *.go --l=3 --m=2 --xln=4.97342 --mode="TM (lmn=321)" --out-dir=./frames/tm-321`
// `go run *.go --l=4 --m=1 --xln=9.96755 --mode="TM (lmn=412)" --out-dir=./frames/tm-412`
var (
	hotHeatmap  = flag.String("hot-heatmap", "../heatmaps/hot.png", "hot heatmap")
	coldHeatmap = flag.String("cold-heatmap", "../heatmaps/cold.png", "cold heatmap")

	l   = flag.Int("l", 0, "order-l")
	m   = flag.Int("m", 0, "order-m")
	xln = flag.Float64(
		"xln",
		0.0,
		"nth root of jl(x) or d(xjl(x))/dx, for TE and TM mode respectively",
	)
	mode = flag.String("mode", "", "has to start with TM|TE")

	outDir = flag.String("out-dir", "", "output dir")
)

const (
	width           = 1280
	height          = 1280
	orbitPeriods    = 24
	framesPerPeriod = 30

	maxAmp       = .3
	rSamples     = 50
	thetaSamples = 241
	phiSamples   = 240

	radius  = 3.0
	camDist = 1.2 * radius
	// Camera orbiting circle axis, which is in the x-y plane, forming this angle with the x-axis.
	cameraOrbitAxisAngle   = math.Pi * 17.0 / 180.0
	cameraOrbitAngleOffset = -math.Pi * 7.0 / 180.0
	tmMode                 = "TM"
	teMode                 = "TE"

	minFading   = 0.2
	maxFading   = 1
	fadingGamma = 0.5
)

// VSH value at theta for all phis.
type vsh struct {
	re []*graphix.Vec3
	im []*graphix.Vec3
}

// Convert the vector in local (r,theta,phi) frame to Cartesian frame, at point (theta,phi).
func sphericalToCartesian(v *graphix.Vec3, theta, phi float64) *graphix.Vec3 {
	ct, st := math.Cos(theta), math.Sin(theta)
	cp, sp := math.Cos(phi), math.Sin(phi)
	return graphix.NewVec3(
		v[0]*st*cp+v[1]*ct*cp-v[2]*sp,
		v[0]*st*sp+v[1]*ct*sp+v[2]*cp,
		v[0]*ct-v[1]*st,
	)
}

func newVSH(
	theta float64,
	phis []float64,
	evaluator func(float64, []float64) ([]*graphix.Vec3, []*graphix.Vec3),
) *vsh {
	// re, im are in local spherical frame.
	re, im := evaluator(theta, phis)
	ret := &vsh{
		re: make([]*graphix.Vec3, len(phis)),
		im: make([]*graphix.Vec3, len(phis)),
	}
	for i, phi := range phis {
		ret.re[i] = sphericalToCartesian(re[i], theta, phi)
		ret.im[i] = sphericalToCartesian(im[i], theta, phi)
	}
	return ret
}

type radVal struct {
	jl, jlDer float64
}

func main() {
	flag.Parse()

	if *l <= 0 {
		panic("l must be greater than 0")
	}

	dtheta := math.Pi / float64(thetaSamples-1)
	thetas := make([]float64, thetaSamples)
	// North and south pole.
	thetas[0] = 0
	thetas[thetaSamples-1] = math.Pi
	for t := 1; t < thetaSamples-1; t++ {
		thetas[t] = float64(t) * dtheta
	}

	dphi := 2.0 * math.Pi / float64(phiSamples)
	phis := make([]float64, phiSamples)
	for p := 0; p < phiSamples; p++ {
		phis[p] = float64(p) * dphi
	}

	// Prepare VSH values for all angles.
	sp := newsph(*l, *m)
	fillVSH := func(vshs []*vsh, evaluator func(float64, []float64) ([]*graphix.Vec3, []*graphix.Vec3)) {
		// North and south pole.
		vshs[0] = newVSH(0, []float64{0}, evaluator)
		vshs[thetaSamples-1] = newVSH(math.Pi, []float64{0}, evaluator)

		for t := 1; t < thetaSamples-1; t++ {
			vshs[t] = newVSH(thetas[t], phis, evaluator)
		}
	}

	vshY := make([]*vsh, thetaSamples)
	fillVSH(vshY, sp.evalY)
	vshPhi := make([]*vsh, thetaSamples)
	fillVSH(vshPhi, sp.evalPhi)
	vshPsi := make([]*vsh, thetaSamples)
	fillVSH(vshPsi, sp.evalPsi)

	radVals := make([]*radVal, rSamples)
	dxln := *xln / float64(rSamples)
	for i := 0; i < rSamples; i++ {
		jl, jlDer := sphericalBessel(*l, dxln*float64(i+1))
		radVals[i] = &radVal{jl: jl, jlDer: jlDer}
	}
	// We know from outset that center of spherer has zero field always.
	gridCnt := rSamples * ((thetaSamples-2)*phiSamples + 2)

	grids := make([]*graphix.Vec3, gridCnt)
	dr := radius / float64(rSamples)

	efields := make([]*graphix.Vec3, gridCnt*framesPerPeriod)
	hfields := make([]*graphix.Vec3, gridCnt*framesPerPeriod)
	idx := 0

	// First record the spherical coordinates into grids [rho,phi,z].
	for nr := 1; nr <= rSamples; nr++ {
		r := dr * float64(nr)
		// North pole.
		grids[idx] = graphix.NewVec3(r, 0, 0)
		idx++
		// Latitudes.
		for nt := 1; nt < thetaSamples-1; nt++ {
			// Longitudes.
			for nphi := 0; nphi < phiSamples; nphi++ {
				grids[idx] = graphix.NewVec3(r, thetas[nt], phis[nphi])
				idx++
			}
		}
		// South pole.
		grids[idx] = graphix.NewVec3(r, math.Pi, 0)
		idx++
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
			for i := range grids {
				ef[i], hf[i] = field(grids, i, omegat, vshY, vshPhi, vshPsi, radVals)
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

	// Convert grids to cartesian.
	for _, grid := range grids {
		r, theta, phi := grid[0], grid[1], grid[2]
		ct, st := math.Cos(theta), math.Sin(theta)
		grid[0], grid[1], grid[2] = r*st*math.Cos(phi), r*st*math.Sin(phi), r*ct
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
	pos := graphix.NewVec3(0, 0, -camDist)
	forward := graphix.NewVec3(0, 0, 1)
	up := graphix.NewVec3(0, 1, 0)
	pr := graphix.NewOrthographic()
	sc := graphix.NewScreen(width, height, -1.2*radius, -1.2*radius, 1.2*radius, 1.2*radius)
	cir := graphix.NewCircularCameraOrbit(
		n,
		pos,
		forward,
		up,
		frameCnt,
		cameraOrbitAngleOffset,
		pr,
		sc,
	)

	for f := 0; f < frameCnt; f++ {
		paths := make([]*zraster.SpacePath, 0, 2*gridCnt)
		idxStart := (f % framesPerPeriod) * gridCnt
		ef, hf := efields[idxStart:], hfields[idxStart:]
		for i := 0; i < gridCnt; i++ {
			// Look up the color of e,h from heatmap based on original theta.
			t := math.Acos(grids[i][2]/grids[i].Norm()) / math.Pi
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

func indices(i int) (int, int, int) {
	// Index within a shell.
	rIdx := i / ((thetaSamples-2)*phiSamples + 2)
	i = i % ((thetaSamples-2)*phiSamples + 2)
	// The first and last grid of the shell is the north and south pole.
	if i == 0 {
		return 0, 0, rIdx
	}
	if i == (thetaSamples-2)*phiSamples+1 {
		return thetaSamples - 1, 0, rIdx
	}
	i -= 1
	return i/phiSamples + 1, i % phiSamples, rIdx
}

func zmul(re, im *graphix.Vec3, xi float64) *graphix.Vec3 {
	c, s := math.Cos(xi), math.Sin(xi)
	return graphix.NewVec3(
		re[0]*c+im[0]*s,
		re[1]*c+im[1]*s,
		re[2]*c+im[2]*s,
	)
}

// Returns E,H field at (rho,phi,z,omegat).
func field(
	grids []*graphix.Vec3,
	i int,
	omegat float64,
	vshY, vshPhi, vshPsi []*vsh,
	radVals []*radVal,
) (*graphix.Vec3, *graphix.Vec3) {
	thetaIdx, phiIdx, rIdx := indices(i)
	r := grids[i][0]
	k := *xln / radius
	x := k * r

	rey, imy := vshY[thetaIdx].re[phiIdx], vshY[thetaIdx].im[phiIdx]
	rephi, imphi := vshPhi[thetaIdx].re[phiIdx], vshPhi[thetaIdx].im[phiIdx]
	repsi, impsi := vshPsi[thetaIdx].re[phiIdx], vshPsi[thetaIdx].im[phiIdx]

	jl, jlDer := radVals[rIdx].jl, radVals[rIdx].jlDer

	v1 := zmul(rephi, imphi, omegat+math.Pi/2)
	v1.Scale(v1, jl)

	v2 := zmul(rey, imy, omegat)
	v2.Scale(v2, float64(*l*(*l+1))*jl/x)

	v3 := zmul(repsi, impsi, omegat)
	v3.Scale(v3, jlDer/x)

	v2.Add(v2, v3)

	if strings.HasPrefix(*mode, teMode) {
		return v1, v2
	} else if strings.HasPrefix(*mode, tmMode) {
		return v2.Scale(v2, -1), v1
	}
	panic(fmt.Sprintf("unsupported mode: %v", *mode))
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
