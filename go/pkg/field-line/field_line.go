package fieldline

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/fogleman/gg"
)

type pixel [2]int

func (p pixel) inBound(w, h int) bool {
	return p[0] >= 0 && p[0] < w && p[1] >= 0 && p[1] < h
}

type trajPoint struct {
	tangentLength float64
	pixel
}

type Symmetry struct {
	transform func(p Vec3) Vec3
	color     [3]float64
}

func newSymmetry(transform func(Vec3) Vec3, color [3]float64) *Symmetry {
	return &Symmetry{
		transform: transform,
		color:     color,
	}
}

type Trajectory struct {
	Start      Vec3
	AtEnd      func(p, v Vec3) bool
	Color      [3]float64
	symmetries []*Symmetry

	// One slice per frame x symmetry.
	points [][]trajPoint

	// Max and min of tangent length along the trajectory.
	maxTanLen float64
	minTanLen float64
}

func (traj *Trajectory) AddSymmetry(transform func(Vec3) Vec3, color [3]float64) {
	traj.symmetries = append(traj.symmetries, newSymmetry(transform, color))
}

type Options struct {
	OutputFile string
	Width      int
	Height     int
	Step       float64
	// The tangent vector at position represented by the argument.
	TangentAt   func(Vec3) Vec3
	LineWidth   float64
	FadingGamma float64
	*CameraOrbit
}

// Runs the field line renderer given the options and trajectory settings. Upon completion
// trajs internal data structure would have been modified.
func Run(opts Options, trajs []Trajectory) {
	if opts.CameraOrbit == nil {
		opts.CameraOrbit = NewCameraOrbit(0.0, 1)
	}
	var wgTrace sync.WaitGroup
	wgTrace.Add(len(trajs))
	h := opts.Step
	h2 := opts.Step * 0.5
	h6 := opts.Step / 6.0
	h3 := opts.Step / 3.0
	for i := range trajs {
		go func(traj *Trajectory) {
			defer wgTrace.Done()
			traj.maxTanLen = -1.0
			traj.minTanLen = -1.0
			identity := newSymmetry(func(p Vec3) Vec3 { return p }, traj.Color)
			traj.symmetries = append([]*Symmetry{identity}, traj.symmetries...)
			traj.points = make([][]trajPoint, len(opts.cameras)*len(traj.symmetries))
			// See multi variable Runge Kutta-4 at https://www.myphysicslab.com/explain/runge-kutta-en.html
			x := traj.Start
			for {
				a := opts.TangentAt(x)
				if traj.AtEnd(x, a) {
					return
				}
				tanlen := a.Norm()
				allOutOfBound := true
				for c, camera := range opts.CameraOrbit.cameras {
					for s, sym := range traj.symmetries {
						idx := c*len(traj.symmetries) + s
						tx := sym.transform(x)
						pt := trajPoint{
							tangentLength: tanlen,
							pixel:         camera.worldToScreen(tx, opts.Width, opts.Height),
						}
						if pt.inBound(opts.Width, opts.Height) {
							allOutOfBound = false
						}
						newPixel := true
						// Include a new point only if it's beyond a pixel away.
						if len(traj.points[idx]) > 0 {
							last := traj.points[idx][len(traj.points[idx])-1]
							if last.pixel[0] == pt.pixel[0] && last.pixel[1] == pt.pixel[1] {
								newPixel = false
							}
						}
						if newPixel {
							traj.points[idx] = append(traj.points[idx], pt)
						}
					}
				}

				if allOutOfBound {
					return
				}

				if traj.maxTanLen < 0.0 || traj.maxTanLen < tanlen {
					traj.maxTanLen = tanlen
				}
				if traj.minTanLen < 0.0 || traj.minTanLen > tanlen {
					traj.minTanLen = tanlen
				}
				xb := x.Add(a.Scale(h2))
				b := opts.TangentAt(xb)
				xc := x.Add(b.Scale(h2))
				c := opts.TangentAt(xc)
				xd := x.Add(c.Scale(h))
				d := opts.TangentAt(xd)
				x = x.Add(a.Scale(h6))
				x = x.Add(b.Scale(h3))
				x = x.Add(c.Scale(h3))
				x = x.Add(d.Scale(h6))
			}
		}(&trajs[i])
	}

	wgTrace.Wait()

	fmt.Println("completed tracing all trajectories")

	// Calculate max and min of tangent lengths.
	max, min := -1.0, -1.0
	for _, traj := range trajs {
		if max < 0 || max < traj.maxTanLen {
			max = traj.maxTanLen
		}
		if min < 0 || min > traj.minTanLen {
			min = traj.minTanLen
		}
	}
	// Degenerate case.
	if max < 0.0 && min < 0.0 {
		max, min = 1.0, 0.0
	}

	var wgRender sync.WaitGroup
	wgRender.Add(len(opts.cameras))
	for c := range opts.cameras {
		go func(cc int) {
			defer wgRender.Done()

			dc := gg.NewContext(opts.Width, opts.Height)
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			dc.SetLineWidth(opts.LineWidth)

			for _, traj := range trajs {
				for j := range traj.symmetries {
					points := traj.points[cc*len(traj.symmetries)+j]
					start := 0
					for {
						// Search for the next in-bound pixel.
						for start < len(points) {
							if points[start].inBound(opts.Width, opts.Height) {
								break
							}
							start++
						}
						if start > len(points) {
							// No more pixels along this trajectory.
							break
						}
						// Search for the next out-of-bound pixel.
						end := start + 1
						for end < len(points) {
							if !points[end].inBound(opts.Width, opts.Height) {
								break
							}
							end++
						}
						// end-1 because we need at least two points to draw a line.
						for p := start; p < end-1; p++ {
							// Determine the alpha of this segment based on the ratio of average tangent length to the max tangent length.
							avg := (points[p].tangentLength + points[p+1].tangentLength) / 2.0
							alpha := math.Pow((avg-min)/(max-min), opts.FadingGamma)
							dc.SetRGBA(traj.symmetries[j].color[0], traj.symmetries[j].color[1], traj.symmetries[j].color[2], alpha)
							dc.DrawLine(
								float64(points[p].pixel[0]),
								float64(points[p].pixel[1]),
								float64(points[p+1].pixel[0]),
								float64(points[p+1].pixel[1]),
							)
							dc.Stroke()
						}
						start = end + 1
					}
				}
			}

			filename := fmt.Sprintf("%v-%03d-of-%03d.png", opts.OutputFile, cc, len(opts.cameras))
			if err := dc.SavePNG(filename); err != nil {
				log.Fatalf("failed to save %v: %v", filename, err)
			}
		}(c)
	}
	wgRender.Wait()
}
