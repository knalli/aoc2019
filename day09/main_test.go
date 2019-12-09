package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"reflect"
	"testing"
)

func Test_runBoostProgram(t *testing.T) {
	type args struct {
		program []int
		base    int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Sample1",
			args: args{
				program: dl.ReadFileAsIntArray("sample1.txt"),
				base:    0,
			},
			want: []int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
		{
			name: "Sample2",
			args: args{
				program: dl.ReadFileAsIntArray("sample2.txt"),
				base:    0,
			},
			want: []int{1219070632396864},
		},
		{
			name: "Sample3",
			args: args{
				program: dl.ReadFileAsIntArray("sample3.txt"),
				base:    0,
			},
			want: []int{1125899906842624},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runBoostProgram(tt.args.program, tt.args.base); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runBoostProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}
