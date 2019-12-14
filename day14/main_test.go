package main

import (
	dl "de.knallisworld/aoc/aoc2019/dayless"
	"reflect"
	"testing"
)

func Test_react(t *testing.T) {
	type args struct {
		file           string
		targetMaterial Material
		baseChemical   string
	}
	tests := []struct {
		name string
		args args
		want Material
	}{
		{
			name: "Sample1",
			args: args{
				file:           "sample1.txt",
				targetMaterial: Material{1, "FUEL"},
				baseChemical:   "ORE",
			},
			want: Material{31, "ORE"},
		},
		{
			name: "Sample2",
			args: args{
				file:           "sample2.txt",
				targetMaterial: Material{1, "FUEL"},
				baseChemical:   "ORE",
			},
			want: Material{165, "ORE"},
		},
		{
			name: "Sample3",
			args: args{
				file:           "sample3.txt",
				targetMaterial: Material{1, "FUEL"},
				baseChemical:   "ORE",
			},
			want: Material{13312, "ORE"},
		},
		{
			name: "Sample4",
			args: args{
				file:           "sample4.txt",
				targetMaterial: Material{1, "FUEL"},
				baseChemical:   "ORE",
			},
			want: Material{180697, "ORE"},
		},
		{
			name: "Sample5",
			args: args{
				file:           "sample5.txt",
				targetMaterial: Material{1, "FUEL"},
				baseChemical:   "ORE",
			},
			want: Material{2210736, "ORE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			puzzle, _ := dl.ReadFileToArray(tt.args.file)
			reactions := newReactions(puzzle)
			if got := react(reactions, tt.args.targetMaterial, tt.args.baseChemical); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("react() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newReaction(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want Reaction
	}{
		{
			args: args{line: "157 ORE => 5 NZVS"},
			want: Reaction{inputs: []Material{{Amount: 157, Chemical: "ORE"}}, output: Material{Amount: 5, Chemical: "NZVS"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newReaction(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newReaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
