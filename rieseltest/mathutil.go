package rieseltest

import (
	"errors"
	"math"
	"fmt"

	// big "math/big"
	big "github.com/ricpacca/gmp"
	// big "github.com/ricpacca/go.flint/fmpz"
)

var zero = new(big.Int).SetInt64(0)
var one = new(big.Int).SetInt64(1)
var two = new(big.Int).SetInt64(2)
var maxInt64 = new(big.Int).SetInt64(math.MaxInt64)

// lowerNonZeroBit returns the position of the lower non zero bit of a given number.
// The position is counted from least significant bit starting from 0.
//
// This function requires:
//		a) n >= 1
//
// Examples:
//		lowerNonZeroBit(11) = 0, since 11 = 0b1011
//		lowerNonZeroBit(18) = 1, since 18 = 0b10010
//		lowerNonZeroBit(24) = 3, since 18 = 0b11000
func lowerNonZeroBit(n int64) (uint, error) {

	// Check preconditions
	if n < 1 {
		return 0, errors.New(fmt.Sprintf("Expected n >= 1, but received n = %v", n))
	}

	ret := uint(0)

	// Count +1 from the right-most bit until we find a 1
	for n&1 == 0 {
		ret += 1
		n >>= 1
	}

	return ret, nil
}

// bitLen returns the number of bits of the given positive number.
//
// This function requires:
//		a) n >= 0
//
// Examples:
//		bitLen(0) = 0
//		bitLen(11) = 4, since 11 = 0b1011
//		bitLen(18) = 5, since 18 = 0b10010
//		bitLen(24) = 5, since 18 = 0b11000
func bitLen(n int64) (uint, error) {

	// Check preconditions
	if n < 0 {
		return 0, errors.New(fmt.Sprintf("Expected n >= 0, but received n = %v", n))
	}

	bitLen := uint(0)

	// Count +1 from the right-most bit until n is 0
	for n != 0 {
		bitLen++
		n >>= 1
	}

	return bitLen, nil
}

// bit returns true if the index-th bit of the given n is 1.
// The index of the right-most bit is 0.
//
// Example with 11 (0b1011):
//		bitLen(11, 0) = true
//		bitLen(11, 1) = true
//		bitLen(11, 2) = false
//		bitLen(11, 3) = true
//		bitLen(11, 4) = false
func bit(n int64, index uint) bool {
	return ((n >> index) & 1) == 1
}

// lround rounds the given float number a to the nearest integer value.
//
// Examples:
//		lround(1.1) = 1
//		lround(-1.1) = -1
//		lround(1.5) = 2
//		lround(-1.5) = -2
func lround(a float64) int64 {
	if a < 0 {
		return int64(math.Ceil(a - 0.5))
	}

	return int64(math.Floor(a + 0.5))
}

// reduce reduces the given x to a square free integer.
// In other words: given x, it computes d, without square factor, and b, such as:
// 		x = d*b^2
//
// NOTE: This function has been copied from the LLR software's "reduce" [Ref3].
func reduce(x int64) (b int64, d int64) {

	// Handle small easy cases
	if x < 4 {
		return x, 1
	}

	d, b = x, 1

	// Divide d by even powers of two.
	for d % 4 == 0 {
		d /= 4
		b *= 2
	}

	for div := int64(3); ; div += 2 {
		sq := div * div

		if sq > d {
			break
		}

		// Divide d by even powers of increasing odd factors.
		for d % sq == 0 {
			d /= sq
			b *= div
		}
	}

	return b, d
}

// isPerfectSquare returns true if n is a perfect square.
//
// This function requires:
//		a) n > 0
func isPerfectSquare(n int64) (bool, error) {

	// Check preconditions
	if n <= 0 {
		return false, errors.New(fmt.Sprintf("Expected n > 0, but received n = %v", n))
	}

	s := lround(math.Sqrt(float64(n)))

	if s*s == n {
		return true, nil
	} else {
		return false, nil
	}
}

