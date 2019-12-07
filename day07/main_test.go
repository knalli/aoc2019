package main

import "testing"

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
				puzzle:        readFileAsIntArray("puzzle1.txt"),
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
