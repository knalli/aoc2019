package main

import (
	"container/list"
	day18 "de.knallisworld/aoc/aoc2019/day18/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"github.com/yourbasic/graph"
	"math"
	"sort"
	"strings"
	"time"
)

const AocDay = 20
const AocDayName = "day20"
const AocDayTitle = "Day 20"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
	maze := newMaze(lines)
	fmt.Printf("%s\n", maze.ToString())
	dl.PrintSolution(fmt.Sprintf("Shortest path: %d", solution1(lines)))

	dl.PrintStepHeader(2)
	dl.PrintSolution(fmt.Sprintf("Shortest path: %d", solutions2(lines, true, false)))

}

const TILE uint8 = '.'
const WALL uint8 = '#'

type Maze struct {
	Start   day18.Point
	Stop    day18.Point
	Values  map[day18.Point]uint8
	Portals map[day18.Point]Portal
}

func (m *Maze) Clone() *Maze {
	values := make(map[day18.Point]uint8)
	for k, v := range m.Values {
		values[day18.Point{X: k.X, Y: k.Y}] = v
	}
	portals := make(map[day18.Point]Portal)
	for k, v := range m.Portals {
		portals[day18.Point{X: k.X, Y: k.Y}] = Portal{Name: v.Name, Target: day18.Point{X: v.Target.X, Y: v.Target.Y}}
	}
	return &Maze{Start: m.Start, Stop: m.Stop, Values: values, Portals: portals}
}

func (m *Maze) ToString() string {
	result := ""
	minX := math.MaxInt32
	minY := math.MaxInt32
	maxX := 0
	maxY := 0
	for p := range m.Values {
		minX = dl.MinInt(minX, p.X)
		minY = dl.MinInt(minY, p.Y)
		maxX = dl.MaxInt(maxX, p.X)
		maxY = dl.MaxInt(maxY, p.Y)
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if v, exist := m.Values[day18.Point{X: x, Y: y}]; exist {
				result += fmt.Sprintf("%c", v)
			} else {
				result += " "
			}
		}
		result += "\n"
	}

	result += fmt.Sprintf("Portals: %d\n", len(m.Portals))
	for from, portal := range m.Portals {
		result += fmt.Sprintf("  %s: %s -> %s\n", portal.Name, from.ToString(), portal.Target.ToString())
	}

	return result
}

type Portal struct {
	Name      string
	Target    day18.Point
	InnerEdge bool
}

type Path struct {
	Target day18.Point
	Value  []day18.Point
}

func solution1(lines []string) int {
	maze := newMaze(lines)
	g, nodes := buildGraph(maze, true)
	path := shortestPath(g, nodes, maze.Start, maze.Stop)
	return len(path) - 1
}

