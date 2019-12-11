package main

import (
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 11
const AocDayName = "day11"
const AocDayTitle = "Day 11"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		robot := newRobot(program)
		robot.Run()
		fmt.Printf("\n%s\n", robot.PanelsAsString())
		dl.PrintSolution(fmt.Sprintf("%d panels has been painted (aka visited)", robot.CountPaintedPanels()))
	}

	{
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		robot := newRobot(program)
		robot.panels[Point{0, 0}] = &Panel{Value: 1, Visited: false}
		robot.Run()
		fmt.Printf("\n%s\n", flipTextBox(robot.PanelsAsString()))
		dl.PrintSolution(fmt.Sprintf("%d panels has been painted (aka visited)", robot.CountPaintedPanels()))
	}

}

func flipTextBox(str string) string {
	lines := strings.Split(str, "\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		j := len(lines) - i // reverse line index (flip horizontal)
		result[j-1] = line
	}
	return strings.Join(result, "\n")
}

type Point struct {
	X, Y int
}

type Robot struct {
	panels    map[Point]*Panel
	position  Point
	direction int32
	memory    []int
	stdin     chan int
	stdout    chan int
}

type Panel struct {
	Value   int
	Visited bool
}

func newRobot(program []int) Robot {
	panels := make(map[Point]*Panel)
	memory := make([]int, len(program))
	copy(memory, program)

	in := make(chan int, 1)
	out := make(chan int)

	position := Point{0, 0}

	return Robot{panels: panels, direction: '^', position: position, memory: memory, stdin: in, stdout: out}
}

func (r *Robot) Run() {

	wait := make(chan bool)
	halt := make(chan error)

	go func(memory []int, in <-chan int, out chan<- int, halt chan<- error) {
		halt <- day09.ExecutionInstructions(memory, in, out, false)
	}(r.memory, r.stdin, r.stdout, halt)

	go func() {
		for {

			select {
			case <-halt:
				wait <- true
				return
			default:
				if _, exist := r.panels[r.position]; !exist {
					r.panels[r.position] = &Panel{}
				}
				panel := r.panels[r.position]
				r.stdin <- panel.Value

				// wait
				paint := <-r.stdout
				nextDirection := <-r.stdout
				panel.Value = paint
				panel.Visited = true
				if nextDirection == 0 {
					r.turnLeft()
				} else if nextDirection == 1 {
					r.turnRight()
				} else {
					panic("invalid direction")
				}
			}
		}
	}()

	<-wait
}

func (r *Robot) PanelsAsString() string {
	result := ""
	bottomLeft := Point{0, 0}
	topRight := Point{0, 0}
	for k := range r.panels {
		bottomLeft.X = dl.MinInt(bottomLeft.X, k.X)
		bottomLeft.Y = dl.MinInt(bottomLeft.Y, k.Y)
		topRight.X = dl.MaxInt(topRight.X, k.X)
		topRight.Y = dl.MaxInt(topRight.Y, k.Y)
	}

	for y := bottomLeft.Y; y <= topRight.Y; y++ {
		for x := bottomLeft.X; x <= topRight.X; x++ {
			if panel, exist := r.panels[Point{x, y}]; exist && panel.Value == 1 {
				result += "#"
			} else {
				result += "."
			}
		}
		result += "\n"
	}
	return result
}

func (r *Robot) turnLeft() {
	switch r.direction {
	case '^':
		r.direction = '<'
	case '>':
		r.direction = '^'
	case 'v':
		r.direction = '>'
	case '<':
		r.direction = 'v'
	}
	r.forward()
}

func (r *Robot) turnRight() {
	switch r.direction {
	case '^':
		r.direction = '>'
	case '>':
		r.direction = 'v'
	case 'v':
		r.direction = '<'
	case '<':
		r.direction = '^'
	}
	r.forward()
}

func (r *Robot) forward() {
	switch r.direction {
	case '^':
		r.position = Point{r.position.X, r.position.Y + 1}
	case '>':
		r.position = Point{r.position.X + 1, r.position.Y}
	case 'v':
		r.position = Point{r.position.X, r.position.Y - 1}
	case '<':
		r.position = Point{r.position.X - 1, r.position.Y}
	}
}

func (r *Robot) CountPaintedPanels() int {
	result := 0
	for _, panel := range r.panels {
		if panel.Visited {
			result++
		}
	}
	return result
}
