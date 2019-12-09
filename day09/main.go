package main

import (
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"time"
)

const AocDay = 9
const AocDayName = "day09"
const AocDayTitle = "Day 09"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		results := runBoostProgram(program, 1)
		dl.PrintSolution(fmt.Sprintf("BOOST keycode: %d length", len(results)))
		for _, r := range results {
			fmt.Printf("%d \n", r)
		}
	}

	{
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		results := runBoostProgram(program, 2)
		dl.PrintSolution(fmt.Sprintf("BOOST keycode: %d length", len(results)))
		for _, r := range results {
			fmt.Printf("%d \n", r)
		}
	}

}

func runBoostProgram(program []int, base int) []int {
	memory := make([]int, len(program))
	copy(memory, program)
	in := make(chan int, 1)
	out := make(chan int)
	in <- base
	go func(in <-chan int, out chan<- int) {
		_ = day09.ExecutionInstructions(memory, in, out, false)
	}(in, out)

	result := make(chan []int)
	go func(out chan int, result chan []int) {
		outputs := make([]int, 0)
		for n := range out {
			outputs = append(outputs, n)
		}
		result <- outputs
	}(out, result)
	return <-result
}
