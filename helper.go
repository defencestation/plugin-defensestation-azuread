package main

import (
	"encoding/json"
)

func StructToInterface(inputStruct interface{}) (interface{}, error) {
	// Marshal the struct to JSON
	data, err := json.Marshal(inputStruct)
	if err != nil {
		return nil, err
	}

	// Create a map and unmarshal the JSON data into it
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}