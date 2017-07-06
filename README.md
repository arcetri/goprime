# GoPrime

GoPrime is a software that can perform a Lucas-Lehmer-Riesel primality test for numbers
of the form __h*2<sup>n</sup>-1__.

## Motivation

The main motivations for why GoPrime has been written are:

- Implement, evaluate, compare and document all the algorithms involved in the LLR test.
- Estimate the time each substep of the LLR test and a full LLR test would take for any given number.

GoPrime is open source to serve as a learning base for all those interested in understanding how the LLR test works.
It can be used also for finding new prime numbers, but finding new prime numbers is not the main purpose for why it
was written.

## Usage

```sh
# Download and install GoPrime
$ go get github.com/arcetri/GoPrime

# Run GoPrime with any h and n
$ goprime 391581 216193
```

If you have errors with these commands, check that you have GoLang (at least v6) installed and configured with:
    
```sh
# Set the $GOPATH and add the $GOPATH/bin to the PATH environment variable if not already done.
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOPATH/bin
```

## Features to be added in the future
- Correctness checks for the computation of terms of the Lucas sequence U(x)
- Checkpoints for the computation of terms of the Lucas sequence U(x)

## Contribute

Please feel invited to contribute by creating a pull request to submit the code you would like to be included.

You can also contact us using the following email address: *goprime-contributors (at) external (dot) cisco (dot) com*.
If you send us an email, please include the phrase "__goprime author question__" somewhere in the subject line or 
your email may be rejected.


## Contribute

Please feel invited to contribute by creating a pull request to submit the code or bug fixes you would like to be 
included in GoPrime.

## License

This project is distributed under the terms of the Apache License v2.0. See file "LICENSE" for further reference.