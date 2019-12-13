package main

import (
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"time"
)

const AocDay = 13
const AocDayName = "day13"
const AocDayTitle = "Day 13"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		result, _ := startGame(puzzle, false)
		blocksTotal := 0
		for _, tile := range result {
			if tile == 2 {
				blocksTotal++
			}
		}
		fmt.Printf("Game: \n%s\n", printGame(result))
		dl.PrintSolution(fmt.Sprintf("Number of block tiles is %d", blocksTotal))
	}

	{
		dl.PrintStepHeader(2)
		puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")
		puzzle[0] = 2 // play free
		result, score := startGame(puzzle, true)
		fmt.Printf("Game: \n%s\n", printGame(result))
		dl.PrintSolution(fmt.Sprintf("Last score is %d", score))
	}

}

type Point struct {
	X, Y int
}

func startGame(program []int, playingMode bool) (map[Point]int, int) {

	data := make(map[Point]int, 0)

	in := make(chan int)     // program stdin
	out := make(chan int)    // program stdout
	fin := make(chan bool)   // game end
	halt := make(chan error) // program halt
	go func() {
		halt <- day09.ExecutionInstructions(program, in, out, false)
	}()

	paddlePosition := Point{0, 0}
	ballPosition := Point{0, 0}
	lastScore := 0
	turn := make(chan bool) // signal players turn
	go func() {
		for {
			select {
			case <-halt:
				fin <- true
				return
			default:
				x := <-out // x
				y := <-out // y
				v := <-out // value/p(l)ayload
				if x == -1 && y == 0 {
					//fmt.Printf("Current score: %d", v)
					lastScore = v
				} else {
					// v is tileId
					p := Point{x, y}
					data[p] = v
					if v == 3 {
						paddlePosition.X = x
						paddlePosition.Y = y
					} else if v == 4 {
						ballPosition.X = x
						ballPosition.Y = y
						if playingMode {
							turn <- true
						}
					}
				}
			}
		}
	}()

	// Control paddle (react on ball)
	if playingMode {
		go func() {
			for {
				select {
				case <-turn:
					if ballPosition.X < paddlePosition.X {
						in <- -1
					} else if ballPosition.X > paddlePosition.X {
						in <- 1
					} else {
						in <- 0
					}
				default:
					continue
				}
			}
		}()
	}

	<-fin // wait for end of all coroutines
	return data, lastScore
}

func printGame(data map[Point]int) string {
	maxX := 0
	maxY := 0
	for p := range data {
		maxX = dl.MaxInt(maxX, p.X)
		maxY = dl.MaxInt(maxY, p.Y)
	}

	result := ""
	for y := maxY; y >= 0; y-- {
		for x := 0; x <= maxX; x++ {
			p := Point{x, y}
			if v, exist := data[p]; exist {
				switch v {
				case 1:
					result += "█"
				case 2:
					result += "░"
				case 3:
					result += "_"
				case 4:
					result += "O"
				default:
					result += " "
				}
			} else {
				result += " "
			}
		}
		result += "\n"
	}
	return result
}
