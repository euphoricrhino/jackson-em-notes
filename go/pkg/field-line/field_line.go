package fieldline

import (
	"log"
	"math"
	"sync"

	"github.com/fogleman/gg"
)

type pixel [2]int

func (p pixel) inBound(w, h int) bool {
	return p[0] >= 0 && p[0] < w && p[1] >= 0 && p[1] < h
}

type orbitPoint struct {
	tangentLength float64
	pixel
}

type Symmetry struct {
	transform  func(p Vec3) Vec3
	color      [3]float64
	outOfBound bool
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

	// One per symmetry.
	points [][]orbitPoint
}

func (traj *Trajectory) AddSymmetry(transform func(Vec3) Vec3, color [3]float64) {
	traj.symmetries = append(traj.symmetries, newSymmetry(transform, color))
}

type Camera struct {
	pos Vec3
	// Normalized screen frame in world coordinate.
	sx Vec3
	sy Vec3
	sz Vec3
}

// Creates a camera at pos, looking at origin, with screenX (after proper orthogonalization) pointing to screen's right direction.
func NewCamera(pos, screenX Vec3) *Camera {
	// Vector from origin to camera.
	sz := pos.Normalize()
	// screenX's projection onto sz.
	xonz := sz.Scale(screenX.Dot(sz))
	// Orthogonalized screenX.
	sx := screenX.Subtract(xonz).Normalize()
	sy := sz.Cross(sx)
	return &Camera{
		pos: pos,
		sx:  sx,
		sy:  sy,
		sz:  sz,
	}
}

func (cam *Camera) worldToScreen(p Vec3, w, h int) pixel {
	q := p.Subtract(cam.pos)
	return pixel{
		w/2 + int(float64(w)*q.Dot(cam.sx)/2.0),
		h/2 - int(float64(h)*q.Dot(cam.sy)/2.0),
	}
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
	*Camera
}

// Runs the field line renderer given the options and trajectory settings. Upon completion
// trajs internal data structure would have been modified.
func Run(opts Options, trajs []Trajectory) {
	if opts.Camera == nil {
		opts.Camera = NewCamera(Vec3{0, 0, 1}, Vec3{1, 0, 0})
	}
	var wg sync.WaitGroup
	wg.Add(len(trajs))
	h := opts.Step
	h2 := opts.Step * 0.5
	h6 := opts.Step / 6.0
	h3 := opts.Step / 3.0
	for i := range trajs {
		go func(traj *Trajectory) {
			defer wg.Done()
			identity := newSymmetry(func(p Vec3) Vec3 { return p }, traj.Color)
			traj.symmetries = append([]*Symmetry{identity}, traj.symmetries...)
			traj.points = make([][]orbitPoint, len(traj.symmetries))
			// See multi variable Runge Kutta-4 at https://www.myphysicslab.com/explain/runge-kutta-en.html
			x := traj.Start
			for {
				a := opts.TangentAt(x)
				if traj.AtEnd(x, a) {
					return
				}
				tanlen := a.Norm()
				allOutOfBound := true
				for j, sym := range traj.symmetries {
					if sym.outOfBound {
						continue
					}
					tx := sym.transform(x)
					op := orbitPoint{
						tangentLength: tanlen,
						pixel:         opts.worldToScreen(tx, opts.Width, opts.Height),
					}
					if !op.inBound(opts.Width, opts.Height) {
						sym.outOfBound = true
						continue
					}
					allOutOfBound = false
					newPixel := true
					// Include a new point only if it's beyond a pixel away.
					if len(traj.points[j]) > 0 {
						last := traj.points[j][len(traj.points[j])-1]
						if last.pixel[0] == op.pixel[0] && last.pixel[1] == op.pixel[1] {
							newPixel = false
						}
					}
					if newPixel {
						traj.points[j] = append(traj.points[j], op)
					}
				}
				if allOutOfBound {
					return
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

	wg.Wait()

	// Calculate max and min of tangent lengths.
	max, min := -1.0, -1.0
	for _, traj := range trajs {
		for i := range traj.points {
			for j := range traj.points[i] {
				if max < 0 || max < traj.points[i][j].tangentLength {
					max = traj.points[i][j].tangentLength
				}
				if min < 0 || min > traj.points[i][j].tangentLength {
					min = traj.points[i][j].tangentLength
				}
			}
		}
	}

	dc := gg.NewContext(opts.Width, opts.Height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	dc.SetLineWidth(opts.LineWidth)
	for _, traj := range trajs {
		for i, pi := range traj.points {
			for j := 0; j < len(pi)-1; j++ {
				// Determine the alpha of this segment based on the ratio of average tangent length to the max tangent length.
				avg := (pi[j].tangentLength + pi[j+1].tangentLength) / 2.0
				alpha := math.Pow((avg-min)/(max-min), opts.FadingGamma)
				dc.SetRGBA(traj.symmetries[i].color[0], traj.symmetries[i].color[1], traj.symmetries[i].color[2], alpha)
				dc.DrawLine(
					float64(pi[j].pixel[0]),
					float64(pi[j].pixel[1]),
					float64(pi[j+1].pixel[0]),
					float64(pi[j+1].pixel[1]),
				)
				dc.Stroke()
			}
		}
	}
	if err := dc.SavePNG(opts.OutputFile); err != nil {
		log.Fatalf("failed to save %v: %v", opts.OutputFile, err)
	}
}
