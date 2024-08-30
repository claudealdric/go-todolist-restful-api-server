package utils

import (
	"encoding/json"
	"slices"
)

func ConvertToJSON(object any) ([]byte, error) {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
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
