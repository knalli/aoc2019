package day18

import "fmt"

type Point struct {
	X, Y int
}

func (p Point) Equal(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}
func (p Point) East() Point {
	return Point{p.X + 1, p.Y}
}
func (p Point) West() Point {
	return Point{p.X - 1, p.Y}
}
func (p Point) North() Point {
	return Point{p.X, p.Y + 1}
}
func (p Point) South() Point {
	return Point{p.X, p.Y - 1}
}
func (p Point) ToString() string {
	return fmt.Sprintf("(%d/%d)", p.X, p.Y)
}

func (p Point) Adjacents() []Point {
	return []Point{p.North(), p.East(), p.South(), p.West()}
}
