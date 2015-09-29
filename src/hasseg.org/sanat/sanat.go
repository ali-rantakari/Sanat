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
  Sanat generate <input_file> <output_format> <output_dir> [-p value]
  Sanat validate <input_file>

Options:
  -p --processors list  The preprocessors to use (comma-separated)
  `
	args, _ := docopt.Parse(usage, nil, true, "Sanat", false)

	// (Optionally) get "group" preprocessor for all the preprocessors
	// we want to run
	//
	preprocessorsArg := args["--processors"]
	var preprocessor preprocessing.Preprocessor
	preprocessor = preprocessing.NewNoOpPreprocessor()
	if preprocessorsArg != nil {
		var err error
		preprocessorNames := util.ComponentsFromCommaSeparatedList(preprocessorsArg.(string))
		preprocessor, err = preprocessing.GroupPreprocessorForProcessorNames(preprocessorNames)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	inputFilePath := args["<input_file>"].(string)

	// Parse translation file
	//
	translationSet, err := parser.TranslationSetFromFile(inputFilePath, preprocessor, parserErrorHandler)
	if err != nil {
		os.Exit(1)
	}

	if args["generate"].(bool) {
		outputDirPath := args["<output_dir>"].(string)
		outputFormat := args["<output_format>"].(string)

		// Write output
		//
		outputFunction, err := output.OutputFunctionForName(outputFormat)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		outputFunction(translationSet, outputDirPath)
	}
}
