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

type Trajectory struct {
	Start Vec3
	AtEnd func(p, v Vec3) bool
	Color [3]float64

	points []orbitPoint
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

func Run(opts Options, trajs []Trajectory) {
	if opts.Camera == nil {
		opts.Camera = NewCamera(Vec3{0, 0, 2}, Vec3{1, 0, 0})
	}
	var wg sync.WaitGroup
	wg.Add(len(trajs))
	h := opts.Step
	h2 := opts.Step * 0.5
	h6 := opts.Step / 6.0
	h3 := opts.Step / 3.0
	for i := range trajs {
		go func(traj *Trajectory) {
			traj.points = nil
			defer wg.Done()
			// See multi variable Runge Kutta-4 at https://www.myphysicslab.com/explain/runge-kutta-en.html
			x := traj.Start
			for {
				a := opts.TangentAt(x)
				op := orbitPoint{
					tangentLength: math.Sqrt(a.Dot(a)),
					pixel:         opts.worldToScreen(x, opts.Width, opts.Height),
				}
				if !op.inBound(opts.Width, opts.Height) || traj.AtEnd(x, a) {
					return
				}
				newPoint := true
				// Include a new point only if it's beyond a pixel away.
				if len(traj.points) > 0 {
					last := traj.points[len(traj.points)-1]
					if last.pixel[0] == op.pixel[0] && last.pixel[1] == op.pixel[1] {
						newPoint = false
					}
				}
				if newPoint {
					traj.points = append(traj.points, op)
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
			if max < 0 || max < traj.points[i].tangentLength {
				max = traj.points[i].tangentLength
			}
			if min < 0 || min > traj.points[i].tangentLength {
				min = traj.points[i].tangentLength
			}
		}
	}

	dc := gg.NewContext(opts.Width, opts.Height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	dc.SetLineWidth(opts.LineWidth)
	for _, traj := range trajs {
		if len(traj.points) == 0 {
			continue
		}
		for i := 0; i < len(traj.points)-1; i++ {
			// Determine the alpha of this segment based on the ratio of average tangent length to the max tangent length.
			avg := (traj.points[i].tangentLength + traj.points[i+1].tangentLength) / 2.0
			alpha := math.Pow((avg-min)/(max-min), opts.FadingGamma)
			dc.SetRGBA(traj.Color[0], traj.Color[1], traj.Color[2], alpha)
			dc.DrawLine(
				float64(traj.points[i].pixel[0]),
				float64(traj.points[i].pixel[1]),
				float64(traj.points[i+1].pixel[0]),
				float64(traj.points[i+1].pixel[1]),
			)
			dc.Stroke()
		}
	}
	if err := dc.SavePNG(opts.OutputFile); err != nil {
		log.Fatalf("failed to save %v: %v", opts.OutputFile, err)
	}
}
