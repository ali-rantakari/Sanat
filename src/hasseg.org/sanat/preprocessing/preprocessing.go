package preprocessing

import (
	"errors"
	"hasseg.org/sanat/model"
)

type PreProcessorFunction func(*model.TranslationSet) error

var PreProcessorFunctionsByName = map[string]PreProcessorFunction{}

func PreProcessorFunctionForName(name string) (PreProcessorFunction, error) {
	ret := PreProcessorFunctionsByName[name]
	if ret != nil {
		return ret, nil
	}

	e := "Unknown preprocessor '" + name + "' — available preprocessors: "
	for preprocessorName, _ := range PreProcessorFunctionsByName {
		e += preprocessorName + " "
	}
	return nil, errors.New(e)
}
