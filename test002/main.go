package main

import (
	"de.knallisworld/aoc/aoc2019/dayless"
	"de.knallisworld/aoc/aoc2019/test001/lib"
	"fmt"
	"time"
)

const AocDay = -2
const AocDayName = "test002"
const AocDayTitle = "Testing 2"

func main() {
	dayless.PrintDayHeader(AocDay, AocDayTitle)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	fmt.Println("Printing local puzzle")
	if s, err := dayless.ReadFileToString(AocDayName + "/puzzle1.txt"); err != nil {
		panic(err)
	} else {
		fmt.Println(*s)
	}
	fmt.Println()

	fmt.Println("Executing shared code")
	for _, s := range lib.TheDayOfTheTentacle() {
		fmt.Println(s)
	}
	fmt.Println()
}
