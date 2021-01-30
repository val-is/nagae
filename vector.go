package nagae

import "math"

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Translate(delta Vec2) {
	v.X += delta.X
	v.Y += delta.Y
}

func (v *Vec2) MultScalar(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2) MultVec(other Vec2) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v Vec2) Hypot() float64    { return math.Hypot(v.X, v.Y) }
func (v Vec2) Angle() float64    { return math.Atan2(v.Y, v.X) }
func (v Vec2) AngleDeg() float64 { return v.Angle() * 180 / math.Pi }

func (v *Vec2) Rotate(angle float64) {
	dist := v.Hypot()
	curAng := v.Angle()
	v.X = dist * math.Cos(curAng+angle)
	v.Y = dist * math.Sin(curAng+angle)
}
func (v *Vec2) RotateDeg(angle float64) { v.Rotate(math.Pi * angle / 180) }
