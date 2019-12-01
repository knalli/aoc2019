package main

import "testing"

func Test_computeFuelByModule(t *testing.T) {
	type args struct {
		mod  int
		deep bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sample1",
			args: args{
				mod:  12,
				deep: false,
			},
			want: 2,
		},
		{
			name: "sample2",
			args: args{
				mod:  14,
				deep: false,
			},
			want: 2,
		},
		{
			name: "sample3",
			args: args{
				mod:  1969,
				deep: false,
			},
			want: 654,
		}, {
			name: "sample4",
			args: args{
				mod:  100756,
				deep: false,
			},
			want: 33583,
		}, {
			name: "sample5",
			args: args{
				mod:  14,
				deep: true,
			},
			want: 2,
		}, {
			name: "sample6",
			args: args{
				mod:  1969,
				deep: true,
			},
			want: 966,
		}, {
			name: "sample7",
			args: args{
				mod:  100756,
				deep: true,
			},
			want: 50346,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeFuelByModule(tt.args.mod, tt.args.deep); got != tt.want {
				t.Errorf("computeFuelByModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
