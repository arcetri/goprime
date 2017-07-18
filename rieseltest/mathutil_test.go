package rieseltest

import (
	"testing"

	big "math/big"
	// big "github.com/arcetri/gmp"
	// big "github.com/arcetri/go.flint/fmpz"
)

func TestLowerNonZeroBit(t *testing.T) {
	var testCases = []struct {
		input int64
		expected uint
	}{
		{11, 0},
		{18, 1},
		{24, 3},
	}

	for _, c := range testCases {
		actual, _ := lowerNonZeroBit(c.input)
		if actual != c.expected {
			t.Errorf("lowerNonZeroBit(%v) == %v, but we expected %v", c.input, actual, c.expected)
		}
	}
}

func TestLowerNonZeroBitError(t *testing.T) {
	var testCases = []struct {
		input int64
	}{
		{0},
		{-1},
		{-24},
	}

	for _, c := range testCases {
		_, err := lowerNonZeroBit(c.input)
		if err == nil {
			t.Errorf("lowerNonZeroBit(%v) should return an error, but it didn't.", c.input)
		}
	}
}

func TestBitLen(t *testing.T) {
	var testCases = []struct {
		input int64
		expected uint
	}{
		{0, 0},
		{11, 4},
		{18, 5},
		{24, 5},
	}

	for _, c := range testCases {
		actual, _ := bitLen(c.input)
		if actual != c.expected {
			t.Errorf("bitLen(%v) == %v, but we expected %v", c.input, actual, c.expected)
		}
	}
}

func TestBitLenError(t *testing.T) {
	var testCases = []struct {
		input int64
	}{
		{-1},
		{-24},
	}

	for _, c := range testCases {
		_, err := bitLen(c.input)
		if err == nil {
			t.Errorf("bitLen(%v) should return an error, but it didn't.", c.input)
		}
	}
}

func TestBit(t *testing.T) {
	var testCases = []struct {
		input    int64
		index    uint
		expected bool
	}{
		{11, 0, true},
		{11, 1, true},
		{11, 2, false},
		{11, 3, true},
		{11, 4, false},
	}

	for _, c := range testCases {
		actual := bit(c.input, c.index)
		if actual != c.expected {
			t.Errorf("bit(%v, %v) == %v, but we expected %v", c.input, c.index, actual, c.expected)
		}
	}
}

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

func TestReduce(t *testing.T) {
	var testCases = []struct {
		x int64
		b int64
		d int64
	}{
		{24, 2, 6},
		{38, 1, 38},
		{27, 3, 3},
		{3, 3, 1},
		{415800, 30, 462},
	}

	for _, c := range testCases {
		b, d := reduce(c.x)
		if c.b != b && c.d != d {
			t.Errorf("reduce(%v) == (%v, %v), but we expected (%v, %v)", c.x, b, d, c.b, c.d)
		}
	}
}

func TestIsPerfectSquare(t *testing.T) {
	perfectSquares := map[int64]bool{4: true, 9: true, 16: true, 25: true, 36: true, 49: true, 64: true, 81: true, 100: true}

	for n := int64(-10); n < 100; n++ {
		is, err := isPerfectSquare(n)

		if n < 0 && err == nil {
			t.Errorf("sqrtOrZero(%v) should return an error, but it didn't", n)
		} else if perfectSquares[n] && err != nil && is == false {
			t.Errorf("sqrtOrZero(%v) = false, but we expected true", n)
		} else if perfectSquares[n] == false && err != nil && is == true {
			t.Errorf("sqrtOrZero(%v) = true, but we expected false", n)
		}
	}
}

func TestModExp(t *testing.T) {
	var testCases = []struct {
		b int64
		e int64
		m int64
		expected int64
	}{
		{111, 123, 53, 35},
		{12, 9, 1, 0},
		{0, 981409, 132421, 0},
		{1000, 1000, 19, 7},
		{12983194, 0, 19, 1},
		{-81792, 73363, 233, -161},
	}

	for _, c := range testCases {
		actual, _ := modExp(c.b, c.e, c.m)
		if actual != c.expected {
			t.Errorf("modExp(%v, %v, %v) == %v, but we expected %v", c.b, c.e, c.m, actual, c.expected)
		}
	}
}

func TestModExpError(t *testing.T) {
	var testCases = []struct {
		b int64
		e int64
		m int64
	}{
		{11, -2, 5},
		{7, -1, 120},
		{24412, 919, 0},
		{1, 132, -1},
	}

	for _, c := range testCases {
		_, err := modExp(c.b, c.e, c.m)
		if err == nil {
			t.Errorf("modExp(%v, %v, %v) should return an error, but it didn't", c.b, c.e, c.m)
		}
	}
}

func TestRieselMod(t *testing.T) {
	var testCases = []struct {
		a string
		h int64
		n int64
		expected string
	}{
		{"191561942608236107294793378393788647952342390272950271", 1, 177, "0"},
		{"383123885216472214589586756787577295904684780545900542", 1, 177, "0"},
		{"1524155677489", 13, 17, "1155404"},
		{"67386085206301762878672178026014945965432808571819096426538514623471147922216230988104" +
			"51014310011205180773952647674517450540237811212921539702288863686237481195029174312884" +
			"121757753830014976", 45, 415, "724449805572665047374217432817204615614" +
			"632579287918225362906109975099326556922295221602684123768506623631360442535074691561376"},
	}

	for _, c := range testCases {
		a, _ := new(big.Int).SetString(c.a, 10)
		R, _ := NewRieselNumber(c.h, c.n)
		rieselMod(a, R)
		expected, _ := new(big.Int).SetString(c.expected, 10)

		if a.Cmp(expected) != 0 {
			t.Errorf("rieselMod(%v, %v) = %v, but we expected %v", c.a, R, a, expected)
		}
	}
}

func TestEfficientJacobi(t *testing.T) {
	var testCases = []struct {
		x int64
		h int64
		n int64
		expected int
	}{
		{9, 507, 217588, 1},
		{131233, 41231029, 217523, -1},
		{66, 742329, 2, -1},
		{131230, 742329, 20523, -1},
	}

	for _, c := range testCases {
		actual, _ := efficientJacobi(c.x, c.h, c.n, nil)

		if actual != c.expected {
			t.Errorf("efficientJacobi(%v, %v, %v, nil) = %v, but we expected %v", c.x, c.h, c.n, actual, c.expected)
		}
	}
}
