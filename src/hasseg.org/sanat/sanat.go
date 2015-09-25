package main

import (
    "fmt"
    "os"
    "github.com/docopt/docopt-go"
    "hasseg.org/sanat/parser"
    "hasseg.org/sanat/output"
    "hasseg.org/sanat/model"
)

func main() {
    usage := "Usage: Sanat <input_file> <output_format> <output_dir>"
    args, _ := docopt.Parse(usage, nil, true, "Sanat", false)

    inputFilePath := args["<input_file>"].(string)
    outputDirPath := args["<output_dir>"].(string)
    outputFormat := args["<output_format>"].(string)
    _ = outputDirPath

    translationSet := parser.NewTranslationSetFromFile(inputFilePath)

    writerMap := make(map[string]func(model.TranslationSet,string))
    writerMap["apple"] = output.WriteAppleStringsFiles
    writerMap["dump"] = output.DumpTranslationSet

    outputFunction := writerMap[outputFormat]
    if outputFunction == nil {
        fmt.Fprint(os.Stderr, "Unknown output format '", outputFormat, "' — allowed formats: ")
        for formatName,_ := range writerMap {
            fmt.Fprint(os.Stderr, formatName + " ")
        }
        fmt.Fprint(os.Stderr, "\n")
        os.Exit(1)
    }
    outputFunction(translationSet, outputDirPath)
}
