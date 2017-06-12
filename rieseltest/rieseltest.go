package rieseltest

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/op/go-logging"
	"github.com/mjibson/go-dsp/fft"
	"os"
	"fmt"
	"time"
	"math"
	"errors"
	"math/big"
)

func IsPrime(h, n int64) bool {
	if lbit, err := lowerNonZeroBit(h); err == nil && lbit > 0 {
		n += int64(lbit)
		h >>= lbit
	}

	N := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(h), new(big.Int).Exp(big.NewInt(2), big.NewInt(n), nil)), big.NewInt(1))

	v1, _ := GenV1(h, n, RODSETH)

	u2 := GenU2(N, h, n, v1)
	uN := GenUN(N, h, n, u2)

	if uN.Cmp(big.NewInt(0)) == 0 {
		return true
	} else {
		log.Debug("n = %v*2^%v-1 lead to u(n)=%v", h, n, uN)
	}

	return false
}

func lround(a float64) int {
	if a < 0 {
		return int(math.Ceil(a - 0.5))
	}
	return int(math.Floor(a + 0.5))
}

// Initialize logger to be used in this package
var log = logging.MustGetLogger("rieseltest")

func init() {
	// Set up the terminal and file loggers with a custom format and backend
	terminalFormat := logging.MustStringFormatter("%{color}[%{shortfunc}]%{color:reset} %{time:15:04:05.000} %{color}%{level:.4s}%{color:reset} -> %{message}")
	terminalBackend := logging.NewLogBackend(os.Stdout, "", 0)

	// Let the file logs be saved on a rolling basis
	fileFormat := logging.MustStringFormatter("[%{pid} - %{shortfunc}] %{time:15:04:05.000} %{level:.4s} -> %{message}")
	fileBackend := logging.NewLogBackend(&lumberjack.Logger{
		Filename:   fmt.Sprintf(".logs/logFile.%d%.2d%.2d.log",
			time.Now().Year(), time.Now().Month(), time.Now().Day()),
		MaxSize:    100, // megabytes
		MaxBackups: 2,
		MaxAge:     3, //days
	}, "", 0)

	// Only notice messages and more severe messages should be sent to the file backend
	fileBackendLeveled := logging.AddModuleLevel(fileBackend)
	fileBackendLeveled.SetLevel(logging.NOTICE, "")
	fileBackendFormatter := logging.NewBackendFormatter(fileBackendLeveled, fileFormat)

	// Only notice messages and more severe messages should be sent to the file backend
	terminalBackendLeveled := logging.AddModuleLevel(terminalBackend)
	terminalBackendLeveled.SetLevel(logging.NOTICE, "")
	terminalBackendFormatter := logging.NewBackendFormatter(terminalBackendLeveled, terminalFormat)

	// Set the two backends for the loggers
	logging.SetBackend(terminalBackendFormatter, fileBackendFormatter)
}

// TODO use larger base (up to half of the CPU numeric precision is good) <-

// Square performs a square on the input polynomial (assumed to be represented with the given base) using FFT
// Useful resources to understand how this works:
// http://numbers.computation.free.fr/Constants/Algorithms/fft.html
// http://icereijo.com/fft-based-multiplication-algorithm/
// http://www.cs.rug.nl/~ando/pdfs/Ando_Emerencia_multiplying_huge_integers_using_fourier_transforms_paper.pdf
// https://math.stackexchange.com/questions/764727/concrete-fft-polynomial-multiplication-example
func square(input []float64, base int) ([]int, error) {
	// TODO verify that Landon's whiteboard test succeeds

	// Pad the coefficient vectors of input to at least length 2n - 1, using zeros
	padWithZeros(&input)

	log.Debugf("Input as real values: %v", input)

	// Compute the FFT of the input
	frequency := fft.FFTReal(input)

	log.Debugf("Frequency (output of FFTReal): %v", frequency)

	// Apply the dot product
	for i, el := range frequency {
		frequency[i] = el * el
	}

	log.Debugf("Frequency squared: %v", frequency)

	// Compute the inverse FFT
	squared := fft.IFFT(frequency)

	log.Debugf("Inverse frequency: %v", squared)

	// Keep only the real part of the output (and round it)
	output := make([]int, len(squared))
	for i, el := range squared {
		output[i] = lround(real(el))

		// Verify that complex part of the output rounds to 0 for each element
		if lround(imag(el)) != 0 {
			return nil, errors.New("The imaginary part of some number didn't round to 0")
		}
	}

	log.Debugf("Output: %v", output)

	propagateCarry(&output, base)

	log.Debugf("Scaled output: %v", output)

	return output, nil
}

