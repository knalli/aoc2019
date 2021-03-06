package main

import (
	"container/list"
	day18 "de.knallisworld/aoc/aoc2019/day18/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"math"
	"strings"
	"time"
)

const AocDay = 24
const AocDayName = "day24"
const AocDayTitle = "Day 24"

const BUGS = "#"
const EMPTY = "."

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
	grid := day18.NewMap(puzzle)
	dl.PrintSolution(fmt.Sprintf("Biodiversity rating is %d", solution1(&grid)))

	dl.PrintStepHeader(2)
	dl.PrintSolution(fmt.Sprintf("After %d minutes, there are %d bugs", 200, solution2(&grid, 200)))

}

func solution1(grid *day18.Map) int {
	fmt.Printf("initial state:\n%s\n\n", grid.ToString())

	type Item struct {
		grid   day18.Map
		minute int
	}

	cache := make(map[string]Item)
	cache[grid.ToString()] = Item{grid: *grid, minute: 0}
	for i := 1; i < math.MaxInt16; i++ {
		grid = tick(grid)
		toString := grid.ToString()
		fmt.Printf("After %d minute:\n%s\n\n", i, toString)
		if _, exist := cache[toString]; exist {
			return biodiversityRating(grid)
		}
		cache[toString] = Item{grid: *grid, minute: i}
	}

	return -1
}

func solution2(grid *day18.Map, minutes int) int {
	fmt.Printf("initial state:\n%s\n\n", grid.ToString())

	l := list.New()
	l.PushBack(grid)
	state := &GridState{
		initial: grid,
		level:   0,
		list:    l,
	}
	for i := 0; i < minutes; i++ {
		state = tickRecursive(state)
	}

	return state.CountBugs()
}

func biodiversityRating(grid *day18.Map) int {
	power := 0
	result := 0
	grid.Each(func(p day18.Point, v string) {
		if v == BUGS {
			result += int(math.Pow(float64(2), float64(power)))
		}
		power++
	})
	return result
}

func tick(grid *day18.Map) *day18.Map {
	next := grid.Clone()
	grid.Each(func(p day18.Point, v string) {
		if v == BUGS {
			adjacentBugs := 0
			for _, adjacent := range p.Adjacents() {
				if !grid.Contains(adjacent) {
					continue
				}
				if *grid.Get(adjacent) == BUGS {
					adjacentBugs++
				}
			}
			if adjacentBugs != 1 {
				next.Set(p, EMPTY)
			}
		} else if v == EMPTY {
			adjacentBugs := 0
			for _, adjacent := range p.Adjacents() {
				if !grid.Contains(adjacent) {
					continue
				}
				if *grid.Get(adjacent) == BUGS {
					adjacentBugs++
				}
			}
			if 1 <= adjacentBugs && adjacentBugs <= 2 {
				next.Set(p, BUGS)
			}
		}
	})
	return &next
}

type GridState struct {
	initial *day18.Map
	level   int
	list    *list.List
}

func (s *GridState) Clone() *GridState {
	initial := s.initial.Clone()
	l := list.New()
	l.PushBackList(s.list)
	return &GridState{
		initial: &initial,
		level:   s.level,
		list:    l,
	}
}

func (s *GridState) CountBugs() int {
	total := 0
	elem := s.list.Front()
	for elem != nil {
		grid := elem.Value.(*day18.Map)
		total += grid.Count(func(v string) bool {
			return v == BUGS
		})
		elem = elem.Next()
	}
	return total
}

type GridScan struct {
	grid  *day18.Map
	inner *day18.Map
	outer *day18.Map
}

func createEmptyGrid() *day18.Map {
	line := strings.Repeat(EMPTY, 5)
	grid := day18.NewMap([]string{
		line,
		line,
		line,
		line,
		line,
	})
	return &grid
}

func tickRecursive(state *GridState) *GridState {
	next := state.Clone()
	next.level++

	l := list.New()
	l.PushBackList(state.list)
	l.PushFront(createEmptyGrid()) // new level
	l.PushBack(createEmptyGrid())  // new level
	l.PushFront(createEmptyGrid()) // only for easier loop
	l.PushBack(createEmptyGrid())  // only for easier loop

	extractNext := func(l *list.List) *day18.Map {
		var elem *day18.Map
		{
			next := l.Front()
			elem = next.Value.(*day18.Map)
			l.Remove(next)
		}
		return elem
	}

	getScans := func(l *list.List) []*GridScan {
		result := make([]*GridScan, 0)
		outer := extractNext(l)
		current := extractNext(l)
		for l.Len() > 0 {
			inner := extractNext(l)
			result = append(result, &GridScan{
				grid:  current,
				inner: inner,
				outer: outer,
			})
			outer = current
			current = inner
		}
		return result
	}

	isBug := func(v string) bool {
		return v == BUGS
	}

	next.list.Init()
	for _, scan := range getScans(l) {
		grid := scan.grid
		inner := scan.inner
		outer := scan.outer
		result := grid.Clone()
		grid.Each(func(p day18.Point, v string) {
			if p.X == 2 && p.Y == 2 {
				return
			}
			if isBug(v) {
				bugs := countRecursiveAdjacents(p, outer, grid, inner, isBug)
				if bugs != 1 {
					result.Set(p, EMPTY)
				}
			} else {
				bugs := countRecursiveAdjacents(p, outer, grid, inner, isBug)
				if bugs == 1 || bugs == 2 {
					result.Set(p, BUGS)
				}
			}
		})
		next.list.PushBack(&result)
	}

	return next
}

func countRecursiveAdjacents(p day18.Point, outer *day18.Map, grid *day18.Map, inner *day18.Map, condition func(v string) bool) int {
	total := 0
	leftOuterPoint := day18.Point{X: 1, Y: 2}
	rightOuterPoint := day18.Point{X: 3, Y: 2}
	topOuterPoint := day18.Point{X: 2, Y: 1}
	bottomOuterPoint := day18.Point{X: 2, Y: 3}
	for _, a := range p.Adjacents() {
		// look for outer
		if a.X == -1 {
			if condition(*outer.Get(leftOuterPoint)) {
				total++
			}
		} else if a.X == grid.Width() {
			if condition(*outer.Get(rightOuterPoint)) {
				total++
			}
		}
		if a.Y == -1 {
			if condition(*outer.Get(topOuterPoint)) {
				total++
			}
		} else if a.Y == grid.Height() {
			if condition(*outer.Get(bottomOuterPoint)) {
				total++
			}
		}
		if !grid.Contains(a) {
			continue
		}
		// look for inner
		if a.X == 2 && a.Y == 2 {
			if p.X < a.X {
				// left side
				inner.EachColumnCell(0, func(_ day18.Point, v2 string) {
					if condition(v2) {
						total++
					}
				})
			} else if p.X > a.X {
				// right side
				inner.EachColumnCell(grid.Width()-1, func(_ day18.Point, v2 string) {
					if condition(v2) {
						total++
					}
				})
			} else if p.Y < a.Y {
				// top side
				inner.EachRowColumn(0, func(_ day18.Point, v2 string) {
					if condition(v2) {
						total++
					}
				})
			} else if p.Y > a.Y {
				// bottom side
				inner.EachRowColumn(grid.Height()-1, func(_ day18.Point, v2 string) {
					if condition(v2) {
						total++
					}
				})
			}
		} else if condition(*grid.Get(a)) {
			total++
		}
	}
	return total
}
