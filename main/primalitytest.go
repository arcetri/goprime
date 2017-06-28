package main

import (
	"goprime/rieseltest"
	"fmt"
)

func main() {

	N, err := rieseltest.NewRieselNumber(507, 217588)

	if err != nil {
		fmt.Println(err)

	} else {
		fmt.Println(rieseltest.IsPrime(N))
	}
}
