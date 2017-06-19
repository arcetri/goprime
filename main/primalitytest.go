package main

import (
	"goprime/rieseltest"
	"fmt"
)

func main() {

	N := rieseltest.NewRieselNumber(507, 217588)
	fmt.Println(rieseltest.IsPrime(N))

	/*rieseltest.genV1Rodseth(2084259, 1257787)
	rieseltest.genV1Riesel(1095, 2992587)
	//rieseltest.genV1Rodseth(1095, 2992587)
	rieseltest.genV1Riesel(3338480145, 257127)
	//rieseltest.genV1Rodseth(3338480145, 257127)
	rieseltest.genV1Riesel(111546435, 257139)
	//rieseltest.genV1Rodseth(111546435, 257139)
	rieseltest.genV1Riesel(2084259, 1257787)
	//rieseltest.genV1Rodseth(2084259, 1257787)
	rieseltest.genV1Riesel(8331405, 1984565)
	//rieseltest.genV1Rodseth(8331405, 1984565)
	rieseltest.genV1Riesel(165, 2207550)
	//rieseltest.genV1Rodseth(165, 2207550)
	rieseltest.genV1Riesel(1155, 1082878)
	//rieseltest.genV1Rodseth(1155, 1082878) */

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
