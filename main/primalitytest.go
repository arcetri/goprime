package main

import (
	"primality/rieseltest"
)

func main() {

	// [ Step 0: use sieve to attempt to eliminate any number
	// that is divisible by a prime less than 2^39 ]

	// Take h and n
	// N = h * 2^n - 1

	// Check Preconditions:
	//    n >= 2
	//    h >= 1
	//    h < 2^n

	// Check if h mod 2 == 1
	// (if not, then move power of two)

	// Generate v(1)

	// From v1 generate u(2) = v(h)
	// this uses the recursion (A)

	// from u(2) generate u(n)
	// this uses the recursion (B)

	// N is prime is u(n) = 0 (mod N)
}