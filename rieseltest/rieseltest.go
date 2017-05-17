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

func genV1rodseth(n, h int64) (int64, error) {

	// Check if h is not a multiple of 3
	if h % 3 != 0 {
		return 4, nil
	}

	// Compute N as a Big num
	N := big.NewInt(n)
	N.Exp(big.NewInt(2), N, nil)
	N.Mul(big.NewInt(h), N)
	N.Sub(N, big.NewInt(1))

	// Check if there is a p which satisfies Rodseth conditions
	var p int64
	for p = 3; p < math.MaxInt64; p++ {
		if big.Jacobi(big.NewInt(p - 2), N) == 1 && big.Jacobi(big.NewInt(p + 2), N) == -1 {
			return p, nil
		}
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
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