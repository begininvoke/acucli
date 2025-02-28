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

// getCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific scan",
	Long: `Get details of a specific scan by ID. Takes scan ID from stdin. Example:

echo "scan_id_here" | acucli scan get`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input == nil || len(input) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no scan ID provided"), "Error")
			return
		}

		scanID := input[0]
		getScanDetails(scanID)
	},
}

func getScanDetails(scanID string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/scans/%s", viper.GetString("URL"), scanID), nil)
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
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
