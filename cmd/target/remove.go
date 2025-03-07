/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

type RemovePostBody struct {
	TargetIDList []string `json:"target_id_list"`
}

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes target",
	Long: `Takes input as stdin. Example:

	echo "9797f3aa-80f7-41a6-9e24-4926b35147cf" | acucli target remove : Removes the target
	acucli target list | echo $(awk '{print $2}') | acucli target remove : Removes all targets`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			makeDeleteRequest(input)
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no input provided"), "Error")
		}
	},
}

func makeDeleteRequest(ids []string) {
	postBody := RemovePostBody{TargetIDList: ids}
	requestJson, err := json.Marshal(postBody)
	if err != nil {
		jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets/delete"), bytes.NewBuffer(requestJson))
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

	// Create a response object
	response := map[string]interface{}{
		"status_code": resp.StatusCode,
		"status":      resp.Status,
		"removed_ids": ids,
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
