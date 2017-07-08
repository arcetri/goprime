// +build linux

// Package rieseltest implements the functions to perform a full
// primality test on Riesel numbers of the form R = h * 2^n - 1.

package rieseltest

import (
	"fmt"
	"math"
	"errors"
	"syscall"

	// Mathematical library implementing the necessary methods
	// big "math/big"
	big "github.com/ricpacca/gmp"
	// big "github.com/ricpacca/go.flint/fmpz"
)

func init() {
	// By default, disable the loggers
	ConfigureLogger(false, 0,false, 0)
}

// IsPrime performs a full Lucas-Lehmer-Riesel primality test on the
// given Riesel number R = h * 2^n - 1, and returns true if the number is prime.
//
// The test works as follows:
// 		1) Generate V(1)
// 		2) From V(1) generate U(2) = V(h)
// 		3) From U(2) generate U(n)
// 		4) N is prime is U(n) == 0 (mod N)
//
// The following conditions must be true for the test to work:
//		a) n >= 2
//		b) h >= 1
func IsPrime(R *RieselNumber) (bool, error) {

	// Check preconditions
	if R.h < 1 {
		return false, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", R.h))
	}
	if R.n < 2 {
		return false, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", R.n))
	}

	// Check if N is a small prime or a multiple of a small prime
	if check, err := screenEasyPrimes(R); err == nil && check != 0 {
		if check == 1 {
			log.Infof("N = %v is a known prime < 257", R)
			return true, nil
		}

		log.Infof("N = %v has a known factor < 257", R)
		return false, nil
	}

	// Step 1: Get a V(1) for the Riesel candidate.
	//
	// The 'RIESEL' and 'RODSETH' methods are equivalent.
	// The 'PENNE' method can be faster but finds a higher V(1),
	// which might slow down the following steps of the test.
	v1, err := GenV1(R, RODSETH)
	if err != nil { return false, err }
	log.Infof("Generated V(1) = %v", v1)

	// Step 2: Use the generated V(1) to generate U(2) = V(h)
	u2, err := GenU2(R, v1)
	if err != nil { return false, err }
	if loggingEnabled { log.Infof("Generated U(2) = V(h). Last 8 digits: %v", getLastDigits(u2)) }

	// Step 3: Use the generated U(2) to generate U(n)
	uN, err := GenUN(R, u2)
	if err != nil { return false, err }
	if loggingEnabled { log.Infof("Generated U(n). Last 8 digits: %v", getLastDigits(uN)) }

	// Step 4: Check if U(n) == 0 (mod N)
	if uN.Cmp(zero) == 0 {
		log.Infof("N = %v is prime!", R)
		return true, nil
	} else {
		log.Infof("N = %v is composite!", R)
		return false, nil
	}
}

// GenV1 available algorithms
const (
	RIESEL uint8 = iota
	RODSETH
	PENNE
)

// GenV1 computes a valid V(1) value for the given Riesel candidate.
//
// This function assumes:
//		a) n >= 2
//		b) h >= 1
//		c) h mod 2 == 1
//
// The generation of V(1) depends on the value of h. There are two cases
// to consider:
// 		1) h mod 3 != 0
// 		2) h mod 3 == 0
//
// CASE 1:	(h mod 3 != 0)
//
// This case is easy. In [Ref1], page 869, one finds that if:
//		a) h mod 6 == +/-1
//		b) N mod 3 != 0
//
// which translates, given that we assumed odd h, into the condition:
//		c) h mod 3 != 0
//
// if this case condition is true then (also page 869):
// 		U(2) = (2 + sqrt(3))^h + (2 - sqrt(3))^h =
//			 = (2 + sqrt(3))^h + (2 + sqrt(3))^(-h)
// [ this is because it is easy to show that: (2 - sqrt(3)) = (2 + sqrt(3))^(-1) ]
//
// Since [Ref1] states:
//		V(i) = alpha^i + alpha^(-i)
//
// and since U(2) = V(h), we can simply let:
//		alpha = (2 + sqrt(3))
//
// thus, in case 1, we return:
//		V(1) = alpha^1 + alpha^(-1) =
//			 = (2 + sqrt(3)) - ((2 - sqrt(3)) = 4
//
// CASE 2:	(h mod 3 == 0):
//
// This case is more complicated and its explanation depends on the algorithm
// chosen (Riesel, Rodseth or Penne).
//
// NOTE: Even though CASE 2 could work for any h, we use it only when h mod 3 == 0,
// because CASE 1 is faster and thus we prefer to use that when h mod 3 != 0.
func GenV1(R *RieselNumber, method uint8) (int64, error) {

	// Check preconditions
	if R.h < 1 {
		return -1, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", R.h))
	}
	if R.n < 2 {
		return -1, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", R.n))
	}
	if R.h % 2 == 0 {
		return -1, errors.New(fmt.Sprintf("Expected odd h, but received h = %v", R.h))
	}

	// Check if h is not a multiple of 3
	if hmod3 := R.h % 3; hmod3 != 0 {

		// Screen easy composites where 3 is a factor.
		// It is easy to show that when:
		//
		// 		(h mod 3 == 1 AND n is even) OR
		// 		(h mod 3 == 2 AND n is odd),
		//
		// then 3 is a factor.
		//
		// This relies on the observation that:
		//		2^(2k) ==  +1 (mod 3)
		//		2^(2k+1) == -1 (mod 3)
		if ((hmod3 == 1) && (R.n&1 == 0)) || ((hmod3 == 2) && (R.n&1 == 1)) {
			return -1, errors.New(fmt.Sprintf("N = %v is a multiple of 3", R))
		}

		// In all these cases, we have that v(1) = 4
		log.Debugf("h = %v is a multiple of 3, thus V(1) = 4", R.h)
		return 4, nil
	}

	// Handle the cases when h is a multiple of 3
	if method == RIESEL {
		return genV1Riesel(R.h, R.n)
	} else if method == RODSETH {
		return genV1Rodseth(R.h, R.n)
	} else if method == PENNE {
		return genV1Penne(R.h, R.n)
	} else {
		return -1, errors.New("The specified method to generate v1 is not valid")
	}
}

