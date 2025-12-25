package filters

// import (
// 	"fmt"
// 	"testing"
// )
//
// func equal(_lhs, _rhs [][]float64) bool {
// 	for y := range len(_lhs) {
// 		for x := range len(_lhs[0]) {
// 			if _lhs[y][x] != _rhs[y][x] {
// 				return false
// 			}
// 		}
// 	}
//
// 	return true
// }
//
// func Test_m(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for target function.
// 		_n   int
// 		want [][]float64
// 	}{
// 		{
// 			name: "M_2",
// 			_n:   1,
// 			want: [][]float64{
// 				{0.0, 0.5},
// 				{0.75, 0.25}},
// 		},
// 		{
// 			name: "M_4",
// 			_n:   2,
// 			want: [][]float64{
// 				{0.0, 0.5, 0.125, 0.625},
// 				{0.75, 0.25, 0.875, 0.375},
// 				{0.1875, 0.6875, 0.0625, 0.5625},
// 				{0.9375, 0.4375, 0.8125, 0.3125}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := m(tt._n)
//
// 			fmt.Println(got)
//
// 			if !equal(tt.want, got) {
// 				t.Errorf("m() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
