package codes

import "testing"

func Test_C1_plusNum(t *testing.T) {
	type test struct {
		input  int
		output int
	}

	tests := map[string]test{
		"No1": {input: 1, output: 11},
		"No2": {input: -1, output: 9},
		"No3": {input: 0, output: 10},
		"No4": {input: -9, output: 1},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			C1_plusNum(&(tc.input))
			if tc.input != tc.output {
				t.Errorf("exceted: %d, got : %d", tc.input, tc.output)
			}
		})
	}
}

func Test_C1_SliceOp(t *testing.T) {
	type test struct {
		input  []int
		output []int
	}
	isEqual := func(a test) bool {
		if len(a.input) != len(a.output) {
			return false
		}
		for i := range a.input {
			if a.input[i] != a.output[i] {
				return false
			}
		}
		return true
	}

	tests := map[string]test{
		"No1": {input: []int{1, 2, 3}, output: []int{2, 4, 6}},
		"No2": {input: []int{}, output: []int{}},
		"No3": {input: []int{-1, -2, 3}, output: []int{-2, -4, 6}},
		"No4": {input: []int{-1, 2, 0}, output: []int{-2, 4, 0}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			C1_SliceOp(&(tc.input))
			if !isEqual(tc) {
				t.Errorf("exceted: %d, got : %d", tc.input, tc.output)
			}
		})
	}
}
