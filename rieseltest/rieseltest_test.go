package rieseltest

import (
	"testing"
	"os"
	"bufio"
	"strings"
	"strconv"

	// big "math/big"
	big "github.com/arcetri/gmp"
	// big "github.com/arcetri/go.flint/fmpz"
)

func TestGenV1SimpleCase(t *testing.T) {

	// Test for all the known Riesel primes with h not multiple of 3 if generating V(1)
	// works as expected (returning 4).
	if file, err := os.Open("testfiles/v1_with_h_NOT_multiple_of_3.out"); err == nil {

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

			R, _ := NewRieselNumber(int64(h), int64(n))
			actual, _ := GenV1(R, RODSETH)

			if actual != int64(4) {
				t.Errorf("genV1Riesel(%v, %v) == %v, but we expected 4", h, n, actual)
			}
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

func TestGenV1Rodseth(t *testing.T) {

	// Test for all the known Riesel primes with h multiple of 3 if generating V(1) works
	// as expected. We use the calc software [Ref5] to generate the test cases.
	if file, err := os.Open("testfiles/v1_with_h_multiple_of_3.out"); err == nil {

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

			expected, err := strconv.Atoi(words[2])
			if err != nil {
				panic(err)
			}

			actual, _ := genV1Rodseth(int64(h), int64(n))
			if actual != int64(expected) {
				t.Errorf("genV1Riesel(%v, %v) == %v, but we expected %v", h, n, actual, expected)
			}
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

func TestGenV1Riesel(t *testing.T) {

	// Test for all the known Riesel primes with h multiple of 3 if generating V(1) works
	// as expected. We use the calc software [Ref5] to generate the test cases.
	if file, err := os.Open("testfiles/v1_with_h_multiple_of_3.out"); err == nil {

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

			expected, err := strconv.Atoi(words[2])
			if err != nil {
				panic(err)
			}

			actual, _ := genV1Riesel(int64(h), int64(n))
			if expected != 29 && expected != 59 && actual != int64(expected) {
				t.Errorf("genV1Riesel(%v, %v) == %v, but we expected %v", h, n, actual, expected)
			}
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

func TestGenU2(t *testing.T) {

	// Test for some known Riesel primes if generating U(2) works as expected.
	// We use the calc software [Ref5] to generate the test cases.
	k, _ := new(big.Int).SetString("11795713261792426815798380820887418752674276338130891767997058623856897193630" +
		"44019617626496399595785831029617354807397622866359370536685450127496104007738667927950678092481408087840326" +
		"76788707308924861149775649032641075553276613952032825786169935399015100878575000138517502280577878565639243" +
		"68684968218923849458923447972573135862423723904333441872330786434091582216595399208036188109565367105073050" +
		"19628466572645823801440681649807420172660094126573000139887515466774346436001698389165752785493408275313705" +
		"70646904337526123314470090772404343585123729255003531576149287549750613692081411866875443165835695629173342" +
		"472423326462", 10)

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
		R, _ := NewRieselNumber(c.h, c.n)
		v1, _ := GenV1(R, RODSETH)
		actual, _ := GenU2(R, v1)

		if actual.Cmp(c.expected) != 0 {
			t.Errorf("GenU2(%v, %v) == %v, but we expected %v", R, v1, actual, c.expected)
		}
	}
}

func TestGenUNSingle(t *testing.T) {
	t.Skip()	// normally skip the test because it takes a long time to run, enable when useful

	// Test for some known Riesel primes if generating U(n) mod N works as expected.
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
		R, _ := NewRieselNumber(c.h, c.n)
		v1, _ := GenV1(R, RODSETH)
		u2, _ := GenU2(R, v1)

		actual, _ := GenUN(R, u2)

		if actual.Cmp(c.expected) != 0 {
			t.Errorf("GenUN(%v, %v) == %v, but we expected %v", R, u2, actual, c.expected)
		}
	}
}

// region benchmarks

func BenchmarkGenV1Riesel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if file, err := os.Open("testfiles/v1_with_h_multiple_of_3_large.out"); err == nil {

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
		if file, err := os.Open("testfiles/v1_with_h_multiple_of_3_large.out"); err == nil {

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
		if file, err := os.Open("testfiles/v1_with_h_multiple_of_3_large.out"); err == nil {

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
		if file, err := os.Open("testfiles/h_n_large_primes.out"); err == nil {

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

				N, _ := NewRieselNumber(int64(h), int64(n))

				v1, err := GenV1(N, RIESEL)
				if err != nil {
					panic(err)
				}

				GenU2(N, v1)
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
		if file, err := os.Open("testfiles/h_n_large_primes.out"); err == nil {

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

				N, _ := NewRieselNumber(int64(h), int64(n))

				v1, err := GenV1(N, RODSETH)
				if err != nil {
					panic(err)
				}

				GenU2(N, v1)
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
		if file, err := os.Open("testfiles/h_n_large_primes.out"); err == nil {

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

				N, _ := NewRieselNumber(int64(h), int64(n))

				v1, err := GenV1(N, PENNE)
				if err != nil {
					panic(err)
				}

				GenU2(N, v1)
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

func BenchmarkGenU2WithV13(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if file, err := os.Open("testfiles/v1_with_h_NOT_multiple_of_3_and_v1_3_large.out"); err == nil {

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

				N, _ := NewRieselNumber(int64(h), int64(n))

				GenU2(N, 3)
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

func BenchmarkGenU2WithV14(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if file, err := os.Open("testfiles/v1_with_h_NOT_multiple_of_3_and_v1_3_large.out"); err == nil {

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

				N, _ := NewRieselNumber(int64(h), int64(n))

				GenU2(N, 4)
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