package codes

import "testing"

func Test_GoAtomicCounter(t *testing.T) {
	type test struct {
		input  int
		output uint64
	}
	tests := map[string]test{
		"No1": {input: 10, output: uint64(100000)},
		"No2": {input: 100, output: uint64(1000000)},
		"No3": {input: 1000, output: uint64(10000000)},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GoAtomicCounter(tc.input)
			if got != tc.output {
				t.Errorf("exceted: %d, got : %d", tc.output, got)
			}
		})
	}
}

func Test_GoMutexCounter(t *testing.T) {
	type test struct {
		input  int
		output uint64
	}
	tests := map[string]test{
		"No1": {input: 10, output: uint64(100000)},
		"No2": {input: 100, output: uint64(1000000)},
		"No3": {input: 1000, output: uint64(10000000)},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GoMutexCounter(tc.input)
			if got != tc.output {
				t.Errorf("exceted: %d, got : %d", tc.output, got)
			}
		})
	}
}
