/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

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

type Target struct {
	Address     string `json:"address"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Criticality int    `json:"criticality"`
}

type PostBody struct {
	Targets []Target `json:"targets"`
	Groups  []string `json:"groups"`
}

var gid string

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add targets",
	Long: `Adding targets
	 It takes urls from stdin, Target Group ID is optional. Example:

	 echo "https://127.0.0.1" | acucli target add : Add the target without target group
	 cat targets.txt | acucli target add --gid=cd3db1f4-6275-478c-8830-8d96d37120f3 : Add targets from a file with target group
	 `,
	Run: func(cmd *cobra.Command, args []string) {
		groups := []string{}
		input := filehelper.ReadStdin()
		inputGID, _ := cmd.Flags().GetString("gid")
		if inputGID != "" {
			groups = append(groups, inputGID)
		}

		if input != nil {
			targets := []Target{}

			for _, line := range input {
				targets = append(targets, Target{Address: line, Description: "", Type: "default", Criticality: 30})
			}
			makeRequest(targets, groups)
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no input provided"), "Error")
		}
	},
}

func makeRequest(t []Target, groups []string) {
	postBody := PostBody{Targets: t, Groups: groups}
	requestJson, err := json.Marshal(postBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets/add"), bytes.NewBuffer(requestJson))
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

	// Read the response body into a variable
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error reading response body")
		return
	}

	// Output only the JSON response
	jsonoutput.OutputRawJSON(responseBody)
}

func init() {
	AddCmd.Flags().StringVarP(&gid, "gid", "g", "", "Group ID (To assign the targets to the group)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
