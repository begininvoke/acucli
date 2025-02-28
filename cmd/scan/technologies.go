/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scan

import (
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

// technologiesCmd represents the technologies command
var TechnologiesCmd = &cobra.Command{
	Use:   "technologies",
	Short: "Get technologies for a scan result",
	Long: `Get technologies for a specific scan result. Takes scan ID and result ID from stdin in the format "scan_id:result_id". Example:

echo "scan_id:result_id" | acucli scan technologies

You can also pipe the output from the results command and extract the result_id:
echo "scan_id" | acucli scan results | jq -r '.results[0].scan_id + ":" + .results[0].result_id' | acucli scan technologies`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input == nil || len(input) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no scan ID and result ID provided"), "Error")
			return
		}

		// Parse the input to get scan ID and result ID
		parts := strings.Split(input[0], ":")
		if len(parts) != 2 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("input must be in the format 'scan_id:result_id'"), "Error")
			return
		}

		scanID := parts[0]
		resultID := parts[1]
		getScanTechnologies(scanID, resultID)
	},
}

func getScanTechnologies(scanID, resultID string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/scans/%s/results/%s/technologies", viper.GetString("URL"), scanID, resultID), nil)
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
	// technologiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// technologiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
