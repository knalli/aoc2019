package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"errors"
	"fmt"
	"strings"
	"time"
)

const AocDay = 5
const AocDayName = "day05"
const AocDayTitle = "Day 05"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle := readFileAsIntArray(AocDayName + "/puzzle1.txt")
		output, halt := compute(puzzle, 1)
		dl.PrintSolution(fmt.Sprintf("Diagnostic code: %d (%s)", output, halt))
	}

	{
		dl.PrintStepHeader(2)
		puzzle := readFileAsIntArray(AocDayName + "/puzzle1.txt")
		output, halt := compute2(puzzle, 5, false)
		dl.PrintSolution(fmt.Sprintf("Diagnostic code: %d (%s)", output, halt))
	}

}

func readFileAsIntArray(file string) []int {
	str, _ := dl.ReadFileToString(file)
	arr := strings.Split(strings.TrimSpace(*str), ",")
	return dl.ParseStringToIntArray(arr)
}

func compute(data []int, input int) (int, error) {
	exchange := input
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
				if mode01 == 1 {
					param1 = data[i+1]
				} else {
					param1 = data[data[i+1]]
				}
				if mode02 == 1 {
					param2 = data[i+2]
				} else {
					param2 = data[data[i+2]]
				}
				data[data[i+3]] = param1 + param2
				i += 4
			}
		case 2:
			{
				var param1, param2 int
				if mode01 == 1 {
					param1 = data[i+1]
				} else {
					param1 = data[data[i+1]]
				}
				if mode02 == 1 {
					param2 = data[i+2]
				} else {
					param2 = data[data[i+2]]
				}
				data[data[i+3]] = param1 * param2
				i += 4
			}
		case 3:
			{
				data[data[i+1]] = exchange
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
				exchange = param1
				i += 2
			}
		case 99:
			return exchange, errors.New("got Halt And Catch Fire")
		}
	}
}

func compute2(data []int, input int, debug bool) (int, error) {
	exchange := input
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
				data[data[i+1]] = exchange
				fmt.Printf("[%05d] IN %d => #%d\n", i, exchange, data[i+1])
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
				exchange = param1
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
			return exchange, errors.New("got Halt And Catch Fire")
		default:
			panic(fmt.Sprintf("invalid op code %d (%d)", opcode, iOpcode))
		}
	}
}
