package main

import (
	"bufio"
	"container/list"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"github.com/fatih/color"
	"math"
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

	if len(os.Args) > 1 && os.Args[1] == "manual" {
		solveManually(puzzle)
	} else {
		solveAutomatically(puzzle)
	}

}

func solveManually(puzzle []int) {
	dl.PrintStepHeader(1)
	// take klein bottle, candy cane, hologram, astrolabe
	// code is 134349952
	droid := newDroid(puzzle, true)

	go func() {
		for line := range droid.out {
			println(line)
			if line == "Command?" {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				droid.in <- strings.TrimSpace(text)
			}
		}
	}()

	droid.Start()
}

func solveAutomatically(puzzle []int) {
	dl.PrintStepHeader(1)
	run(puzzle, false)
}

type Droid struct {
	program []int
	debug   bool
	in      chan string
	out     chan string
	err     chan error
	Stop    func()
}

func (d *Droid) Start() {
	in := make(chan int, 100)
	out := make(chan int)
	sgn := make(chan int)

	var wg sync.WaitGroup
	wg.Add(3)
	go func(in chan int, out chan int, sgn chan int) {
		for e := range dl.ExecuteIntcode(d.program, in, out, sgn, false) {
			d.err <- e
		}
		close(in)
		close(sgn)
		close(d.err)
		wg.Done()
	}(in, out, sgn)

	d.Stop = func() {
		sgn <- 9
	}

	// Transport input (string channel to int channel)
	go func() {
		for line := range d.in {
			if d.debug {
				fmt.Printf("%s\n", color.New(color.FgRed).Sprint(line))
			}
			// send command (ASCII to int)
			for _, c := range line {
				in <- int(c)
			}
			in <- '\n'
		}
		wg.Done()
	}()

	var buffer []int // line buffer
	go func() {
		for c := range out {
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
		}
		close(d.in)
		close(d.out)
		wg.Done()
	}()

	wg.Wait()
}

func newDroid(program []int, debug bool) *Droid {
	return &Droid{
		program: program,
		debug:   debug,
		in:      make(chan string, 10),
		out:     make(chan string, 10),
		err:     make(chan error, 10),
	}
}

type Room struct {
	Name  string
	Doors map[dl.Direction]*Room
	Path  []string
	Items []string
}

func (r *Room) Equal(o *Room) bool {
	return r.Name == o.Name
}

var sentinelRoom = Room{Name: "SENTINEL"}

