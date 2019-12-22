package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"testing"
)

func Test_solution1(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Sample1",
			args: args{filename: "sample1.txt"},
			want: 23,
		},
		{
			name: "Sample2",
			args: args{filename: "sample2.txt"},
			want: 58,
		},
		{
			name: "Puzzle1",
			args: args{filename: "puzzle1.txt"},
			want: 616,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, _ := dl.ReadFileToArray(tt.args.filename)
			if got := solution1(lines); got != tt.want {
				t.Errorf("solution1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_solutions2(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Puzzle 1",
			args: args{filename: "puzzle1.txt"},
			want: 7498,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, _ := dl.ReadFileToArray(tt.args.filename)
			if got := solutions2(lines, false, false); got != tt.want {
				t.Errorf("solutions2() = %v, want %v", got, tt.want)
			}
		})
	}
}
