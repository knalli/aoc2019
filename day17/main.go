package main

import (
	"container/list"
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 17
const AocDayName = "day17"
const AocDayTitle = "Day 17"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		world := RenderAsciiToString(runProgram1(program, false))
		fmt.Printf("%s\n", world)
		world = world[:len(world)-3] // remove last linebreak
		intersections := findIntersections(splitWorld(world))
		sum := 0
		for _, intersection := range intersections {
			sum += intersection.X * intersection.Y
		}
		dl.PrintSolution(fmt.Sprintf("There are %d intersections and the sum of the alignment parameters is %d.", len(intersections), sum))
	}

	{
		debug := false
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		world := RenderAsciiToString(runProgram1(program, false))
		fmt.Printf("%s\n", world)
		world = world[:len(world)-3] // remove last linebreak
		movements := findMovements(splitWorld(world))
		fmt.Printf("Movements: %s\n\n", movements)
		a, b, c, main, groups := findRepeatingGroups(movements, 3, 20)
		fmt.Printf("Main routine: %s\n", main)
		fmt.Printf("Function A: %s\n", groups[0])
		fmt.Printf("Function C: %s\n", groups[1])
		fmt.Printf("Function B: %s\n", groups[2])
		if debug {
			fmt.Printf("A=%d,B=%d,C=%d\n", a, b, c)
		}

		inputs := make([]int, 0)
		for _, c := range main {
			inputs = append(inputs, int(c))
		}
		inputs = append(inputs, '\n')
		for _, group := range groups {
			for i, v := range strings.Split(group, ",") {
				if i > 0 {
					inputs = append(inputs, ',')
				}
				if v == "L" || v == "R" {
					u := v[0]
					inputs = append(inputs, int(u))
				} else {
					for _, v2 := range v {
						inputs = append(inputs, int(v2))
					}
				}
			}
			inputs = append(inputs, '\n')
		}

		inputs = append(inputs, 'n')
		inputs = append(inputs, '\n')

		//println("Instructions")
		//println(RenderAsciiToString(inputs))

		println("Run program...")
		program[0] = 2
		output := runProgram2(program, inputs, false)
		if debug {
			println(RenderAsciiToString(output))
		}
		dust := dl.MaxInt(output[len(output)-2], output[len(output)-1]) // workaround (last output is sometimes 0)
		dl.PrintSolution(fmt.Sprintf("Collected dust is %d", dust))
	}

}

func findRepeatingGroups(line string, groups int, maxCommandLength int) (int, int, int, string, []string) {
	l := len(line)
	lengths := make([]int, groups)
	for i := 0; i < groups; i++ {
		lengths[i] = 1
	}
	a := 0
	b := 0
	c := 0

	findNextFree := func(line string) int {
		parts := strings.Split(line, ",")
		for i, part := range parts {
			if part != "A" && part != "B" && part != "C" {
				return i
			}
		}
		return -1
	}

	for {
		temp := "" + line
		a++
		if a+b+c > l {
			a = 1
			b++
		}
		if a+b+c > l {
			b = 1
			c++
		}
		if a+b+c > l {
			break
		}
		//fmt.Printf("%02d %02d %02d\n", a, b, c)

		var subA, subB, subC string

		{
			parts := strings.Split(temp, ",")
			avail := findNextFree(temp)
			if avail+2*a > len(parts) {
				continue
			}
			subA = strings.Join(parts[avail:avail+2*a], ",")
			if strings.ContainsAny(subA, "ABC") {
				continue
			}
			if len(subA) > maxCommandLength {
				continue
			}
			temp = strings.Replace(temp, subA, "A", -1)
		}

		{
			parts := strings.Split(temp, ",")
			avail := findNextFree(temp)
			if avail+2*b > len(parts) {
				continue
			}
			subB = strings.Join(parts[avail:avail+2*b], ",")
			if strings.ContainsAny(subB, "ABC") {
				continue
			}
			if len(subB) > maxCommandLength {
				continue
			}
			temp = strings.Replace(temp, subB, "B", -1)
		}

		{
			parts := strings.Split(temp, ",")
			avail := findNextFree(temp)
			if avail+2*c > len(parts) {
				continue
			}
			subC = strings.Join(parts[avail:avail+2*c], ",")
			if strings.ContainsAny(subC, "ABC") {
				continue
			}
			if len(subC) > maxCommandLength {
				continue
			}
			temp = strings.Replace(temp, subC, "C", -1)
		}

		if -1 == findNextFree(temp) {
			if len(temp) > maxCommandLength {
				fmt.Printf("Found main function, but exceeding command length: %s\n", temp)
				continue
			}
			return a, b, c, temp, []string{subA, subB, subC}
		}
	}

	return 0, 0, 0, line, nil
}

func findTopRepeatingString(line string, minLength int, maxLength int, max int) []string {
	result := make([]string, 0)
	temp := "" + line
	for i := 0; i < max; i++ {
		r := findLongestRepeatingString(temp, minLength, maxLength)
		if r == "" {
			break
		}
		result = append(result, r)
		temp = strings.Replace(temp, r, "_", -1)
	}
	return result
}

func findLongestRepeatingString(line string, minLength int, maxLength int) string {

	length := len(line)
	found := make(map[string]int)
	longestStrike := ""

	for ; minLength < length && minLength <= maxLength; minLength++ {
		known := make(map[string]bool)
		for pos := 0; pos+minLength < length; pos++ {
			needle := line[pos : pos+minLength]
			if _, exist := known[needle]; exist {
				continue
			} else {
				known[needle] = true
			}
			count := 0
			for search := pos + minLength; search+minLength < length; search++ {
				sub := line[search : search+minLength]
				if needle == sub {
					count++
				}
			}
			if count > 0 {
				found[needle] = count
				if len(needle) > len(longestStrike) {
					longestStrike = needle
				}
			}
		}
	}

	return longestStrike
}

