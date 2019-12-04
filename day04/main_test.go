package main

import "testing"

func Test_hasSameAdjacents(t *testing.T) {
	type args struct {
		numbers  []int
		required int
		exact    bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "111111",
			args: args{
				numbers:  []int{1, 1, 1, 1, 1, 1},
				required: 2,
				exact:    false,
			},
			want: true,
		},
		{
			name: "223450",
			args: args{
				numbers:  []int{2, 2, 3, 4, 5, 0},
				required: 2,
				exact:    false,
			},
			want: true,
		},
		{
			name: "123789",
			args: args{
				numbers:  []int{1, 2, 3, 7, 8, 9},
				required: 2,
				exact:    false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSameAdjacents(tt.args.numbers, tt.args.required, tt.args.exact); got != tt.want {
				t.Errorf("hasSameAdjacents() = %v, want %v", got, tt.want)
			}
		})
	}
}
