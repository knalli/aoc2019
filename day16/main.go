package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"strings"
	"sync"
	"time"
)

const AocDay = 16
const AocDayName = "day16"
const AocDayTitle = "Day 16"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle, _ := dl.ReadFileToString(AocDayName + "/puzzle1.txt")
		result := fft(*puzzle, "0, 1, 0, -1", 100)
		dl.PrintSolution(fmt.Sprintf("The final result's first eight digits are '%s'", result[0:8]))
	}

	{
		dl.PrintStepHeader(2)
		puzzle, _ := dl.ReadFileToString(AocDayName + "/puzzle1.txt")
		var offset int
		{
			p := *puzzle
			offset = dl.ParseInt(p[0:7])
		}
		puzzle2 := strings.Repeat(*puzzle, 10000)
		dl.PrintSolution(fmt.Sprintf("The message offset is '%d'", offset))
		result := fftComplex(puzzle2[offset:], "0, 1, 0, -1", 100, fft0Summing, func(ints []int) string {
			return intArrayAsString(ints[0:8])
		})
		dl.PrintSolution(fmt.Sprintf("The final result's eight digits at offset '%d' are '%s'", offset, result))
	}

}

func fft(input string, basePattern string, phases int) string {
	inputInt := make([]int, len(input))
	for i, c := range input {
		inputInt[i] = int(c - '0')
	}
	basePatternInt := make([]int, 0)
	for _, part := range strings.Split(basePattern, ",") {
		basePatternInt = append(basePatternInt, dl.ParseInt(strings.TrimSpace(part)))
	}
	numbers := fft0Single(inputInt, basePatternInt, phases)

	result := ""
	for _, n := range numbers {
		result += fmt.Sprintf("%d", n)
	}
	return result
}

func fftComplex(input string, basePattern string, phases int, f func(input []int, basePattern []int, phases int) []int, r func([]int) string) string {
	inputInt := make([]int, len(input))
	for i, c := range input {
		inputInt[i] = int(c - '0')
	}
	basePatternInt := make([]int, 0)
	for _, part := range strings.Split(basePattern, ",") {
		basePatternInt = append(basePatternInt, dl.ParseInt(strings.TrimSpace(part)))
	}
	numbers := f(inputInt, basePatternInt, phases)

	return r(numbers)
}

func fft0Multi(input []int, basePattern []int, phases int) []int {

	n := len(input)
	result := make([]int, n)
	copy(result, input)

	workersMax := 256
	workerWorkload := n / workersMax

	worker := func(wg *sync.WaitGroup, offset int, limit int, output []int) {
		//fmt.Printf("Worker %09d-%09d [%d]\n", offset, offset+limit, limit)
		for o := offset; o < offset+limit; o++ {
			t := 0
			// fmt.Printf("Output %04d...\n", o)
			for i := 0; i < n; i++ {
				b := basePattern[((i+1)/(o+1))%4]
				t += result[i] * b
			}
			output[o] = dl.AbsInt(t) % 10
		}
		wg.Done()
	}

	for phase := 1; phase <= phases; phase++ {
		fmt.Printf("Phase %03d...\n", phase)
		output := make([]int, n)

		var wg sync.WaitGroup
		wg.Add(workersMax + 1)
		for core := 0; core < workersMax; core++ {
			go worker(&wg, core*workerWorkload, workerWorkload, output)
		}
		go worker(&wg, workersMax*workerWorkload, n-(workersMax*workerWorkload), output) // rest
		wg.Wait()

		copy(result, output)
	}
	return result
}

func intArrayAsString(arr []int) string {
	result := ""
	for _, n := range arr {
		result += fmt.Sprintf("%d", n)
	}
	return result
}

func fft0Single(input []int, basePattern []int, phases int) []int {

	n := len(input)
	result := make([]int, n)
	copy(result, input)

	for phase := 1; phase <= phases; phase++ {
		//fmt.Printf("Phase %03d...\n", phase)
		output := make([]int, n)

		for o := 0; o < n; o++ {
			t := 0
			for i := 0; i < n; i++ {
				b := basePattern[((i+1)/(o+1))%4]
				t += result[i] * b
			}
			output[o] = dl.AbsInt(t) % 10
		}
		copy(result, output)
	}
	return result
}

// summing up (does only work if offset is in the last 50% part of the input)
func fft0Summing(input []int, basePattern []int, phases int) []int {

	n := len(input)
	result := make([]int, n)
	copy(result, input)

	for phase := 1; phase <= phases; phase++ {
		//fmt.Printf("Phase %03d...\n", phase)

		for i := len(result) - 1; i > 0; i-- {
			result[i-1] += result[i]
			result[i-1] %= 10
		}
	}
	return result
}
