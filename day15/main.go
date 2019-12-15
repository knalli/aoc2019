package main

import (
	"container/list"
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"math"
	"time"
)

const AocDay = 15
const AocDayName = "day15"
const AocDayTitle = "Day 15"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		movements := findOptimalPathToOxygen(program, false)
		dl.PrintSolution(fmt.Sprintf("Number of movements is %d", movements))
	}

	{
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		world, oxygenPos := discoverMap(program, false, false)
		minute := 0
		world[oxygenPos] = OXYGEN
		for {
			currentFree := 0
			for _, v := range world {
				switch v {
				case FREE, FREE_VISITED:
					currentFree++
				}
			}
			if currentFree == 0 {
				break
			}
			flooding := make([]Point, 0)
			for pos, posValue := range world {
				if posValue == OXYGEN {
					for _, adjacent := range []Point{pos.North(), pos.East(), pos.South(), pos.West()} {
						if adjacentValue, exist := world[adjacent]; exist {
							if adjacentValue != WALL {
								flooding = append(flooding, adjacent)
							}
						}
					}
				}
			}
			for _, point := range flooding {
				world[point] = OXYGEN
			}
			minute++
		}
		renderMap(world, Point{math.MinInt32, math.MinInt32}, Point{math.MinInt32, math.MinInt32}, list.New(), list.New())
		dl.PrintSolution(fmt.Sprintf("%d minutes until all is filled with oxygen", minute))
	}

}

type Point struct {
	X, Y int
}

func (p *Point) East() Point {
	return Point{p.X + 1, p.Y}
}
func (p *Point) West() Point {
	return Point{p.X - 1, p.Y}
}
func (p *Point) North() Point {
	return Point{p.X, p.Y + 1}
}
func (p *Point) South() Point {
	return Point{p.X, p.Y - 1}
}
func (p *Point) ToString() string {
	return fmt.Sprintf("(%d/%d)", p.X, p.Y)
}

const NORTH = 1
const SOUTH = 2
const WEST = 3
const EAST = 4
const WALL = '#'
const FREE = '.'
const OXYGEN = 'O'
const FREE_VISITED = '-'
const DROID_HIT_WALL = 0
const DROID_MOVED = 1
const DROID_MOVED_AND_FOUND = 2

func renderMap(world map[Point]int16, pos Point, oxygen Point, track *list.List, trackPos *list.List) {
	bl := Point{0, 0}
	tr := Point{0, 0}
	for p := range world {
		bl.X = dl.MinInt(bl.X, p.X)
		bl.Y = dl.MinInt(bl.Y, p.Y)
		tr.X = dl.MaxInt(tr.X, p.X)
		tr.Y = dl.MaxInt(tr.Y, p.Y)
	}
	fmt.Println()
	dir := "^"
	{
		x := track.Back()
		if x != nil {
			switch x.Value {
			case NORTH:
				dir = "^"
			case SOUTH:
				dir = "v"
			case WEST:
				dir = "<"
			case EAST:
				dir = ">"
			}
		}
	}
	for y := tr.Y; y >= bl.Y; y-- {
		line := ""
		for x := bl.X; x <= tr.X; x++ {
			if oxygen.X == x && oxygen.Y == y {
				line += fmt.Sprintf("%s", Blue("O"))
			} else if pos.X == x && pos.Y == y {
				line += fmt.Sprintf("%s", Blue(dir))
			} else if v, exist := world[Point{x, y}]; exist {
				tracked := false
				p := trackPos.Front()
				for p != nil {
					point := p.Value.(Point)
					if point.X == x && point.Y == y {
						tracked = true
						break
					}
					p = p.Next()
				}
				if tracked {
					line += fmt.Sprintf("%s", Red(string(v)))
				} else {
					switch v {
					case WALL:
						line += string(v)
					case FREE:
						line += string(FREE)
					case FREE_VISITED:
						line += string(FREE_VISITED)
					case OXYGEN:
						line += fmt.Sprintf("%s", Blue(string(OXYGEN)))
					}
				}
			} else {
				line += " "
			}
		}
		fmt.Printf("%s\n", line)
	}
	fmt.Println()
	fmt.Printf("Directions: ")
	movements := 0
	x := track.Front()
	for x != nil {
		s := ""
		switch x.Value {
		case NORTH:
			s = "^"
		case SOUTH:
			s = "v"
		case WEST:
			s = "<"
		case EAST:
			s = ">"
		}
		fmt.Printf("%s ", s)
		x = x.Next()
		movements++
	}
	fmt.Printf(" total = %d\n", movements)
	fmt.Println()
}

