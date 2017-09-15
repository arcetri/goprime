package main

import (
	"github.com/arcetri/goprime/rieseltest"
	"github.com/op/go-logging"
	"fmt"
	"flag"
	"os"
	"strconv"
)

var logLevels = map[int]logging.Level{
	1: logging.WARNING,
	2: logging.INFO,
	3: logging.DEBUG,
}

func main() {

	// Define the Usage message
	flag.Usage = func() {
		fmt.Print("GoPrime, a software to test the primality of numbers of the form h*2^n-1.\n\n")
		fmt.Print("Usage:\n")
		fmt.Print("  goprime [h] [n]\n\n")
		fmt.Print("Optional flags:\n")
		flag.PrintDefaults()
	}

	// Read command line arguments
	//		'-t level' for outputting logs of the specified level and higher to the terminal
	//		'-f level' for outputting logs of the specified level and higher to a file
	terminalLoggerPtr := flag.Int("t", 0, "Level of logs to be written to stdout " +
		"{0 = None (default); 1 = Warning; 2 = Info; 3 = Debug}.")
	fileLoggerPtr := flag.Int("f", 0, "Level of logs to be written to log files " +
		"{0 = None (default); 1 = Warning; 2 = Info; 3 = Debug}.")
	flag.Parse()

	// Check for validity of command line arguments
	if *fileLoggerPtr < 0 || *fileLoggerPtr > 3 || *terminalLoggerPtr < 0 || *terminalLoggerPtr > 3 {
		fmt.Print("Unexpected command line flags.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Configure logger according to the command line arguments
	rieseltest.ConfigureLogger(*fileLoggerPtr != 0, logLevels[*fileLoggerPtr],
		*terminalLoggerPtr != 0, logLevels[*terminalLoggerPtr])

	// Read h and n arguments
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Try to convert them to int64 values
	h, err := strconv.ParseInt(flag.Args()[0], 10, 64)
	if err != nil { panic(err) }
	n, err := strconv.ParseInt(flag.Args()[1], 10, 64)
	if err != nil { panic(err) }

	// Create RieselNumber instance with the specified h and n
	N, err := rieseltest.NewRieselNumber(int64(h), int64(n))

	// N, err := rieseltest.NewRieselNumber(507, 217588)
	// N, err := rieseltest.NewRieselNumber(502573, 7181987)	// largest known Riesel prime

	if err != nil {
		fmt.Println(err)

	} else {

		// Test the specified Riesel number for primality
		result, err := rieseltest.IsPrime(N)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}
