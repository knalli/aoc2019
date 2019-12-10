package main

import "C"
import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"sort"
	"time"
)

const AocDay = 10
const AocDayName = "day10"
const AocDayTitle = "Day 10"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	run := func(title string, file string, showSteps bool) {
		lines, _ := dl.ReadFileToArray(file)
		grid := buildGrid(lines)
		counts := calculateGridFieldToAsteroidVisions(grid)
		bestVision := findBestVision(grid)
		visions := counts[bestVision]
		f := grid.data[bestVision.Y][bestVision.X]
		dl.PrintStepHeader(0)
		fmt.Printf("%s:\n%s\n", title, grid.ToStringColorized(bestVision, Point{-1, -1}, visions))
		dl.PrintStepHeader(1)
		dl.PrintSolution(fmt.Sprintf("Best coordinates of %s with %d other asteriods", bestVision.ToString(), f.Detected))
		destroyed := make([]Point, 0)
		vaporize(grid, bestVision, func(total int, p Point) {
			destroyed = append(destroyed, p)
			if showSteps {
				fmt.Printf("  💥 #%d (%d,%d)\n", total, p.X, p.Y)
				counts := calculateGridFieldToAsteroidVisions(grid)
				visions := counts[bestVision]
				fmt.Printf("%s:\n%s\n", title, grid.ToStringColorized(bestVision, p, visions))
				time.Sleep(100 * time.Millisecond)
			}
		})
		dl.PrintStepHeader(2)
		dl.PrintSolution(fmt.Sprintf("%d asteroids destroyed", len(destroyed)))
		if len(destroyed) >= 200 {
			n200 := destroyed[199]
			dl.PrintSolution(fmt.Sprintf("200th asteroid is %s, anwser = %d", n200.ToString(), n200.X*100+n200.Y))
		}
	}

	run("Sample 1", AocDayName+"/sample1.txt", false)
	run("Sample 5", AocDayName+"/sample5.txt", false)
	run("Sample 6", AocDayName+"/sample6.txt", false)
	run("Puzzle 1", AocDayName+"/puzzle1.txt", false)

}

type Field struct {
	Filled   bool
	Detected int
}

type Point struct {
	X, Y int
}

func (p *Point) Clone() Point {
	return Point{p.X, p.Y}
}

func (p *Point) Equals(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p *Point) Add(o Point) Point {
	return Point{p.X + o.X, p.Y + o.Y}
}

func (p *Point) Sub(o Point) Point {
	return Point{p.X - o.X, p.Y - o.Y}
}

func (p *Point) ToString() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

type Grid struct {
	data [][]Field
}

func (g *Grid) ToString(highlights ...Point) string {
	result := ""
	isHighlighted := func(x, y int) bool {
		for _, h := range highlights {
			if h.X == x && h.Y == y {
				return true
			}
		}
		return false
	}

	header := true
	if header {
		result += " "
		for x := 0; x < len(g.data[0]); x++ {
			result += fmt.Sprintf("%d", x%10)
		}
		result += "\n"
	}
	for y, line := range g.data {
		if y > 0 {
			result += "\n"
		}
		if header {
			result += fmt.Sprintf("%d", y%10)
		}
		for x, f := range line {
			if f.Filled {
				if isHighlighted(x, y) {
					result += Red("#").String()
				} else {
					result += "#"
				}
			} else {
				if isHighlighted(x, y) {
					result += Red(".").String()
				} else {
					result += "."
				}
			}
		}
	}
	return result
}

func (g *Grid) ToStringColorized(bestVision Point, destroy Point, visions []Point) string {
	result := ""
	highlight := func(x, y int, s string) Value {
		if bestVision.X == x && bestVision.Y == y {
			return Blue(s)
		}
		if destroy.X == x && destroy.Y == y {
			return Red(s)
		}
		for _, h := range visions {
			if h.X == x && h.Y == y {
				return Yellow(s)
			}
		}
		return Reset(s)
	}
	for y, line := range g.data {
		if y > 0 {
			result += "\n"
		}
		for x, f := range line {
			if f.Filled {
				result += highlight(x, y, " #").String()
			} else {
				result += highlight(x, y, " .").String()
			}
		}
	}
	return result
}

