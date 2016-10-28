// Command nconv converts numbers between different bases.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type NumberFormat int

const (
	NUMBER_DECIMAL = iota
	NUMBER_BINARY
	NUMBER_OCTAL
	NUMBER_HEX
)

var (
	convertError       = errors.New("Failed to convert")
	unknownFormatError = errors.New("Unknown number format")
	systemBases        = map[NumberFormat]int{
		NUMBER_DECIMAL: 10,
		NUMBER_BINARY:  2,
		NUMBER_OCTAL:   8,
		NUMBER_HEX:     16,
	}
	systemFormats = map[NumberFormat]string{
		NUMBER_DECIMAL: "%d",
		NUMBER_BINARY:  "%b",
		NUMBER_OCTAL:   "%o",
		NUMBER_HEX:     "%x",
	}
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tnconv [dboh] # operate on standard input\n")
	flag.PrintDefaults()
}

func convertNumber(inputType, outputType NumberFormat, in string) (string, error) {
	inputBase := systemBases[inputType]
	intermediate, err := strconv.ParseInt(in, inputBase, 64)
	if err != nil {
		return "", convertError
	}
	outputFormat := systemFormats[outputType]
	return fmt.Sprintf(outputFormat, intermediate), nil
}

func parseAll(input io.Reader, output io.Writer, inputType, outputType NumberFormat) {
	lineReader := bufio.NewReader(input)

	for {
		line, _, err := lineReader.ReadLine()
		if err != nil {
			return
		}
		trimmed := strings.TrimSpace(string(line))
		conv, err := convertNumber(inputType, outputType, trimmed)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fmt.Fprintln(output, conv)
	}
}

func getTypes(types string) (NumberFormat, NumberFormat, error) {
	var in, out, t NumberFormat
	for _, c := range strings.ToLower(types) {
		switch c {
		case 'h':
			t = NUMBER_HEX
		case 'b':
			t = NUMBER_BINARY
		case 'o':
			t = NUMBER_OCTAL
		case 'd':
			t = NUMBER_DECIMAL
		default:
			return in, out, unknownFormatError
		}
		in = t
		in, out = out, in
	}
	return in, out, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	switch flag.NArg() {
	case 1:
		in, out, err := getTypes(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad types: %s\n", err.Error())
			return
		}
		parseAll(os.Stdin, os.Stdout, in, out)
	default:
		usage()
	}
}
