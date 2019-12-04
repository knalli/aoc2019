package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 4
const AocDayName = "day04"
const AocDayTitle = "Day 04"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	dl.PrintStepHeader(1)
	{
		puzzle, _ := dl.ReadFileToString(AocDayName + "/puzzle1.txt")
		parts := strings.Split(*puzzle, "-")
		total := countPasswords(dl.ParseInt(parts[0]), dl.ParseInt(parts[1]), func(n []int) bool {
			return hasSameAdjacents(n, 2, false) && hasIncreasingDigits(n)
		})
		dl.PrintSolution(fmt.Sprintf("There are %d different passwords", total))
	}

	dl.PrintStepHeader(2)
	{
		puzzle, _ := dl.ReadFileToString(AocDayName + "/puzzle1.txt")
		parts := strings.Split(*puzzle, "-")
		total := countPasswords(dl.ParseInt(parts[0]), dl.ParseInt(parts[1]), func(n []int) bool {
			return hasSameAdjacents(n, 2, true) && hasIncreasingDigits(n)
		})
		dl.PrintSolution(fmt.Sprintf("There are %d different passwords", total))
	}

}

func countPasswords(from, to int, predicate func([]int) bool) int {
	items := make([]int, 0)
	for i := from; i <= to; i++ {
		n0 := i % 10
		n1 := (i % 100) / 10
		n2 := (i % 1000) / 100
		n3 := (i % 10000) / 1000
		n4 := (i % 100000) / 10000
		n5 := (i % 1000000) / 100000
		n := []int{n5, n4, n3, n2, n1, n0}
		if predicate(n) {
			items = append(items, i)
		}
	}
	items = unique(items)
	return len(items)
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func hasSameAdjacents(numbers []int, required int, exact bool) bool {
	for i := 0; i < 10; i++ {
		strike := 0
		highest := 1
		for _, n := range numbers {
			if n == i {
				strike++
			} else {
				if highest < strike {
					highest = strike
				}
				strike = 0
			}
		}
		if strike > 0 && highest < strike {
			highest = strike
		}
		if exact {
			if highest == required {
				return true
			}
		} else {
			if highest >= required {
				return true
			}
		}
	}
	return false
}

func hasIncreasingDigits(numbers []int) bool {
	min := 1
	for _, n := range numbers {
		if n < min {
			return false
		}
		min = n
	}
	return true
}
