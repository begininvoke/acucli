/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package report

import (
	"bytes"
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

type RemoveReportRequest struct {
	ReportIDList []string `json:"report_id_list"`
}

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove reports",
	Long: `Remove reports by ID. Takes report IDs from stdin. Example:

echo "report_id_here" | acucli report remove
cat report_ids.txt | acucli report remove : Removes multiple reports`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input == nil || len(input) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no report IDs provided"), "Error")
			return
		}

		removeReports(input)
	},
}

func removeReports(reportIDs []string) {
	request := RemoveReportRequest{
		ReportIDList: reportIDs,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/reports/delete"), bytes.NewBuffer(requestJson))
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

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
	// RemoveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// RemoveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
