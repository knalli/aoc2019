package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"math"
	"strings"
	"time"
)

const AocDay = 3
const AocDayName = "day03"
const AocDayTitle = "Day 03"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(0)
	fmt.Printf("💪 Computing...\n")
	lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
	wires := getWires(lines)
	intersections := getIntersections(getPaths(lines))
	fmt.Printf("👉 Found %d intersections at all\n", len(intersections))

	{
		dl.PrintStepHeader(1)
		centralCrossingPoint, centralCrossingDistance := getShortestDistance(Point{0, 0}, intersections)
		dl.PrintSolution(fmt.Sprintf("Most central intersection is %s with Manhatten Distance %d", centralCrossingPoint.ToString(), centralCrossingDistance))
	}

	{
		dl.PrintStepHeader(2)
		minStepsTotal, minStepsPoint := getShortestPath(Point{0, 0}, wires, intersections)
		dl.PrintSolution(fmt.Sprintf("Most central intersection is %s with Lowest Total Of Steps %d", minStepsPoint.ToString(), minStepsTotal))
	}

}

type Point struct {
	X, Y int
}

func (p *Point) Equals(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p *Point) ToString() string {
	return fmt.Sprintf("%d/%d", p.X, p.Y)
}

func (p Point) Clone() Point {
	return Point{p.X, p.Y}
}

type Wire struct {
	Segments []WireSegment
}

type WireSegment struct {
	Start     Point
	Stop      Point
	Cost      int
	Length    int
	Direction uint8
	path      []Point
}

func (w *WireSegment) IsVertical() bool {
	return w.Direction == 'U' || w.Direction == 'D'
}

func (w *WireSegment) isInverted() bool {
	return w.Direction == 'D' || w.Direction == 'L'
}

func (w *WireSegment) Path() []Point {
	if w.path == nil {
		result := make([]Point, 0)
		current := w.Start.Clone()
		result = append(result, current)
		for !w.Stop.Equals(current) {
			switch w.Direction {
			case 'U':
				current = Point{current.X, current.Y + 1}
				break
			case 'R':
				current = Point{current.X + 1, current.Y}
				break
			case 'D':
				current = Point{current.X, current.Y - 1}
				break
			case 'L':
				current = Point{current.X - 1, current.Y}
				break
			}
			result = append(result, current)
		}
		w.path = result
	}
	return w.path
}

func newWireByString(line string) Wire {
	parts := strings.Split(line, ",")
	segments := make([]WireSegment, 0)
	current := Point{X: 0, Y: 0}
	cost := 0
	for _, part := range parts {
		dir := part[0]
		length := dl.ParseInt(part[1:])
		var stop Point
		switch dir {
		case 'U':
			stop = Point{
				X: current.X,
				Y: current.Y + length,
			}
			break
		case 'R':
			stop = Point{
				X: current.X + length,
				Y: current.Y,
			}
			break
		case 'D':
			stop = Point{
				X: current.X,
				Y: current.Y - length,
			}
			break
		case 'L':
			stop = Point{
				X: current.X - length,
				Y: current.Y,
			}
			break
		}

		segments = append(segments, WireSegment{
			Start:     current,
			Stop:      stop,
			Length:    length,
			Cost:      cost,
			Direction: dir,
		})
		current = stop.Clone()
		cost += length
	}
	return Wire{
		Segments: segments,
	}
}

func getPath(line string) []Point {
	var x = 0
	var y = 0
	var parts = strings.Split(line, ",")
	var result = make([]Point, 0)
	result = append(result, Point{x, y})
	for _, line := range parts {
		length := dl.ParseInt(line[1:])
		switch line[0] {
		case 'U':
			for d := 0; d < length; d++ {
				y += 1
				result = append(result, Point{x, y})
			}
			break
		case 'R':
			for d := 0; d < length; d++ {
				x += 1
				result = append(result, Point{x, y})
			}
			break
		case 'D':
			for d := 0; d < length; d++ {
				y -= 1
				result = append(result, Point{x, y})
			}
			break
		case 'L':
			for d := 0; d < length; d++ {
				x -= 1
				result = append(result, Point{x, y})
			}
			break
		default:
			panic("invalid instruction.Direction")
		}
	}
	return result
}

func getPaths(lines []string) [][]Point {
	var result = make([][]Point, len(lines))
	for i, line := range lines {
		result[i] = getPath(line)
	}
	return result
}

func getWires(lines []string) []Wire {
	var result = make([]Wire, len(lines))
	for i, line := range lines {
		result[i] = newWireByString(line)
	}
	return result
}

func getIntersections(paths [][]Point) []Point {
	all := make(map[string]int)
	for _, path := range paths {
		local := make(map[string]bool)
		for _, c := range path {
			key := fmt.Sprintf("%d/%d", c.X, c.Y)
			local[key] = true
		}
		for key := range local {
			if _, exist := all[key]; exist {
				all[key] = all[key] + 1
			} else {
				all[key] = 1
			}
		}
	}

	crossed := 0
	for _, m := range all {
		if m > 1 {
			crossed++
		}
	}
	var result = make([]Point, crossed)
	var i = 0
	for s, m := range all {
		if m > 1 {
			split := strings.Split(s, "/")
			result[i] = Point{dl.ParseInt(split[0]), dl.ParseInt(split[1])}
			i++
		}
	}
	return result
}

func getManhattenDistance(a Point, b Point) int {
	return int(math.Abs(float64(a.X-b.X))) + int(math.Abs(float64(a.Y-b.Y)))
}

func getShortestDistance(base Point, crossings []Point) (Point, int) {
	min := math.MaxInt64
	minPoint := base
	for _, p := range crossings {
		d := getManhattenDistance(base, p)
		if !(p.X == 0 && p.Y == 0) && d < min {
			min = d
			minPoint = p
		}
	}
	return minPoint, min
}

type dataItem struct {
	Dist      int
	Direction uint8
}

func getShortestPath(base Point, wires []Wire, intersections []Point) (int, Point) {
	allData := make(map[int]map[string]dataItem)
	for i, wire := range wires {
		data := make(map[string]dataItem)
		for _, segment := range wire.Segments {
			for pointIdx, point := range segment.Path() {
				// ignore base (0/0)
				if point.Equals(base) {
					continue
				}
				pointKey := point.ToString()
				var d int
				if segment.isInverted() {
					d = segment.Cost + (segment.Length - pointIdx)
				} else {
					d = segment.Cost + pointIdx
				}
				if v, exist := data[pointKey]; !exist || d < v.Dist {
					data[pointKey] = dataItem{
						Dist:      d,
						Direction: segment.Direction,
					}
					// TODO what if overridden direction is important?
				}
			}
		}
		allData[i] = data
	}

	minSteps := math.MaxInt64
	minStepsPoint := Point{0, 0}
	for _, intersection := range intersections {
		key := intersection.ToString()
		wireSegment0 := allData[0][key]
		wireSegment1 := allData[1][key]
		wireSegment0IsVertical := wireSegment0.Direction == 'U' || wireSegment0.Direction == 'D'
		wireSegment1IsVertical := wireSegment1.Direction == 'U' || wireSegment1.Direction == 'D'
		if (wireSegment0IsVertical && !wireSegment1IsVertical) || (!wireSegment0IsVertical && wireSegment1IsVertical) {
			steps := wireSegment0.Dist + wireSegment1.Dist
			if steps < minSteps {
				minSteps = steps
				minStepsPoint = intersection.Clone()
			}
		}
	}

	return minSteps, minStepsPoint
}