func solutions2(lines []string, debug bool, trace bool) int {
	maze := newMaze(lines)

	shortestDistances := make(map[day18.Point][]Path)
	{
		g, nodes := buildGraph(maze, false)
		search := []day18.Point{maze.Start, maze.Stop}
		for p := range maze.Portals {
			search = append(search, p)
		}
		for _, p1 := range search {
			for _, p2 := range search {
				if p1.Equal(p2) {
					continue
				}
				path := shortestPath(g, nodes, p1, p2)
				if len(path) > 1 {
					if _, exist := shortestDistances[p1]; !exist {
						shortestDistances[p1] = []Path{}
					}
					shortestDistances[p1] = append(shortestDistances[p1], Path{Target: p2, Value: path})
				}
			}
		}
	}
	getPossiblePaths := func(p day18.Point) []Path {
		var result []Path
		for _, path := range shortestDistances[p] {
			result = append(result, path)
		}
		return result
	}

	type LevelPoint struct {
		Level int
		Point day18.Point
	}

	type LevelPointDistance struct {
		LevelPoint LevelPoint
		Distance   int
	}

	visited := make(map[LevelPoint]bool)
	parents := make(map[LevelPoint]LevelPoint)
	distance := make(map[LevelPoint]int)
	queue := list.New()

	resolvePath := func(start, goal day18.Point) int {
		result := make([]LevelPointDistance, 0)
		dist := distance[LevelPoint{Point: goal}]
		result = append(result, LevelPointDistance{LevelPoint: LevelPoint{Point: goal}})
		for !result[len(result)-1].LevelPoint.Point.Equal(start) {
			last := result[len(result)-1]
			next := parents[last.LevelPoint]
			ndist := -1 // dist is not correct
			result = append(result, LevelPointDistance{LevelPoint: LevelPoint{Level: next.Level, Point: next.Point}, Distance: ndist})
		}
		// reverse
		sort.SliceStable(result, func(i, j int) bool {
			return true
		})

		if debug {
			lastLevel := 0
			lastPortal := "AA"
			for i, x := range result {
				if i == 0 {
					continue
				}
				nextPortal := maze.Portals[x.LevelPoint.Point].Name
				if lastLevel == x.LevelPoint.Level {
					fmt.Printf("%sWalk from %s to %s (%d steps)\n", strings.Repeat(" ", lastLevel), lastPortal, nextPortal, x.Distance)
				} else if lastLevel < x.LevelPoint.Level {
					fmt.Printf("%sRecurse into level %d through %s (1 step)\n", strings.Repeat(" ", lastLevel), x.LevelPoint.Level, nextPortal)
				} else {
					fmt.Printf("%sReturn to level %d through %s (1 step)\n", strings.Repeat(" ", lastLevel), x.LevelPoint.Level, nextPortal)
				}
				lastLevel = x.LevelPoint.Level
				lastPortal = nextPortal
			}
		}
		return dist
	}

	queue.PushFront(LevelPoint{Point: maze.Start})
	visited[LevelPoint{Point: maze.Start}] = true
	distance[LevelPoint{Point: maze.Start}] = 0
	for queue.Len() > 0 {
		var current LevelPoint
		{
			next := queue.Front()
			current = next.Value.(LevelPoint)
			queue.Remove(next)
		}
		for _, path := range getPossiblePaths(current.Point) {
			child := LevelPoint{Level: current.Level, Point: path.Target}
			if child.Point.Equal(maze.Start) {
				// skip start
				continue
			} else if child.Point.Equal(maze.Stop) {
				if current.Level == 0 {
					parents[child] = current
					nextDistance := distance[current] + len(path.Value) - 1
					if d, exist := distance[child]; !exist || d > nextDistance {
						distance[child] = nextDistance
					}
					return resolvePath(maze.Start, maze.Stop)
				} else {
					// skip
					continue
				}
			}

			portal := maze.Portals[child.Point]

			// rule: on level 0 no outer portals
			if current.Level == 0 && !portal.InnerEdge {
				continue
			}

			if _, exist := visited[child]; exist {
				continue
			}

			if trace {
				fmt.Printf("%sWalk from %s to %s (%d steps)\n", strings.Repeat(" ", current.Level), maze.Portals[current.Point].Name, portal.Name, len(path.Value)-1)
			}
			parents[child] = current
			nextDistance := distance[current] + 1
			if d, exist := distance[child]; !exist || d > nextDistance {
				distance[child] = nextDistance
			}

			if portal.InnerEdge {
				// recurse (aka add)
				if trace {
					fmt.Printf("%sRecurse into level %d through %s (1 step)\n", strings.Repeat(" ", current.Level), current.Level+1, portal.Name)
				}
				levelNext := LevelPoint{Level: current.Level + 1, Point: portal.Target}
				queue.PushBack(levelNext)
				parents[levelNext] = child
				nextDistance := distance[child] + len(path.Value) - 1
				if d, exist := distance[levelNext]; !exist || d > nextDistance {
					distance[levelNext] = nextDistance
				}
				visited[levelNext] = true
			} else {
				// return (aka remove)
				if trace {
					fmt.Printf("%sReturn to level %d through %s (1 step)\n", strings.Repeat(" ", current.Level), current.Level-1, portal.Name)
				}
				levelNext := LevelPoint{Level: current.Level - 1, Point: portal.Target}
				queue.PushBack(levelNext)
				nextDistance := distance[child] + len(path.Value) - 1
				if d, exist := distance[levelNext]; !exist || d > nextDistance {
					distance[levelNext] = nextDistance
				}
				parents[levelNext] = child
				visited[levelNext] = true
			}
			visited[child] = true
		}
	}
	return resolvePath(maze.Start, maze.Stop)
}