func north(world map[Point]int16, pos Point, blocked bool) (int, bool) {
	if !blocked {
		if _, exist := world[pos.North()]; !exist {
			// prefer only if empty
			return NORTH, false
		}
	}
	if v, exist := world[pos.East()]; !exist || (v != WALL && v != FREE_VISITED) {
		return EAST, false
	}
	if v, exist := world[pos.West()]; !exist || (v != WALL && v != FREE_VISITED) {
		return WEST, false
	}
	if !blocked {
		if v, exist := world[pos.North()]; !exist || (v != WALL && v != FREE_VISITED) {
			return NORTH, false
		}
	}
	return SOUTH, true
}

func east(world map[Point]int16, pos Point, blocked bool) (int, bool) {
	if !blocked {
		if _, exist := world[pos.East()]; !exist {
			// prefer only if empty
			return EAST, false
		}
	}
	if v, exist := world[pos.South()]; !exist || (v != WALL && v != FREE_VISITED) {
		return SOUTH, false
	}
	if v, exist := world[pos.North()]; !exist || (v != WALL && v != FREE_VISITED) {
		return NORTH, false
	}
	if !blocked {
		if v, exist := world[pos.East()]; !exist || (v != WALL && v != FREE_VISITED) {
			return EAST, false
		}
	}
	return WEST, true
}

func south(world map[Point]int16, pos Point, blocked bool) (int, bool) {
	if !blocked {
		if _, exist := world[pos.South()]; !exist {
			// prefer only if empty
			return SOUTH, false
		}
	}
	if v, exist := world[pos.West()]; !exist || (v != WALL && v != FREE_VISITED) {
		return WEST, false
	}
	if v, exist := world[pos.East()]; !exist || (v != WALL && v != FREE_VISITED) {
		return EAST, false
	}
	if !blocked {
		if v, exist := world[pos.South()]; !exist || (v != WALL && v != FREE_VISITED) {
			return SOUTH, false
		}
	}
	return NORTH, true
}

func west(world map[Point]int16, pos Point, blocked bool) (int, bool) {
	if !blocked {
		if _, exist := world[pos.West()]; !exist {
			// prefer only if empty
			return WEST, false
		}
	}
	if v, exist := world[pos.North()]; !exist || (v != WALL && v != FREE_VISITED) {
		return NORTH, false
	}
	if v, exist := world[pos.South()]; !exist || (v != WALL && v != FREE_VISITED) {
		return SOUTH, false
	}
	if !blocked {
		if v, exist := world[pos.West()]; !exist || (v != WALL && v != FREE_VISITED) {
			return WEST, false
		}
	}
	return EAST, true
}

func handle(world map[Point]int16, pos Point, dir int, blocked bool) (int, bool) {
	switch dir {
	case NORTH:
		return north(world, pos, blocked)
	case SOUTH:
		return south(world, pos, blocked)
	case WEST:
		return west(world, pos, blocked)
	case EAST:
		return east(world, pos, blocked)
	}
	panic("invalid handle")
}

