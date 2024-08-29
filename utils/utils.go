package utils

import (
	"encoding/json"
)

func ConvertToJSON(object any) ([]byte, error) {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
