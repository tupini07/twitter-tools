package print_utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

var isLoggingToFile = false
var outputFilePath = ""

func SetupLogger(desiredOutputFilePath string) {
	if len(desiredOutputFilePath) == 0 {
		isLoggingToFile = false
	} else {
		isLoggingToFile = true
		outputFilePath = desiredOutputFilePath

		color.NoColor = true

		// create file if necessary (or at least clear it)
		file, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if !errors.Is(err, os.ErrNotExist) {
			file.Truncate(0)
		}
		file.Close()
	}
}

func Printf(format string, a ...any) {
	if isLoggingToFile {
		appendToLogFile(fmt.Sprintf(format, a...))
	} else {
		fmt.Printf(format, a...)
	}
}

func Fprint(w io.Writer, a ...any) {
	if isLoggingToFile {
		appendToLogFile(fmt.Sprint(a...))
	} else {
		fmt.Fprint(w, a...)
	}
}

func Fatalf(format string, v ...any) {
	if isLoggingToFile {
		appendToLogFile(fmt.Sprintf(format, v...))
		os.Exit(1)
	} else {
		log.Fatalf(format, v...)
	}
}

func Fatal(a ...any) {
	if isLoggingToFile {
		appendToLogFile(fmt.Sprint(a...))
		os.Exit(1)
	} else {
		log.Fatal(a...)
	}
}

func appendToLogFile(msg string) {
	f, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(msg); err != nil {
		log.Fatal(err)
	}
}
