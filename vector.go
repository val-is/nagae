package nagae

type Vec2 struct {
	x, y float64
}

var ZeroVector = Vec2{0, 0}

func (v *Vec2) Translate(delta Vec2) {
	v.x += delta.x
	v.y += delta.y
}

func (v *Vec2) MultScalar(scalar float64) {
	v.x *= scalar
	v.y *= scalar
}
