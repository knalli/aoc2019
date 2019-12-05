package main

import "testing"

func Test_compute(t *testing.T) {
	type args struct {
		data  []int
		input int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Sample1",
			args: args{
				data:  []int{1002, 4, 3, 4, 33},
				input: 1,
			},
			want:    1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compute(tt.args.data, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("compute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("compute() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compute2(t *testing.T) {
	type args struct {
		data  []int
		input int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// Using position mode, consider whether the input is equal to 8; output 1 (if it is) or 0 (if it is not).
		{
			name: "Sample2.1",
			args: args{
				data:  []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
				input: 8,
			},
			want:    1,
			wantErr: true,
		},
		{
			name: "Sample2.1",
			args: args{
				data:  []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
				input: 9,
			},
			want:    0,
			wantErr: true,
		},
		// Using position mode, consider whether the input is less than 8; output 1 (if it is) or 0 (if it is not).
		{
			name: "Sample2.2",
			args: args{
				data:  []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
				input: 7,
			},
			want:    1,
			wantErr: true,
		},
		{
			name: "Sample2.2",
			args: args{
				data:  []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
				input: 8,
			},
			want:    0,
			wantErr: true,
		},
		// Using immediate mode, consider whether the input is equal to 8; output 1 (if it is) or 0 (if it is not).
		{
			name: "Sample2.3",
			args: args{
				data:  []int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
				input: 8,
			},
			want:    1,
			wantErr: true,
		},
		{
			name: "Sample2.3",
			args: args{
				data:  []int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
				input: 9,
			},
			want:    0,
			wantErr: true,
		},
		// Using immediate mode, consider whether the input is less than 8; output 1 (if it is) or 0 (if it is not).
		{
			name: "Sample2.4",
			args: args{
				data:  []int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
				input: 7,
			},
			want:    1,
			wantErr: true,
		},
		{
			name: "Sample2.4",
			args: args{
				data:  []int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
				input: 9,
			},
			want:    0,
			wantErr: true,
		},
		// Jump
		{
			name: "Sample2.5",
			args: args{
				data:  []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
				input: 0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Sample2.5",
			args: args{
				data:  []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
				input: 1,
			},
			want:    1,
			wantErr: true,
		},
		// Complete sample
		{
			name: "Sample2.6",
			args: args{
				data: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
					1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
					999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
				input: 7,
			},
			want:    999,
			wantErr: true,
		},
		{
			name: "Sample2.6",
			args: args{
				data: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
					1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
					999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
				input: 8,
			},
			want:    1000,
			wantErr: true,
		},
		{
			name: "Sample2.6",
			args: args{
				data: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
					1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
					999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
				input: 9,
			},
			want:    1001,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compute2(tt.args.data, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("compute2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("compute2() got = %v, want %v", got, tt.want)
			}
		})
	}
}
