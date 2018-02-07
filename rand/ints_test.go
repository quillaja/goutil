package rand

import (
	"math/rand"
	"testing"
	"time"
)

const testN = 10000 // run each test this many times.

func TestIntNM(t *testing.T) {
	type args struct {
		low  int
		high int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "positive 1", args: args{low: 0, high: 20}},
		{name: "positive 2", args: args{low: 9, high: 100}},
		{name: "positive 3", args: args{low: 0, high: 1}},
		{name: "negative 1", args: args{low: -10, high: -2}},
		{name: "negative 2", args: args{low: -520, high: -401}},
		{name: "negitive 3", args: args{low: -2, high: -1}},
		{name: "span 1", args: args{low: -5, high: 5}},
		{name: "span 2", args: args{low: -10236, high: 1}},
		{name: "span 3", args: args{low: -1, high: 523498}},
	}
	rand.Seed(time.Now().UnixNano())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < testN; i++ {
				got := IntNM(tt.args.low, tt.args.high)
				if !(tt.args.low <= got && got < tt.args.high) {
					t.Errorf("%d <= %d < %d", tt.args.low, got, tt.args.high)

				}
			}
		})
	}
}
