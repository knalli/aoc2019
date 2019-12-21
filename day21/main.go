package main

import (
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"time"
)

const AocDay = 21
const AocDayName = "day21"
const AocDayTitle = "Day 21"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		jumpInstructions := []string{
			"NOT A J",

			"NOT B T",
			"OR T J",

			"NOT C T",
			"OR T J",

			"AND D J",

			"WALK",
		}
		result := runProgram1(program, jumpInstructions, false)
		fmt.Printf(renderAsciiToString(result))
		reportedDamage := result[len(result)-1]
		dl.PrintSolution(fmt.Sprintf("Reported damage is %d", reportedDamage))
	}

	{
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		jumpInstructions := []string{
			"NOT A J",

			"NOT B T",
			"OR T J",

			"NOT C T",
			"OR T J",

			"AND D J",

			// Look ahead: after it jumped, is the next jump possible or not required
			// either D+4=H is a hole and D+1=E is ground, or D+4=H is ground
			"NOT H T",
			"AND E T",
			"OR H T",
			"AND T J",

			"RUN",
		}
		result := runProgram1(program, jumpInstructions, false)
		fmt.Printf(renderAsciiToString(result))
		reportedDamage := result[len(result)-1]
		dl.PrintSolution(fmt.Sprintf("Reported damage is %d", reportedDamage))
	}

}

func renderAsciiToString(n []int) string {
	result := ""
	for _, c := range n {
		result += string(c)
	}
	return result
}

func runProgram1(program []int, inputLines []string, debug bool) []int {

	result := make([]int, 0)

	in := make(chan int)     // program stdin
	out := make(chan int)    // program stdout
	fin := make(chan bool)   // program end
	halt := make(chan error) // program halt

	go func() {
		halt <- day09.ExecutionInstructions(program, in, out, debug)
	}()

	go func() {
		for _, line := range inputLines {
			for _, c := range line {
				in <- int(c)
			}
			in <- '\n'
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
