#!/bin/bash

if [ "$1" = "gmp" ]; then
	perl -i -pe 's/\tbig \"math\/big\"/\t\/\/ big \"math\/big\"/s' *.go
	perl -i -pe 's/\t\/\/ big \"github\.com\/arcetri\/gmp\"/\tbig \"github\.com\/arcetri\/gmp\"/s' *.go
	perl -i -pe 's/\tbig \"github\.com\/arcetri\/go.flint\/fmpz\"/\t\/\/ big \"github\.com\/arcetri\/go.flint\/fmpz\"/s' *.go
elif [ "$1" = "flint" ]; then
        perl -i -pe 's/\tbig \"math\/big\"/\t\/\/ big \"math\/big\"/s' *.go
        perl -i -pe 's/\tbig \"github\.com\/arcetri\/gmp\"/\t\/\/ big \"github\.com\/arcetri\/gmp\"/s' *.go
        perl -i -pe 's/\t\/\/ big \"github\.com\/arcetri\/go.flint\/fmpz\"/\tbig \"github\.com\/arcetri\/go.flint\/fmpz\"/s' *.go
elif [ "$1" = "big" ]; then
        perl -i -pe 's/\t\/\/ big \"math\/big\"/\tbig \"math\/big\"/s' *.go
        perl -i -pe 's/\tbig \"github\.com\/arcetri\/gmp\"/\t\/\/ big \"github\.com\/arcetri\/gmp\"/s' *.go
        perl -i -pe 's/\tbig \"github\.com\/arcetri\/go.flint\/fmpz\"/\t\/\/ big \"github\.com\/arcetri\/go.flint\/fmpz\"/s' *.go
else
	echo "Invalid argument"
fi
