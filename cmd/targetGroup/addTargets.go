/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package targetGroup

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
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

// addTargetsCmd represents the addTargets command
var AddTargetsCmd = &cobra.Command{
	Use:   "addTargets",
	Short: "Add targets to a target group",
	Long: `Add targets from stdin and the target group via id flag. Example:
	echo targets.txt| acucli targetGroup addTargets --id 0637a8b0-900d-44e8-9a04-edef6ac25e23 : Add targets from file adds the defined target group id
		`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ = cmd.Flags().GetString("id")
		if id == "" {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("target group ID is required"), "Error")
			return
		}

		input := filehelper.ReadStdin()
		if input != nil && len(input) > 0 {
			pBody := postBody{}
			pBody.Add = input
			pBody.Remove = []string{}
			addTargets(pBody, id)
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no target IDs provided"), "Error")
		}
	},
}

func addTargets(pBody postBody, id string) {
	requestJson, err := json.Marshal(pBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
		return
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/target_groups/%s/targets", viper.GetString("URL"), id), bytes.NewBuffer(requestJson))
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
		"status_code":   resp.StatusCode,
		"status":        resp.Status,
		"added_targets": pBody.Add,
		"group_id":      id,
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
	AddTargetsCmd.Flags().StringVarP(&id, "id", "", "", "Group Target ID")
	AddTargetsCmd.MarkFlagRequired("id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTargetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTargetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
