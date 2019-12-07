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
		highestSignal := findBestSignal(
			puzzle,
			[]int{0, 1, 2, 3, 4},
			evaluateAmplifiers,
			dl.MaxInt,
		)
		dl.PrintSolution(fmt.Sprintf("The highest signal is: %d", highestSignal))
	}

	{
		dl.PrintStepHeader(2)
		puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		highestSignal := findBestSignal(
			puzzle,
			[]int{5, 6, 7, 8, 9},
			evaluateAmplifiersInLoop,
			dl.MaxInt,
		)
		dl.PrintSolution(fmt.Sprintf("The highest signal using aplifier feedback is: %d", highestSignal))
	}

}

func findBestSignal(
	puzzle []int,
	phaseSequence []int, evaluator func(puzzle []int, sequence []int, input int) int,
	bestReducer func(a, b int) int,
) int {

	best := 0
	for _, sequence := range dl.Permutations(phaseSequence) {
		best = bestReducer(best, evaluator(puzzle, sequence, 0))
	}
	return best
}

func evaluateAmplifiers(instructions []int, phaseSequence []int, input int) int {
	// [0] -> A -> [1] -> B -> [2] -> … -> [n] -> [E] -> [n+1]
	// n+1 is the actual output
	pipeline := make([]chan int, len(phaseSequence)+1)
	for i := range pipeline {
		pipeline[i] = make(chan int)
	}

	for i := range phaseSequence {
		// [i] -> Amplifier -> [i+1]
		go func(in <-chan int, out chan<- int) {
			memory := make([]int, len(instructions))
			copy(memory, instructions)
			if ret := day07.ExecutionInstructions(memory, in, out, false); ret == nil {
				panic("Program not halted correctly")
			}
		}(pipeline[i], pipeline[i+1])
	}

	// START
	go func() {
		for i, phase := range phaseSequence {
			pipeline[i] <- phase
		}
		pipeline[0] <- input
	}()

	// output has one element
	return <-pipeline[len(pipeline)-1]
}

func evaluateAmplifiersInLoop(instructions []int, phaseSequence []int, input int) int {
	// [0] -> A -> [1] -> B -> [2] -> C -> [3] -> … -> [n] -> [E] -> [n+1]
	// [0] must be buffered
	pipeline := make([]chan int, len(phaseSequence)+1)
	for i := range pipeline {
		if i == 0 {
			pipeline[i] = make(chan int, 1)
		} else {
			pipeline[i] = make(chan int)
		}
	}
	// last pipeline's channel [n+1] will be copied into result later
	results := make(chan int, 1)

	for i := range phaseSequence {
		// [i] -> Amplifier -> [i+1]
		go func(in <-chan int, out chan<- int) {
			memory := make([]int, len(instructions))
			copy(memory, instructions)
			if ret := day07.ExecutionInstructions(memory, in, out, false); ret == nil {
				panic("Program not halted correctly")
			}
		}(pipeline[i], pipeline[i+1])
	}

	// Take everything of the last and put in the channels "first" and "copy"
	//   1) connects last amplifier's output with the first amplifier's input
	//   2) gets a result (channel) of the last amplifier's output separately
	go func(last <-chan int, first chan<- int, copy chan<- int) {
		day07.FanOut(last, first, copy)
		close(copy)
	}(pipeline[len(pipeline)-1], pipeline[0], results)

	// Reduce "results" to the max value
	max := make(chan int)
	go day07.ReduceIntChannel(results, max, dl.MaxIntArrayValue)

	// START
	go func() {
		for i, phase := range phaseSequence {
			pipeline[i] <- phase
		}
		pipeline[0] <- input
	}()

	return <-max
}