// genV1Rodseth computes a valid V(1) value for the given Riesel candidate.
//
// This function assumes:
//		a) n >= 2
//		b) h >= 1
//		c) h mod 2 == 1
//
// This method is the easiest one, and simply sets V(1) to the first value P that satisfies:
//		a) Jacobi(P-2, h*2^n-1) == 1
//		b) Jacobi(P+2, h*2^n-1) == -1
//
// Read [Ref2] for a theoretical explanation of this method, and the comments
// below for an explanation of the optimizations implemented.
func genV1Rodseth(h, n int64) (int64, error) {

	// Check preconditions
	if h < 1 {
		return -1, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", h))
	}
	if n < 2 {
		return -1, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", n))
	}
	if h % 2 == 0 {
		return -1, errors.New(fmt.Sprintf("Expected odd h, but received h = %v", h))
	}

	// OPTIMIZATION: Store a cache of already computed Jacobi symbols.
	//
	// The cache will work as follows:
	// For every P we check if Jacobi(P-2, N) == +1.
	// If Jacobi(P-2, N) != +1, we move on to the next P.
	// If Jacobi(P-2, N) == +1, we check if Jacobi(P+2, N) == -1. If that is verified, we have found
	// a valid P and we return it, but if that fails, it means that Jacobi(P+2, N) == +1.
	//
	// Four iterations later, when we are considering P'=P+4, our first check will be
	// Jacobi(P'-2, N) == +1. Notice that P'-2 == P+2. Therefore we might have already computed
	// Jacobi(P'-2, N), and in that case, we would have had Jacobi(P'-2, N) == +1, so our first
	// check is done without the need to compute another Jacobi symbol.
	//
	// Thus, the information that we need is whether we had already computed that value and failed
	// the check. If we reached a following iteration, that means that the check would have failed
	// anyway (otherwise our loop would have finished), so eventually we only need to know if that
	// value was computed (and had failed) or it was not computed.
	//
	// The cache will be used to store this information for the Jacobi symbols of the P+2 cases.
	cache := make(map[int64]bool)

	// This cache is used by the efficientJacobi function when we do repeated calls of it.
	efficientJacobiCache := make(map[int64]int)

	// This function is used to compute the Jacobi(P-2, N) case.
	jacobi_minus := func (x, h, n int64) (int, error) {

		// Check in the cache if that symbol was already computed before.
		// If it was computed for a P'=P-4, that means that it must have been Jacobi(P'+2, N) == 1,
		// otherwise the loop would not have reached this point but it would have instead returned a
		// valid V(1). Thus, if it had been computed before, we simply return 1.
		//
		// Otherwise, if it was not computed (this might have happened if the first condition failed),
		// then we simply compute it and return the result. We do not even need to store the fact that
		// we have computed it in the cache, because we will not need to compute Jacobi(P-2) again for
		// other increasing values of P.
		if _, ok := cache[x]; ok {
			log.Debugf("Retrieved Jacobi(%v, N) from the cache", x)

			// We won't need this value anymore, so we can remove it
			// from the cache to keep it smaller.
			delete(cache, x)
			return 1, nil

		} else {
			return efficientJacobi(x, h, n, efficientJacobiCache)
		}
	}

	// This function is used to compute the Jacobi(P+2, N) case.
	jacobi_plus := func (x, h, n int64) (int, error) {

		// Since we might need to compute Jacobi(P+2, N) again for a P'=P+4 later,
		// we store the fact that we have computed it in the cache before returning it.
		cache[x] = true

		return efficientJacobi(x, h, n, efficientJacobiCache)
	}

	// Check if there is a P which satisfies Rodseth conditions.
	//
	// QUESTIONS:
	// 1) Will it ever be the case that the first P verifying these conditions is P > 2^64 - 1?
	// 2) Could we safely consider only odd values of P?
	for P := int64(3); P < math.MaxInt64; P++ {
		var j_minus int
		var j_plus int
		var err error

		// Compute Jacobi(P - 2, N) and check for condition 1
		if j_minus, err = jacobi_minus(P - 2, h, n); err == nil && j_minus == 1 {
			log.Debugf("Jacobi(%v - 2, N) == 1: 1st condition passed", P)

			// Compute Jacobi(P + 2, N) and check for condition 2
			if j_plus, err = jacobi_plus(P + 2, h, n); err == nil && j_plus == -1 {
				log.Debugf("Jacobi(%v + 2, N) == -1: 2nd condition passed", P)
				return P, nil
			}
		}

		// Return an error when we have found a factor of N
		if err != nil {
			return -1, err
		}
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
}

// genV1Riesel computes a valid V(1) value for the given Riesel candidate in the case
// where h is a multiple of 3.
//
// This function requires:
//		a) n >= 2
//		b) h >= 1
//		c) h mod 2 == 1
//
// The implementation of this method is very similar to the implementation of the
// LLR software from Jean Penne, but contains a few optimizations.
//
// Please read the references [Ref1] and TODO for a comprehensive explanation of the theory
// underlying this method.
func genV1Riesel(h, n int64) (int64, error) {

	// Check preconditions
	if h < 1 {
		return -1, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", h))
	}
	if n < 2 {
		return -1, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", n))
	}
	if h % 2 == 0 {
		return -1, errors.New(fmt.Sprintf("Expected odd h, but received h = %v", h))
	}

	// OPTIMIZATION:
	// Store a cache of already computed Jacobi symbols.
	cache := make(map[int64]int)

	for v := int64(3); ; v++ {

		// D is the square free part of (v^2 - 4)
		_, D := reduce(v * v - 4)

		// Set dred = odd part of D
		//
		// Since D is square free, the exact power of 2 dividing it can only be 0 or 1.
		// Therefore, when D is even, we can get its odd part by simply dividing it by 2.
		// When D is odd, we don't need to divide it (since dividing it by 2^0 = 1 is same as not dividing).
		//
		// We compute dred for later computing Jacobi(D, N), because it is more efficient to
		// compute Jacobi(N, D) thanks to the properties of the Jacobi symbol. However, if D
		// is even, we need to reduce it to an odd dred. That is not a problem, because we will
		// later take into account the sign difference between Jacobi(D, N) and Jacobi(N, dred)
		// and compute the former through the latter.
		var dred int64
		if (D & 1) == 1 {
			dred = D
		} else {
			dred = D >> 1
		}

		// Since we use dred instead of D, we have to take into account the sign difference
		// between Jacobi(D, N) and Jacobi(dred, N).
		//
		// When D is odd, dred == D, and therefore there is no sign difference between
		// Jacobi(D, N) and Jacobi(dred, N).
		//
		// When n > 2, sign = +1 as well. In fact, when D is even and n > 2, we have that:
		//
		// sign = Jacobi(D, N) / Jacobi(dred, N) = Jacobi(2, N) = (-1)^((N-1)*(N+1)/8) =
		// (-1)^((h*2^n-2)*(h*2^n)/8) = (-1)^((h*2^(n-1)-1)*h*2^(n-2))
		//
		// --> sign = +1 if D is odd OR n > 2
		//     sign = -1 if D is even and n == 2
		var sign bool
		if (D & 1 == 1) || (n > 2) {
			sign = true
		} else { 			// D is even and n == 2
			sign = false
		}

		if hmodd := int64(h) % dred; hmodd == 0 && sign == true {

			// When h == 0 (mod dred), then N == -1 (mod dred), and therefore,
			// keeping in mind that n >= 2, we can write:
			//
			// Jacobi(dred, N) = Jacobi(N, dred) * (-1)^((N-1)*(dred-1)/4) =
			// = Jacobi(-1, dred) * (-1)^((h*2^n-2)/2*(dred-1)/2) =
			// = (-1)^((dred-1)/2) * (-1)^((h*2^(n-1)-1)*(dred-1)/2) =
			//
			// Now, since (h*2^(n-1)-1) is always odd, we have that:
			// 		when (dred-1)/2 is even:
			// 			(-1)^((dred-1)/2) == +1
			// 			(-1)^((h*2^(n-1)-1)*(dred-1)/2) == +1
			// 		when (dred-1)/2 is odd:
			// 			(-1)^((dred-1)/2) == -1
			// 			(-1)^((h*2^(n-1)-1)*(dred-1)/2) == -1
			//
			// In both of the cases, since the signs are concordant, we have that:
			// (-1)^((dred-1)/2) * (-1)^((h*2^(n-1)-1)*(dred-1)/2) = +1
			//
			// Thus, Jacobi(dred, N) = 1.
			//
			// We can immediately verify if the current candidate for V(1) is valid:
			//      when sign == false, then Jacobi(D, N) = -1  ->  Check VALID
			//      when sign == true, then Jacobi(D, N) = 1  -> Check NOT VALID
			log.Debugf("[C1] Jacobi(%v, N) = 1: %v is not a valid candidate for V(1)", D, v)
			continue

		} else if hmodd != 0 {

			// Count the sign difference between Jacobi(dred, N) and Jacobi(N, dred).
			// Using the expression in the above comment, we know that:
			//
			// Jacobi(dred, N) / Jacobi(N, dred) = (-1)^((N-1)/2*(dred-1)/2) =
			// = (-1)^((h*2^(n-1)-1)*(dred-1)/2)
			//
			// Now, since (h*2^(n-1)-1) is always odd, the sign difference depends
			// on whether (dred-1)/2 is even or odd.
			//
			// Reminding that dred is odd, we observe that when:
			//		(dred mod 4 == 1) -> (dred-1)/2 is even -> (-1)^((h*2^(n-1)-1)*(dred-1)/2) = +1
			//		(dred mod 4 == 3) -> (dred-1)/2 is odd -> (-1)^((h*2^(n-1)-1)*(dred-1)/2) = -1
			if (dred % 4) == 3 { sign = !sign }

			// Compute Jacobi(N, dred) or get it from the cache if it was already computed.
			var jNd int
			if val, ok := cache[dred]; ok {
				log.Debugf("Retrieved Jacobi(N, %v) from the cache", dred)
				jNd = val

			} else {

				// Jacobi(N, dred) = Jacobi(h * 2^n - 1, dred) =
				// Jacobi(((h mod dred) * (2^n mod dred) - 1) mod dred, dred)
				twonModd, err := modExp(2, n, dred)
				if err != nil {
					return 0, errors.New("Something went wrong: n or dred were negative")
				}

				Nmodd := (hmodd * twonModd - 1 + dred) % dred

				// Check if dred divides N (just in case)
				if (Nmodd == 0) && (dred != 1) {
					return -1, errors.New(fmt.Sprintf("%v divides N: N does not need to be tested " +
						"further.", dred))
				}

				jNd = big.Jacobi(new(big.Int).SetInt64(Nmodd), new(big.Int).SetInt64(dred))

				// Check if GCD(N, dred) != 1
				// If GCD(N, dred) != 1, then N has a divisor > 1, and does not need to be tested further.
				if jNd == 0 {
					return -1, errors.New("N has a known factor, it does not need to be tested further.")
				}

				// Store the computed Jacobi(N, dred) in the cache for potential later use
				cache[dred] = jNd
			}

			// Verify that the resulting Jacobi(D, N) == -1
			//
			// This is the first condition of the Riesel theorem. If the candidate
			// V(1) does not satisfy it, we can move to the next candidate.
			if ((sign == true) && (jNd == 1)) || ((sign == false) && (jNd == -1)) {
				log.Debugf("[C1] Jacobi(%v, N) = 1: %v is not a valid candidate for V(1)", D, v)
				continue
			}
		}

		// If we reached this point, it means that Jacobi(D, N) == -1.
		//
		// Now, we check if we are in the case where alpha = epsilon^2, by verifying if
		// V(1) − 2 is a perfect square. If we are in this case, we do not need to
		// verify any further condition, and our candidate V(1) is valid already.
		//
		// If V(1) − 2 is not a perfect square, then we are in the case when
		// alpha = epsilon, and we need to verify a further condition instead.
		//
		// Read TODO paper for a better explanation of this
		if issquare, err := isPerfectSquare(v - 2); issquare && err == nil {
			log.Debugf("%v-2 is a perfect square -> alpha = epsilon^2 -> %v is a valid V(1) candidate", v, v)
			return v, nil

		} else if err != nil {
			return 0, errors.New("Something went wrong: (v - 2) was negative")

		} else {

			// Now we need to verify the second condition of Riesel's theorem, i.e.:
			// Jacobi(r, N) * (a^2 - b^2 * D) / r == -1  [C2]
			//
			// As TODO explains, it can be shown that r = 4a, and therefore we can write:
			// Jacobi(r, N) = Jacobi(4a, N) = Jacobi(4, N) * Jacobi(a, N) = Jacobi(a, N)
			// where either a = v + 2 or a = v - 2.
			//
			// Since it is equivalent to choose any of them, we choose a = v - 2 which is smaller.
			//
			// Now, to compute Jacobi (a, N), we will use some optimizations similar to the ones that
			// we used to compute Jacobi(N, D).
			//
			// If a is even, we want to reduce it to an odd number. We keep count of how many times we
			// divide a by 2 to be able to later take into account the sign difference between
			// Jacobi (a, N) and Jacobi (ared, N).
			a := v - 2
			ared := a
			i := 0

			// NOTE: a is not square free as D was, so it might be necessary to
			// divide it by two more than once.
			for (ared & 1) == 0 {
				ared >>= 1		// a = a / 2
				i++
			}

			// TODO Ref shows that:
			// Jacobi(r, N) * (a^2 - b^2 * D) / r = Jacobi(r, N) * sgn(a^2 − b^2 * D)
			//
			// Moreover, for our a, we have that:
			// a^2 - b^2 * D = (v - 2)^2 - D * b^2 =
			// = (v - 2)^2 - (v^2 - 4) = -4 * v + 8 == -4 * aminus < 0
			//
			// Therefore we start with sign = false == -1.
			sign = false

			// Now, let's check the sign difference between Jacobi(a, N) and Jacobi(ared, N).
			//
			// When a is odd, ared == a, and therefore there is no sign difference between
			// Jacobi(a, N) and Jacobi(ared, N). Moreover:
			//
			// sign = Jacobi(a, N) / Jacobi(ared, N) = (Jacobi(2, N))^i = ... =
			// = ((-1)^((h*2^(n-1)-1)*h*2^(n-2)))^i
			//
			// --> sign = sign if: a is odd OR i is even OR n > 2
			//     sign = !sign otherwise
			//
			// NOTE: the ((i & 1) == 0) check covers the "a is odd" case too
			// (since in that case i == 0).
			if !(((i & 1) == 0) || (n > 2)) {
				sign = !sign
			}

			if hmoda := h % ared; hmoda == 0 && sign == true {

				// When h == 0 (mod ared), then N == -1 (mod ared), and therefore,
				// keeping in mind that n >= 2, we can write:
				//
				// Jacobi(ared, N) = Jacobi(N, ared) * (-1)^((N-1)*(ared-1)/4) =
				// = Jacobi(-1, ared) * (-1)^((h*2^n-2)/2*(ared-1)/2) =
				// = (-1)^((ared-1)/2) * (-1)^((h*2^(n-1)-1)*(ared-1)/2) =
				//
				// Now, since (h*2^(n-1)-1) is always odd, we have that:
				// 		when (ared-1)/2 is even:
				// 			(-1)^((ared-1)/2) == +1
				// 			(-1)^((h*2^(n-1)-1)*(ared-1)/2) == +1
				// 		when (ared-1)/2 is odd:
				// 			(-1)^((ared-1)/2) == -1
				// 			(-1)^((h*2^(n-1)-1)*(ared-1)/2) == -1
				//
				// In both of the cases, since the signs are concordant, we have that:
				// (-1)^((ared-1)/2) * (-1)^((h*2^(n-1)-1)*(ared-1)/2) = +1
				//
				// Thus, Jacobi(ared, N) = 1.
				//
				// We can immediately verify if the current candidate for V(1) is valid:
				//      when sign == false, then C2 = -1  ->  Check VALID
				//      when sign == true, then C2 = 1  -> Check NOT VALID
				log.Debugf("[C2] Jacobi(%v, N) * sign = 1: %v is NOT a valid candidate for V(1)", ared, v)
				continue

			} else if hmoda != 0 {

				// Check the sign difference between Jacobi(ared, N) and Jacobi(N, ared)
				// Using the expression in the above comment, we know that:
				//
				// Jacobi(ared, N) / Jacobi(N, ared) = (-1)^((N-1)/2*(ared-1)/2) =
				// = (-1)^((h*2^(n-1)-1)*(ared-1)/2)
				//
				// Since (h*2^(n-1)-1) is always odd, the result depends on whether (ared-1)/2 is even or odd.
				//
				// Reminding that ared is odd, we observe that when:
				//		(ared mod 4 == 1) -> (ared-1)/2 is even -> (-1)^((h*2^(n-1)-1)*(ared-1)/2) = +1
				//		(ared mod 4 == 3) -> (ared-1)/2 is odd -> (-1)^((h*2^(n-1)-1)*(ared-1)/2) = -1
				if (ared % 4) == 3 { sign = !sign }

				// Compute Jacobi(N,ared) or get it from the cache if it was already computed.
				var jNa int
				if val, ok := cache[ared]; ok {
					log.Debugf("Retrieved Jacobi(N, %v) from the cache", ared)
					jNa = val

				} else {

					// Jacobi(N, ared) = Jacobi(h * 2^n - 1, ared) =
					// Jacobi(((h mod ared) * (2^n mod ared) - 1) mod ared, ared)
					twonModa, err := modExp(2, n, ared)
					if err != nil {
						return 0, errors.New("Something went wrong: n or ared were negative")
					}

					Nmoda := (hmoda * twonModa - 1 + ared) % ared

					// Check if dred divides N (just in case)
					if (Nmoda == 0) && (ared != 1) {
						return -1, errors.New(fmt.Sprintf("%v divides N: N does not need to be tested " +
							"further.", ared))
					}

					jNa = big.Jacobi(new(big.Int).SetInt64(Nmoda), new(big.Int).SetInt64(ared))

					// Check if GCD(N, ared) != 1
					// If GCD(N, ared) != 1, then N has a divisor > 1, and does not need to be tested further.
					if jNa == 0 {
						return -1, errors.New("N has a known factor, it does not need to be tested further.")
					}

					// Store the computed Jacobi(N, ared) in the cache for potential later use
					cache[ared] = jNa
				}

				// Verify that C2 is verified.
				// If the candidate V(1) does not satisfy it, we move to the next candidate.
				if ((sign == true) && (jNa == 1)) || ((sign == false) && (jNa == -1)) {
					log.Debugf("[C2] Jacobi(%v, N) * sign = 1: %v is not a valid candidate for V(1)", ared, v)
					continue
				}
			}

			// If we reached this point, it means that both C1 and C2 are
			// verified and that the candidate V(1) is valid.
			return v, nil
		}
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
}

// genV1Penne computes a valid V(1) value for the given Riesel candidate.
//
// This function requires:
//		a) n >= 2
//		b) h >= 1
//		c) h mod 2 == 1
//
// This method is the fastest one, but may not find the smallest valid V(1).
// Not finding the smallest valid V(1) may slow down the computation of U(2) later.
//
// Read TODO and [Ref3] for a comprehensive explanation of this method.
func genV1Penne(h, n int64) (int64, error) {

	// Check preconditions
	if h < 1 {
		return -1, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", h))
	}
	if n < 2 {
		return -1, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", n))
	}
	if h % 2 == 0 {
		return -1, errors.New(fmt.Sprintf("Expected odd h, but received h = %v", h))
	}

	// OPTIMIZATION:
	// Store a cache of already computed Jacobi symbols.
	cache := make(map[int64]int)

	for x := int64(1); ; x++ {

		// D is the square free part of (x^2 + 4)
		_, D := reduce(x * x + 4)

		// Set dred = odd part of D
		//
		// Since D is square free, the exact power of 2 dividing it can only be 0 or 1.
		// Therefore, when D is even, we can get its odd part by simply dividing it by 2.
		// When D is odd, we don't need to divide it (since dividing it by 2^0 = 1 is same as not dividing).
		//
		// We compute dred for later computing Jacobi(D, N), because it is more efficient to
		// compute Jacobi(N, D) thanks to the properties of the Jacobi symbol. However, if D
		// is even, we need to reduce it to an odd dred. That is not a problem, because we will
		// later take into account the sign difference between Jacobi(D, N) and Jacobi(N, dred)
		// and compute the former through the latter.
		var dred int64
		if (D & 1) == 1 {
			dred = D
		} else {
			dred = D >> 1
		}

		// Since we use dred instead of D, we have to take into account the sign difference
		// between Jacobi(D, N) and Jacobi(dred, N).
		//
		// When D is odd, dred == D, and therefore there is no sign difference between
		// Jacobi(D, N) and Jacobi(dred, N).
		//
		// When n > 2, sign = +1 as well. In fact, when D is even and n > 2, we have that:
		//
		// sign = Jacobi(D, N) / Jacobi(dred, N) = Jacobi(2, N) = (-1)^((N-1)*(N+1)/8) =
		// (-1)^((h*2^n-2)*(h*2^n)/8) = (-1)^((h*2^(n-1)-1)*h*2^(n-2))
		//
		// --> sign = +1 if D is odd OR n > 2
		//     sign = -1 if D is even and n == 2
		var sign bool
		if (n > 2) || (D& 1 == 1) {
			sign = true
		} else {
			sign = false
		}

		if hmodd := h % dred; hmodd == 0 && sign == true {

			// When h == 0 (mod dred), then N == -1 (mod dred), and therefore,
			// keeping in mind that n >= 2, we can write:
			//
			// Jacobi(dred, N) = Jacobi(N, dred) * (-1)^((N-1)*(dred-1)/4) =
			// = Jacobi(-1, dred) * (-1)^((h*2^n-2)/2*(dred-1)/2) =
			// = (-1)^((dred-1)/2) * (-1)^((h*2^(n-1)-1)*(dred-1)/2) =
			//
			// Now, since (h*2^(n-1)-1) is always odd, we have that:
			// 		when (dred-1)/2 is even:
			// 			(-1)^((dred-1)/2) == +1
			// 			(-1)^((h*2^(n-1)-1)*(dred-1)/2) == +1
			// 		when (dred-1)/2 is odd:
			// 			(-1)^((dred-1)/2) == -1
			// 			(-1)^((h*2^(n-1)-1)*(dred-1)/2) == -1
			//
			// In both of the cases, since the signs are concordant, we have that:
			// (-1)^((dred-1)/2) * (-1)^((h*2^(n-1)-1)*(dred-1)/2) = +1
			//
			// Thus, Jacobi(dred, N) = 1.
			//
			// We can immediately verify if the current candidate for V(1) is valid:
			//      when sign == false, then Jacobi(D, N) = -1  ->  Check VALID
			//      when sign == true, then Jacobi(D, N) = 1  -> Check NOT VALID
			continue

		} else if hmodd != 0 {

			// Count the sign difference between Jacobi(dred, N) and Jacobi(N, dred).
			// Using the expression in the above comment, we know that:
			//
			// Jacobi(dred, N) / Jacobi(N, dred) = (-1)^((N-1)/2*(dred-1)/2) =
			// = (-1)^((h*2^(n-1)-1)*(dred-1)/2)
			//
			// Now, since (h*2^(n-1)-1) is always odd, the sign difference depends
			// on whether (dred-1)/2 is even or odd.
			//
			// Reminding that dred is odd, we observe that when:
			//		(dred mod 4 == 1) -> (dred-1)/2 is even -> (-1)^((h*2^(n-1)-1)*(dred-1)/2) = +1
			//		(dred mod 4 == 3) -> (dred-1)/2 is odd -> (-1)^((h*2^(n-1)-1)*(dred-1)/2) = -1
			if (dred % 4) == 3 { sign = !sign }

			// Compute Jacobi(N,dred) or get it from the cache if it was already computed.
			var jNd int
			if val, ok := cache[dred]; ok {
				jNd = val

			} else {

				// Jacobi(N, dred) = Jacobi(h * 2^n - 1, dred) =
				// Jacobi(((h mod dred) * (2^n mod dred) - 1) mod dred, dred)
				twonModd, err := modExp(2, n, dred)
				if err != nil {
					return 0, errors.New("Something went wrong: n or dred were negative")
				}

				Nmodd := (hmodd * twonModd - 1 + dred) % dred

				// Check if dred divides N (just in case)
				if (Nmodd == 0) && (dred != 1) {
					return -1, errors.New(fmt.Sprintf("%v divides N: N does not need to be tested further.", dred))
				}

				jNd = big.Jacobi(new(big.Int).SetInt64(Nmodd), new(big.Int).SetInt64(dred))

				// Check if GCD(N, dred) != 1
				// If GCD(N, dred) != 1, then N has a divisor > 1, and does not need to be tested further.
				if jNd == 0 {
					return -1, errors.New("N has a known factor, it does not need to be tested further.")
				}

				// Store the computed Jacobi(N, dred) in the cache for potential later use
				cache[dred] = jNd
			}

			// Verify that the resulting Jacobi(D, N) == -1
			//
			// This is the first condition of the Riesel theorem. If the candidate
			// V(1) does not satisfy it, we can move to the next candidate.
			if ((sign == true) && (jNd == 1)) || ((sign == false) && (jNd == -1)) {
				continue
			}
		}

		// In this case, when D is the square-free part of X^2+4, it can be proved that
		// alpha is always of the form epsilon^2, so, Jacobi(D, N) == -1 is sufficient for v to
		// be valid for Riesel's theorem 5.
		return x * x + 2, nil
	}

	// If there is no valid v1, return an error
	return -1, errors.New("It was not possible to find a valid v1 for the given n and h")
}

// GenU2 computes U(2) for the given Riesel candidate, where:
// 		U(2) = V(h)
//
// This function requires:
//		a) n >= 2
//		b) h >= 1
//		c) h mod 2 == 1
//		d) v1 >= 3
//
// We calculate any V(x) as follows:	(Ref1, top of page 873)
//
//		V(0) = alpha^0 + alpha^(-0) = 2
//		V(1) = alpha^1 + alpha^(-1) = GenV1(h,n)
//		V(x+2) = V(1)*V(x+1) - V(x)
//
// It can be shown that the following are true:
//
//		V(2*x) = V(x)^2 - 2
//		V(2*x+1) = V(x+1) * V(x) - V(1)
//
// To prevent V(x) from growing too large, we will replace all V(x) with (V(x) mod N).
func GenU2(R *RieselNumber, v1 int64) (*big.Int, error) {

	// Check preconditions
	if R.h < 1 {
		return nil, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", R.h))
	}
	if R.n < 2 {
		return nil, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", R.n))
	}
	if R.h % 2 == 0 {
		return nil, errors.New(fmt.Sprintf("Expected odd h, but received h = %v", R.h))
	}
	if v1 < 3 {
		return nil, errors.New(fmt.Sprintf("Expected v1 >= 3, but received v1 = %v", v1))
	}

	// Compute V(1) as big.Int for later arithmetic
	v1_big := new(big.Int)
	v1_big.SetInt64(v1)

	// vTwoXPlusOne computes: V(2*x+1) = V(x+1) * V(x) - V(1)
	// and sends the result to the given channel c
	vTwoXPlusOne := func(vXPlus1, vX *big.Int, c chan *big.Int) {
		tmp := new(big.Int)
		tmp.Mul(vX, vXPlus1)
		tmp.Sub(tmp, v1_big)
		rieselMod(tmp, R)
		c <- tmp
	}

	// vTwoX computes: V(2*x) = V(x)^2 - 2
	// and sends the result to the given channel c
	vTwoX := func(vX *big.Int, c chan *big.Int) {
		tmp := new(big.Int)
		tmp.Mul(vX, vX)
		tmp.Sub(tmp, two)
		rieselMod(tmp, R)
		c <- tmp
	}

	// r represents V(x) at every iteration,
	// with x starting from 1. Thus we set it to:
	//
	// r = V(1)
	r := new(big.Int).Set(v1_big)

	// If h == 1, we simply return: V(1) mod N
	if R.h == 1 {
		rieselMod(r, R)
		return r, nil
	}

	// s represents V(x+1) at every iteration,
	// with x starting from 1. Thus we set it to:
	//
	// s = V(1)^2 - 2 = V(2)
	s := new(big.Int)
	s.Mul(r, r)
	s.Sub(s, two)

	// These two channels will be used for the parallel computation
	// of r and s at every iteration.
	c_r := make(chan *big.Int)
	c_s := make(chan *big.Int)

	bitLen, err := bitLen(R.h)
	if err != nil {
		return nil, errors.New("Something went wrong: h was negative")
	}

	// Cycle from second highest bit to second lowest bit of h.
	for i := bitLen - 2; i > 0; i-- {

		// Starting from:
		//		r = V(x)
		//		s = V(x+1)
		if bit(R.h, uint(i)) {

			// If the current bit is a 1, set:
			// 		r = V(2*x+1)
			// 		s = V(2*x+2)
			//
			// These two operations are done in parallel
			go vTwoXPlusOne(s, r, c_r)
			go vTwoX(s, c_s)

			// Receive the resulting r and s from the respective
			// channels when the threads are done
			r = <- c_r
			s = <- c_s

			if loggingEnabled {
				log.Debugf("r = %v", getLastDigits(r))
				log.Debugf("s = %v", getLastDigits(s))
			}

		} else {

			// If the current bit is a 0, set:
			// 		s = V(2*x+1)
			// 		r = V(2*x)
			//
			// These two operations are done in parallel
			go vTwoXPlusOne(s, r, c_s)
			go vTwoX(r, c_r)

			// Receive the resulting r and s from the respective
			// channels when the threads are done
			s = <- c_s
			r = <- c_r

			if loggingEnabled {
				log.Debugf("_r = %v", getLastDigits(r))
				log.Debugf("_s = %v", getLastDigits(s))
			}
		}
	}

	// Since we know that h is odd, the final bit(0) is 1. Thus:
	// 		r = V(2*x+1)
	r.Mul(r, s)
	r.Sub(r, v1_big)
	rieselMod(r, R)

	if loggingEnabled { log.Debugf(".r = %v", getLastDigits(r)) }

	// At this point r = V(h)
	return r, nil
}

// GenUN computes U(n) for the given Riesel candidate.
//
// This function requires:
//		a) n >= 2
//		b) h >= 1
//		c) h mod 2 == 1
//		d) u >= 3
//
// We calculate each term of the Lucas sequence sequentially as follows:	(Ref1, page 871)
//		U(x+1) = U(x)^2 - 2
//
// To prevent U(x) from growing too large, we will replace all U(x) with (U(x) mod N).
func GenUN(R *RieselNumber, u *big.Int) (*big.Int, error) {

	// Check preconditions
	if R.h < 1 {
		return nil, errors.New(fmt.Sprintf("Expected h >= 1, but received h = %v", R.h))
	}
	if R.n < 2 {
		return nil, errors.New(fmt.Sprintf("Expected n >= 2, but received n = %v", R.n))
	}
	if R.h % 2 == 0 {
		return nil, errors.New(fmt.Sprintf("Expected odd h, but received h = %v", R.h))
	}
	if u.Sign() < 0 {
		return nil, errors.New(fmt.Sprintf("Expected u > 0, but received u = %v", u))
	}

	begin := new(syscall.Tms)
	syscall.Times(begin)

	// TODO add checkpoints and correctness checks here
	for i := int64(3); i <= R.n; i++ {

		// u = (u^2 - 2) mod N
		u.Mul(u, u)
		u.Sub(u, two)
		rieselMod(u, R)

		if i % 1000 == 0 {

			current := new(syscall.Tms)
			syscall.Times(current)

			s := u.String()
			log.Infof("Reached U(%v). Last 8 digits = %v. Utime = %.2f. Stime = %.2f.", i, s[len(s)-8:],
				float64(current.Utime - begin.Utime) / 100.0, float64(current.Stime - begin.Stime) / 100.0)
		}
	}

	return u, nil
}