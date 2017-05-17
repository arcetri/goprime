package rieseltest

import (
	"testing"
	"fmt"
	"reflect"
)

func TestLround(t *testing.T) {
	var testCases = []struct {
		input float64
		expected int
	}{
		{1.4999999999999, 1},
		{1.5, 2},
		{2, 2},
		{0, 0},
		{-1.5, -2},
	}

	for i, c := range testCases {
		actual := lround(c.input)
		if actual != c.expected {
			t.Errorf("[%v] Something(%v) == %v, but we expected %v", i, c.input, actual, c.expected)
		}
	}
}

func TestSomething(t *testing.T) {
	/*var testCases = []struct {
		s, expected string
	}{
		{"yo", "yes"},
		{"ya", "no"},
	}

	for _, c := range testCases {
		actual := Something(c.s)
		if actual != c.expected {
			t.Errorf("Something(%q) == %q, but we expected %q", c.s, actual, c.expected)
		}
	}*/
}

func BenchmarkSomething(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func TestSquare(t *testing.T) {

	var testCases = []struct {
		input []float64
		base int
		expected []int
	}{
		{[]float64 {2, 3}, 10, []int {4, 2, 0, 1}},
		{[]float64 {5, 4, 3, 2}, 10,[]int {5, 2, 0, 9, 9, 4, 5}},

		{[]float64 {985, 39, 701, 991, 506, 649, 805, 28, 255, 313, 936, 887, 598, 493, 69, 628, 287,
			118, 157, 134, 791, 602, 195, 59, 254, 381, 508, 573, 525, 233, 765, 606, 446, 268, 260, 682,
			596, 820, 939, 628, 830, 117, 269, 7},
			1000,
			[]int {225, 800, 568, 330, 527, 947, 673, 70, 532, 231, 132, 380, 154, 532, 467, 286,
				7, 336, 279, 385, 941, 814, 539, 744, 630, 667, 353, 35, 599, 852, 696, 515, 228, 352,
				873, 346, 123, 328, 83, 231, 850, 216, 389, 780, 674, 911, 200, 363, 318, 871, 927, 783,
				625, 859, 577, 652, 71, 469, 455, 382, 444, 810, 467, 957, 29, 673, 59, 421, 267, 437,
				516, 930, 204, 377, 498, 318, 856, 811, 253, 188, 228, 584, 567, 35, 74, 840, 52}},
	}

	for i, c := range testCases {
		actual, _ := square(c.input, c.base)

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("[%v] Something(%v) == %v, but we expected %v", i, c.input, actual, c.expected)
		}
	}
}
