package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"hasseg.org/sanat/output"
	"hasseg.org/sanat/parser"
	"hasseg.org/sanat/preprocessing"
	"hasseg.org/sanat/util"
	"os"
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

	preprocessorsArg := args["--processors"]
	var preprocessorNames []string
	if preprocessorsArg != nil {
		preprocessorNames = util.ComponentsFromCommaSeparatedList(preprocessorsArg.(string))
	}

	// Parse translation file
	//
	translationSet, err := parser.TranslationSetFromFile(inputFilePath, parserErrorHandler)
	if err != nil {
		os.Exit(1)
	}

	// (Optionally) process the translations
	//
	if 0 < len(preprocessorNames) {
		for _, preprocessorName := range preprocessorNames {
			preprocessorFunction, err := preprocessing.PreProcessorFunctionForName(preprocessorName)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			if processErr := preprocessorFunction(&translationSet); processErr != nil {
				fmt.Fprintln(os.Stderr, processErr.Error())
				os.Exit(1)
			}
		}
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
