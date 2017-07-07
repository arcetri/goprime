#!/bin/bash

if [ "$1" = "gmp" ]; then
	perl -i -pe 's/\tbig \"math\/big\"/\t\/\/ big \"math\/big\"/s' *.go
	perl -i -pe 's/\t\/\/ big \"github\.com\/ricpacca\/gmp\"/\tbig \"github\.com\/ricpacca\/gmp\"/s' *.go
	perl -i -pe 's/\tbig \"github\.com\/ricpacca\/go.flint\/fmpz\"/\t\/\/ big \"github\.com\/ricpacca\/go.flint\/fmpz\"/s' *.go
elif [ "$1" = "flint" ]; then
        perl -i -pe 's/\tbig \"math\/big\"/\t\/\/ big \"math\/big\"/s' *.go
        perl -i -pe 's/\tbig \"github\.com\/ricpacca\/gmp\"/\t\/\/ big \"github\.com\/ricpacca\/gmp\"/s' *.go
        perl -i -pe 's/\t\/\/ big \"github\.com\/ricpacca\/go.flint\/fmpz\"/\tbig \"github\.com\/ricpacca\/go.flint\/fmpz\"/s' *.go
elif [ "$1" = "big" ]; then
        perl -i -pe 's/\t\/\/ big \"math\/big\"/\tbig \"math\/big\"/s' *.go
        perl -i -pe 's/\tbig \"github\.com\/ricpacca\/gmp\"/\t\/\/ big \"github\.com\/ricpacca\/gmp\"/s' *.go
        perl -i -pe 's/\tbig \"github\.com\/ricpacca\/go.flint\/fmpz\"/\t\/\/ big \"github\.com\/ricpacca\/go.flint\/fmpz\"/s' *.go
else
	echo "Invalid argument"
fi
