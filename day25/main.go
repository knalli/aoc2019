package main

import (
	"bufio"
	day09 "de.knallisworld/aoc/aoc2019/day09/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const AocDay = 25
const AocDayName = "day25"
const AocDayTitle = "Day 25"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	puzzle := dl.ReadFileAsIntArray(AocDayName + "/puzzle1.txt")

	if len(os.Args) > 0 && os.Args[1] == "manual" {
		solveManually(puzzle)
	} else {
		solveAutomatically(puzzle)
	}

}

func solveManually(puzzle []int) {
	dl.PrintStepHeader(1)
	// take klein bottle, candy cane, hologram, astrolabe
	// code is 134349952
	droid := newDroid(puzzle)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for line := range droid.out {
			println(line)
			if line == "Command?" {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				droid.in <- strings.TrimSpace(text)
			}
		}
		wg.Done()
	}()

	droid.Run()
	wg.Wait()
}

// TODO automatic
func solveAutomatically(puzzle []int) {
	dl.PrintStepHeader(1)

	tree := &DecisionNode{
		decisions: make(map[string]*DecisionNode),
	}

	var walk func(node *DecisionNode)
	walk = func(node *DecisionNode) {
		if node.parent == nil {
			run(puzzle, []string{}, node)
		} else {
			for _, option := range node.options {
				if _, exist := node.decisions[option]; !exist {
					cmds := make([]string, len(node.path))
					copy(cmds, node.path)
					for _, optionItem := range strings.Split(option, ",") {
						cmds = append(cmds, optionItem)
					}
					run(puzzle, cmds, node)
				}
			}
		}
		for _, child := range node.decisions {
			walk(child)
		}
	}

	walk(tree)
	dl.PrintSolution("Not solved yet")

	dl.PrintStepHeader(2)
	dl.PrintSolution("Not solved yet")
}

type DecisionNode struct {
	parent    *DecisionNode
	path      []string
	options   []string
	decisions map[string]*DecisionNode
	Final     bool
	Cyclic    bool
}

type Droid struct {
	program []int
	stopped bool
	in      chan string
	out     chan string
}

