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

// removeCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a scan",
	Long: `Remove a scan by ID. Takes scan ID from stdin. Example:

echo "scan_id_here" | acucli scan remove
cat scan_ids.txt | acucli scan remove : Removes multiple scans`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input == nil || len(input) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no scan ID provided"), "Error")
			return
		}

		results := make(map[string]interface{})
		for _, scanID := range input {
			statusCode, responseBody := removeScan(scanID)
			results[scanID] = map[string]interface{}{
				"status_code": statusCode,
				"response":    responseBody,
			}
		}

		// Output only the JSON response
		jsonoutput.OutputJSON(results)
	},
}

func removeScan(scanID string) (int, string) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/scans/%s", viper.GetString("URL"), scanID), nil)
	if err != nil {
		return 500, fmt.Sprintf("Error creating request: %v", err)
	}

	// Perform the request using the custom client
	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return 500, fmt.Sprintf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Sprintf("Error reading response body: %v", err)
	}

	// If there's a response body, try to parse it as JSON
	if len(body) > 0 {
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err == nil {
			jsonBytes, _ := json.Marshal(jsonBody)
			return resp.StatusCode, string(jsonBytes)
		}
	}

	return resp.StatusCode, string(body)
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
