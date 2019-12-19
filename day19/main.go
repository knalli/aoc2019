package main

import (
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	day18 "de.knallisworld/aoc/aoc2019/day18/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"time"
)

const AocDay = 19
const AocDayName = "day19"
const AocDayTitle = "Day 19"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		grid := scan(program, 50, 50)
		fmt.Println(renderGrid(grid))
		dl.PrintSolution(fmt.Sprintf("There are %d points affected by the tractor beam.", countAffected(grid)))
	}

	{
		dl.PrintStepHeader(2)
		program := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		grid := scan(program, 100, 100)
		x, y := findBestPositionForSquare(program, grid, 100)
		dl.PrintSolution(fmt.Sprintf("Found possible point at %d/%d, result = %d", x, y, x*10000+y))
	}

}

func countAffected(grid [][]int) int {
	total := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == 1 {
				total++
			}
		}
	}
	return total
}

func run(program []int, x, y int) int {
	in := make(chan int, 2)  // program stdin
	out := make(chan int, 1) // program stdout
	in <- x
	in <- y
	_ = day09.ExecutionInstructions(program, in, out, false)
	o := <-out
	close(in)
	return o
}

func scan(program []int, width int, height int) [][]int {
	result := make([][]int, height)

	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		for x := 0; x < width; x++ {
			result[y][x] = run(program, x, y)
		}
	}

	return result
}

func renderGrid(grid [][]int) string {
	result := ""
	for y, line := range grid {
		if y != 0 {
			result += "\n"
		}
		for _, c := range line {
			if c == 0 {
				result += "."
			} else if c == 1 {
				result += "#"
			} else {
				result += "?"
			}
		}
	}
	return result
}

func findBestPositionForSquare(program []int, grid [][]int, dim int) (int, int) {
	height := len(grid)
	width := len(grid[0])
	tl := day18.Point{X: 0, Y: 0}
	hbr := day18.Point{X: width - 1, Y: 0}
	for y := 0; y < height; y++ {
		if grid[y][hbr.X] == 1 {
			hbr.Y = y
			break
		}
	}
	lbr := day18.Point{X: width - 1, Y: height - 1}
	for y := height - 1; y >= 0; y-- {
		if grid[y][lbr.X] == 1 {
			lbr.Y = y
			break
		}
	}
	fmt.Printf("%s -> %s / %s -> %s\n", tl.ToString(), hbr.ToString(), tl.ToString(), lbr.ToString())

	topLineFactor := float64(hbr.Y) / float64(hbr.X)
	bottomLineFactor := float64(lbr.Y) / float64(lbr.X)
	// matching the top/left point (because in the middle)
	centerLineFactor := topLineFactor + (bottomLineFactor-topLineFactor)/2

	foundX := dl.BinarySearch(lbr.X, 100000, func(x int) bool {
		y := int(centerLineFactor * (float64(x)))
		p0 := run(program, x, y)
		p1 := run(program, x+(dim-1), y)
		p2 := run(program, x, y+(dim-1))
		p3 := run(program, x+(dim-1), y+(dim-1))
		return p0 == 1 && p1 == 1 && p2 == 1 && p3 == 1
	}, false)
	foundY := int(centerLineFactor * (float64(foundX)))
	fmt.Printf("Found via aligned function: %d/%d\n", foundX, foundY)

	// fine tune
	corrected := false
	for y := foundY - 2*dim; y < foundY+dim; y++ {
		for x := foundX - 2*dim; x < foundX+dim; x++ {
			if run(program, x, y) != 1 {
				continue
			}
			if run(program, x+(dim-1), y) != 1 {
				continue
			}
			if run(program, x, y+(dim-1)) != 1 {
				continue
			}
			if run(program, x+(dim-1), y+(dim-1)) != 1 {
				continue
			}
			fmt.Printf("Corrected: %d/%d\n", x, y)
			foundY = y
			foundX = x
			corrected = true
			break
		}
		if corrected {
			break
		}
	}

	return foundX, foundY

}
