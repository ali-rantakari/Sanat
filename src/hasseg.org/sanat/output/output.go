package output

import (
	"errors"
	"hasseg.org/sanat/model"
)

type OutputFunction func(model.TranslationSet, string)

var OutputFunctionsByName = map[string]OutputFunction{
	"apple":   WriteAppleStringsFiles,
	"android": WriteAndroidStringsFiles,
	"dump":    DumpTranslationSet,
}

func OutputFunctionForName(name string) (OutputFunction, error) {
	ret := OutputFunctionsByName[name]
	if ret != nil {
		return ret, nil
	}

	e := "Unknown output format '" + name + "' — allowed formats: "
	for formatName, _ := range OutputFunctionsByName {
		e += formatName + " "
	}
	return nil, errors.New(e)
}
