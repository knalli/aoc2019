package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"errors"
	"fmt"
	"strings"
	"time"
)

const AocDay = 2
const AocDayName = "day02"
const AocDayTitle = "Day 02"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	input := readFileAsIntArray(AocDayName + "/puzzle1.txt")
	//compute([]int{1,0,0,0,99})
	//compute([]int{2,3,0,3,99})
	// 1202 program alert avoid
	input[1] = 12
	input[2] = 2
	err := compute(input)
	if err != nil {
		fmt.Printf("🔥 %s\n", err.Error())
	}
	dl.PrintSolution(fmt.Sprintf("Solution 1: Program position #0 = %d", input[0]))

	err = nil
	dl.PrintStepHeader(2)
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			local := readFileAsIntArray(AocDayName + "/puzzle1.txt")
			local[1] = noun
			local[2] = verb
			err = compute(local)
			if err != nil && local[0] == 19690720 {
				dl.PrintSolution(fmt.Sprintf("Solution 2: noun=%d, verb=%d => 100*noun+verb = %d", noun, verb, 100*noun+verb))
				break
			}
			err = nil
		}
		if err != nil {
			break
		}
	}

}

func readFileAsIntArray(file string) []int {
	str, _ := dl.ReadFileToString(file)
	arr := strings.Split(*str, ",")
	return dl.ParseStringToIntArray(arr)
}

func compute(data []int) error {
	i := 0
	for {
		switch data[i] {
		case 1:
			data[data[i+3]] = data[data[i+1]] + data[data[i+2]]
			i += 4
			break
		case 2:
			data[data[i+3]] = data[data[i+1]] * data[data[i+2]]
			i += 4
			break
		case 99:
			return errors.New("got Halt And Catch Fire")
		}
	}
}
