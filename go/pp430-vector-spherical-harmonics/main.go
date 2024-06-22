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
	"sync"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/graphix/zraster"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	l         = flag.Int("l", 0, "l")
	m         = flag.Int("m", 0, "m")
	heatmap   = flag.String("heatmap", "", "")
	gamma     = flag.Float64("gamma", 1, "gamma")
	outDir    = flag.String("out-dir", "", "")
	component = flag.String("component", "", "")
)

const (
	componentY   = "Y"
	componentPsi = "Psi"
	componentPhi = "Phi"
)

func main() {
	flag.Parse()

	const thetaSamples = 240
	const phiSamples = 240
	dtheta := math.Pi / thetaSamples
	dphi := math.Pi * 2 / phiSamples

	hm, _ := graphix.LoadHeatmap(*heatmap, 1)
	sp := newsph(*l, *m)
	gridCnt := 2 + (thetaSamples-1)*phiSamples
	grids := make([]*graphix.Vec3, gridCnt)
	fields := make([]*graphix.Vec3, gridCnt)

	var evaluator func(float64, []float64) ([]*graphix.Vec3, []*graphix.Vec3)

	switch *component {
	case componentY:
		evaluator = sp.evalY
	case componentPsi:
		evaluator = sp.evalPsi
	case componentPhi:
		evaluator = sp.evalPhi
	default:
		panic(fmt.Sprintf("unknown component: %v", *component))
	}
	maxf, minf := math.Inf(-1), math.Inf(1)
	// North and south pole.
	grids[0] = graphix.NewVec3(1, 0, 0)
	re, _ := evaluator(0, []float64{0})
	fields[0] = re[0]
	maxf = max(maxf, re[0].Norm())
	minf = min(minf, re[0].Norm())

	grids[1] = graphix.NewVec3(1, math.Pi, 0)
	re, _ = evaluator(math.Pi, []float64{0})
	fields[1] = re[0]
	maxf = max(maxf, re[0].Norm())
	minf = min(minf, re[0].Norm())

	phis := make([]float64, phiSamples)
	for p := 0; p < phiSamples; p++ {
		phis[p] = float64(p) * dphi
	}

	var wg sync.WaitGroup
	workers := runtime.NumCPU()
	wg.Add(workers)

	maxff := make([]float64, workers)
	minff := make([]float64, workers)
	for w := range maxff {
		maxff[w] = math.Inf(-1)
		minff[w] = math.Inf(1)
	}
	for w := 0; w < workers; w++ {
		go func(wk int) {
			for t := 1; t < thetaSamples; t++ {
				if t%workers != wk {
					continue
				}
				theta := dtheta * float64(t)
				re, _ := evaluator(theta, phis)
				for p, phi := range phis {
					uidx := (t-1)*phiSamples + p + 2
					grids[uidx] = graphix.NewVec3(1, theta, phi)
					fields[uidx] = re[p]
					maxff[wk] = max(maxff[wk], re[p].Norm())
					minff[wk] = min(minff[wk], re[p].Norm())
				}
			}
			wg.Done()
		}(w)
	}

	wg.Wait()
	for w := range maxff {
		maxf = max(maxf, maxff[w])
		minf = min(minf, minff[w])
	}

	const framesPerPeriod = 30

	const maxAmp = 0.05
	const minAmp = 0.02
	paths := make([][]*zraster.SpacePath, framesPerPeriod)
	for f := range paths {
		paths[f] = make([]*zraster.SpacePath, gridCnt+3)
		amp := (minAmp+maxAmp)/2 + (maxAmp-minAmp)/2*math.Cos(float64(f)*math.Pi*2/framesPerPeriod)
		const minFading = 0.0
		const maxFading = 1.0
		for i := 0; i < gridCnt; i++ {
			absf := fields[i].Norm()
			rr, theta, phi := grids[i][0], grids[i][1], grids[i][2]
			ct, st := math.Cos(theta), math.Sin(theta)
			cp, sp := math.Cos(phi), math.Sin(phi)
			x, y, z := rr*st*cp, rr*st*sp, rr*ct
			offset := graphix.NewVec3(
				fields[i][0]*st*cp+fields[i][1]*ct*cp-fields[i][2]*sp,
				fields[i][0]*st*sp+fields[i][1]*ct*sp+fields[i][2]*cp,
				fields[i][0]*ct-fields[i][1]*st,
			)
			offset.Scale(offset, amp/maxf)
			lambda := (absf - minf) / (maxf - minf)
			zlambda := (rr*ct + 1) / 2
			r, g, b, _ := hm[int(zlambda*float64(len(hm)-1))].RGBA()

			c := color.NRGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(255 * (minFading + (maxFading-minFading)*math.Pow(lambda, *gamma))),
			}
			paths[f][i] = &zraster.SpacePath{
				Segments: []*zraster.SpaceVertex{{
					Pos:   graphix.NewVec3(x, y, z),
					Color: c,
				}},
				End:       graphix.NewVec3(x+offset[0], y+offset[1], z+offset[2]),
				LineWidth: 1,
			}
		}
		paths[f][gridCnt] = &zraster.SpacePath{
			Segments: []*zraster.SpaceVertex{{
				Pos:   graphix.NewVec3(0, 0, 0),
				Color: color.NRGBA{R: 0xff, A: 0xff},
			}},
			End:       graphix.NewVec3(1.2, 0, 0),
			LineWidth: 2,
		}
		paths[f][gridCnt+1] = &zraster.SpacePath{
			Segments: []*zraster.SpaceVertex{{
				Pos:   graphix.NewVec3(0, 0, 0),
				Color: color.NRGBA{G: 0xff, A: 0xff},
			}},
			End:       graphix.NewVec3(0, 1.2, 0),
			LineWidth: 2,
		}
		paths[f][gridCnt+2] = &zraster.SpacePath{
			Segments: []*zraster.SpaceVertex{{
				Pos:   graphix.NewVec3(0, 0, 0),
				Color: color.NRGBA{B: 0xff, A: 0xff},
			}},
			End:       graphix.NewVec3(0, 0, 1.2),
			LineWidth: 2,
		}
	}

	const orbitPeriods = 24
	frameCnt := orbitPeriods * framesPerPeriod

	cameraOrbitAxisAngle := math.Pi * 17.0 / 180.0
	cameraOrbitAngleOffset := -math.Pi * 7.0 / 180.0

	// Configuration for the camera orbit.
	n := graphix.NewVec3(math.Cos(cameraOrbitAxisAngle), math.Sin(cameraOrbitAxisAngle), 0)
	pos := graphix.NewVec3(0, 0, -2)
	forward := graphix.NewVec3(0, 0, 1)
	up := graphix.NewVec3(0, 1, 0)
	pr := graphix.NewPerspective(1.5)
	sc := graphix.NewScreen(1280, 1280, -1.2, -1.2, 1.2, 1.2)
	cir := graphix.NewCircularCameraOrbit(n, pos, forward, up, frameCnt, cameraOrbitAngleOffset, pr, sc)

	for f := 0; f < frameCnt; f++ {
		img := zraster.Run(zraster.Options{
			Camera:  cir.GetCamera(f),
			Paths:   paths[f%framesPerPeriod],
			Workers: runtime.NumCPU(),
		})
		renderCaption(img)
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

func renderCaption(img draw.Image) {
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetLineWidth(1)
	draw2d.SetFontFolder("/System/Library/Fonts")
	draw2d.SetFontNamer(func(_ draw2d.FontData) string { return "Monaco.ttf" })
	textColor := color.RGBA{0, 0xcc, 0xcc, 0xff}
	gc.SetFillColor(textColor)
	gc.SetStrokeColor(textColor)
	gc.SetDPI(288)
	gc.SetFontSize(5)
	text := ""
	switch *component {
	case componentY:
		text = fmt.Sprintf("Y (l=%v, m=%v)", *l, *m)
	case componentPsi:
		text = fmt.Sprintf("Ψ=r∇Y (l=%v, m=%v)", *l, *m)
	case componentPhi:
		text = fmt.Sprintf("Φ=r×∇Y (l=%v, m=%v)", *l, *m)
	}
	gc.FillStringAt(text, 40.0, 40.0)
}
