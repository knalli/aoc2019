package main

import (
	day07 "de.knallisworld/aoc/aoc2019/day07/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"time"
)

const AocDay = 7
const AocDayName = "day07"
const AocDayTitle = "Day 07"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		max := findHighestSignal(puzzle, []int{0, 1, 2, 3, 4})
		dl.PrintSolution(fmt.Sprintf("The highest signal is: %d", max))
	}

	{
		dl.PrintStepHeader(2)
		puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		max := findHighestAmplifiedSignal(puzzle, []int{5, 6, 7, 8, 9})
		dl.PrintSolution(fmt.Sprintf("The highest signal using aplifier feedback is: %d", max))
	}

}

func findHighestSignal(puzzle []int, phaseSequence []int) int {
	max := 0
	for _, phaseSequences := range dl.Permutations(phaseSequence) {
		out := <-compute(puzzle, phaseSequences, 0)
		if out > max {
			max = out
		}
	}
	return max
}

func findHighestAmplifiedSignal(puzzle []int, phaseSequence []int) int {
	max := 0
	for _, phaseSequences := range dl.Permutations(phaseSequence) {
		out := computeAmplified(puzzle, phaseSequences, 0)
		if out > max {
			max = out
		}
	}
	return max
}

func compute(instructions []int, phaseSequence []int, input int) chan int {
	ios := make([]chan int, len(phaseSequence)+1)
	for i := range ios {
		ios[i] = make(chan int)
	}

	for i := range phaseSequence {
		in := ios[i]
		out := ios[i+1]
		go func() {
			data := make([]int, len(instructions))
			for i, v := range instructions {
				data[i] = v
			}
			halt := day07.ExecutionInstructions(data, in, out, false)
			if halt == nil {
				panic("Program not halted correctly")
			}
			close(out)
		}()
	}

	for i, phase := range phaseSequence {
		ios[i] <- phase
	}
	ios[0] <- input

	return ios[len(phaseSequence)]
}

func computeAmplified(instructions []int, phaseSequence []int, input int) int {
	ios := make([]chan int, len(phaseSequence)+1)
	for i := range ios {
		ios[i] = make(chan int, 10)
	}

	for i := range phaseSequence {
		go func(in <-chan int, out chan int, i int) {
			data := make([]int, len(instructions))
			for i, v := range instructions {
				data[i] = v
			}
			halt := day07.ExecutionInstructions(data, in, out, false)
			if halt == nil {
				panic("Program not halted correctly")
			}
			close(out)
		}(ios[i], ios[i+1], i)
	}

	results := make(chan int, 10)

	// Take everything of the last and put in the channels "first" and "copy"
	go func(channels []chan int, copy chan int) {
		first := ios[0]
		last := ios[len(channels)-1]
		day07.FanOut(last, first, copy)
		close(copy)
	}(ios, results)

	// Reduce "results" to the max value
	max := make(chan int)
	go day07.ReduceMax(results, max)

	// START
	go func() {
		for i, phase := range phaseSequence {
			ios[i] <- phase
		}
		ios[0] <- input
	}()

	return <-max
}