// modExp performs a modular exponentiation to compute (base^exponent mod modulus)
// using the right-to-left binary method.
//
// This function requires:
//		a) modulus > 0
//		b) exponent >= 0
//
// This implementation is based on the pseudocode found on:
// 		en.wikipedia.org/wiki/Modular_exponentiation
func modExp(base, exponent, modulus int64) (int64, error) {

	// Check preconditions
	if modulus <= 0 {
		return 0, errors.New(fmt.Sprintf("Expected modulus > 0, but received modulus = %v", modulus))
	}
	if exponent < 0 {
		return 0, errors.New(fmt.Sprintf("Expected exponent >= 0, but received exponent = %v", exponent))
	}

	if base == 0 || modulus == 1 { return 0, nil }
	result := int64(1)
	base = base % modulus

	for exponent > 0 {
		if exponent&1 == 1 {
			result = (result * base) % modulus
		}

		exponent >>= 1
		base = (base * base) % modulus
	}

	return result, nil
}

// rieselMod computes (a mod N), where N = (h * 2^n - 1) in an efficient
// way using the shift and add method.
//
// Read [Ref4] for more information on this method.
func rieselMod(a *big.Int, R *RieselNumber) {
	if R.N.Cmp(maxInt64) == -1 {
		a.Mod(a, R.N)

	} else {
		for a.Cmp(R.N) == 1 {
			if int64(a.BitLen()) <= R.n {
				break	// QUESTION: is this code ever reached?
			}

			j := new(big.Int)
			j.Rsh(a, uint(R.n))

			k := new(big.Int)
			k.Sub(a, k.Lsh(j, uint(R.n)))

			if R.h == 1 {
				a.Add(k, j)
			} else {
				tquo := new(big.Int)
				tmod := new(big.Int)

				tquo.DivMod(j, R.hBig, tmod)
				a.Add(a.Add(tmod.Lsh(tmod, uint(R.n)), k), tquo)
			}
		}

		if a.Cmp(R.N) == 0 {
			a.SetInt64(0)
		}
	}
}

