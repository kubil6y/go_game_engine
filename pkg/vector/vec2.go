// source: https://github.com/joonazan/vec2 (modified)
package vector

import "math"

type Vec2 struct {
	X, Y float32
}

func NewVec2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

func NewZeroVec2() Vec2 {
	return Vec2{X: 0, Y: 0}
}

func (self *Vec2) Add(other Vec2) {
	self.X += other.X
	self.Y += other.Y
}

func (self *Vec2) Sub(other Vec2) {
	self.X -= other.X
	self.Y -= other.Y
}

func (self *Vec2) Mul(other float32) {
	self.X *= other
	self.Y *= other
}

func (self Vec2) Dot(other Vec2) float32 {
	return self.X*other.X + self.Y*other.Y
}

func (self Vec2) Cross(other Vec2) float32 {
	return self.X*other.Y - self.Y*other.X
}

func (self Vec2) Crossf(other float32) Vec2 {
	return Vec2{-self.Y * other, self.X * other}
}

func (self Vec2) LengthSquared() float32 {
	return self.X*self.X + self.Y*self.Y
}

func (self Vec2) Length() float32 {
	return float32(math.Sqrt(float64(self.LengthSquared())))
}

func (self *Vec2) Normalize() {
	self.Mul(1 / self.Length())
}

func (self Vec2) Normalized() Vec2 {
	return Mul(self, 1/self.Length())
}

func (v Vec2) Plus(v2 Vec2) Vec2 {
	return Add(v, v2)
}

func (v Vec2) Minus(v2 Vec2) Vec2 {
	return Sub(v, v2)
}

func (v Vec2) Times(r float32) Vec2 {
	return Mul(v, r)
}

func Add(v, u Vec2) Vec2 {
	return Vec2{v.X + u.X, v.Y + u.Y}
}

func Sub(v, u Vec2) Vec2 {
	return Vec2{v.X - u.X, v.Y - u.Y}
}

func Mul(v Vec2, r float32) Vec2 {
	return Vec2{v.X * r, v.Y * r}
}