func run(program []int, cmds []string, node *DecisionNode) {
	droid := newDroid(program)

	dirs := make(map[string][]string, 4)
	{
		dirs["north"] = []string{"west", "north", "east", "south"}
		dirs["east"] = []string{"north", "east", "south", "west"}
		dirs["south"] = []string{"east", "south", "west", "north"}
		dirs["west"] = []string{"south", "west", "north", "east"}
	}

	unusedCommands := make([]string, len(cmds))
	copy(unusedCommands, cmds)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		doorParsing := false
		itemParsing := false
		var doorOptions []string
		var itemOptions []string
		var itemBlacklist string
		//var itemDropList []string

		visitedRooms := make(map[string]int)

		itemBlacklist += "giant electromagnet,escape pod,photons,molten lava,infinite loop,"

		lastDir := "north"
		for {
			select {
			case line := <-droid.out:
				fmt.Printf("<< %s\n", line)

				if len(line) == 0 {
					continue
				}

				if len(line) > 0 && line[0:3] == "== " {

					reg, _ := regexp.Compile("== (.*) ==")
					r := reg.FindStringSubmatch(line)
					roomName := r[1]

					visitedRooms[roomName]++
					if visitedRooms[roomName] > 2 {
						node.Final = true
						node.Cyclic = true
						return
					}

					doorOptions = []string{}
					itemOptions = []string{}
				}

				if line == "Doors here lead:" {
					doorParsing = true
					doorOptions = []string{}
				} else if doorParsing && len(line) > 0 && line[0:2] == "- " {
					doorOptions = append(doorOptions, line[2:])
				} else {
					doorParsing = false
				}
				if line == "Items here:" {
					itemParsing = true
					itemOptions = []string{}
				} else if itemParsing && len(line) > 0 && line[0:2] == "- " {
					item := line[2:]
					if !strings.Contains(itemBlacklist, item+",") {
						itemOptions = append(itemOptions, item)
					}
				} else {
					itemParsing = false
				}

				/*
					if strings.Contains(line, "You can't move!") {
						reg, _ := regexp.Compile(".*The (.*) is stuck to you\\..*")
						r := reg.FindStringSubmatch(line)
						itemDropList = append(itemDropList, r[1])
					}
				*/

				if strings.Contains(line, "ejected back to the") {
					switch lastDir {
					case "north":
						lastDir = "south"
					case "east":
						lastDir = "west"
					case "south":
						lastDir = "north"
					case "west":
						lastDir = "east"
					}
				}

				if line == "Unrecognized command." || line == "You can't go that way." || strings.Contains(line, "You can't move!") || strings.Contains(line, "You take the infinite loop") {
					return
				}

				if line == "Command?" {
					if len(unusedCommands) > 0 {
						cmd := unusedCommands[0]
						unusedCommands = unusedCommands[1:]
						droid.in <- cmd
						continue
					} else {
						if len(node.options) == 0 {
							options := make([]string, 0)
							for _, doorOption := range doorOptions {
								options = append(options, doorOption)
							}
							/*
								for _, itemOption := range itemOptions {
									for _, doorOption := range doorOptions {
										option := make([]string, 0)
										option = append(option, "take "+itemOption)
										option = append(option, doorOption)
										options = append(options, strings.Join(option, ","))
									}
								}
							*/
							// TODO more item permutations?
							node.options = options
						}

						for _, option := range node.options {
							if _, exist := node.decisions[option]; !exist {

								nextCmds := make([]string, len(cmds))
								copy(nextCmds, cmds)
								nextCmds = append(nextCmds, option)

								next := &DecisionNode{
									parent:    node,
									path:      nextCmds,
									options:   nil,
									decisions: make(map[string]*DecisionNode),
								}
								node.decisions[option] = next
								node = next
								droid.in <- option
								break
							}
						}

					}

					/*
						if len(itemDropList) > 0 {
							item := itemDropList[0]
							itemDropList = itemDropList[1:]
							itemBlacklist += item + ","
							droid.in <- "drop " + item
						} else if line == "Unrecognized command." || line == "You can't go that way." || strings.Contains(line, "You can't move!") {
							return
						} else {
							if len(itemOptions) > 0 {
								droid.in <- "take " + itemOptions[0]
								itemOptions = itemOptions[1:]
							} else {
								for _, dir := range dirs[lastDir] {
									if strings.Contains(strings.Join(doorOptions, ","), dir) {
										lastDir = dir
										droid.in <- dir
										break
									}
								}
							}
						}
					*/
				}
			}
		}
		wg.Done()
	}()
	droid.Run()
	wg.Wait()
}

func (d *Droid) Run() {
	in := make(chan int, 100)
	out := make(chan int)
	halt := make(chan error)
	go func(in <-chan int, out chan<- int) {
		halt <- day09.ExecutionInstructions(d.program, in, out, false)
	}(in, out)

	go func() {
		for {
			select {
			case line := <-d.in:
				fmt.Printf(">> %s\n", color.New(color.FgRed).Sprint(line))
				// send command (ASCII to int)
				for _, c := range line {
					//fmt.Printf("%d\n", c)
					in <- int(c)
				}
				in <- '\n'
			}
		}
	}()

	var buffer []int
	go func() {
		for {
			select {
			case <-halt:
				close(halt)
				close(in)
				close(d.out)
				close(d.in)
				d.stopped = true
				return
			case c := <-out:
				//fmt.Printf("Receiving byte... %d\n", c)
				// receive command/output (in to ASCII)
				if c == '\n' {
					str := renderAsciiToString(buffer)
					//fmt.Printf("Receiving buffered line: %s\n", str)
					d.out <- str
					buffer = []int{}
				} else {
					buffer = append(buffer, c)
				}
			default:
			}
		}
	}()
}

func newDroid(program []int) *Droid {
	return &Droid{
		program: program,
		in:      make(chan string, 10),
		out:     make(chan string),
	}
}
func renderAsciiToString(n []int) string {
	result := ""
	for _, c := range n {
		result += string(c)
	}
	return result
}
