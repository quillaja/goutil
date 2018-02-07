package rand

import "testing"

func TestFloat64NM(t *testing.T) {
	type args struct {
		low  float64
		high float64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "positive 1", args: args{low: 0.1, high: 21.1}},
		{name: "positive 2", args: args{low: 9.9, high: 100.01}},
		{name: "positive 3", args: args{low: 0.0, high: 1.0}},
		{name: "negative 1", args: args{low: -10.34, high: -2.23}},
		{name: "negative 2", args: args{low: -520.78, high: -401.23}},
		{name: "negitive 3", args: args{low: -2.0, high: -1.0}},
		{name: "span 1", args: args{low: -5.0, high: 5.0}},
		{name: "span 2", args: args{low: -10236.2, high: 1.0}},
		{name: "span 3", args: args{low: -1.9, high: 523498.1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < testN; i++ {
				got := Float64NM(tt.args.low, tt.args.high)
				if !(tt.args.low <= got && got < tt.args.high) {
					t.Errorf("%g <= %g < %g", tt.args.low, got, tt.args.high)
				}
			}
		})
	}
}
