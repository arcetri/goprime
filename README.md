# GoPrime

GoPrime is a software that can perform a Lucas-Lehmer-Riesel primality test for numbers of the form 
__h*2<sup>n</sup>-1__.

## Motivation

The main motivations for why GoPrime has been written are:

- Implement, evaluate, compare and document all the algorithms involved in the LLR test.
- Estimate the time each substep of the LLR test and a full LLR test would take for any given number.

GoPrime is open source to serve as a generic learning base for all those interested in understanding how the LLR test 
works. It can be used also for finding new prime numbers, but finding new prime numbers is not the purpose for 
why it was written (there currently is much better software for that purpose).

## Results

### Generating _V(1)_

Generating _V(1)_ is the fastest substep, but at the same time the most difficult to understand part of the LLR test.
We implemented three different algorithms to generate _V(1)_ ([Riesel][riesel], [Rödseth][rodseth] and [Penné][penne]).

In general, we found the Rödseth algorithm to be the most straightforward to implement and we recommend to use it, 
given that it performs well in comparison with the other methods.

Further details may be found in our code comments. 

### Generating _U(2)_

Generating _U(2)_ = _V(h)_ requires to compute approximately log<sub>2</sub>(h) terms of the {_V<sub>i<sub>_} sequence.
Each iteration of this substep works by computing _V(2x+1)_ and _V(2x)_ until we reach _V(h)_.

We found out that, during every iteration, the operations of computing _V(2x+1)_ and _V(2x)_ can be easily 
parallelized and do not need to be done sequentially. Implementing this optimization reduced the total
computation time of this substep of about 50%.

### Generating _U(n)_

Generating _U(n)_ is the most time consuming substep of the LLR test, as it requires to compute _n_ terms of the 
{_U<sub>i<sub>_} sequence where each term depends on the previous term (which makes it hard to parallelize). 

We evaluated the speed of this substep, with a special focus on comparing the time required to square large 
numbers with three different libraries ([Go math/big][big], [FLINT][flint] and [GMP][gmp]). The squaring routine is the 
most crucial part of this substep, because it is where most of the computation time is spent.

It appears that, in squaring large numbers, GoLang math/big is very slow, and that FLINT is slightly faster than GMP.

You may wish to explore other squaring solutions. We expect that approaches based on [Crandall's transform][crandall], 
George Woltman's [Gwnums library][gwnums], [Colin Percival paper][percival] or hardware-specific hand tuned code 
(such as using C with inline assembly to access special hardware instructions) can achieve results at least one 
order of magnitude faster than what we observed so far.

## Usage

### Goprime

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

### Goprime-c

GoPrime-c is a C translation of the GoPrime software. In our tests, it appeared that GoPrime-c is about 7% faster 
than the corresponding GoLang code of GoPrime. 

```sh
# Download and install GoPrime
$ go get github.com/arcetri/GoPrime
$ cd $GOPATH/src/github.com/arcetri/GoPrime/c
$ make install

# Run GoPrime with any h and n
$ goprime-c 391581 216193
```

## Future work
- Evaluate other methods to perform the squaring in the "Generating U(n)" substep.
- Add correctness checks to be regularly performed during the "Generating U(n)" substep.
- Add checkpoints to be regularly saved during the "Generating U(n)" substep.
- Improve the goprime-c code using the GoPrime code and comments as an example.

## Contribute

Please feel invited to contribute by creating a pull request to submit the code or bug fixes you would like to be 
included in GoPrime.

You can also contact us using the following email address: *goprime-contributors (at) external (dot) cisco (dot) com*.
If you send us an email, please include the phrase "__goprime author question__" somewhere in the subject line or 
your email may be rejected.

## License

This project is distributed under the terms of the Apache License v2.0. See file "LICENSE" for further reference.

[rodseth]: <http://folk.uib.no/nmaoy/papers/luc.pdf>
[riesel]: <http://www.ams.org/journals/mcom/1969-23-108/S0025-5718-1969-0262163-1/S0025-5718-1969-0262163-1.pdf>
[penne]: <http://jpenne.free.fr/index2.html>
[flint]: <http://www.flintlib.org/>
[gmp]: <https://gmplib.org>
[big]: <https://golang.org/pkg/math/big/>
[gwnums]: <https://www.mersenne.org/download/>
[crandall]: <http://www.ams.org/journals/mcom/1994-62-205/S0025-5718-1994-1185244-1/S0025-5718-1994-1185244-1.pdf>
[percival]: <http://www.daemonology.net/papers/fft.pdf>