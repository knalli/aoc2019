package main

import (
	"container/list"
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	tm "github.com/buger/goterm"
	. "github.com/logrusorgru/aurora"
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
		videoStreaming := true
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		world := RenderAsciiToString(runProgram1(program, false))
		fmt.Printf("%s\n", world)
		world = world[:len(world)-3] // remove last linebreak
		movements := findMovements(splitWorld(world))
		fmt.Printf("Movements: %s\n\n", movements)
		main, _, functions := DeflateString(movements, ",", 3, 20)
		fmt.Printf("Main routine: %s\n", main)
		fmt.Printf("Function A: %s\n", functions[0])
		fmt.Printf("Function C: %s\n", functions[1])
		fmt.Printf("Function B: %s\n", functions[2])

		inputs := make([]int, 0)
		for _, c := range main {
			inputs = append(inputs, int(c))
		}
		inputs = append(inputs, '\n')
		for _, group := range functions {
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

		if videoStreaming {
			for i := 0; i < len(strings.Split(movements, ","))/2; i++ {
				inputs = append(inputs, 'y')
				inputs = append(inputs, '\n')
			}
		} else {
			inputs = append(inputs, 'n')
			inputs = append(inputs, '\n')
		}

		//println("Instructions")
		//println(RenderAsciiToString(inputs))

		fmt.Println("Run program...")
		program[0] = 2
		outputChan := make(chan string)
		wait := make(chan bool)
		if videoStreaming {
			go func() {
				tm.Clear()
				received := 0
				//noinspection ALL
				tm.Printf("Live video feed (syncing):\n")
				for line := range outputChan {
					if line == "" {
						received++
						time.Sleep(60 * time.Millisecond)
						tm.MoveCursor(1, 1)
						//noinspection ALL
						tm.Printf("Live video feed (received %03d images):\n", received)
					}
					//noinspection ALL
					tm.Println(colorize(line))
					tm.Flush()
				}
				//noinspection ALL
				tm.Println("\nLive video feed ended")
				wait <- true
			}()
		}
		output := runProgram2(program, inputs, outputChan, videoStreaming, false)
		if videoStreaming {
			<-wait
		}
		dust := dl.MaxInt(output[len(output)-2], output[len(output)-1]) // workaround (last output is sometimes 0)
		dl.PrintSolution(fmt.Sprintf("Collected dust is %d", dust))
	}

}

func colorize(line string) string {
	line = strings.Replace(line, "#", fmt.Sprintf("%s", Yellow("#")), -1)
	line = strings.Replace(line, ".", fmt.Sprintf("%s", Black(".")), -1)
	for _, s := range []string{"^", ">", "v", "<"} {
		line = strings.Replace(line, s, fmt.Sprintf("%s", Red(Bold(s))), -1)
	}
	return line
}

func DeflateString(line string, separator string, blockNum int, blockSize int) (string, []int, []string) {

	l := len(line)

	type Pair struct {
		Length  int
		Data    string
		Filling string
	}

	pairs := make([]Pair, blockNum)
	for i := 0; i < blockNum; i++ {
		pairs[i] = Pair{0, "", ""}
	}

	if blockNum >= 'L' {
		// for this specific solution enough, but 'L' will break the "replace-approach" b/c dir movements
		// alternative: store pair ranges
		panic("specified blockNum not supported")
	}

	blockFillingKeys := ""
	for i := 0; i < len(pairs); i++ {
		name := string('A' + i)
		blockFillingKeys += name
		pairs[i].Filling = name
	}

	findNextFree := func(line string) int {
		parts := strings.Split(line, separator)
		for i, part := range parts {
			if !strings.ContainsAny(part, blockFillingKeys) {
				return i
			}
		}
		return -1
	}

	isBlockColliding := func(begin int, end int, block string) bool {
		return strings.ContainsAny(block, blockFillingKeys)
	}

	sumLengths := func() int {
		total := 0
		for _, pair := range pairs {
			total += pair.Length
		}
		return total
	}

	for {

		pairs[0].Length++
		for i := 0; i < len(pairs)-1; i++ {
			if sumLengths() > l {
				pairs[i].Length = 1
				pairs[i+1].Length++
			}
		}
		if sumLengths() > l {
			break
		}

		{
			valid := true
			for i := 0; i < len(pairs); i++ {
				valid = valid && pairs[i].Length > 0
			}
			if !valid {
				continue
			}
		}
		//fmt.Printf("%02d %02d %02d\n", a, b, c)

		temp := "" + line // copy/reset
		for i := range pairs {
			pairs[i].Data = ""
		}

		validBlocks := false
		for i := 0; i < blockNum; i++ {
			pair := &pairs[i]
			parts := strings.Split(temp, separator)

			blockBegin := findNextFree(temp)
			if blockBegin < 0 {
				continue
			}
			blockEnd := blockBegin + 2*pair.Length
			if blockEnd > len(parts) {
				continue
			}

			sub := strings.Join(parts[blockBegin:blockEnd], separator)
			if isBlockColliding(blockBegin, blockEnd, sub) || len(sub) > blockSize {
				continue
			}

			pair.Data = sub
			temp = strings.Replace(temp, sub, pair.Filling, -1)
			validBlocks = true
		}
		if !validBlocks {
			continue
		}

		if -1 == findNextFree(temp) {
			if len(temp) > blockSize {
				// mt.Printf("Found main function, but exceeding command length: %s\n", temp)
				continue
			}
			blockLengths := make([]int, blockNum)
			blockData := make([]string, blockNum)
			for i := 0; i < blockNum; i++ {
				blockLengths[i] = pairs[i].Length
				blockData[i] = pairs[i].Data
			}
			return temp, blockLengths, blockData
		}
	}

	return line, []int{0, 0, 0}, nil
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

func runProgram2(program []int, inputs []int, output chan string, videoStreaming bool, debug bool) []int {

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
				close(output)
				return
			default:
				b := <-out
				if videoStreaming && b == '\n' {
					output <- RenderAsciiToString(result)
					result = make([]int, 0)
				} else {
					result = append(result, b)
				}
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
