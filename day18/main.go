package main

import (
	"container/list"
	day18 "de.knallisworld/aoc/aoc2019/day18/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const AocDay = 18
const AocDayName = "day18"
const AocDayTitle = "Day 18"

const PLAYER = "@"
const WALL = "#"
const OPEN = "."

func main() {
	dl.PrintDayHeader(AocDay, AocDayTitle)
	defer dl.TimeTrack(time.Now(), AocDayName)

	{
		dl.PrintStepHeader(1)
		puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		m := day18.NewMap(puzzle)
		fmt.Printf("Puzzle:\n%s\n", m.ToString())
		steps := findShortestPathCollectingAllKeys(m, false)
		dl.PrintSolution(fmt.Sprintf("Shortest path (steps) is %d", steps))
	}

	{
		dl.PrintStepHeader(2)
		puzzle, _ := dl.ReadFileToArray(AocDayName + "/puzzle1.txt")
		m := day18.NewMap(puzzle)

		// update map
		p := m.Filter(func(v string) bool {
			return v == PLAYER
		})[0]
		for _, a := range []day18.Point{
			p,
			p.North(),
			p.East(),
			p.South(),
			p.West(),
		} {
			m.Set(a, WALL)
		}
		for _, a := range []day18.Point{
			p.North().East(),
			p.South().East(),
			p.South().West(),
			p.North().West(),
		} {
			m.Set(a, PLAYER)
		}

		fmt.Printf("Puzzle:\n%s\n", m.ToString())
		steps := findShortestPathCollectingAllKeys(m, false)
		dl.PrintSolution(fmt.Sprintf("Shortest path (steps) is %d", steps))
	}

}

type State struct {
	Positions []day18.Point
	Vault     *day18.Map
	Cost      int
}

func (s *State) Move(i int, p day18.Point) *State {
	s.Positions[i] = p
	return s
}

func (s *State) ConsumeKey(k string) *State {
	key := s.Vault.FindFirst(func(v string) bool {
		return v == k
	})
	if key != nil {
		s.Vault.Set(*key, OPEN)
	}
	door := s.Vault.FindFirst(func(v string) bool {
		return string(v[0]+32) == k
	})
	if door != nil {
		s.Vault.Set(*door, OPEN)
	}
	return s
}

func (s *State) AddCost(n int) *State {
	s.Cost += n
	return s
}

func (s *State) Clone() *State {
	vault := s.Vault.Clone()
	return &State{
		Positions: func() []day18.Point {
			result := make([]day18.Point, len(s.Positions))
			for i, p := range s.Positions {
				result[i] = day18.Point{
					X: p.X,
					Y: p.Y,
				}
			}
			return result
		}(),
		Vault: vault,
		Cost:  s.Cost,
	}
}

func findShortestPathCollectingAllKeys(vault day18.Map, debug bool) int {

	starts := vault.Filter(func(v string) bool {
		return v == PLAYER
	})

	// Optimize: static map of shortest path between all keys (incl. start to keys)
	key2keyDistances := make(map[day18.Point]map[day18.Point][]day18.Point)
	{
		// ensure no path to itself are collected
		bfsFiltered := func(vault day18.Map, p1 day18.Point, p2 day18.Point) []day18.Point {
			r := bfs(vault, p1, p2)
			if len(r) < 2 {
				return nil
			} else {
				return r[1:] // remove itself (first)
			}
		}
		for _, start := range starts {
			key2keyDistances[start] = make(map[day18.Point][]day18.Point)
		}
		vault.Each(func(p1 day18.Point, v1 string) {
			c1 := v1[0]
			if !('a' <= c1 && c1 <= 'z') {
				return
			}
			for _, start := range starts {
				if r := bfsFiltered(vault, start, p1); r != nil {
					key2keyDistances[start][p1] = r
				}
			}
			m := make(map[day18.Point][]day18.Point)
			vault.Each(func(p2 day18.Point, v2 string) {
				c2 := v2[0]
				if !('a' <= c2 && c2 <= 'z') {
					return
				}
				if r := bfsFiltered(vault, p1, p2); r != nil {
					m[p2] = r
				}
			})
			key2keyDistances[p1] = m
		})
	}

	isFree := func(vault *day18.Map, path []day18.Point) bool {
		for i := 0; i < len(path)-1; i++ { // without last (actual target)
			p := path[i]
			v := *vault.Get(p)
			if v != "." {
				return false
			}
		}
		return true
	}

	keysLeft := func(vault *day18.Map) bool {
		available := 0
		vault.EachValue(func(v string) {
			c := v[0]
			if 'a' <= c && c <= 'z' {
				available++
			}
		})
		return available > 0
	}

	buildCacheKey := func(state *State) string {
		cacheKey := ""
		for _, position := range state.Positions {
			cacheKey += position.ToString() + ";"
		}
		points := state.Vault.Filter(func(v string) bool {
			c := v[0]
			return 'a' <= c && c <= 'z'
		})
		keys := make([]string, len(points))
		for i, p := range points {
			keys[i] = *state.Vault.Get(p)
		}
		sort.SliceStable(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		cacheKey += "_" + strings.Join(keys, "")
		return cacheKey
	}

	queue := list.New()
	{
		clone := vault.Clone()

		// clear player position
		for _, start := range starts {
			clone.Set(start, OPEN)
		}

		queue.PushFront(&State{
			Positions: starts,
			Vault:     clone,
			Cost:      0,
		})
	}

	resultCache := make(map[string]int)

	minCost := math.MaxInt32

	for queue.Len() > 0 {
		var state *State
		{
			next := queue.Front()
			state = next.Value.(*State)
			queue.Remove(next)
		}

		if !keysLeft(state.Vault) {
			minCost = dl.MinInt(minCost, state.Cost)
			continue
		}

		for i, position := range state.Positions {
			for target, path := range key2keyDistances[position] {
				if isFree(state.Vault, path) {
					key := *state.Vault.Get(target)
					if key == OPEN {
						continue // already used
					}
					next := state.Clone()
					next.AddCost(len(path))
					next.ConsumeKey(key)
					next.Move(i, target)

					// Improve speed and memory performance (drastically)
					// If at this point the costs are not better then a previously run with the same left keys,
					// we can terminate this branch. It will not get better..
					cacheKey := buildCacheKey(next)
					if cachedCost, exist := resultCache[cacheKey]; !exist || cachedCost > next.Cost {
						resultCache[cacheKey] = next.Cost
						// only push for later if it could improve the costs
						queue.PushBack(next)
					}
				}
			}
		}
	}

	return minCost
}

func bfs(m day18.Map, start day18.Point, goal day18.Point) []day18.Point {
	visited := make(map[day18.Point]bool)
	parents := make(map[day18.Point]day18.Point)
	queue := list.New()
	m.Each(func(p day18.Point, v string) {
		if v == WALL {
			visited[p] = true
		} else {
			visited[p] = false
		}
	})

	resolvePath := func() []day18.Point {
		result := make([]day18.Point, 0)
		result = append(result, goal)
		for !result[len(result)-1].Equal(start) {
			result = append(result, parents[result[len(result)-1]])
		}
		// reverse
		sort.SliceStable(result, func(i, j int) bool {
			return true
		})
		return result
	}

	queue.PushFront(start)
	visited[start] = true
	for queue.Len() > 0 {
		var node day18.Point
		{
			next := queue.Front()
			node = next.Value.(day18.Point)
			queue.Remove(next)
		}
		if node.Equal(goal) {
			return resolvePath()
		}
		for _, child := range node.Adjacents() {
			if m.Contains(child) { // valid child
				if childVisited := visited[child]; !childVisited {
					queue.PushBack(child)
					parents[child] = node
					visited[child] = true
				}
			}
		}
	}
	return nil
}