func buildGraph(maze *Maze, addPortalJumps bool) (*graph.Mutable, map[day18.Point]int) {
	nodes := make(map[day18.Point]int)
	for p, v := range maze.Values {
		if v == TILE {
			nodes[p] = len(nodes)
		}
	}
	g := graph.New(len(nodes))
	for p, v := range maze.Values {
		if v == TILE {
			if addPortalJumps {
				if other, exist := maze.Portals[p]; exist {
					g.AddCost(nodes[p], nodes[other.Target], 1)
				}
			}
			right := p.East()
			if v, exist := maze.Values[right]; exist && v == TILE {
				g.AddBothCost(nodes[p], nodes[right], 1)
			}
			bottom := p.South()
			if v, exist := maze.Values[bottom]; exist && v == TILE {
				g.AddBothCost(nodes[p], nodes[bottom], 1)
			}
		}
	}
	return g, nodes
}

func shortestPath(g *graph.Mutable, nodes map[day18.Point]int, start day18.Point, goal day18.Point) []day18.Point {
	path, _ := graph.ShortestPath(g, nodes[start], nodes[goal])

	reverseMap := make(map[int]day18.Point)
	for k, v := range nodes {
		reverseMap[v] = k
	}

	result := make([]day18.Point, 0)
	for _, v := range path {
		result = append(result, reverseMap[v])
	}
	return result
}

func newMaze(lines []string) *Maze {
	values := make(map[day18.Point]uint8)
	portals := make(map[day18.Point]Portal)
	tempPortals := make(map[string]day18.Point)
	maxHeight := len(lines)
	maxWidth := len(lines[0])
	for _, line := range lines {
		maxWidth = dl.MaxInt(maxWidth, len(line))
	}
	isUpperCaseLetter := func(c uint8) bool {
		return 'A' <= c && c <= 'Z'
	}
	for y, line := range lines {
		for x := range line {
			p := day18.Point{X: x, Y: y}
			c := line[x]
			if c == TILE || c == WALL {
				values[p] = c
			}

			if c == TILE {
				var name string
				if y-2 >= 0 && isUpperCaseLetter(lines[y-1][x]) && isUpperCaseLetter(lines[y-2][x]) {
					name = string([]uint8{lines[y-2][x], lines[y-1][x]})
				} else if x+2 < maxWidth && isUpperCaseLetter(lines[y][x+1]) && isUpperCaseLetter(lines[y][x+2]) {
					// right
					name = string([]uint8{lines[y][x+1], lines[y][x+2]})
				} else if y+2 < maxHeight && isUpperCaseLetter(lines[y+1][x]) && isUpperCaseLetter(lines[y+2][x]) {
					// bottom
					name = string([]uint8{lines[y+1][x], lines[y+2][x]})
				} else if x-2 >= 0 && isUpperCaseLetter(lines[y][x-2]) && isUpperCaseLetter(lines[y][x-1]) {
					// left
					name = string([]uint8{lines[y][x-2], lines[y][x-1]})
				}
				if name == "" {
					continue
				}
				if other, exist := tempPortals[name]; exist {
					portals[p] = Portal{Name: name, Target: other, InnerEdge: 4 < p.X && p.X < maxWidth-4 && 4 < p.Y && p.Y < maxHeight-4}
					portals[other] = Portal{Name: name, Target: p, InnerEdge: 4 < other.X && other.X < maxWidth-4 && 4 < other.Y && other.Y < maxHeight-4}
					delete(tempPortals, name)
				} else {
					tempPortals[name] = p
				}
			}
		}
	}
	return &Maze{Start: tempPortals["AA"], Stop: tempPortals["ZZ"], Values: values, Portals: portals}
}
