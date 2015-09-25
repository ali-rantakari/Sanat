package main

import (
    "flag"
    //"fmt"
    //"os"
    "hasseg.org/sanat/parser"
    "hasseg.org/sanat/output"
)

func main() {
    inputPathPtr := flag.String("input_path", "", "Input file")
    flag.Parse()

    translationSet := parser.NewTranslationSetFromFile(*inputPathPtr)
    output.WriteAppleStringsFile(translationSet, "fi")
}
