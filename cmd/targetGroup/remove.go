/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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

type RemovePostBody struct {
	TargetGroupIDList []string `json:"group_id_list"`
}

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a target group",
	Long: `Takes id input from stdin, removes it without deleting the targets inside. Example:
	
	echo "39468242-2706-43ce-8278-3edf30ed1889" | acucli targetGroup remove : Removes a single target group
	cat toremove.txt | acucli targetGroup remove : Removes multiple
	`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil && len(input) > 0 {
			makeDeleteRequest(input)
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no target group IDs provided"), "Error")
		}
	},
}

func makeDeleteRequest(ids []string) {
	postBody := RemovePostBody{TargetGroupIDList: ids}
	requestJson, err := json.Marshal(postBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/target_groups/delete"), bytes.NewBuffer(requestJson))
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
		"removed_ids": ids,
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
	// RemoveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// RemoveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
