package rieseltest

import (
	"testing"
	"reflect"
	"os"
	"bufio"
	"strings"
	"strconv"
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

func TestGenV1Rodseth(t *testing.T) {
	var testCases = []struct {
		h, n, expected int64
	}{
		// only primes
		{1095, 2992587, 5},
		{3338480145, 257127, 21},
		{111546435, 257139, 27},
		{2084259, 1257787, 15},
		{8331405, 1984565, 9},
		{165, 2207550, 17},
		{1155, 1082878, 11},
	}

	for _, c := range testCases {
		actual, _ := GenV1rodseth(c.h, c.n)
		if actual != c.expected {
			t.Errorf("Something(%v, %v) == %v, but we expected %v", c.h, c.n, actual, c.expected)
		}
	}
}

func benchmarkGenV1rodseth(h int64, n int64, b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenV1rodseth(h, n)
	}
}

func BenchmarkV1Rodseth1(b *testing.B) { benchmarkGenV1rodseth(1095, 2992587, b) }
func BenchmarkV1Rodseth2(b *testing.B) { benchmarkGenV1rodseth(3338480145, 257127, b) }
func BenchmarkV1Rodseth3(b *testing.B) { benchmarkGenV1rodseth(111546435, 257139, b) }
func BenchmarkV1Rodseth4(b *testing.B) { benchmarkGenV1rodseth(2084259, 1257787, b) }
func BenchmarkV1Rodseth5(b *testing.B) { benchmarkGenV1rodseth(8331405, 1984565, b) }
func BenchmarkV1Rodseth6(b *testing.B) { benchmarkGenV1rodseth(165, 2207550, b) }
func BenchmarkV1Rodseth7(b *testing.B) { benchmarkGenV1rodseth(1155, 1082878, b) }

func TestGenV1Riesel(t *testing.T) {
	var testCases = []struct {
		h, n, expected int64
	}{
		// only primes
		{1095, 2992587, 5},
		{3338480145, 257127, 21},
		{111546435, 257139, 27},
		{2084259, 1257787, 15},
		{8331405, 1984565, 9},
		{165, 2207550, 17},
		{1155, 1082878, 11},
	}

	for _, c := range testCases {
		actual, _ := GenV1riesel(c.h, c.n)
		if actual != c.expected {
			t.Errorf("Something(%v, %v) == %v, but we expected %v", c.h, c.n, actual, c.expected)
		}
	}
}

func benchmarkGenV1riesel(h int64, n int64, b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenV1riesel(h, n)
	}
}

func BenchmarkV1Riesel1(b *testing.B) { benchmarkGenV1riesel(1095, 2992587, b) }
func BenchmarkV1Riesel2(b *testing.B) { benchmarkGenV1riesel(3338480145, 257127, b) }
func BenchmarkV1Riesel3(b *testing.B) { benchmarkGenV1riesel(111546435, 257139, b) }
func BenchmarkV1Riesel4(b *testing.B) { benchmarkGenV1riesel(2084259, 1257787, b) }
func BenchmarkV1Riesel5(b *testing.B) { benchmarkGenV1riesel(8331405, 1984565, b) }
func BenchmarkV1Riesel6(b *testing.B) { benchmarkGenV1riesel(165, 2207550, b) }
func BenchmarkV1Riesel7(b *testing.B) { benchmarkGenV1riesel(1155, 1082878, b) }

func BenchmarkGenV1rieselFull(b *testing.B) {

	for i := 0; i < b.N; i++ {
		if file, err := os.Open("h-n.verified.txt"); err == nil {

			// create a new scanner and read the file line by line
			s := bufio.NewScanner(file)

			for s.Scan() {
				line := s.Text()
				words := strings.Split(line, " ")
				h, err := strconv.Atoi(words[0])
				if err != nil {
					panic(err)
				}

				n, err := strconv.Atoi(words[1])
				if err != nil {
					panic(err)
				}

				GenV1riesel(int64(h), int64(n))
			}

			// check for errors
			if err = s.Err(); err != nil {
				panic(err)
			}

			file.Close()

		} else {
			panic(err)
		}
	}
}

func BenchmarkGenV1rodsethFull(b *testing.B) {

	for i := 0; i < b.N; i++ {
		if file, err := os.Open("h-n.verified.txt"); err == nil {

			// create a new scanner and read the file line by line
			s := bufio.NewScanner(file)

			for s.Scan() {
				line := s.Text()
				words := strings.Split(line, " ")
				h, err := strconv.Atoi(words[0])
				if err != nil {
					panic(err)
				}

				n, err := strconv.Atoi(words[1])
				if err != nil {
					panic(err)
				}

				GenV1rodseth(int64(h), int64(n))
			}

			// check for errors
			if err = s.Err(); err != nil {
				panic(err)
			}

			file.Close()

		} else {
			panic(err)
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
