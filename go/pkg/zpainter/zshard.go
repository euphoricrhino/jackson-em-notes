package zpainter

import (
	"image/color"

	"github.com/golang/freetype/raster"
)

type zShard struct {
	parent     *ZPainter
	r, g, b, a uint32
	p1         [3]float64
	p2         [3]float64
	dd         float64

	zbuf [][]*zColor
}

// Sets the current stroke color.
func (zs *zShard) SetColor(c color.Color) {
	zs.r, zs.g, zs.b, zs.a = c.RGBA()
}

// Sets the end points' coordinate for the current stroke.
func (zs *zShard) SetEndpoints(p1, p2 [3]float64) {
	zs.p1, zs.p2 = p1, p2
	dx, dy := p1[0]-p2[0], p1[1]-p2[1]
	zs.dd = dx*dx + dy*dy
}

// Paint only updates the zbuffer without actually painting any pixels on the image.
// ZPainter.Commit will sort the zbuffer and apply pixels in the z-order.
func (zs *zShard) Paint(ss []raster.Span, done bool) {
	rect := zs.parent.img.Bounds()
	for _, s := range ss {
		if s.Y < rect.Min.Y {
			continue
		}
		if s.Y >= rect.Max.Y {
			return
		}
		if s.X0 < rect.Min.X {
			s.X0 = rect.Min.X
		}
		if s.X1 > rect.Max.X {
			s.X1 = rect.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}
		r, g, b, a := zs.r*s.Alpha, zs.g*s.Alpha, zs.b*s.Alpha, zs.a*s.Alpha
		for x := s.X0; x < s.X1; x++ {
			zs.updateZBuf(x, s.Y, r, g, b, a)
		}
	}
}

func (zs *zShard) updateZBuf(x, y int, r, g, b, a uint32) {
	// Compute the z depth of the pixel by linearly interpolating between the two end points.
	// Due to rasterizing, (x,y) may not be on the line connecting the two end points.
	// The projection point from (x,y) to this line is used to interpolate z depth.
	z := zs.p1[2]
	// Degenerate case.
	if zs.dd == 0 {
		if z < zs.p2[2] {
			z = zs.p2[2]
		}
	} else {
		dx1, dy1 := float64(x)-zs.p1[0], float64(y)-zs.p1[1]
		dd1 := dx1*dx1 + dy1*dy1
		dx2, dy2 := float64(x)-zs.p2[0], float64(y)-zs.p2[1]
		dd2 := dx2*dx2 + dy2*dy2
		t := ((dd1-dd2)/zs.dd + 1.0) * 0.5
		z = zs.p1[2] + t*(zs.p2[2]-zs.p1[2])
	}

	i := (y-zs.parent.img.Rect.Min.Y)*zs.parent.width + (x - zs.parent.img.Rect.Min.X)
	zs.zbuf[i] = append(zs.zbuf[i], &zColor{
		r: r,
		g: g,
		b: b,
		a: a,
		z: z,
	})
}
