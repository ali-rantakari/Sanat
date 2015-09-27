package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"hasseg.org/sanat/output"
	"hasseg.org/sanat/parser"
	"os"
)

func parserErrorHandler(lineNumber int, message string) {
	fmt.Fprintln(os.Stderr, "ERROR on line", lineNumber, message)
}

func main() {
	usage := "Usage: Sanat <input_file> <output_format> <output_dir>"
	args, _ := docopt.Parse(usage, nil, true, "Sanat", false)

	inputFilePath := args["<input_file>"].(string)
	outputDirPath := args["<output_dir>"].(string)
	outputFormat := args["<output_format>"].(string)
	_ = outputDirPath

	translationSet, err := parser.NewTranslationSetFromFile(inputFilePath, parserErrorHandler)
	if err != nil {
		os.Exit(1)
	}

	outputFunction, err := output.OutputFunctionForName(outputFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	outputFunction(translationSet, outputDirPath)
}
