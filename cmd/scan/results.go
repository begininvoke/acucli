/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

// resultsCmd represents the results command
var ResultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Get scan result history",
	Long: `Get scan result history for a specific scan ID. Takes scan ID from stdin. Example:

echo "scan_id_here" | acucli scan results`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input == nil || len(input) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no scan ID provided"), "Error")
			return
		}

		scanID := input[0]
		getScanResults(scanID)
	},
}

func getScanResults(scanID string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/scans/%s/results", viper.GetString("URL"), scanID), nil)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating request")
		return
	}

	// Perform the request using the custom client
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

	// Check if the response is valid JSON
	var responseBody interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error parsing JSON")
		return
	}

	// Output only the JSON response
	jsonoutput.OutputRawJSON(body)
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resultsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resultsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
