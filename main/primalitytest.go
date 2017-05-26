package main

import (
	"primality/rieseltest"
)

func main() {

	rieseltest.GenV1riesel(2084259, 1257787)
	//rieseltest.GenV1rodseth(2084259, 1257787)
	rieseltest.GenV1riesel(1095, 2992587)
	//rieseltest.GenV1rodseth(1095, 2992587)
	rieseltest.GenV1riesel(3338480145, 257127)
	//rieseltest.GenV1rodseth(3338480145, 257127)
	rieseltest.GenV1riesel(111546435, 257139)
	//rieseltest.GenV1rodseth(111546435, 257139)
	rieseltest.GenV1riesel(2084259, 1257787)
	//rieseltest.GenV1rodseth(2084259, 1257787)
	rieseltest.GenV1riesel(8331405, 1984565)
	//rieseltest.GenV1rodseth(8331405, 1984565)
	rieseltest.GenV1riesel(165, 2207550)
	//rieseltest.GenV1rodseth(165, 2207550)
	rieseltest.GenV1riesel(1155, 1082878)
	//rieseltest.GenV1rodseth(1155, 1082878)

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