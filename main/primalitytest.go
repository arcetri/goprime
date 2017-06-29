package main

import (
	"goprime/rieseltest"
	"fmt"
	"github.com/op/go-logging"
)

func main() {

	rieseltest.ConfigureLogger(false, logging.INFO, true, logging.INFO)
	N, err := rieseltest.NewRieselNumber(507, 217588)

	// N, err := rieseltest.NewRieselNumber(502573, 7181987)	// largest known Riesel prime

	if err != nil {
		fmt.Println(err)

	} else {
		rieseltest.IsPrime(N)
	}
}
