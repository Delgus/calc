package calc

import (
	"testing"
)

func TestSumPrice(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "case1", args: args{x: 0.001, y: 0.003}, want: 0.00},
		{name: "case2", args: args{x: 0.002, y: 0.003}, want: 0.01},
		{name: "case3", args: args{x: 0.00099, y: 0.00401}, want: 0.01},
		{name: "case4", args: args{x: 0.00099, y: 0.00400}, want: 0.00},
		{name: "case5", args: args{x: 0.99, y: 0.01}, want: 1.00},
		{name: "case6", args: args{x: 0.99, y: 0.001}, want: 0.99},
		{name: "case7", args: args{x: 0.004, y: 0.99}, want: 0.99},
		{name: "case8", args: args{x: 0.005, y: 0.99}, want: 1.00},
		{name: "case9", args: args{x: -0.001, y: -0.003}, want: 0.00},
		{name: "case10", args: args{x: -0.002, y: -0.003}, want: -0.01},
		{name: "case11", args: args{x: -0.00099, y: -0.00401}, want: -0.01},
		{name: "case12", args: args{x: -0.00099, y: -0.00400}, want: 0.00},
		{name: "case13", args: args{x: -0.99, y: -0.01}, want: -1.00},
		{name: "case14", args: args{x: -0.99, y: -0.001}, want: -0.99},
		{name: "case15", args: args{x: -0.004, y: -0.99}, want: -0.99},
		{name: "case16", args: args{x: -0.005, y: -0.99}, want: -1.00},
		{name: "case17", args: args{x: 9999.9999, y: 9999.9999}, want: 20000.00},
		{name: "case18", args: args{x: -9999.9999, y: -9999.9999}, want: -20000.00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumPrice(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("unexpected result %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumQuantity(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumQuantity(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("SumQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}
