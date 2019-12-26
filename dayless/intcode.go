package dayless

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func ExecuteIntcode(program []int, in <-chan int, out chan<- int, signal <-chan int, debug bool) <-chan error {
	ip := 0
	relativeBase := 0
	memory := make([]int, MaxInt(math.MaxInt16, 2*len(program)))
	copy(memory, program)

	err := make(chan error, 1)

	stopped := false
	isStopped := func() bool {
		if !stopped {
			select {
			case s := <-signal:
				err <- errors.New(fmt.Sprintf("killed by signal #%d", s))
				stopped = true
			default:
			}
		}
		return stopped
	}

	for !isStopped() {
		instruction := fmt.Sprintf("%05d", memory[ip])
		opcode, _ := strconv.Atoi(instruction[3:])
		arg := func(i int) int {
			switch instruction[3-i] {
			case '1': // immediate mode (its value)
				return ip + i
			case '2': // relative mode
				return relativeBase + memory[ip+i]
			default: // 1, position mode
				return memory[ip+i]
			}
		}
		debugArg := func(i int) string {
			switch instruction[3-i] {
			case '1':
				return fmt.Sprintf("%d", ip+i)
			case '2':
				return fmt.Sprintf("#%d+%d", relativeBase, memory[ip+i])
			default:
				return fmt.Sprintf("#%d", memory[ip+i])
			}
		}
		switch opcode {
		case 1:
			if debug {
				fmt.Printf("[%s] ADD %s %s -> %s\n", instruction, debugArg(1), debugArg(2), debugArg(3))
			}
			memory[arg(3)] = memory[arg(1)] + memory[arg(2)]
			ip += 4
		case 2:
			if debug {
				fmt.Printf("[%s] MUL %s %s -> %s\n", instruction, debugArg(1), debugArg(2), debugArg(3))
			}
			memory[arg(3)] = memory[arg(1)] * memory[arg(2)]
			ip += 4
		case 3:
			if debug {
				fmt.Printf("[%s] REA STDIN -> %s\n", instruction, debugArg(1))
			}
			select {
			case s := <-signal:
				err <- errors.New(fmt.Sprintf("killed by signal #%d", s))
				stopped = true
				break
			case v := <-in:
				memory[arg(1)] = v
				ip += 2
			}
		case 4:
			if debug {
				fmt.Printf("[%s] WRT %s -> STDOUT\n", instruction, debugArg(1))
			}
			out <- memory[arg(1)]
			ip += 2
		case 5: // jump-if-true
			if debug {
				fmt.Printf("[%s] JIT %s != 0 -> %s\n", instruction, debugArg(1), debugArg(2))
			}
			if memory[arg(1)] != 0 {
				ip = memory[arg(2)]
			} else {
				ip += 3
			}
		case 6: // jump-if-false
			if debug {
				fmt.Printf("[%s] JIF %s == 0 -> %s\n", instruction, debugArg(1), debugArg(2))
			}
			if memory[arg(1)] == 0 {
				ip = memory[arg(2)]
			} else {
				ip += 3
			}
		case 7: // less-than
			if debug {
				fmt.Printf("[%s] JLT %s < %s -> %s\n", instruction, debugArg(1), debugArg(2), debugArg(3))
			}
			if memory[arg(1)] < memory[arg(2)] {
				memory[arg(3)] = 1
			} else {
				memory[arg(3)] = 0
			}
			ip += 4
		case 8: // equals
			if debug {
				fmt.Printf("[%s] JEQ %s == %s -> %s\n", instruction, debugArg(1), debugArg(2), debugArg(3))
			}
			if memory[arg(1)] == memory[arg(2)] {
				memory[arg(3)] = 1
			} else {
				memory[arg(3)] = 0
			}
			ip += 4
		case 9: // adjusts the relative base
			if debug {
				fmt.Printf("[%s] ARB %s -> RB\n", instruction, debugArg(1))
			}
			relativeBase += memory[arg(1)]
			ip += 2
		case 99:
			if debug {
				fmt.Printf("[%s] HLT\n", instruction)
			}
			stopped = true
			break
		default:
			err <- errors.New(fmt.Sprintf("invalid op code %d (%s)", opcode, instruction))
			break
		}
	}

	close(err)
	close(out)

	return err
}