func padWithZeros(number *[]float64) {
	// TODO consider speed improvement by adding n instead of n-1 zeros (factor 2)

	for i, max := 0, len(*number); i < max - 1; i++ {
		*number = append(*number, 0)
	}
}

func propagateCarry(number *[]int, base int) {
	carry := 0

	// Propagate the carry
	for i := 0; i < len(*number); i++ {
		(*number)[i] += carry
		carry = (*number)[i] / base
		(*number)[i] -= carry * base
	}

	if carry != 0 {
		// Propagate the extra carry
		for carry != 0. {
			x := carry
			carry = x / base
			*number = append(*number, x - carry * base)
		}

	} else {
		// Remove the trailing zeros
		i := len(*number)
		for ; (*number)[i - 1] == 0.; i-- {}
		*number = (*number)[:i]
	}
}

const (
	RIESEL uint8 = iota
	RODSETH
	PENNE
)

func GenV1(h, n int64, method uint8) (int64, error) {

	// TODO MAKE H ODD

	// Check if h is not a multiple of 3
	if hmod3 := h % 3; hmod3 != 0 {

		// When k == 1 mod 3 AND n is even or k == 2 mod 3 and n is odd, then 3 is a factor
		if !(h == 1 && n == 2) && (((hmod3 == 1) && (n&1 == 0)) || ((hmod3 == 2) && (n&1 == 1))) {
			return -1, errors.New(fmt.Sprintf("N = %v*2^%v-1 is a multiple of 3", h, n))
		}

		// In the other cases, we have that v(1) = 4
		return 4, nil
	}

	if method == RIESEL {
		return genV1Riesel(h, n)
	} else if method == RODSETH {
		return genV1Rodseth(h, n)
	} else if method == PENNE {
		return genV1Penne(h, n)
	} else {
		return -1, errors.New("The specified method to generate v1 is not valid")
	}
}

func efficient_jacobi(x, h, n int64, NEqualsMod8, checkDone bool) int {
	equalsModQ := func(x int64, modValues ...int64) bool {
		for _, m := range modValues {
			if x == m { return true }
		}
		return false
	}

	sign := true

	// While x is even
	if x & 1 == 0 {

		hMod8 := h % 8

		// When h == 0 (mod 8), then N == dred - 1 (mod dred),
		// which means that Jacobi(dred,N) = 1, not valid
		if checkDone == false  && hMod8 != 0 {
			twoNMod8 := ModExp(2, n, 8)
			NMod8 := (hMod8*twoNMod8 - 1 + 8) % 8
			NEqualsMod8 = equalsModQ(NMod8 , 3, 5)
			checkDone = true
		}

		for x & 1 == 0 {
			if NEqualsMod8 { sign = !sign }
			x /= 2
		}
	}

	hModX := h % x
	var jNx int

	// TODO add explanation for this
	if hModX == 0 {
		jNx = 1

	} else {
		twoNModX := ModExp(2, n, x)
		NModX := (hModX*twoNModX - 1 + x) % x

		jNx = big.Jacobi(big.NewInt(NModX), big.NewInt(x))
		if (x % 4) == 3 { sign = !sign }
	}

	if sign == true {
		return jNx
	} else {
		return -jNx
	}
}

