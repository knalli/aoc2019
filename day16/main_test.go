package main

import "testing"

func Test_fft(t *testing.T) {
	type args struct {
		input       string
		basePattern string
		phases      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Sample1",
			args: args{
				input:       "12345678",
				basePattern: "0, 1, 0, -1",
				phases:      1,
			},
			want: "48226158",
		},
		{
			name: "Sample1",
			args: args{
				input:       "12345678",
				basePattern: "0, 1, 0, -1",
				phases:      4,
			},
			want: "01029498",
		},
		{
			name: "Sample1",
			args: args{
				input:       "12345678",
				basePattern: "0, 1, 0, -1",
				phases:      100,
			},
			want: "23845678",
		},
		{
			name: "SampleLarge1",
			args: args{
				input:       "80871224585914546619083218645595",
				basePattern: "0, 1, 0, -1",
				phases:      100,
			},
			want: "24176176",
		},
		{
			name: "SampleLarge2",
			args: args{
				input:       "19617804207202209144916044189917",
				basePattern: "0, 1, 0, -1",
				phases:      100,
			},
			want: "73745418",
		},
		{
			name: "SampleLarge3",
			args: args{
				input:       "69317163492948606335995924319873",
				basePattern: "0, 1, 0, -1",
				phases:      100,
			},
			want: "52432133",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fft(tt.args.input, tt.args.basePattern, tt.args.phases)[0:8]; got != tt.want {
				t.Errorf("fft() = %v, want %v", got, tt.want)
			}
		})
	}
}
