/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

type responseBody struct {
	Address                  string `json:"address"`
	Agents                   []any  `json:"agents"`
	ContinuousMode           bool   `json:"continuous_mode"`
	Criticality              int    `json:"criticality"`
	DefaultScanningProfileID string `json:"default_scanning_profile_id"`
	DeletedAt                any    `json:"deleted_at"`
	Description              string `json:"description"`
	Fqdn                     string `json:"fqdn"`
	FqdnHash                 string `json:"fqdn_hash"`
	FqdnStatus               string `json:"fqdn_status"`
	FqdnTmHash               string `json:"fqdn_tm_hash"`
	IssueTrackerID           any    `json:"issue_tracker_id"`
	LastScanDate             string `json:"last_scan_date"`
	LastScanID               string `json:"last_scan_id"`
	LastScanSessionID        string `json:"last_scan_session_id"`
	LastScanSessionStatus    string `json:"last_scan_session_status"`
	ManualIntervention       bool   `json:"manual_intervention"`
	SeverityCounts           struct {
		Critical int `json:"critical"`
		High     int `json:"high"`
		Info     int `json:"info"`
		Low      int `json:"low"`
		Medium   int `json:"medium"`
	} `json:"severity_counts"`
	TargetID     string `json:"target_id"`
	Threat       int    `json:"threat"`
	Type         any    `json:"type"`
	Verification any    `json:"verification"`
}

var id string

// targetCmd represents the target command
var TargetCmd = &cobra.Command{
	Use:   "target",
	Short: "Endpoint for target operations",
	Long:  `Retrieve target information from id flag`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ = cmd.Flags().GetString("id")
		if id == "" {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("target ID is required"), "Error")
			return
		}
		GetTargetRequest(id)
	},
}

func GetTargetRequest(id string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", viper.GetString("URL"), "/targets/", id), nil)
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
	var respBody responseBody
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error parsing JSON")
		return
	}

	// Output only the JSON response
	jsonoutput.OutputRawJSON(body)
}

func init() {
	TargetCmd.Flags().StringVarP(&id, "id", "", "", "Target ID")
	TargetCmd.MarkFlagRequired("id")

	TargetCmd.AddCommand(ListCmd)
	TargetCmd.AddCommand(AddCmd)
	TargetCmd.AddCommand(RemoveCmd)
	TargetCmd.AddCommand(GetConfigCmd)
	TargetCmd.AddCommand(SetConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
