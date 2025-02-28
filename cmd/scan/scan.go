/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scan

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

type postBody struct {
	TargetID    string   `json:"target_id"`
	ProfileID   string   `json:"profile_id"`
	Schedule    Schedule `json:"schedule"`
	Incremental bool     `json:"incremental"`
}

type Schedule struct {
	Disable       bool    `json:"disable"`
	TimeSensitive bool    `json:"time_sensitive"`
	StartDate     *string `json:"start_date"`
}

var scanProfileId string

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Command to start the scan",
	Long: `Command to start the scan, takes target ids from stdin and id of scan profile in flag, example:

cat target_ids.txt | acucli scan --scanProfileID=47973ea9-018b-4294-9903-bb1cf3b1e886`,
	Run: func(cmd *cobra.Command, args []string) {
		targets := filehelper.ReadStdin()
		if targets == nil || len(targets) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no target IDs provided"), "Error")
			return
		}

		scanProfileID, _ := cmd.Flags().GetString("scanProfileID")
		if scanProfileID == "" {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("scan profile ID is required"), "Error")
			return
		}

		results := make(map[string]interface{})
		for _, target := range targets {
			statusCode, responseBody := startScan(target, scanProfileID)
			results[target] = map[string]interface{}{
				"status_code": statusCode,
				"response":    responseBody,
			}
		}

		// Output only the JSON response
		jsonoutput.OutputJSON(results)
	},
}

func startScan(targetID string, scanProfileID string) (int, string) {
	postBody := postBody{ProfileID: scanProfileID, Incremental: false, Schedule: Schedule{Disable: false, TimeSensitive: false, StartDate: nil}}
	postBody.TargetID = targetID

	requestJson, err := json.Marshal(postBody)
	if err != nil {
		return 500, fmt.Sprintf("Error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/scans"), bytes.NewBuffer(requestJson))
	if err != nil {
		return 500, fmt.Sprintf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return 500, fmt.Sprintf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Sprintf("Error reading response body: %v", err)
	}

	return resp.StatusCode, string(body)
}

func init() {
	ScanCmd.Flags().StringVarP(&scanProfileId, "scanProfileID", "", "", "scanProfile ID")
	ScanCmd.MarkFlagRequired("scanProfileID")

	ScanCmd.AddCommand(ListCmd)
	ScanCmd.AddCommand(GetCmd)
	ScanCmd.AddCommand(RemoveCmd)
	ScanCmd.AddCommand(ResultsCmd)
	ScanCmd.AddCommand(VulnerabilitiesCmd)
	ScanCmd.AddCommand(TechnologiesCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
