package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"time"
)

const AocDay = 1
const AocDayName = "day01"
const AocDayTitle = "Day 01"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	dl.PrintSolution("Not solved yet")

	dl.PrintStepHeader(2)
	dl.PrintSolution("Not solved yet")

}
