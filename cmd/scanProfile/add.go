/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scanProfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the exported Scan Profile",
	Long: `Imports the exported scan profile. It takes the config json file from stdin (you can export the scan with --export flag). Example:
	cat scanProfile.json | acucli scanProfile add`,
	Run: func(cmd *cobra.Command, args []string) {
		data := filehelper.ReadStdin()
		if data == nil || len(data) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no scan profile data provided"), "Error")
			return
		}

		var stringBuilder strings.Builder
		for _, str := range data {
			stringBuilder.WriteString(str) // Add each string to the builder.
		}
		combinedString := stringBuilder.String()

		// Convert the combined string to a byte slice.
		byteSlice := []byte(combinedString)
		var scanProfile ScanProfile
		err := json.Unmarshal(byteSlice, &scanProfile)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error parsing scan profile JSON")
			return
		}

		makeRequest(scanProfile)
	},
}

func makeRequest(scanProfile ScanProfile) {
	requestJson, err := json.Marshal(scanProfile)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/scanning_profiles"), bytes.NewBuffer(requestJson))
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error making request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error reading response body")
		return
	}

	// Create a response object
	response := map[string]interface{}{
		"status_code": resp.StatusCode,
		"status":      resp.Status,
	}

	// If there's a response body, include it
	if len(body) > 0 {
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err == nil {
			response["response_body"] = jsonBody
		} else {
			response["response_body"] = string(body)
		}
	}

	// Output only the JSON response
	jsonoutput.OutputJSON(response)
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
