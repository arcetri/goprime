package rieseltest

import (
	big "github.com/ricpacca/gmp"
	"fmt"
)

var zero = new(big.Int).SetInt64(0)
var one = new(big.Int).SetInt64(1)
var two = new(big.Int).SetInt64(2)

type RieselNumber struct {
	h int64
	hBig *big.Int
	n int64
	nBig *big.Int
	N *big.Int
}

func NewRieselNumber(h, n int64) *RieselNumber {
	r := new(RieselNumber)
	r.h = h
	r.n = n
	r.hBig = new(big.Int).SetInt64(h)
	r.nBig = new(big.Int).SetInt64(n)

	N := new(big.Int)
	N.Exp(two, r.nBig, nil)
	N.Mul(r.hBig, N)
	N.Sub(N, one)

	r.N = N
	return r
}

func (N *RieselNumber) String() string {
	return fmt.Sprintf("%v * 2^%v - 1", N.h, N.n)
}