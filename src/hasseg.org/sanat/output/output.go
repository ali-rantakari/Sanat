package output

import (
	"errors"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/output/android"
	"hasseg.org/sanat/output/apple"
	"hasseg.org/sanat/output/dump"
	"hasseg.org/sanat/output/windows"
)

type OutputFunction func(model.TranslationSet, string)

var OutputFunctionsByName = map[string]OutputFunction{
	"apple":   apple.WriteStringsFiles,
	"android": android.WriteStringsFiles,
	"windows": windows.WriteStringsFiles,
	"dump":    dump.DumpTranslationSet,
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
