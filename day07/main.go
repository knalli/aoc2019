package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"errors"
	"fmt"
	"strings"
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
		puzzle := readFileAsIntArray(AocDayName + "/puzzle1.txt")
		max := findHighestSignal(puzzle, []int{0, 1, 2, 3, 4})
		dl.PrintSolution(fmt.Sprintf("The highest signal is: %d", max))
	}

	dl.PrintStepHeader(2)
	dl.PrintSolution("Not solved yet")

}

func findHighestSignal(puzzle []int, phaseSequence []int) int {
	max := 0
	for _, phaseSequences := range permutations(phaseSequence) {
		out := <-compute(puzzle, phaseSequences, 0)
		if out > max {
			max = out
		}
	}
	return max
}

func permutations(arr []int) [][]int {
	result := make([][]int, 0)
	for p := make([]int, len(arr)); p[0] < len(p); nextPerm(p) {
		result = append(result, getPerm(arr, p))
	}
	return result
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

func readFileAsIntArray(file string) []int {
	str, _ := dl.ReadFileToString(file)
	arr := strings.Split(strings.TrimSpace(*str), ",")
	return dl.ParseStringToIntArray(arr)
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
			halt := executionInstructions(data, in, out, false)
			if halt == nil {
				panic("Program not halted correctly")
			}
			close(in)
		}()
	}

	for i, phase := range phaseSequence {
		ios[i] <- phase
	}
	ios[0] <- input

	return ios[len(phaseSequence)]
}

// copy of day5/compute2
func executionInstructions(data []int, in <-chan int, out chan int, debug bool) error {
	i := 0
	for {
		iOpcode := data[i]
		opcode := iOpcode
		if opcode > 10000 {
			opcode %= 10000
		}
		if opcode > 1000 {
			opcode %= 1000
		}
		if opcode > 100 {
			opcode %= 100
		}
		mode01 := (iOpcode / 100) % 10
		mode02 := (iOpcode / 1000) % 10
		//mode03 := (iOpcode / 10000) % 10
		switch opcode {
		case 1:
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				data[data[i+3]] = param1 + param2
				if debug {
					fmt.Printf("[%05d] ADD %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				i += 4
			}
		case 2:
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				data[data[i+3]] = param1 * param2
				if debug {
					fmt.Printf("[%05d] MUL %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				i += 4
			}
		case 3:
			{
				next := <-in
				data[data[i+1]] = next
				if debug {
					fmt.Printf("[%05d] IN %d => #%d\n", i, next, data[i+1])
				}
				i += 2
			}
		case 4:
			{
				var param1 int
				if mode01 == 1 {
					param1 = data[i+1]
				} else {
					param1 = data[data[i+1]]
				}
				out <- param1
				if debug {
					fmt.Printf("[%05d] OUT #%d\n", i, data[i+1])
				}
				i += 2
			}
		case 5: // jump-if-true
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if debug {
					fmt.Printf("[%05d] JIT %s %s\n", i, param1s, param2s)
				}
				if param1 != 0 {
					i = param2
				} else {
					i += 3
				}
			}
		case 6: // jump-if-false
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if debug {
					fmt.Printf("[%05d] JIF %s %s\n", i, param1s, param2s)
				}
				if param1 == 0 {
					i = param2
				} else {
					i += 3
				}
			}
		case 7: // less-than
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if debug {
					fmt.Printf("[%05d] SLT %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				if param1 < param2 {
					data[data[i+3]] = 1
				} else {
					data[data[i+3]] = 0
				}
				i += 4
			}
		case 8: // equals
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if debug {
					fmt.Printf("[%05d] SIE %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				if param1 == param2 {
					data[data[i+3]] = 1
				} else {
					data[data[i+3]] = 0
				}
				i += 4
			}
		case 99:
			return errors.New("got Halt And Catch Fire")
		default:
			panic(fmt.Sprintf("invalid op code %d (%d)", opcode, iOpcode))
		}
	}
}