func run(program []int, debug bool) {

	rooms := make(map[string]*Room)
	blacklistItems := make([]string, 0)
	securityRoomPath := make([]string, 0)
	itemMap := make(map[string]string) // itemName to roomName
	foundItems := make([]string, 0)

	type State struct {
		Commands  []string
		Direction dl.Direction
		LastRoom  *Room
	}

	queue := list.New()
	queue.PushBack(&State{
		Commands:  []string{},
		Direction: dl.NORTH,
		LastRoom:  nil,
	})

	for queue.Len() > 0 {

		var state *State
		{
			next := queue.Front()
			state = next.Value.(*State)
			queue.Remove(next)
		}

		commandReplay := make([]string, len(state.Commands))
		copy(commandReplay, state.Commands)
		commandPath := make([]string, 0)
		var prevDirection dl.Direction = state.Direction
		var prevRoom *Room = state.LastRoom
		inventory := make([]string, 0)

		unexpectedEnd := true
		droid := newDroid(program, debug)
		go func() {
			var roomName string
			var doorParsing bool
			var itemParsing bool
			var doorOptions []string
			var itemOptions []string

			var lastEmptyLines int
			var lastMessages []string

			lineProblemDetectors := []func(line string) (bool, bool){
				func(line string) (bool, bool) {
					l := len(lastMessages)
					if l > 0 {
						if strings.Contains(lastMessages[l-1], "Bye!") {
							return true, false
						}
					}
					return false, false
				},
				func(line string) (bool, bool) {
					l := len(lastMessages)
					if l > 0 {
						if strings.Contains(lastMessages[l-1], "It is suddenly completely dark! You are eaten by a Grue!") {
							return true, false
						}
					}
					return false, false
				},
				func(line string) (bool, bool) {
					l := len(lastMessages)
					if l > 1 {
						if lastMessages[l-2] == lastMessages[l-1] {
							// recursive output
							return true, true
						}
					}
					return false, false
				},
			}

			advancedProblemDetectors := []func() (bool, bool){
				func() (bool, bool) {
					if prevRoom != nil && prevRoom.Name == roomName {
						if len(lastMessages) > 1 && strings.Contains(lastMessages[len(lastMessages)-2], "You can't move") {
							return true, true
						}
					}
					return false, false
				},
			}

			for line := range droid.out {

				if len(line) == 0 {
					lastEmptyLines++
				} else {
					lastEmptyLines = 0
					lastMessages = append(lastMessages, line)
				}

				if debug {
					fmt.Println(line)
				}

				if len(line) > 0 && line[0:3] == "== " {
					reg, _ := regexp.Compile("== (.*) ==")
					r := reg.FindStringSubmatch(line)
					roomName = strings.TrimSpace(r[1])
					doorOptions = []string{}
					itemOptions = []string{}
					lastMessages = []string{}
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
					if !strings.Contains(strings.Join(blacklistItems, ","), item+",") {
						itemOptions = append(itemOptions, item)
					}
				} else {
					itemParsing = false
				}

				// room output ready (without command interaction)
				if line == "Command?" || (len(roomName) > 0 && !doorParsing && len(doorOptions) > 0 && !itemParsing && lastEmptyLines > 2) {

					// store knowledge about map (room)
					if _, exist := rooms[roomName]; !exist {
						room := &Room{
							Name:  roomName,
							Path:  commandPath,
							Doors: make(map[dl.Direction]*Room),
							Items: itemOptions,
						}
						doorOptionsStr := strings.Join(doorOptions, "")
						// using the sentinel room, this directions are being marked as "visited"
						for _, dir := range []string{"north", "east", "south", "west"} {
							if !strings.Contains(doorOptionsStr, dir) {
								room.Doors[dl.NewDirection(dir)] = &sentinelRoom
							}
						}
						rooms[roomName] = room
					}

					room := rooms[roomName]
					if prevRoom != nil {
						prevRoom.Doors[prevDirection] = room
						room.Doors[prevDirection.Return()] = prevRoom
					}

					if roomName == "Pressure-Sensitive Floor" {
						securityRoomPath = make([]string, len(commandPath))
						copy(securityRoomPath, commandPath)
						fmt.Println("🎉 Found security floor")
						unexpectedEnd = false
						go droid.Stop()
						break
					}
				}

				detected, killRequired := false, false
				for _, detector := range lineProblemDetectors {
					if d, k := detector(line); d {
						detected = d
						killRequired = k
						break
					}
				}
				if detected {
					// problem with last item
					item := inventory[len(inventory)-1]
					blacklistItems = append(blacklistItems, item)
					inventory = inventory[0 : len(inventory)-1]

					fmt.Printf("👉 Putting item '%s' on blacklist...\n", item)

					// restart
					queue.PushBack(&State{
						Commands:  commandPath,
						Direction: prevDirection,
						LastRoom:  prevRoom,
					})

					if killRequired {
						go droid.Stop()
					}
					unexpectedEnd = false
					break
				}

				if line == "Command?" {
					room := rooms[roomName]

					if len(itemOptions) > 0 {
						option := itemOptions[0]
						if _, exist := itemMap[option]; !exist {
							itemMap[option] = room.Name
						}
						itemOptions = itemOptions[1:]
						if !strings.Contains(strings.Join(blacklistItems, ","), option+",") {
							droid.in <- "take " + option
							inventory = append(inventory, option)
						}
						continue
					}

					if len(commandReplay) > 0 {
						cmd := commandReplay[0]
						commandReplay = commandReplay[1:]
						commandPath = append(commandPath, cmd)
						// state room&dir already set correctly
						droid.in <- cmd
						continue
					}

					// At this point AFTER command replay, the state of prevRoom is correct
					detected, killRequired := false, false
					for _, detector := range advancedProblemDetectors {
						if d, k := detector(); d {
							detected = d
							killRequired = k
							break
						}
					}
					if detected {
						// problem with last item
						item := inventory[len(inventory)-1]
						blacklistItems = append(blacklistItems, item)
						inventory = inventory[0 : len(inventory)-1]

						fmt.Printf("👉 Putting item '%s' on blacklist...\n", item)

						// restart
						queue.PushBack(&State{
							Commands:  commandPath,
							Direction: prevDirection,
							LastRoom:  prevRoom,
						})

						if killRequired {
							go droid.Stop()
						}
						unexpectedEnd = false
						break
					}

					unvisited := make([]dl.Direction, 0)
					for _, dir := range []dl.Direction{
						dl.NORTH,
						dl.EAST,
						dl.SOUTH,
						dl.WEST,
					} {
						if _, exist := room.Doors[dir]; !exist {
							unvisited = append(unvisited, dir)
						}
					}

					// branching required (restart droid)
					for len(unvisited) > 1 {
						nextDirection := unvisited[len(unvisited)-1]
						unvisited = unvisited[0 : len(unvisited)-1]

						branchedCommandPath := make([]string, len(commandPath)+1)
						copy(branchedCommandPath, commandPath)
						branchedCommandPath[len(branchedCommandPath)-1] = nextDirection.ToString()
						queue.PushBack(&State{
							Commands:  branchedCommandPath,
							Direction: nextDirection,
							LastRoom:  room,
						})
					}

					if len(unvisited) == 1 {
						cmd := unvisited[0].ToString()
						prevRoom = room
						prevDirection = unvisited[0]
						commandPath = append(commandPath, cmd)
						droid.in <- cmd
					} else {
						fmt.Println("☠️ Droid has found a dead end.")
						droid.Stop()
						unexpectedEnd = false
					}
				}
			}
			if unexpectedEnd {
				// problem with last item
				item := inventory[len(inventory)-1]
				blacklistItems = append(blacklistItems, item)
				inventory = inventory[0 : len(inventory)-1]
				fmt.Printf("👉 Putting item '%s' on blacklist...\n", item)

				// Restart
				queue.PushBack(&State{
					Commands:  commandPath,
					Direction: prevDirection,
					LastRoom:  prevRoom,
				})
			}
			for _, inv := range inventory {
				if !strings.Contains(strings.Join(foundItems, ","), inv+",") {
					foundItems = append(foundItems, inv)
				}
			}
		}()
		droid.Start()
	}

	fmt.Printf("👉 Found path to security room: %s\n", strings.Join(securityRoomPath, ", "))
	fmt.Printf("👉 Found items: %s\n", strings.Join(foundItems, ", "))
	fmt.Printf("👉 The following items are on the blacklist: %s\n", strings.Join(blacklistItems, ", "))

	droid := newDroid(program, debug)

	// Collect each item: Go to the room, go back again to origin.
	commands := make([]string, 0)
	for _, itemName := range foundItems {
		directions := rooms[itemMap[itemName]].Path
		for _, direction := range directions {
			commands = append(commands, direction)
		}
		commands = append(commands, "take "+itemName)
		for i := len(directions) - 1; i >= 0; i-- {
			commands = append(commands, dl.NewDirection(directions[i]).Return().ToString())
		}
	}
	// Add path to Security Checkpoint
	checkpointRoom := rooms["Security Checkpoint"]
	for _, direction := range checkpointRoom.Path {
		commands = append(commands, direction)
	}
	sensitiveRoom := rooms["Pressure-Sensitive Floor"]
	securityDirection := dl.NORTH
	for dir, room := range checkpointRoom.Doors {
		if room == sensitiveRoom {
			securityDirection = dir
			break
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		skip := true

		permutationsOfActiveItems := unorderedAndUniquePermutations(foundItems)

		// at the begin, all items are being carried
		carriedItems := make([]string, len(foundItems))
		copy(carriedItems, foundItems)

		mode := "init"
		currentItemSet := make([]string, 0)
		for line := range droid.out {
			if line == "Command?" && len(commands) > 0 {
				command := commands[0]
				commands = commands[1:]
				droid.in <- command
				continue
			}
			if skip && !strings.Contains(line, "Security Checkpoint") {
				continue
			} else {
				skip = false
			}

			if debug {
				fmt.Println(line)
			}

			if strings.Contains(line, "Security Checkpoint") {
				mode = "drop"
			}

			if line == "Command?" {
				if mode == "drop" {
					if len(carriedItems) > 0 {
						droid.in <- "drop " + carriedItems[0]
						carriedItems = carriedItems[1:]
						continue
					} else {
						mode = "take"
						currentItemSet = permutationsOfActiveItems[0]
						permutationsOfActiveItems = permutationsOfActiveItems[1:]
					}
				}

				if mode == "take" {
					if len(currentItemSet) > 0 {
						droid.in <- "take " + currentItemSet[0]
						carriedItems = append(carriedItems, currentItemSet[0])
						currentItemSet = currentItemSet[1:]
						continue
					} else {
						mode = "walk"
					}
				}

				if mode == "walk" {
					droid.in <- securityDirection.ToString()
					mode = "test"
				}
			}

			if strings.Contains(line, "on the keypad") {
				r, _ := regexp.Compile(" You should be able to get in by typing (\\d+) on the keypad at the main airlock")
				code := dl.ParseInt(r.FindStringSubmatch(line)[1])
				dl.PrintSolution(fmt.Sprintf("The code is %d", code))
			}
		}
		wg.Done()
	}()
	droid.Start()
	wg.Wait()
}

func unorderedAndUniquePermutations(options []string) [][]string {
	n := uint(len(options))
	m := uint(math.Pow(float64(2), float64(n))) - 1
	result := make([][]string, 0)

	for i := uint(1); i <= m; i++ {
		sub := make([]string, 0)
		for k := n; k > 0; k-- {
			if (i>>(k-1))&1 == 1 {
				sub = append(sub, options[n-k])
			}
		}
		result = append(result, sub)
	}

	return result
}

func renderAsciiToString(n []int) string {
	result := ""
	for _, c := range n {
		result += string(c)
	}
	return result
}
