package utils

import (
	"encoding/json"
	"runtime"
	"slices"
	"strings"
)

func ConvertToJSON(object any) ([]byte, error) {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func GetCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcNameParts := strings.Split(fullFuncName, "/")
	return funcNameParts[len(funcNameParts)-1]
}

func SliceFind[T any](s []T, f func(T) bool) (T, bool) {
	var valueToReturn T
	i := slices.IndexFunc(s, f)
	if i == -1 {
		return valueToReturn, false
	}
	valueToReturn = s[i]
	return valueToReturn, true
}
