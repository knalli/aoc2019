package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const AocDay = 12
const AocDayName = "day12"
const AocDayTitle = "Day 12"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		moons := parseLines(lines)
		steps := 1000
		never := func(tick int) bool {
			return false
		}
		runTime(steps, moons, false, never)
		totalEnergy := computeTotalEnergy(moons, false)
		dl.PrintSolution(fmt.Sprintf("The total energy after %d steps: %d", steps, totalEnergy))
	}

	{
		dl.PrintStepHeader(2)
		lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		moons := parseLines(lines)
		steps := math.MaxInt32

		maxStep := 0
		mx := 0
		my := 0
		mz := 0
		initial := make([]Moon, len(moons))
		copy(initial, moons)

		history := func(tick int) bool {
			// for each x,y,z (separated): found the first match with the initial state
			if mx == 0 {
				fx := true
				for i, moon := range moons {
					if moon.Position.X != initial[i].Position.X || moon.Velocity.X != initial[i].Velocity.X {
						fx = false
						break
					}
				}
				if fx && tick != 0 {
					mx = tick
					fmt.Printf("👉 Found circulation x=%d\n", mx)
				}
			}

			if my == 0 {
				fy := true
				for i, moon := range moons {
					if moon.Position.Y != initial[i].Position.Y || moon.Velocity.Y != initial[i].Velocity.Y {
						fy = false
						break
					}
				}
				if fy && tick != 0 {
					my = tick
					fmt.Printf("👉 Found circulation y=%d\n", my)
				}
			}

			if mz == 0 {
				fz := true
				for i, moon := range moons {
					if moon.Position.Z != initial[i].Position.Z || moon.Velocity.Z != initial[i].Velocity.Z {
						fz = false
						break
					}
				}
				if fz && tick != 0 {
					mz = tick
					fmt.Printf("👉 Found circulation z=%d\n", mz)
				}
			}

			// Given the pieces of x,y,z: the least common multiple is the first match of all
			if mx != 0 && my != 0 && mz != 0 {
				maxStep = dl.LeastCommonMultiple(mx, my, mz)
				return true
			}

			return false
		}
		runTime(steps, moons, false, history)
		dl.PrintSolution(fmt.Sprintf("First circulation state found after %d steps", maxStep))
	}

}

func computeTotalEnergy(moons []Moon, debug bool) int {
	result := 0
	for _, moon := range moons {
		result += moon.TotalEnergy()
		if debug {
			fmt.Printf("%s\n", moon.EnergyString())
		}
	}
	return result
}

func runTime(ticks int, moons []Moon, debug bool, onTick func(tick int) bool) {
	for tick := 0; tick <= ticks; tick++ {
		if debug {
			fmt.Printf("After %d steps:\n", tick)
			for _, moon := range moons {
				fmt.Printf("%s\n", moon.ToString())
			}
		}
		if onTick(tick) {
			break
		}
		if tick == ticks {
			if debug {
				fmt.Println()
			}
			break
		}
		processMoons(moons)
		if debug {
			fmt.Println()
		}
	}
}

type Vector struct {
	X, Y, Z int
}

func (v *Vector) Add(o Vector) Vector {
	return Vector{X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z}
}

func (v *Vector) SumOfValues() int {
	return dl.AbsInt(v.X) + dl.AbsInt(v.Y) + dl.AbsInt(v.Z)
}

func (v *Vector) toString() string {
	return fmt.Sprintf("(%d,%d,%d)", v.X, v.Y, v.Z)
}

func (v *Vector) ToLongString() string {
	return fmt.Sprintf("<x=%3d, y=%3d, z=%3d>", v.X, v.Y, v.Z)
}

type Moon struct {
	Position Vector
	Velocity Vector
}

func (m *Moon) ToString() string {
	return fmt.Sprintf("pos=%s, vel=%s", m.Position.ToLongString(), m.Velocity.ToLongString())
}

func (m *Moon) PotentialEnergy() int {
	return m.Position.SumOfValues()
}

func (m *Moon) KineticEnergy() int {
	return m.Velocity.SumOfValues()
}

func (m *Moon) TotalEnergy() int {
	return m.PotentialEnergy() * m.KineticEnergy()
}

func newMoon(x, y, z int) Moon {
	return Moon{Position: Vector{x, y, z}, Velocity: Vector{0, 0, 0}}
}

func parseLines(lines []string) []Moon {
	result := make([]Moon, len(lines))
	for i, line := range lines {
		result[i] = parseLine(line)
	}
	return result
}

func parseLine(line string) Moon {
	parts := strings.Split(line[1:len(line)-1], ", ")
	x, _ := strconv.Atoi(parts[0][2:])
	y, _ := strconv.Atoi(parts[1][2:])
	z, _ := strconv.Atoi(parts[2][2:])
	return newMoon(x, y, z)
}

func processMoons(moons []Moon) {
	// apply gravity
	for a := range moons {
		for b := range moons {
			if moons[a].Position.X < moons[b].Position.X {
				moons[a].Velocity.X++
			} else if moons[a].Position.X > moons[b].Position.X {
				moons[a].Velocity.X--
			}
			if moons[a].Position.Y < moons[b].Position.Y {
				moons[a].Velocity.Y++
			} else if moons[a].Position.Y > moons[b].Position.Y {
				moons[a].Velocity.Y--
			}
			if moons[a].Position.Z < moons[b].Position.Z {
				moons[a].Velocity.Z++
			} else if moons[a].Position.Z > moons[b].Position.Z {
				moons[a].Velocity.Z--
			}
		}
	}
	// apply velocity
	for i := range moons {
		moons[i].Position = moons[i].Position.Add(moons[i].Velocity)
	}
}

func (m *Moon) EnergyString() string {
	return fmt.Sprintf("pot: %2d + %2d + %2d = %2d;   kin: %2d + %2d + %2d = %2d;   total: %2d * %2d = %3d",
		m.Position.X, m.Position.Y, m.Position.Z, m.Position.SumOfValues(), m.Velocity.X, m.Velocity.Y, m.Velocity.Z, m.Velocity.SumOfValues(), m.Position.SumOfValues(), m.Velocity.SumOfValues(), m.TotalEnergy())
}
