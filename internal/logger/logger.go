// package logger is a basic internal logger used to debug errors. The logger
// is only active if the -debug flag is provided to a webtex binary.
package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
}

// Log outputs a message to STDOUT. [fmt] and [msg] are used the same way as
// the fmt package.
func Log(fmt string, msg ...any) {
	logger.SetOutput(os.Stdout)
	logger.Printf(fmt, msg...)

	// Probably misguided attempt to amortize switching the output, since I care
	// mainly about error output => there ought to be more error log statements.
	defer logger.SetOutput(os.Stderr)
}

// Error outputs a message to STDERR. [fmt] and [msg] are used in the same way
// as the fmt package.
func Error(fmt string, msg ...any) {
	logger.Printf(fmt, msg...)
}
