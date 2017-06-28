package rieseltest

import (
	big "github.com/ricpacca/gmp"
	"fmt"
	"errors"
)

// A RieselNumber represents a number in the form h*2^n-1
//
// QUESTION: Will it ever be the case that h > 2^63 - 1?
type RieselNumber struct {
	h int64
	hBig *big.Int	// multi-precision h
	n int64
	nBig *big.Int	// multi-precision n
	N *big.Int	// h*2^n-1
}

// NewRieselNumber constructs a new RieselNumber instance with the given h and n
//
// This function assumes:
//		a) h >= 1
//		b) n >= 2
//
// When h is even, we will reduce it to odd and add number of times we had
// to divide it by two to n.
func NewRieselNumber(h, n int64) (*RieselNumber, error) {

	// Check preconditions
	if h < 1 {
		return nil, errors.New(fmt.Sprintf("Expected h > 0, but received h = %v", h))
	}
	if n < 2 {
		return nil, errors.New(fmt.Sprintf("Expected n > 1, but received n = %v", n))
	}

	r := new(RieselNumber)
	r.h = h
	r.n = n

	// Make h odd by moving powers of two over 2^n
	if lbit, err := lowerNonZeroBit(r.h); err == nil && lbit > 0 {
		r.n += int64(lbit)
		r.h >>= lbit
	}

	r.hBig = new(big.Int).SetInt64(r.h)
	r.nBig = new(big.Int).SetInt64(r.n)

	N := new(big.Int)
	N.Exp(two, r.nBig, nil)
	N.Mul(r.hBig, N)
	N.Sub(N, one)

	r.N = N
	return r, nil
}

// Custom "toString" functionality to print instances of RieselNumber as h*2^n-1
func (R *RieselNumber) String() string {
	return fmt.Sprintf("%v * 2^%v - 1", R.h, R.n)
}