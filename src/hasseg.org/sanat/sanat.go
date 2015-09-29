package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"hasseg.org/sanat/output"
	"hasseg.org/sanat/parser"
	"hasseg.org/sanat/preprocessing"
	"hasseg.org/sanat/util"
)

func parserErrorHandler(lineNumber int, message string) {
	fmt.Fprintln(os.Stderr, "ERROR on line", lineNumber, message)
}

func main() {
	// Arguments
	//
	usage := `Sanat.

Usage:
  Sanat <input_file> <output_format> <output_dir> [-p value]

Options:
  -p --processors list  The preprocessors to use (comma-separated)
  `
	args, _ := docopt.Parse(usage, nil, true, "Sanat", false)

	inputFilePath := args["<input_file>"].(string)
	outputDirPath := args["<output_dir>"].(string)
	outputFormat := args["<output_format>"].(string)

	// (Optionally) get "group" preprocessor for all the preprocessors
	// we want to run
	//
	preprocessorsArg := args["--processors"]
	var preProcessor preprocessing.PreProcessor
	preProcessor = preprocessing.NewNoOpPreProcessor()
	if preprocessorsArg != nil {
		var err error
		preprocessorNames := util.ComponentsFromCommaSeparatedList(preprocessorsArg.(string))
		preProcessor, err = preprocessing.GroupPreProcessorForProcessorNames(preprocessorNames)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	// Parse translation file
	//
	translationSet, err := parser.TranslationSetFromFile(inputFilePath, preProcessor, parserErrorHandler)
	if err != nil {
		os.Exit(1)
	}

	// Write output
	//
	outputFunction, err := output.OutputFunctionForName(outputFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	outputFunction(translationSet, outputDirPath)
}