func genV1Rodseth(h, n int64) (int64, error) {
	NEqualsMod8 := false
	checkDone := false

	queue := make(chan int, 4)

	// Only need 4 values
	jacobi_minus := func (x, h, n int64) int {
		if len(queue) < 4 {
			return efficient_jacobi(x, h, n, NEqualsMod8, checkDone)
		} else {
			return <-queue
		}
	}

	jacobi_plus := func (x, h, n int64) int {
		j := efficient_jacobi(x, h, n, NEqualsMod8, checkDone)
		queue <- j
		return j
	}

	// Check if there is a p which satisfies Rodseth conditions
	var p int64
	for p = 3; p < math.MaxInt64; p++ {
		if jacobi_minus(p - 2, h, n) == 1 && jacobi_plus(p + 2, h, n) == -1 {
			return p, nil
		}
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
}

func genV1Riesel(h, n int64) (int64, error) {

	// Check if there is a v which satisfies Riesel conditions
	var v int64
	cache := make(map[int64]int)

	for v = 3; ; v++ {
		_, d, err := reduce(v * v - 4)

		if err != nil {
			fmt.Println(err)
			continue
		}

		var dred int64
		if (d & 1) == 1 {
			dred = d
		} else {
			dred = d >> 1
		}

		var sign bool
		if (n > 2) || (d & 1 == 1) {
			sign = true
		} else {
			sign = false
		}

		hmodd := h % dred

		// When h == 0 (mod dred), then N == dred - 1 (mod dred),
		// which means that Jacobi(dred,N) = 1, not valid
		if hmodd == 0 {
			continue
		}

		if (n > 1) && ((((dred - 1) / 2) & 1) == 1) {
			sign = !sign
		}

		var jNd int
		if val, ok := cache[dred]; ok {
			jNd = val

		} else {
			twonModd := ModExp(2, n, dred)
			Nmodd := (hmodd * twonModd - 1 + dred) % dred

			if (Nmodd == 0) && (dred != 1) {
				return -1, errors.New("N has a known factor, it does not need to be tested further.")
			}

			jNd = big.Jacobi(big.NewInt(Nmodd), big.NewInt(dred))
			// TODO add the check that a and b are coprime

			cache[dred] = jNd
		}

		if ((sign == true) && (jNd == 1)) || ((sign == false) && (jNd == -1)) {
			continue
		}

		if X := issquare(v - 2); X != 0 {
			return v, nil
		} else {
			ared := v - 2
			i := 0

			for (ared & 1) == 0 {
				ared >>= 1
				i++
			}

			sign = false
			if !((n > 2) || (((v - 2) & 1) == 1) || ((i & 1) == 0)) {
				sign = true
			}

			hmoda := h % ared
			if (n > 1) && ((((ared - 1) / 2) & 1) == 1) {
				sign = !sign
			}


			var jNa int
			if val, ok := cache[ared]; ok {
				jNa = val

			} else {
				twonModa := ModExp(2, n, ared)
				Nmoda := (hmoda * twonModa - 1 + ared) % ared

				if (Nmoda == 0) && (ared != 1) {
					return -1, errors.New("N has a known factor, it does not need to be tested further.")
				}

				jNa = big.Jacobi(big.NewInt(Nmoda), big.NewInt(ared))
				// TODO add the check that a and b are coprime

				cache[ared] = jNa
			}

			if ((sign == true) && (jNa == 1)) || ((sign == false) && (jNa == -1)) {
				continue
			}

			return v, nil
		}
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
}

func genV1Penne(h, n int64) (int64, error) {
	// Check if there is a v which satisfies Riesel conditions
	var x int64
	cache := make(map[int64]int)

	for x = 1; ; x++ {
		_, d, err := reduce(x * x + 4)

		if err != nil {
			fmt.Println(err)
			continue
		}

		var dred int64
		if (d & 1) == 1 {
			dred = d
		} else {
			dred = d >> 1
		}

		var sign bool
		if (n > 2) || (d & 1 == 1) {
			sign = true
		} else {
			sign = false
		}

		hmodd := h % dred

		// When h == 0 (mod dred), then N == dred - 1 (mod dred),
		// which means that Jacobi(dred,N) = 1, not valid
		if hmodd == 0 {
			continue
		}

		if (n > 1) && ((((dred - 1) / 2) & 1) == 1) {
			sign = !sign
		}

		var jNd int
		if val, ok := cache[dred]; ok {
			jNd = val

		} else {
			twonModd := ModExp(2, n, dred)
			Nmodd := (hmodd * twonModd - 1 + dred) % dred

			if (Nmodd == 0) && (dred != 1) {
				return -1, errors.New("N has a known factor, it does not need to be tested further.")
			}

			jNd = big.Jacobi(big.NewInt(Nmodd), big.NewInt(dred))
			// TODO add the check that a and b are coprime

			cache[dred] = jNd
		}

		if ((sign == true) && (jNd == 1)) || ((sign == false) && (jNd == -1)) {
			continue
		}

		return x * x + 2, nil
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
}

func reduce(x int64) (b int64, d int64, e error) {
	// Reduce a discriminant to a square free integer.
	// Given x, compute d, without square factor, and b, such as x = d*b^2

	d, b = x, 1

	if x < 4 {
		return -1, -1, errors.New("x needs to be greater than 4")
	}

	for d % 4 == 0 {			// Divide by even power of two.
		d /= 4
		b *= 2
	}

	for div := int64(3); ; div += 2 {

		sq := div * div

		if sq > d {
			break
		}

		for d % sq == 0 {	// Divide by even powers of odd factors.
			d /= sq
			b *= div
		}
	}

	return b, d, nil
}

func issquare(n int64) (s int64) {
	// This function returns the square root of an integer square, or zero.
	s = int64(lround(math.Floor(math.Sqrt(float64(n)))))

	if s*s == n {
		return s
	} else {
		return 0
	}
}

func ModExp(base, exponent, modulus int64) int64 {
	if modulus == 1 { return 0 }
	result := int64(1)
	base = base % modulus

	for ; exponent > 0; exponent >>= 1 {

		if exponent&1 == 1 {
			result = (result * base) % modulus
		}

		base = (base * base) % modulus
	}

	return result
}

func bitLen(n int64) int {
	bitLen := 0
	for ; n != 0; n = n >> 1 {
		bitLen++
	}

	return bitLen
}

func lowerNonZeroBit(n int64) (uint, error) {
	ret := uint(0)

	if n == 0 {
		return 0, errors.New("Given a zero input, which does not have any set bit")
	}

	for n&1 == 0 && n > 0 {
		ret += 1
		n >>= 1
	}

	return ret, nil
}

func bit(n int64, index uint) bool {
	return ((n >> index) & 1) == 1
}

func rieselMod(N, v *big.Int, h, n int64) *big.Int {
	// Returns r % (h * 2^n - 1)

	// Check h, n positive integers
	// Make h odd if not:

	// TODO benchmark and why?
	if N.Cmp(big.NewInt(math.MaxInt64)) == -1 {
		ret := new(big.Int).Mod(v, N)
		return ret

	} else {
		ret := v
		for ret.Cmp(N) == 1 {

			if int64(ret.BitLen()) - 1 < n {
				break
			}

			j := new(big.Int).Rsh(ret, uint(n))
			k := new(big.Int).Sub(ret, new(big.Int).Lsh(j, uint(n)))

			if h == 1 {
				ret.Add(k, j)
			} else {
				tquo := new(big.Int)
				tmod := new(big.Int)
				tquo.DivMod(j, big.NewInt(h), tmod)

				ret.Add(new(big.Int).Add(new(big.Int).Lsh(tmod, uint(n)), k), tquo)
			}
		}

		if ret.Sign() == -1 {
			ret.Add(ret, N)
			return ret
		} else if ret.Cmp(N) == 0 {
			return big.NewInt(0)
		} else {
			return ret
		}
	}
}

func rieselModCh(N, v *big.Int, h, n int64, c chan *big.Int) {
	c <- rieselMod(N, v, h, n)
}

func GenU2(N *big.Int, h, n, v1 int64) *big.Int {

	// Check sanity of arguments / preconditions, like that h is positive odd, v positive, n gt or eq to 2

	v1_big := big.NewInt(v1)

	efficient_2n_plus_one := func(a, b *big.Int, c chan *big.Int) {
		rieselModCh(N, new(big.Int).Sub(new(big.Int).Mul(a, b), v1_big), h, n, c)
	}

	efficient_2n := func(a *big.Int, c chan *big.Int) {
		rieselModCh(N, new(big.Int).Sub(new(big.Int).Mul(a, a), big.NewInt(2)), h, n, c)
	}

	r := v1_big

	if h == 1 {
		return rieselMod(N, r, h, n)
	}

	// s := v1^2 - 2
	s := new(big.Int).Mul(r, r)
	s = s.Sub(s, big.NewInt(2))

	c_r := make(chan *big.Int)
	c_s := make(chan *big.Int)

	// BitLen counts also the last one 0 so it's like an array from 0 to 4 means 5, and we want to start from 3
	for i := bitLen(h) - 2; i > 0; i-- {

		if bit(h,uint(i)) {
			go efficient_2n_plus_one(r, s, c_r)
			go efficient_2n(s, c_s)

			r = <- c_r
			s = <- c_s

		} else {
			go efficient_2n_plus_one(r, s, c_s)
			go efficient_2n(r, c_r)

			s = <- c_s
			r = <- c_r
		}
	}

	r = rieselMod(N, new(big.Int).Sub(new(big.Int).Mul(r, s), v1_big), h, n)

	// v(i + 1) = v(1) * v(i) - v(i - 1)
	// v(2i) = v(i)^2 - 2
	// v(2i + 1) = v(i) * v(i + 1) - v(1)
	return r
}

func GenUN(N *big.Int, h, n int64, u *big.Int) *big.Int {
	// Check sanity of arguments / preconditions, like that h is positive odd, v positive, n gt or eq to 2

	for i := int64(3); i <= n; i++ {
		u = rieselMod(N, new(big.Int).Sub(new(big.Int).Mul(u, u), big.NewInt(2)), h, n)
	}

	return u
}