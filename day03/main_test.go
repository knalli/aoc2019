package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"testing"
)

func readSample(file string) []string {
	res, _ := dl.ReadFileToArray(file)
	return res
}

func Test_getShortestDistance(t *testing.T) {
	type args struct {
		base  Point
		lines []string
	}
	tests := []struct {
		name  string
		args  args
		point Point
		dist  int
	}{
		{
			name: "Puzzle",
			args: args{
				base:  Point{},
				lines: readSample("puzzle1.txt"),
			},
			dist: 651,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, dist := getShortestDistance(tt.args.base, getIntersections(getPaths(tt.args.lines)))
			if dist != tt.dist {
				t.Errorf("getShortestDistance() got = %v, want %v", dist, tt.dist)
			}
		})
	}
}

func Test_getShortestPath(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		dist int
	}{
		{
			name: "Sample1",
			args: args{
				lines: readSample("sample1.txt"),
			},
			dist: 610,
		},
		{
			name: "Sample2",
			args: args{
				lines: readSample("sample2.txt"),
			},
			dist: 410,
		},
		{
			name: "Sample3",
			args: args{
				lines: readSample("sample3.txt"),
			},
			dist: 14,
		},
		{
			name: "Puzzle1",
			args: args{
				lines: readSample("puzzle1.txt"),
			},
			dist: 7534,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wires := getWires(tt.args.lines)
			dist, _ := getShortestPath(Point{0, 0}, wires, getIntersections(getPaths(tt.args.lines)))
			if dist != tt.dist {
				t.Errorf("Test_getShortestPath() got = %v, want %v", dist, tt.dist)
			}
		})
	}
}
