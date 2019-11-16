package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"de.knallisworld/aoc/aoc2019/test001/lib"
	"errors"
	"fmt"
	"github.com/yourbasic/graph"
	"time"
)

const AocDay = -1
const AocDayName = "test001"
const AocDayTitle = "Testing 1"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	fmt.Println("Printing local puzzle")
	if s, err := dl.ReadFileToString(AocDayName + "/puzzle1.txt"); err != nil {
		panic(err)
	} else {
		fmt.Println(*s)
	}
	dl.PrintSolution("123")

	g := graph.New(5)
	g.Add(0, 1)
	g.Add(1, 2)
	g.Add(2, 3)
	graph.BFS(g, 0, func(v, w int, c int64) {
		fmt.Printf("%d -> %d %d\n", v, w, c)
	})

	dl.PrintStepHeader(2)
	fmt.Println("Executing shared code")
	for _, s := range lib.TheDayOfTheTentacle() {
		fmt.Println(s)
	}
	dl.PrintSolution("42")

	dl.PrintError(errors.New("NullPointerException :)"))

}
