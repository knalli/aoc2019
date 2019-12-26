package main

import (
	day18 "de.knallisworld/aoc/aoc2019/day18/lib"
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"reflect"
	"testing"
)

func Test_bfs(t *testing.T) {
	type args struct {
		filename string
		start    string
		goal     string
	}
	tests := []struct {
		name string
		args args
		want []day18.Point
	}{
		{
			name: "Sample1 (@ -> A)",
			args: args{
				filename: "sample1.txt",
				start:    PLAYER,
				goal:     "A",
			},
			want: []day18.Point{
				{5, 1},
				{4, 1},
				{3, 1},
			},
		},
		{
			name: "Sample1 (@ -> b)",
			args: args{
				filename: "sample1.txt",
				start:    PLAYER,
				goal:     "b",
			},
			want: []day18.Point{
				{5, 1},
				{4, 1},
				{3, 1},
				{2, 1},
				{1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, _ := dl.ReadFileToArray(tt.args.filename)
			m := day18.NewMap(lines)
			start := m.FindFirst(func(v string) bool {
				return v == tt.args.start
			})
			goal := m.FindFirst(func(v string) bool {
				return v == tt.args.goal
			})
			if got := bfs(m, *start, *goal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bfs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findShortestPathCollectingAllKeys(t *testing.T) {
	type args struct {
		filename string
		patcher  func(m day18.Map)
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Sample 1",
			args: args{
				filename: "sample1.txt",
			},
			want: 8,
		},
		{
			name: "Sample 2",
			args: args{
				filename: "sample2.txt",
			},
			want: 86,
		},
		{
			name: "Sample 3",
			args: args{
				filename: "sample3.txt",
			},
			want: 132,
		},
		{
			name: "Sample 4",
			args: args{
				filename: "sample4.txt",
			},
			want: 136,
		},
		{
			name: "Sample 5",
			args: args{
				filename: "sample5.txt",
			},
			want: 81,
		},
		{
			name: "Puzzle 1",
			args: args{
				filename: "puzzle1.txt",
			},
			want: 3856,
		},
		{
			name: "Puzzle 1 (4 Players)",
			args: args{
				filename: "puzzle1.txt",
				patcher: func(m day18.Map) {
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
				},
			},
			want: 1660,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, _ := dl.ReadFileToArray(tt.args.filename)
			m := day18.NewMap(lines)
			if tt.args.patcher != nil {
				tt.args.patcher(m)
			}
			if got := findShortestPathCollectingAllKeys(m, false); got != tt.want {
				t.Errorf("findShortestPathCollectingAllKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
