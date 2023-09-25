package fieldline

import "math"

type camera struct {
	pos Vec3
	// Normalized screen frame in world coordinate.
	sx Vec3
	sy Vec3
	sz Vec3
}

func (cam *camera) worldToScreen(p Vec3, w, h int) pixel {
	q := p.Subtract(cam.pos)
	return pixel{
		w/2 + int(float64(w)*q.Dot(cam.sx)/2.0),
		h/2 - int(float64(h)*q.Dot(cam.sy)/2.0),
	}
}

// Camera orbit settings.
type CameraOrbit struct {
	// One camera per frame, equally distributing along the orbiting circle.
	cameras []*camera
}

// NewCameraOrbit sets up the camera orbit. The camera orbit is a circle whose normal vector is the rotated y axis by rollDegree around z axis.
// The orbit has 'frames' positions equally distributed around the circular orbit.
func NewCameraOrbit(rollDegree float64, frames int) *CameraOrbit {
	co := &CameraOrbit{
		cameras: make([]*camera, frames),
	}
	rollRad := rollDegree * math.Pi / 180.0
	ry := Vec3{-math.Sin(rollRad), math.Cos(rollRad), 0.0}
	rz := Vec3{0, 0, 1}
	rx := ry.Cross(rz)
	dtheta := math.Pi * 2.0 / float64(frames)
	for f := 0; f < frames; f++ {
		theta := float64(f) * dtheta
		ct, st := math.Cos(theta), math.Sin(theta)
		pos := rz.Scale(ct).Add(rx.Scale(st))
		sx := rx.Scale(ct).Add(rz.Scale(-st))
		co.cameras[f] = &camera{
			pos: pos,
			sx:  sx,
			sy:  pos.Cross(sx),
			sz:  pos,
		}
	}
	return co
}
