package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func parseJson(dataFile string) (map[string]interface{}, string) {
	// Open and parse data.json
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, fmt.Sprintf("Failed to open data file: %v", err)
	}
	defer file.Close()

	var data interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Sprintf("Failed to parse data file: %v", err)
	}

	// Validate structure: must have keys 'c' and 'p', each a list of objects with 'id' (number) and 'a' (bool)
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, "Data file is not a JSON object"
	}
	return m, ""
}

func validateJson(dataFile string) (map[string]interface{}, string) {
	m, errMsg := parseJson(dataFile)
	if m == nil {
		return nil, errMsg
	}
	for _, key := range []string{"c", "p"} {
		val, exists := m[key]
		if !exists {
			return nil, fmt.Sprintf("Invalid data file: Missing key '%s'", key)
		}
		arr, ok := val.([]interface{})
		if !ok {
			return nil, fmt.Sprintf("Invalid data file: Key '%s' is not a list", key)
		}
		for i, item := range arr {
			obj, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Sprintf("Invalid data file: Element %d in '%s' is not an object", i, key)
			}
			id, idOk := obj["id"]
			a, aOk := obj["a"]
			if !idOk || !aOk {
				return nil, fmt.Sprintf("Invalid data file: Element %d in '%s' missing 'id' or 'a'", i, key)
			}
			// id must be a number
			switch id.(type) {
			case float64, int, int64:
			default:
				return nil, fmt.Sprintf("Invalid data file: 'id' in element %d of '%s' is not a number", i, key)
			}
			// a must be a bool
			if _, ok := a.(bool); !ok {
				return nil, fmt.Sprintf("Invalid data file: 'a' in element %d of '%s' is not a bool", i, key)
			}
		}
	}
	return m, ""
}