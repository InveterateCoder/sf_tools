package misc

import (
	"sf_tools/internal/awshelp"
	"strings"
)

func Selector(name string) bool {
	settings := ParseSFMap()
	for key := range settings.States {
		if key == name {
			return true
		}
	}
	return false
}

func MutateState(detail *awshelp.StateDetails) {
	stateSettings := ParseSFMap().States[detail.Name]
	input := make(map[string]any)
	for key, path := range stateSettings.Input {
		input[key] = pickVal(detail.Input, path)
	}
	detail.Input = input
	output := make(map[string]any)
	for key, path := range stateSettings.Output {
		output[key] = pickVal(detail.Output, path)
	}
	detail.Output = output
}

func pickVal(val any, path string) any {
	valMap, ok := val.(map[string]any)
	if !ok {
		return nil
	}
	pathKeys := strings.Split(path, ".")
	if len(pathKeys) == 0 {
		return nil
	}
	result, ok := valMap[pathKeys[0]]
	if !ok {
		return nil
	}
	for _, pathKey := range pathKeys[1:] {
		pathItem, ok := result.(map[string]any)
		if !ok {
			return nil
		}
		pathResult, ok := pathItem[pathKey]
		if !ok {
			return nil
		}
		result = pathResult
	}
	return result
}