// efficientJacobi efficiently computes the Jacobi symbol for (x, h * 2^n - 1)
//
// This method is faster than directly calling big.Int.Jacobi(a, b) when b is of
// the form h*2^n-1, because it contains some optimizations derived from the
// properties of the Jacobi symbol for that specific case.
//
// The cache argument is used when we call this method multiple times with
// the same values of h and n, to avoid performing the same operations over and over.
//
// If you are calling this method once, you can set it to 'nil' in the method invocation.
//
// If you are calling this method multiple times with the same h and n, but different x, you can
// initialize a cache item as a variable in the caller function (through the newEfficientJacobiCache
// function), and pass its pointer to the repeated invocations of efficientJacobi.
func efficientJacobi(x, h, n int64, cache map[int64]int) (int, error) {

	// true == +1
	sign := true

	// While x is even, we have:
	// 		Jacobi(x, N) == Jacobi(2, N) * Jacobi(x/2, N)
	//
	// First of all, we can write:
	//		Jacobi(2,N) = (-1)^((N-1)*(N+1)/8) =
	// 		= (-1)^((h*2^n-2)*(h*2^n)/8) = (-1)^((h*2^(n-1)-1)*h*2^(n-2))
	//
	// And notice that Jacobi(2,N) == -1 only when (n == 2).
	// In fact, when n != 2, the term 2^(n-2) is always even and any number
	// multiplied by an even number is even, leading to a positive power in our case.
	//
	// Our sign will be as follows:
	// sign = Jacobi(x, N) / Jacobi(x_reduced, N) = (Jacobi(2, N)) * ... * (Jacobi(2, N))
	//
	// Thus, we can divide x by 2 until it is odd and flip the sign in all those cases
	// when (Jacobi(2, N)) == -1.
	for (x & 1) == 0 {
		x >>= 1		// a = a / 2
		if n == 2 { sign = !sign }
	}

	// At this point, we know that x is odd, and we can
	// proceed with computing Jacobi(x, N).
	if hModX := h % x; hModX == 0 {

		// When h == 0 (mod x), then N == -1 (mod x), and therefore,
		// keeping in mind that n >= 2, we can write:
		//
		// Jacobi(x, N) = Jacobi(N, x) * (-1)^((N-1)*(x-1)/4) =
		// = Jacobi(-1, x) * (-1)^((h*2^n-2)/2*(x-1)/2) =
		// = (-1)^((x-1)/2) * (-1)^((h*2^(n-1)-1)*(x-1)/2) =
		//
		// Now, since (h*2^(n-1)-1) is always odd, we have that:
		// 		when (x-1)/2 is even:
		// 			(-1)^((x-1)/2) == +1
		// 			(-1)^((h*2^(n-1)-1)*(x-1)/2) == +1
		// 		when (x-1)/2 is odd:
		// 			(-1)^((x-1)/2) == -1
		// 			(-1)^((h*2^(n-1)-1)*(x-1)/2) == -1
		//
		// In both of the cases, since the signs are concordant, we have that:
		// (-1)^((x-1)/2) * (-1)^((h*2^(n-1)-1)*(x-1)/2) = +1
		//
		// Thus, Jacobi(x, N) = 1.
		//
		// We just need to take into consideration the sign difference due to
		// having previously reduced x to odd.
		if sign == true {
			return 1, nil
		} else {
			return -1, nil
		}

	} else {

		// This is the general case.
		//
		// Given that we know that x is odd, we can write:
		// 		Jacobi(x, N) = Jacobi(N, x) * (-1)^((N-1)/2*(x-1)/2) =
		// 		= Jacobi(((h mod x) * (2^n mod x) - 1) mod x, x) * (-1)^((h*2^(n-1)-1)*(x-1)/2)
		//
		// Let's start by computing Jacobi(N, x) or getting it from the computedJNX
		// cache if it was already computed.
		var jNx int
		if val, ok := cache[x]; ok {
			jNx = val

		} else {

			// Jacobi(N, x) = Jacobi(((h mod x) * (2^n mod x) - 1) mod x, x).
			twoNModX, err := modExp(2, n, x)
			if err != nil {
				return 0, errors.New("Something went wrong: n or x were negative")
			}

			NModX := (hModX * twoNModX - 1 + x) % x

			// Check if x divides N (just in case)
			if (NModX == 0) && (x != 1) {
				return 0, errors.New("N has a known factor, it does not need to be tested further.")
			}

			// Now we can compute Jacobi(N, x) on smaller numbers
			jNx = big.Jacobi(new(big.Int).SetInt64(NModX), new(big.Int).SetInt64(x))

			// Check if GCD(N, x) != 1
			// If GCD(N, x) != 1, then N has a divisor > 1, and does not need to be tested further.
			if jNx == 0 {
				return 0, errors.New("N has a known factor, it does not need to be tested further.")
			}

			// Store the computed Jacobi(N, x) in the cache for potential later use
			if cache != nil {
				cache[x] = jNx
			}
		}

		// The last thing that we have to do is to compute the value of (-1)^((N-1)/2*(x-1)/2).
		//
		// Now, since (h*2^(n-1)-1) is always odd, the result depends on whether (x-1)/2 is even or odd.
		//
		// Reminding that x is odd, we observe that when:
		//		(x mod 4 == 1) -> (x-1)/2 is even -> (-1)^((h*2^(n-1)-1)*(x-1)/2) = +1
		//		(x mod 4 == 3) -> (x-1)/2 is odd -> (-1)^((h*2^(n-1)-1)*(x-1)/2) = -1
		if (x % 4) == 3 { sign = !sign }

		// Jacobi(x, N) = Jacobi(N, x) * sign
		if sign == true {
			return jNx, nil
		} else {
			return -jNx, nil
		}
	}
}

