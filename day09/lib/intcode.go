package day09

import (
	"errors"
	"fmt"
	"math"
)

// see also day07
func ExecutionInstructions(data2 []int, in <-chan int, out chan<- int, debug bool) error {
	i := 0
	relativeBase := 0
	data := make([]int, math.MaxInt16)
	copy(data, data2)
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
		mode03 := (iOpcode / 10000) % 10
		switch opcode {
		case 1:
			{
				var param1, param2 int
				var param1s, param2s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else if mode02 == 2 {
					param2 = data[relativeBase+data[i+2]]
					param2s = fmt.Sprintf("#%d", relativeBase+data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if mode03 == 2 {
					data[relativeBase+data[i+3]] = param1 + param2
				} else {
					data[data[i+3]] = param1 + param2
				}
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
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else if mode02 == 2 {
					param2 = data[relativeBase+data[i+2]]
					param2s = fmt.Sprintf("#%d", relativeBase+data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if mode03 == 2 {
					data[relativeBase+data[i+3]] = param1 * param2
				} else {
					data[data[i+3]] = param1 * param2
				}
				if debug {
					fmt.Printf("[%05d] MUL %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				i += 4
			}
		case 3:
			{
				next := <-in
				if mode01 == 2 {
					data[relativeBase+data[i+1]] = next
				} else {
					data[data[i+1]] = next
				}
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
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
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
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else if mode02 == 2 {
					param2 = data[relativeBase+data[i+2]]
					param2s = fmt.Sprintf("#%d", relativeBase+data[i+2])
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
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else if mode02 == 2 {
					param2 = data[relativeBase+data[i+2]]
					param2s = fmt.Sprintf("#%d", relativeBase+data[i+2])
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
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else if mode02 == 2 {
					param2 = data[relativeBase+data[i+2]]
					param2s = fmt.Sprintf("#%d", relativeBase+data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if debug {
					fmt.Printf("[%05d] SLT %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				if param1 < param2 {
					if mode03 == 2 {
						data[relativeBase+data[i+3]] = 1
					} else {
						data[data[i+3]] = 1
					}
				} else {
					if mode03 == 2 {
						data[relativeBase+data[i+3]] = 0
					} else {
						data[data[i+3]] = 0
					}
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
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				if mode02 == 1 {
					param2 = data[i+2]
					param2s = fmt.Sprintf("%d", data[i+2])
				} else if mode02 == 2 {
					param2 = data[relativeBase+data[i+2]]
					param2s = fmt.Sprintf("#%d", relativeBase+data[i+2])
				} else {
					param2 = data[data[i+2]]
					param2s = fmt.Sprintf("#%d", data[i+2])
				}
				if debug {
					fmt.Printf("[%05d] SIE %s %s => #%d\n", i, param1s, param2s, data[i+3])
				}
				if param1 == param2 {
					if mode03 == 2 {
						data[relativeBase+data[i+3]] = 1
					} else {
						data[data[i+3]] = 1
					}
				} else {
					if mode03 == 2 {
						data[relativeBase+data[i+3]] = 0
					} else {
						data[data[i+3]] = 0
					}
				}
				i += 4
			}
		case 9: // adjusts the relative base
			{
				var param1 int
				var param1s string
				if mode01 == 1 {
					param1 = data[i+1]
					param1s = fmt.Sprintf("%d", data[i+1])
				} else if mode01 == 2 {
					param1 = data[relativeBase+data[i+1]]
					param1s = fmt.Sprintf("#%d", relativeBase+data[i+1])
				} else {
					param1 = data[data[i+1]]
					param1s = fmt.Sprintf("#%d", data[i+1])
				}
				relativeBase += param1
				if debug {
					fmt.Printf("[%05d] REB %s\n", i, param1s)
				}
				i += 2
			}
		case 99:
			close(out)
			return errors.New("got Halt And Catch Fire")
		default:
			panic(fmt.Sprintf("invalid op code %d (%d)", opcode, iOpcode))
		}
	}
}
