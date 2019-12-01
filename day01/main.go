package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"time"
)

const AocDay = 1
const AocDayName = "day01"
const AocDayTitle = "Day 01"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
	dl.PrintSolution(fmt.Sprintf("The sum of the fuel requirements: %d", computeTotalFuel(dl.ParseStringToIntArray(lines), false)))

	dl.PrintStepHeader(2)
	dl.PrintSolution(fmt.Sprintf("The sum of the fuel requirements: %d", computeTotalFuel(dl.ParseStringToIntArray(lines), true)))

}

func computeFuelByModule(mod int, deep bool) int {
	if !deep {
		// solution 1
		return mod/3 - 2
	} else {
		// solution 2
		fuel := mod/3 - 2
		if fuel > 0 {
			return fuel + computeFuelByModule(fuel, true)
		} else {
			return 0
		}
	}
}

func computeTotalFuel(mods []int, deep bool) int {
	res := 0
	for _, mod := range mods {
		res += computeFuelByModule(mod, deep)
	}
	return res
}
