package rieseltest

import (
	"testing"
	"reflect"
	"os"
	"bufio"
	"strings"
	"strconv"
	big "github.com/ncw/gmp"
	"math"
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
		actual, _ := genV1Rodseth(c.h, c.n)
		if actual != c.expected {
			t.Errorf("Something(%v, %v) == %v, but we expected %v", c.h, c.n, actual, c.expected)
		}
	}
}

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
		actual, _ := genV1Riesel(c.h, c.n)
		if actual != c.expected {
			t.Errorf("Something(%v, %v) == %v, but we expected %v", c.h, c.n, actual, c.expected)
		}
	}
}

func TestGenU2Single(t *testing.T) {
	k, _ := new(big.Int).SetString("117957132617924268157983808208874187526742763381308917679970586238568971936304401961762649639959578583102961735480739762286635937053668545012749610400773866792795067809248140808784032676788707308924861149775649032641075553276613952032825786169935399015100878575000138517502280577878565639243686849682189238494589234479725731358624237239043334418723307864340915822165953992080361881095653671050730501962846657264582380144068164980742017266009412657300013988751546677434643600169838916575278549340827531370570646904337526123314470090772404343585123729255003531576149287549750613692081411866875443165835695629173342472423326462", 10)

	var testCases = []struct {
		h, n int64
		expected *big.Int
	}{
		// only primes
		{15, 5, big.NewInt(91)},
		{9, 7, big.NewInt(473)},
		{375, 9, big.NewInt(157186)},
		{105, 8, big.NewInt(3749)},
		{57, 8, big.NewInt(4037)},
		{8565, 15, big.NewInt(244410061)},
		{315, 10, big.NewInt(4015)},
		{507, 217588, k},
	}

	for _, c := range testCases {
		v1, _ := GenV1(c.h, c.n, RODSETH)
		N := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(c.h), new(big.Int).
			Exp(big.NewInt(2), big.NewInt(c.n), nil)), big.NewInt(1))

		actual := GenU2(N, c.h, c.n, v1)
		if actual.Cmp(c.expected) != 0 {
			t.Errorf("GenU2(%v, %v, %v) == %v, but we expected %v", c.h, c.n, v1, actual, c.expected)
		}
	}
}

func TestGenUNSingle(t *testing.T) {
	t.Skip()

	zero := big.NewInt(0)

	var testCases = []struct {
		h, n int64
		expected *big.Int
	}{
		// only primes
		{15, 5, zero},
		{9, 7, zero},
		{375, 9, zero},
		{105, 8, zero},
		{57, 8, zero},
		{8565, 15, zero},
		{315, 10, zero},
	}

	for _, c := range testCases {
		v1, _ := GenV1(c.h, c.n, RODSETH)
		N := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(c.h), new(big.Int).
			Exp(big.NewInt(2), big.NewInt(c.n), nil)), big.NewInt(1))

		u2 := GenU2(N, c.h, c.n, v1)
		actual := GenUN(N, c.h, c.n, u2)

		if actual.Cmp(c.expected) != 0 {
			t.Errorf("GenUN(%v, %v, %v) == %v, but we expected %v", c.h, c.n, u2, actual, c.expected)
		}
	}
}

func BenchmarkGenV1Riesel(b *testing.B) {

	for i := 0; i < b.N; i++ {

		if file, err := os.Open("filter_odd_multiple_of_3_h.out"); err == nil {

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

				genV1Riesel(int64(h), int64(n))
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

func BenchmarkGenV1Rodseth(b *testing.B) {

	for i := 0; i < b.N; i++ {
		if file, err := os.Open("filter_odd_multiple_of_3_h.out"); err == nil {

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

				genV1Rodseth(int64(h), int64(n))
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

func BenchmarkGenV1Penne(b *testing.B) {

	for i := 0; i < b.N; i++ {

		if file, err := os.Open("filter_odd_multiple_of_3_h.out"); err == nil {

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

				genV1Penne(int64(h), int64(n))
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

func BenchmarkGenU2Riesel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if file, err := os.Open("filter_odd_h_large.out"); err == nil {

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

				v1, err := GenV1(int64(h), int64(n), RIESEL)
				if err != nil {
					panic(err)
				}

				N := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(int64(h)), new(big.Int).
					Exp(big.NewInt(2), big.NewInt(int64(n)), nil)), big.NewInt(1))

				GenU2(N, int64(h), int64(n), v1)
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

func BenchmarkGenU2Rodseth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if file, err := os.Open("filter_odd_h_large.out"); err == nil {

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

				v1, err := GenV1(int64(h), int64(n), RODSETH)
				if err != nil {
					panic(err)
				}

				N := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(int64(h)), new(big.Int).
					Exp(big.NewInt(2), big.NewInt(int64(n)), nil)), big.NewInt(1))

				GenU2(N, int64(h), int64(n), v1)
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

func BenchmarkGenU2Penne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if file, err := os.Open("filter_odd_h_large.out"); err == nil {

			// create a new scanner and read the file line by line
			s := bufio.NewScanner(file)

			for s.Scan() {
				line := s.Text()
				words := strings.Split(line, " ")
				h_raw, err := strconv.Atoi(words[0])
				if err != nil {
					panic(err)
				}

				n_raw, err := strconv.Atoi(words[1])
				if err != nil {
					panic(err)
				}

				h, n := int64(h_raw), int64(n_raw)

				v1, err := GenV1(h, n, PENNE)
				if err != nil {
					panic(err)
				}

				N := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(h), new(big.Int).
					Exp(big.NewInt(2), big.NewInt(n), nil)), big.NewInt(1))

				GenU2(N, h, n, v1)
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

func BenchmarkComputation(b *testing.B) {

	for i := 0; i < b.N; i++ {
		big.NewInt(math.MaxInt64)
	}
}
