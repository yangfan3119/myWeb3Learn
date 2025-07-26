package codes

import "testing"

func Test_ChanTransmitNum(t *testing.T) {
	type test struct {
		input  int
		output string
	}
	tests := map[string]test{
		"No1": {input: 10, output: "1,2,3,4,5,6,7,8,9,10,"},
		"No2": {input: 0, output: ""},
		"No3": {input: 5, output: "1,2,3,4,5,"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := ChanTransmitNum(tc.input)
			if got != tc.output {
				t.Errorf("exceted: %s, got : %s", tc.output, got)
			}
		})
	}
}
