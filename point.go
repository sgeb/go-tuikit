package tuikit

import log "github.com/sgeb/go-sglog"

type Point struct {
	X, Y int
}

var PointZero Point
var PointHidden Point = Point{-1, -1}

func NewPoint(x, y int) Point {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return Point{
		X: x,
		Y: y,
	}
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point) Mul(k int) Point {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return Point{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k int) Point {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return Point{p.X / k, p.Y / k}
}

// Eq reports whether p and q are equal.
func (p Point) Eq(q Point) bool {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return p.X == q.X && p.Y == q.Y
}

// Hidden reports whether p is hidden.
func (p Point) Hidden() bool {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return p.X < 0 || p.Y < 0
}
