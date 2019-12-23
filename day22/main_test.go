package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"reflect"
	"testing"
)

func Test_shuffle(t *testing.T) {
	type args struct {
		factorySize int
		techniques  []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Sample 1",
			args: args{
				factorySize: 10,
				techniques: []string{
					"deal with increment 7",
					"deal into new stack",
					"deal into new stack",
				},
			},
			want: []int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},
		},
		{
			name: "Sample 2",
			args: args{
				factorySize: 10,
				techniques: []string{
					"cut 6",
					"deal with increment 7",
					"deal into new stack",
				},
			},
			want: []int{3, 0, 7, 4, 1, 8, 5, 2, 9, 6},
		},
		{
			name: "Sample 3",
			args: args{
				factorySize: 10,
				techniques: []string{
					"deal with increment 7",
					"deal with increment 9",
					"cut -2",
				},
			},
			want: []int{6, 3, 0, 7, 4, 1, 8, 5, 2, 9},
		},
		{
			name: "Sample 4",
			args: args{
				factorySize: 10,
				techniques: []string{
					"deal into new stack",
					"cut -2",
					"deal with increment 7",
					"cut 8",
					"cut -4",
					"deal with increment 7",
					"cut 3",
					"deal with increment 9",
					"deal with increment 3",
					"cut -1",
				},
			},
			want: []int{9, 2, 5, 8, 1, 4, 7, 0, 3, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shuffle(tt.args.factorySize, tt.args.techniques, 1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shuffle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_solve2(t *testing.T) {
	type args struct {
		lines []string
		c     int64
		n     int64
		p     int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Solution",
			args: args{
				lines: func() []string {
					l, _ := dl.ReadFileToArray("puzzle1.txt")
					return l
				}(),
				c: 119315717514047,
				n: 101741582076661,
				p: 2020,
			},
			want: 55574110161534,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve2(tt.args.lines, tt.args.c, tt.args.n, tt.args.p); got != tt.want {
				t.Errorf("solve2() = %v, want %v", got, tt.want)
			}
		})
	}
}
