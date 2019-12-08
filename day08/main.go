package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"math"
	"time"
)

const AocDay = 8
const AocDayName = "day08"
const AocDayTitle = "Day 08"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(0)
		puzzle, _ := dl.ReadFileToString(AocDayName + "/sample1.txt")
		layers := buildLayers(*puzzle, 3, 2)
		min := math.MaxInt64
		layerIdx := -1
		for i, layer := range layers {
			c := layer.CountDigits(0)
			if c < min {
				min = c
				layerIdx = i
			}
		}
		result := layers[layerIdx].CountDigits(1) * layers[layerIdx].CountDigits(2)
		dl.PrintSolution(fmt.Sprintf("Layer #%d has the fewest 0 digits. The result of 1 digits multiplied with 2 digits: %d", layerIdx, result))
	}

	{
		dl.PrintStepHeader(0)
		puzzle, _ := dl.ReadFileToString(AocDayName + "/sample2.txt")
		layers := buildLayers(*puzzle, 2, 2)
		min := math.MaxInt64
		layerIdx := -1
		for i, layer := range layers {
			c := layer.CountDigits(0)
			if c < min {
				min = c
				layerIdx = i
			}
		}
		result := layers[layerIdx].CountDigits(1) * layers[layerIdx].CountDigits(2)
		dl.PrintSolution(fmt.Sprintf("Layer #%d has the fewest 0 digits. The result of 1 digits multiplied with 2 digits: %d", layerIdx, result))
	}

	{
		dl.PrintStepHeader(1)
		puzzle, _ := dl.ReadFileToString(AocDayName + "/puzzle1.txt")
		layers := buildLayers(*puzzle, 25, 6)
		min := math.MaxInt64
		layerIdx := -1
		for i, layer := range layers {
			c := layer.CountDigits(0)
			if c < min {
				min = c
				layerIdx = i
			}
		}
		result := layers[layerIdx].CountDigits(1) * layers[layerIdx].CountDigits(2)
		dl.PrintSolution(fmt.Sprintf("Layer #%d has the fewest 0 digits. The result of 1 digits multiplied with 2 digits: %d", layerIdx, result))
	}

	{
		dl.PrintStepHeader(2)
		puzzle, _ := dl.ReadFileToString(AocDayName + "/puzzle1.txt")
		layers := buildLayers(*puzzle, 25, 6)
		final := reduceLayers(layers, func(v int) bool {
			return v != 2
		})
		dl.PrintSolution(fmt.Sprintf("Final message: \n%s", final.ToString(25)))
	}

}

func buildLayers(line string, width int, height int) []layer {
	layerLength := width * height
	layerCount := len(line) / layerLength

	result := make([]layer, layerCount)
	for i := 0; i < layerCount; i++ {
		offset := i * layerLength
		sub := line[offset : offset+layerLength]
		result[i] = make(layer, layerLength)
		for j, c := range sub {
			result[i][j] = int(c - 48)
		}
	}

	return result
}

func reduceLayers(layers []layer, filter func(v int) bool) layer {
	result := make(layer, len(layers[0]))
	for i := 0; i < len(result); i++ {
		found := false
		for _, l := range layers {
			v := l[i]
			if filter(v) {
				result[i] = v
				found = true
				break
			}
		}
		// fallback
		if !found {
			result[i] = layers[len(layers)-1][i]
		}
	}
	return result
}

type layer []int

func (l layer) CountDigits(d int) int {
	r := 0
	for _, n := range l {
		if n == d {
			r++
		}
	}
	return r
}

func (l layer) ToString(width int) string {
	result := ""
	for i, n := range l {
		if i%width == 0 {
			result += "\n"
		}
		if n == 1 {
			result += "█"
		} else {
			result += " "
		}
	}
	return result
}