func findOptimalPathToOxygen(program []int, debug bool) int {
	world := make(map[Point]int16, 0)

	in := make(chan int)   // program stdin
	out := make(chan int)  // program stdout
	fin := make(chan bool) // program end

	track := list.New()
	trackPos := list.New()
	backtracking := false
	dir := NORTH
	pos := Point{0, 0}
	oxygen := Point{math.MaxInt32, math.MaxInt32}

	go func() {
		_ = day09.ExecutionInstructions(program, in, out, false)
	}()

	go func() {
		for {

			// render
			if debug {
				renderMap(world, pos, oxygen, track, trackPos)
				time.Sleep(1 * time.Millisecond)
			}

			status := <-out
			switch status {
			case DROID_HIT_WALL:
				if debug {
					fmt.Printf("Droid %s hit the wall at", pos.ToString())
				}
				switch dir {
				case NORTH:
					next := pos.North()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				case SOUTH:
					next := pos.South()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				case WEST:
					next := pos.West()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				case EAST:
					next := pos.East()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				}
				in <- dir
			case DROID_MOVED, DROID_MOVED_AND_FOUND:
				if debug {
					fmt.Printf("Droid %s moved", pos.ToString())
				}
				switch dir {
				case NORTH:
					pos = pos.North()
				case SOUTH:
					pos = pos.South()
				case WEST:
					pos = pos.West()
				case EAST:
					pos = pos.East()
				}
				if debug {
					fmt.Printf(" to %s\n", pos.ToString())
				}

				if backtracking {
					if _, exist := world[pos]; !exist {
						backtracking = false
					}
				}

				if !backtracking {
					world[pos] = FREE
					track.PushBack(dir)
					trackPos.PushBack(pos)
				} else {
					track.Remove(track.Back())
					trackPos.Remove(trackPos.Back())
				}

				if status == DROID_MOVED_AND_FOUND {
					if debug {
						fmt.Println("FOUND!")
					}
					oxygen = pos
					fin <- true
					return
				} else {
					switch dir {
					case NORTH:
						d, b := north(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					case SOUTH:
						d, b := south(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					case WEST:
						d, b := west(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					case EAST:
						d, b := east(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					}
					in <- dir
				}
			}
		}
	}()

	world[pos] = '.'
	trackPos.PushBack(pos)
	in <- dir

	<-fin // wait for end of all coroutines

	renderMap(world, pos, oxygen, track, trackPos)

	return track.Len()
}

func discoverMap(program []int, stopAtOxy bool, debug bool) (map[Point]int16, Point) {
	world := make(map[Point]int16, 0)

	in := make(chan int)   // program stdin
	out := make(chan int)  // program stdout
	fin := make(chan bool) // program end

	track := list.New()
	trackPos := list.New()
	backtracking := false
	dir := NORTH
	pos := Point{0, 0}
	oxygenPos := Point{math.MaxInt32, math.MaxInt32}

	go func() {
		_ = day09.ExecutionInstructions(program, in, out, false)
	}()

	go func() {
		for {

			// render
			if debug {
				renderMap(world, pos, oxygenPos, track, trackPos)
				time.Sleep(100 * time.Millisecond)
			}

			status := <-out
			switch status {
			case DROID_HIT_WALL:
				if debug {
					fmt.Printf("Droid %s hit the wall at", pos.ToString())
				}
				switch dir {
				case NORTH:
					next := pos.North()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				case SOUTH:
					next := pos.South()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				case WEST:
					next := pos.West()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				case EAST:
					next := pos.East()
					if debug {
						fmt.Printf("%s\n", next.ToString())
					}
					world[next] = WALL
					d, b := handle(world, pos, track.Back().Value.(int), true)
					dir = d
					if b {
						world[pos] = FREE_VISITED // close
						backtracking = b
					}
				}
				in <- dir
			case DROID_MOVED, DROID_MOVED_AND_FOUND:
				if debug {
					fmt.Printf("Droid %s moved", pos.ToString())
				}
				switch dir {
				case NORTH:
					pos = pos.North()
				case SOUTH:
					pos = pos.South()
				case WEST:
					pos = pos.West()
				case EAST:
					pos = pos.East()
				}
				if debug {
					fmt.Printf(" to %s\n", pos.ToString())
				}

				if backtracking {
					if _, exist := world[pos]; !exist {
						backtracking = false
					}
				}

				if !backtracking {
					world[pos] = FREE
					track.PushBack(dir)
					trackPos.PushBack(pos)
				} else {
					track.Remove(track.Back())
					trackPos.Remove(trackPos.Back())
				}

				if !stopAtOxy && track.Len() == 0 {
					fin <- true
					return
				}

				if status == DROID_MOVED_AND_FOUND {
					if debug {
						fmt.Println("FOUND!")
					}
					oxygenPos = pos
				}

				if debug {
					renderMap(world, pos, oxygenPos, track, trackPos)
				}

				if stopAtOxy && status == DROID_MOVED_AND_FOUND {
					fin <- true
					return
				} else {
					d := dir
					if backtracking {
						d = track.Back().Value.(int)
					}
					switch d {
					case NORTH:
						d, b := north(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					case SOUTH:
						d, b := south(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					case WEST:
						d, b := west(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					case EAST:
						d, b := east(world, pos, false)
						dir = d
						if b {
							backtracking = b
						}
					}
					if backtracking && world[pos] == FREE {
						world[pos] = FREE_VISITED // close
					}
					in <- dir
				}
			}
		}
	}()

	world[pos] = '.'
	trackPos.PushBack(pos)
	in <- dir

	<-fin // wait for end of all coroutines

	renderMap(world, pos, oxygenPos, track, trackPos)

	return world, oxygenPos
}
