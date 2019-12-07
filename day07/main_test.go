package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"reflect"
	"testing"
)

func Test_compute(t *testing.T) {
	type args struct {
		instructions   []int
		phaseSequences [] int
		input          int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Sample1",
			args: args{
				instructions:   []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
				phaseSequences: []int{4, 3, 2, 1, 0},
				input:          0,
			},
			want: 43210,
		},
		{
			name: "Sample2",
			args: args{
				instructions:   []int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
				phaseSequences: []int{0, 1, 2, 3, 4},
				input:          0,
			},
			want: 54321,
		},
		{
			name: "Sample3",
			args: args{
				instructions:   []int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0},
				phaseSequences: []int{1, 0, 4, 3, 2},
				input:          0,
			},
			want: 65210,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := <-compute(tt.args.instructions, tt.args.phaseSequences, tt.args.input); got != tt.want {
				t.Errorf("compute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findHighestSignal(t *testing.T) {
	type args struct {
		puzzle        []int
		phaseSequence []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Puzzle1",
			args: args{
				puzzle:        dl.ReadFileAsIntArray("puzzle1.txt"),
				phaseSequence: []int{0, 1, 2, 3, 4},
			},
			want: 46248,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findHighestSignal(tt.args.puzzle, tt.args.phaseSequence); got != tt.want {
				t.Errorf("findHighestSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_computeAmplified(t *testing.T) {
	type args struct {
		instructions  []int
		phaseSequence []int
		input         int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Sample1",
			args: args{
				instructions:  []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
				phaseSequence: []int{9, 8, 7, 6, 5},
				input:         0,
			},
			want: 139629729,
		},
		{
			name: "Sample2",
			args: args{
				instructions:  []int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10},
				phaseSequence: []int{9, 7, 8, 5, 6},
				input:         0,
			},
			want: 18216,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeAmplified(tt.args.instructions, tt.args.phaseSequence, tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeAmplified() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findHighestAmplifiedSignal(t *testing.T) {
	type args struct {
		puzzle        []int
		phaseSequence []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Puzzle1",
			args: args{
				puzzle:        dl.ReadFileAsIntArray("puzzle1.txt"),
				phaseSequence: []int{5, 6, 7, 8, 9},
			},
			want: 54163586,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findHighestAmplifiedSignal(tt.args.puzzle, tt.args.phaseSequence); got != tt.want {
				t.Errorf("findHighestSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}
