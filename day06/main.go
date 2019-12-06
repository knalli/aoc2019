package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const AocDay = 6
const AocDayName = "day06"
const AocDayTitle = "Day 06"

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(0)
		lines, _ := dl.ReadFileToArray(AocDayName + "/sample1.txt")
		orbit := newOrbit(lines)
		dl.PrintSolution(fmt.Sprintf("Total number of direct and indirect orbits: %d", orbit.getTotalOrbits()))
	}
	{
		dl.PrintStepHeader(1)
		lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		orbit := newOrbit(lines)
		dl.PrintSolution(fmt.Sprintf("Total number of direct and indirect orbits: %d", orbit.getTotalOrbits()))
	}

	{
		dl.PrintStepHeader(0)
		lines, _ := dl.ReadFileToArray(AocDayName + "/sample2.txt")
		orbit := newOrbit(lines)
		transfer, length := orbit.resolveOrbitalTransfer(orbit.getOrbiting("YOU"), orbit.getOrbiting("SAN"))
		dl.PrintSolution(fmt.Sprintf("The orbital transfer length is %d via rendenvous point %s", length, transfer))
	}
	{
		dl.PrintStepHeader(2)
		lines, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		orbit := newOrbit(lines)
		transfer, length := orbit.resolveOrbitalTransfer(orbit.getOrbiting("YOU"), orbit.getOrbiting("SAN"))
		dl.PrintSolution(fmt.Sprintf("The orbital transfer length is %d via rendenvous point %s", length, transfer))
	}

}

type Orbit struct {
	data map[string]string
}

func (o *Orbit) getTotalOrbitsByObject(name string, includeIndirects bool) int {
	if _, exist := o.data[name]; !exist {
		panic("invalid name")
	}
	result := 0
	if direct := o.data[name]; direct != "" {
		result += 1
		if includeIndirects {
			result += o.getTotalOrbitsByObject(direct, includeIndirects)
		}
	}
	return result
}

func (o *Orbit) getTotalOrbits() int {
	result := 0
	cache := make(map[string]int, 0)
	for obj := range o.data {
		if _, exist := cache[obj]; !exist {
			cache[obj] = o.getTotalOrbitsByObject(obj, true)
		}
		result += cache[obj]
	}
	return result
}

// Returns the next orbiting object
func (o *Orbit) getOrbiting(name string) string {
	if k, exist := o.data[name]; exist {
		return k
	}
	panic("invalid name")
}

func (o *Orbit) getPath(name string) []string {
	result := make([]string, 0)
	result = append(result, name)
	next := name
	for o.data[next] != "" {
		next = o.data[next]
		result = append(result, next)
	}
	return result
}

func (o *Orbit) resolveOrbitalTransfer(from string, to string) (position string, transferLength int) {
	fromPath := o.getPath(from)
	toPath := o.getPath(to)
	ReverseSlice(fromPath)
	ReverseSlice(toPath)
	// count same
	base := 0
	for base < len(fromPath) && base < len(toPath) {
		if fromPath[base] != toPath[base] {
			break
		}
		base++
	}

	if base == 0 {
		panic("invalid objects in orbit")
	}

	return fromPath[base-1], (len(fromPath) - base) + (len(toPath) - base)
}

func newOrbit(lines []string) Orbit {
	data := make(map[string]string, 0)
	for _, line := range lines {
		parts := strings.Split(line, ")")
		obj1 := parts[0]
		obj2 := parts[1]
		if _, exist := data[obj1]; !exist {
			data[obj1] = ""
		}
		data[obj2] = obj1
	}
	return Orbit{data: data}
}

func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