// screenEasyPrimes checks whether R is a small known prime < 257 or whether it
// has a small known prime < 257 as a factor.
//
// It returns:
//		+1 if R is a prime < 257
//		-1 if R has a prime < 257 as a factor
//		 0 otherwise
func screenEasyPrimes(R *RieselNumber) (int, error) {

	// Check preconditions
	if R == nil {
		return 0, errors.New("Received R == nil")
	}

	// Catch the degenerate case of h*2^n-1 == 1
	if R.h == 1 && R.n == 1 {
		log.Debugf("N = %v = 1 is not prime", R)
		return -1, nil	// 1*2^1-1 == 1 is not prime
	}

	// Catch the degenerate case of n == 2
	//
	// n == 2 and 0 < h < 2^n  ->  0 < h < 4
	// Since h is odd  ->  h == 1 or h == 3
	if R.h == 1 && R.n == 2 {
		log.Debugf("N = %v = 3 is prime", R)
		return +1, nil		// 1*2^2-1 == 3 is prime
	}
	if R.h == 3 && R.n == 2 {
		log.Debugf("N = %v = 11 is prime", R)
		return +1, nil		// 3*2^2-1 == 11 is prime
	}

	// Catch small primes < 257
	//
	// We check for only a few primes because the other primes < 257
	// violate the checks above.
	if R.h == 1 {
		if R.n == 3 || R.n == 5 || R.n == 7 {
			log.Debugf("N = %v is prime", R)
			return +1, nil		// 3, 7, 31, 127 are prime
		}
	}
	if R.h == 3 {
		if R.n == 2 || R.n == 3 || R.n == 4 || R.n == 6 {
			log.Debugf("N = %v is prime", R)
			return +1, nil		// 11, 23, 47, 191 are prime
		}
	}
	if R.h == 5 && R.n == 4 {
		log.Debugf("N = %v = 79 is prime", R)
		return +1, nil 	// 79 is prime
	}
	if R.h == 7 && R.n == 5 {
		log.Debugf("N = %v = 223 is prime", R)
		return +1, nil		// 223 is prime
	}
	if R.h == 15 && R.n == 4 {
		log.Debugf("N = %v = 239 is prime", R)
		return +1, nil		// 239 is prime
	}

	// Check for 3 <= prime factors < 29
	// The product of primes up to 28, excluding 2 is:
	// 		111546435
	if new(big.Int).GCD(nil, nil, R.N, new(big.Int).SetInt64(111546435)).Cmp(one) != 0 {
		log.Debugf("N = %v is not prime: a small 3 <= prime < 29 divides it", R)
		return -1, nil	// a small 3 <= prime < 29 divides N
	}

	// Check for 29 <= prime factors < 47
	// The product of primes from 28 to 46 is:
	//		5864229
	if new(big.Int).GCD(nil, nil, R.N, new(big.Int).SetInt64(58642669)).Cmp(one) != 0 {
		log.Debugf("N = %v is not prime: a small 29 <= prime < 47 divides it", R)
		return -1, nil 	// a small 29 <= prime < 47 divides N
	}

	// Check for 47 <= prime factors < 257, if N is large
	// 2^282 > pfact(256)/pfact(46) > 2^281
	if bits := R.N.BitLen(); bits - 1 >= 281 {

		pprod256, success := new(big.Int).
			SetString("4912291013238638017062389731791584291410159591853190162192019099864799926800582498341",
			10)

		if success != true {
			log.Warning("Was not able to initialize pprod256, the big.Int product of " +
				"primes between 47 and 256. Skipping this pre-check.")

		} else if new(big.Int).GCD(nil, nil, R.N, pprod256).Cmp(one) != 0 {
			log.Debugf("N = %v is not prime: a small 47 <= prime < 257 divides it", R)
			return -1, nil	// a small 47 <= prime < 257 divides N
		}
	}

	return 0, nil
}