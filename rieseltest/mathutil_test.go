package rieseltest

import "testing"

func TestLround(t *testing.T) {
	var testCases = []struct {
		input float64
		expected int64
	}{
		{1.4999999999999, 1},
		{1.5, 2},
		{2, 2},
		{0, 0},
		{-1.5, -2},
	}

	for _, c := range testCases {
		actual := lround(c.input)
		if actual != c.expected {
			t.Errorf("lround(%v) == %v, but we expected %v", c.input, actual, c.expected)
		}
	}
}