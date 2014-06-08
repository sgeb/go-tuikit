package tuikit

import "strconv"

//----------------------------------------------------------------------------
// Point
//----------------------------------------------------------------------------

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y int
}

// PointZero is the zero Point.
var PointZero Point

// NewPoint is shorthand for Point{X, Y}.
func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

// String returns a string representation of p like "(x3,y4)".
func (p Point) String() string {
	return "(x" + strconv.Itoa(p.X) + ",y" + strconv.Itoa(p.Y) + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point) Mul(k int) Point {
	return Point{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k int) Point {
	return Point{p.X / k, p.Y / k}
}

// Eq reports whether p and q are equal.
func (p Point) Eq(q Point) bool {
	return p.X == q.X && p.Y == q.Y
}

// In reports whether p is in r.
func (p Point) In(r Rect) bool {
	return r.X <= p.X && p.X < (r.X+r.Width) &&
		r.Y <= p.Y && p.Y < (r.Y+r.Height)
}

//----------------------------------------------------------------------------
// Size
//----------------------------------------------------------------------------

// A Size is a Width, Height pair.
type Size struct {
	Width, Height int
}

// SizeZero is the zero Size.
var SizeZero Size

// NewSize is shorthand for Size{Width, Height}.
func NewSize(w, h int) Size {
	return Size{w, h}
}

// String returns a string representation of s like "(w3,h4)".
func (s Size) String() string {
	return "(w" + strconv.Itoa(s.Width) + ",h" + strconv.Itoa(s.Height) + ")"
}

// Eq reports whether s and t are equal.
func (s Size) Eq(t Size) bool {
	return s.Width == t.Width && s.Height == t.Height
}

// Empty reports whether the size contains no points.
func (s Size) Empty() bool {
	return s.Width <= 0 || s.Height <= 0
}

//----------------------------------------------------------------------------
// Rect
//----------------------------------------------------------------------------

// A Rect is the composition of an origin Point and a Size.
type Rect struct {
	Point
	Size
}

// RectZero is the zero Rect.
var RectZero Rect

// NewRect is shorthand for Rect{Point{X, Y}, Size{Width, Height}}.
func NewRect(x, y, w, h int) Rect {
	return Rect{NewPoint(x, y), NewSize(w, h)}
}

// String returns a string representation of r like "(x3,y4)-(w6,h5)".
func (r Rect) String() string {
	return r.Point.String() + "-" + r.Size.String()
}

// Eq reports whether r and s are equal.
func (r Rect) Eq(s Rect) bool {
	return r.Point.Eq(s.Point) && r.Size.Eq(s.Size)
}

// Max returns the most bottom-right Point still inside r.
func (r Rect) Max() Point {
	if r.Empty() {
		return r.Point
	}

	return NewPoint(r.X+r.Width-1, r.Y+r.Height-1)
}

// Add returns the rectangle r translated by p.
func (r Rect) Add(p Point) Rect {
	return NewRect(r.X+p.X, r.Y+p.Y, r.Width, r.Height)
}

// Sub returns the rectangle r translated by -p.
func (r Rect) Sub(p Point) Rect {
	return NewRect(r.X-p.X, r.Y-p.Y, r.Width, r.Height)
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rect) Inset(n int) Rect {
	if r.Width < 2*n {
		r.X = r.X + r.Width/2
		r.Width = 0
	} else {
		r.X += n
		r.Width -= 2 * n
	}

	if r.Height < 2*n {
		r.Y = r.Y + r.Height/2
		r.Height = 0
	} else {
		r.Y += n
		r.Height -= 2 * n
	}

	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rect) Intersect(s Rect) Rect {
	// Calculate Max() before any changes
	rm := r.Max()
	sm := s.Max()

	if r.X < s.X {
		r.X = s.X
	}

	if r.Y < s.Y {
		r.Y = s.Y
	}

	if rm.X > sm.X {
		r.Width = sm.X - r.X + 1
	}

	if rm.Y > sm.Y {
		r.Height = sm.Y - r.Y + 1
	}

	if !r.Valid() {
		return RectZero
	}

	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rect) Union(s Rect) Rect {
	// Calculate Max() before any changes
	rm := r.Max()
	sm := s.Max()

	if r.X > s.X {
		r.X = s.X
	}

	if r.Y > s.Y {
		r.Y = s.Y
	}

	if rm.X < sm.X {
		r.Width = sm.X - r.X + 1
	}

	if rm.Y < sm.Y {
		r.Height = sm.Y - r.Y + 1
	}

	return r
}

// Valid reports whether r is a valid Rect.
func (r Rect) Valid() bool {
	return r.Width >= 0 && r.Height >= 0
}

// Empty reports whether the rectangle contains no points.
func (r Rect) Empty() bool {
	return r.Size.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rect) Overlaps(s Rect) bool {
	rm := r.Max()
	sm := s.Max()
	return r.X <= sm.X && s.X <= rm.X && r.Y <= sm.Y && s.Y <= rm.Y
}

// In reports whether every point in r is in s.
func (r Rect) In(s Rect) bool {
	if r.Empty() {
		return true
	}

	rm := r.Max()
	sm := s.Max()
	return s.X <= r.X && rm.X <= sm.X && s.Y <= r.Y && rm.Y <= sm.Y
}