func findMovements(world [][]uint8) string {
	pos := findPosition(world, func(_, _ int, v uint8) bool {
		return v == '^' || v == '>' || v == 'v' || v == '<'
	})
	path := list.New()
	path.PushBack(pos)
	movements := list.New()
	dir := world[pos.Y][pos.X]

	getNextPoint := func(origin Point, dir uint8) *Point {
		switch dir {
		case '^':
			next := origin.South() // mirror
			if next.Y < 0 {
				return nil
			}
			return &next
		case '>':
			next := origin.East()
			if next.X == len(world[0]) {
				return nil
			}
			return &next
		case 'v':
			next := origin.North() // mirror
			if next.Y == len(world) {
				return nil
			}
			return &next
		case '<':
			next := origin.West()
			if next.X < 0 {
				return nil
			}
			return &next
		default:
			return nil
		}
	}
	getNextPointWithFallback := func(origin Point, dir uint8) (*Point, uint8) {
		next := getNextPoint(origin, dir)
		if next != nil && world[next.Y][next.X] == '#' {
			return next, dir
		}
		alternates := map[uint8][]uint8{}
		alternates['^'] = []uint8{'<', '>'}
		alternates['>'] = []uint8{'^', 'v'}
		alternates['v'] = []uint8{'<', '>'}
		alternates['<'] = []uint8{'^', 'v'}
		for _, ndir := range alternates[dir] {
			next = getNextPoint(origin, ndir)
			if next == nil {
				continue
			}
			if world[next.Y][next.X] == '#' {
				return next, ndir
			}
		}
		return nil, dir
	}

	intersections := findIntersections(world)

	isIntersection := func(p Point) bool {
		for _, i := range intersections {
			if i.Equal(p) {
				return true
			}
		}
		return false
	}

	strike := uint8(0)
	for {
		path.PushBack(pos)

		if !isIntersection(pos) {
			world[pos.Y][pos.X] = 'R'
		}
		npos, ndir := getNextPointWithFallback(pos, dir)

		if npos == nil {
			if strike > 0 {
				movements.PushBack(fmt.Sprintf("%d", strike))
			}
			break
		} else if dir != ndir {
			if strike > 0 {
				movements.PushBack(fmt.Sprintf("%d", strike))
			}
			if (dir == '^' && ndir == '>') || (dir == '>' && ndir == 'v') || (dir == 'v' && ndir == '<') || (dir == '<' && ndir == '^') {
				movements.PushBack("R")
			} else {
				movements.PushBack("L")
			}
			strike = 1
		} else {
			strike++
		}
		pos = *npos
		dir = ndir
	}

	result := ""
	m := movements.Front()
	for m != nil {
		if result != "" {
			result += ","
		}
		result += m.Value.(string)
		m = m.Next()
	}
	return result
}

type Point struct {
	X, Y int
}

func (p *Point) Equal(o Point) bool {
	return p.X == o.X && p.Y == o.Y
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

func runProgram1(program []int, debug bool) []int {

	world := make([]int, 0)

	in := make(chan int)     // program stdin
	out := make(chan int)    // program stdout
	fin := make(chan bool)   // program end
	halt := make(chan error) // program halt

	go func() {
		halt <- day09.ExecutionInstructions(program, in, out, debug)
	}()

	go func() {
		for {
			select {
			case <-halt:
				fin <- true
				return
			default:
				world = append(world, <-out)
			}
		}
	}()

	<-fin // wait for end of all coroutines

	return world
}

func runProgram2(program []int, inputs []int, debug bool) []int {

	result := make([]int, 0)

	in := make(chan int)     // program stdin
	out := make(chan int)    // program stdout
	fin := make(chan bool)   // program end
	halt := make(chan error) // program halt

	go func() {
		halt <- day09.ExecutionInstructions(program, in, out, debug)
	}()

	go func() {
		for _, n := range inputs {
			in <- n
		}
	}()

	go func() {
		for {
			select {
			case <-halt:
				fin <- true
				return
			default:
				result = append(result, <-out)
			}
		}
	}()

	<-fin // wait for end of all coroutines

	return result
}

func RenderAsciiToString(n []int) string {
	result := ""
	for _, c := range n {
		result += string(c)
	}
	return result
}

func splitWorld(s string) [][]uint8 {
	result := make([][]uint8, 0)
	for _, line := range strings.Split(s, "\n") {
		result2 := make([]uint8, len(line))
		for i, c := range line {
			result2[i] = uint8(c)
		}
		result = append(result, result2)
	}
	return result
}

func findPosition(world [][]uint8, f func(x, y int, v uint8) bool) Point {
	for y := 0; y < len(world); y++ {
		for x := 0; x < len(world[y]); x++ {
			if f(x, y, world[y][x]) {
				p := Point{x, y}
				return p
			}
		}
	}
	panic("could not find point")
}

func findIntersections(world [][]uint8) []Point {
	result := make([]Point, 0)
	for y := 1; y < len(world)-1; y++ { // skip boundary
		for x := 1; x < len(world[y])-1; x++ { // skip boundary
			if world[y][x] != '#' {
				continue
			}
			p := Point{x, y}
			top := p.South()
			if world[top.Y][top.X] != '#' {
				continue
			}
			left := p.West()
			if world[left.Y][left.X] != '#' {
				continue
			}
			right := p.East()
			if world[right.Y][right.X] != '#' {
				continue
			}
			bottom := p.North()
			if world[bottom.Y][bottom.X] != '#' {
				continue
			}
			result = append(result, p)
		}
	}
	return result
}
