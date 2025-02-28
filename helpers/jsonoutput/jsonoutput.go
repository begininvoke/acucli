package jsonoutput

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// OutputJSON takes any data structure and outputs it as formatted JSON without any additional text
func OutputJSON(data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// In case of error, return a JSON error object instead of plain text
		errorResponse := map[string]string{"error": fmt.Sprintf("Error marshaling JSON: %v", err)}
		jsonBytes, _ = json.Marshal(errorResponse)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonBytes, "", "  ")
	if err != nil {
		// In case of error, return a JSON error object instead of plain text
		errorResponse := map[string]string{"error": fmt.Sprintf("Error formatting JSON: %v", err)}
		jsonBytes, _ = json.Marshal(errorResponse)
		fmt.Println(string(jsonBytes))
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}

// OutputRawJSON outputs pre-marshaled JSON bytes without any additional text
func OutputRawJSON(jsonBytes []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, jsonBytes, "", "  ")
	if err != nil {
		// In case of error, return a JSON error object instead of plain text
		errorResponse := map[string]string{"error": fmt.Sprintf("Error formatting JSON: %v", err)}
		jsonBytes, _ = json.Marshal(errorResponse)
		fmt.Println(string(jsonBytes))
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}

// OutputErrorAsJSON converts an error to a JSON error response
func OutputErrorAsJSON(err error, message string) {
	errorResponse := map[string]string{"error": fmt.Sprintf("%s: %v", message, err)}
	jsonBytes, _ := json.Marshal(errorResponse)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, jsonBytes, "", "  ")
	fmt.Println(string(prettyJSON.Bytes()))
}
