package utils

import "math"

type Vector2D struct {
	X float64
	Y float64
}

func (v *Vector2D) Add(v2 Vector2D) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector2D) Subtract(v2 Vector2D) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vector2D) Limit(max float64) {
	magSq := v.MagnitudeSquared()
	if magSq > max*max {
		v.Divide(math.Sqrt(magSq))
		v.Multiply(max)
	}
}

func (v *Vector2D) Normalize() {
	mag := math.Sqrt(v.X*v.X + v.Y*v.Y)
	v.X /= mag
	v.Y /= mag
}

func (v *Vector2D) SetMagnitude(z float64) {
	v.Normalize()
	v.X *= z
	v.Y *= z
}

func (v *Vector2D) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2D) Divide(z float64) {
	v.X /= z
	v.Y /= z
}

func (v *Vector2D) Multiply(z float64) {
	v.X *= z
	v.Y *= z
}

func (v Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v2.X-v.X, 2) + math.Pow(v2.Y-v.Y, 2))
}

func Rotate(v Vector2D, ang int) Vector2D {
	aR := AngleToRadians(ang)
	oldX := v.X
	oldY := v.Y
	v.X = oldX*math.Cos(aR) - oldY*math.Sin(aR)
	v.Y = oldX*math.Sin(aR) + oldY*math.Cos(aR)
	return v
}

func Sign(p1 Vector2D, p2 Vector2D, p3 Vector2D) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

//Check si un point est un triangle
func PointInTriangle(pt Vector2D, v1 Vector2D, v2 Vector2D, v3 Vector2D) bool {

	d1 := Sign(pt, v1, v2)
	d2 := Sign(pt, v2, v3)
	d3 := Sign(pt, v3, v1)

	has_neg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	has_pos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(has_neg && has_pos)
}

// Convertie un angle en radian
func AngleToRadians(angle int) float64 {
	return (math.Pi / 180) * float64(angle)
}
