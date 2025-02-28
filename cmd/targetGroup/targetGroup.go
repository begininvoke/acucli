/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package targetGroup

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

type idResponseBody struct {
	TargetIDList []string `json:"target_id_list"`
}

var id string

// targetGroupCmd represents the targetGroup command
var TargetGroupCmd = &cobra.Command{
	Use:   "targetGroup",
	Short: "Get targets from target group id",
	Long: `Takes id of the target group and prints the targets with their id. Example:
	acucli targetGroup --id=cd3db1f4-6275-478c-8830-8d96d37120f3 : Prints the targets of the target group
	`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("target group ID is required"), "Error")
			return
		}
		GetTargetGroupRequest(id)
	},
}

func GetTargetGroupRequest(id string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/target_groups/%s/targets", viper.GetString("URL"), id), nil)
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
	var respBody idResponseBody
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error parsing JSON")
		return
	}

	// Output only the JSON response
	jsonoutput.OutputRawJSON(body)
}

func init() {
	TargetGroupCmd.Flags().StringVarP(&id, "id", "", "", "Target ID")
	TargetGroupCmd.MarkFlagRequired("id")

	TargetGroupCmd.AddCommand(RemoveCmd)
	TargetGroupCmd.AddCommand(AddCmd)
	TargetGroupCmd.AddCommand(ListCmd)
	TargetGroupCmd.AddCommand(AddTargetsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetGroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetGroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
