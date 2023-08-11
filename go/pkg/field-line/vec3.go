package fieldline

import "math"

// Common operations for 3 vectors.
type Vec3 [3]float64

func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		v[0] + u[0],
		v[1] + u[1],
		v[2] + u[2],
	}
}

func (v Vec3) Subtract(u Vec3) Vec3 {
	return Vec3{
		v[0] - u[0],
		v[1] - u[1],
		v[2] - u[2],
	}
}

func (v Vec3) Scale(s float64) Vec3 {
	return Vec3{v[0] * s, v[1] * s, v[2] * s}
}

func (v Vec3) Dot(u Vec3) float64 {
	return v[0]*u[0] + v[1]*u[1] + v[2]*u[2]
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		v[1]*u[2] - u[1]*v[2],
		v[2]*u[0] - u[2]*v[0],
		v[0]*u[1] - u[0]*v[1],
	}
}

func (v Vec3) Normalize() Vec3 {
	return v.Scale(1.0 / math.Sqrt(v.Dot(v)))
}