func (g *Grid) Each(f func(x int, y int, f *Field)) {
	for y := 0; y < len(g.data); y++ {
		for x := 0; x < len(g.data[y]); x++ {
			f(x, y, &g.data[y][x])
		}
	}
}

func (g *Grid) CountAsteroids() int {
	total := 0
	g.Each(func(_ int, _ int, f *Field) {
		if f.Filled {
			total++
		}
	})
	return total
}

func (g *Grid) Filled(x, y int) bool {
	return g.data[y][x].Filled
}

func (g *Grid) IsAnythingExactlyBetween(fromX, fromY, toX, toY int) bool {
	gcd := dl.GreatestCommonDivisor(toX-fromX, toY-fromY)
	stepX := (toX - fromX) / gcd
	stepY := (toY - fromY) / gcd
	walkX := fromX + stepX
	walkY := fromY + stepY
	for walkX != toX || walkY != toY {
		if g.Filled(walkX, walkY) {
			return true
		}
		walkX += stepX
		walkY += stepY
	}

	return false
}

func buildGrid(lines []string) *Grid {
	data := make([][]Field, len(lines))
	for y, line := range lines {
		fields := make([]Field, len(line))
		for x, c := range line {
			fields[x] = Field{Filled: c == '#', Detected: 0}
		}
		data[y] = fields
	}
	return &Grid{data}
}

func calculateGridFieldToAsteroidVisions(grid *Grid) map[Point][]Point {
	result := make(map[Point][]Point)
	grid.Each(func(fromX int, fromY int, fromField *Field) {
		fromField.Detected = 0
		if !fromField.Filled {
			return
		}
		subResult := make([]Point, 0)

		grid.Each(func(toX int, toY int, toField *Field) {
			if fromX == toX && fromY == toY {
				return
			}
			if !toField.Filled {
				return
			}

			if !grid.IsAnythingExactlyBetween(fromX, fromY, toX, toY) {
				fromField.Detected++
				subResult = append(subResult, Point{toX, toY})
			}
		})
		result[Point{fromX, fromY}] = subResult
	})
	return result
}

func findBestVision(grid *Grid) Point {
	max := -1
	maxFieldCoordinates := make([]int, 2)
	grid.Each(func(x int, y int, f *Field) {
		if max < f.Detected {
			max = f.Detected
			maxFieldCoordinates[0] = x
			maxFieldCoordinates[1] = y
		}
	})
	return Point{maxFieldCoordinates[0], maxFieldCoordinates[1]}
}

func iterateClockwise(points []Point, center Point, f func(p Point)) {
	list := make([]Point, len(points))
	for i, p := range points {
		list[i] = p
	}

	sign := func(a int) int {
		if a < 0 {
			return -1
		} else if a > 0 {
			return 1
		} else {
			return 0
		}
	}
	cross := func(v1, v2 Point) int {
		return v1.X*v2.Y - v1.Y*v2.X
	}

	// https://www.reddit.com/r/adventofcode/comments/e8r1jx/day_10_part_2_discrete_anglecomparing_function_ie/
	sort.SliceStable(list, func(i, j int) bool {
		a := list[i].Sub(center)
		b := list[j].Sub(center)

		da := a.X < 0
		db := b.X < 0

		if da != db {
			return !da
		} else if a.X == 0 && b.X == 0 {
			return sign(a.Y) > sign(b.Y)
		} else {
			return sign(cross(b, a)) < 0
		}
	})

	for _, p := range list {
		f(p)
	}
}

func vaporize(grid *Grid, offset Point, notify func(total int, p Point)) {
	left := grid.CountAsteroids()
	total := 0
	for left > 1 {
		counts := calculateGridFieldToAsteroidVisions(grid)
		iterateClockwise(counts[offset], offset, func(p Point) {
			total++
			notify(total, p)
			grid.data[p.Y][p.X].Filled = false
			left--
		})
	}
}
