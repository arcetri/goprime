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

type BigNum struct {
	base int
	digits string
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

func genV1(h, n int64) (int64, error) {

	// Check if h is not a multiple of 3
	if hmod3 := h % 3; hmod3 != 0 {

		// When k == 1 mod 3 AND n is even or k == 2 mod 3 and n is odd, then 3 is a factor
		if ((hmod3 == 1) && (n&1 == 0)) || ((hmod3 == 2) && (n&1 == 1)) {
			return -1, errors.New("N = h*2^n-1 is a multiple of 3")
		}

		// In the other cases, we have that v(1) = 4
		return 4, nil
	}

	return GenV1rodseth(h, n)
}


var NEqualsMod8 = false
var checkDone = false

func efficient_jacobi(x, h, n int64) int {
	// TODO caching

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
			twoNMod8 := new(big.Int).Exp(big.NewInt(2), big.NewInt(n), big.NewInt(8)).Int64()
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
		twoNModX := new(big.Int).Exp(big.NewInt(2), big.NewInt(n), big.NewInt(x)).Int64()
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

type stack []int

func (s stack) Push(v int) stack {
	return append(s, v)
}

func (s stack) PopFromTop() (stack, int) {
	// FIXME: What do we do if the stack is empty, though?
	return  s[1:], s[0]
}

func GenV1rodseth(h, n int64) (int64, error) {
	var s stack

	// Only need 4 values
	jacobi_minus := func (x, h, n int64) int {
		if len(s) < 4 {
			return efficient_jacobi(x, h, n)
		} else {
			var ret int
			s, ret = s.PopFromTop()
			return ret
		}
	}

	jacobi_plus := func (x, h, n int64) int {
		j := efficient_jacobi(x, h, n)
		s = s.Push(j)
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

func GenV1riesel(h, n int64) (int64, error) {

	// Check if there is a p which satisfies Rodseth conditions
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
			twonModd := new(big.Int).Exp(big.NewInt(2), big.NewInt(n), big.NewInt(dred)).Int64()
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
				twonModa := new(big.Int).Exp(big.NewInt(2), big.NewInt(n), big.NewInt(ared)).Int64()
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

func issquare (n int64) (s int64) {
	// This function returns the square root of an integer square, or zero.
	s = int64(lround(math.Floor(math.Sqrt(float64(n)))))

	if s*s == n {
		return s
	} else {
		return 0
	}
}

/*func computeV(n, h int) int {

	if h == 1 {
		if n % 4 == 3 {
			return 3
		} else if n % 2 == 1 {
			return 4
		}
	}

	if h == 3 {
		if m := n % 4; m == 0 || m == 3 {
			return 5778
		}
	}


}

func computeU(prevU, h, n int) {

	squaredU := square(prevU)
}*/