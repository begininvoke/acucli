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

type PostBody struct {
	Name string `json:"name"`
}

type ResponseBody struct {
	Name    string `json:"name"`
	GroupID string `json:"group_id"`
}

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new target group",
	Long: `Adds target group. Example:
	echo "test2" | acucli targetGroup add : add a target group named "test2"
	cat targetGroups.txt | acucli targetGroup add : Add multiple target groups

`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			results := make(map[string]interface{})

			for _, targetGroupName := range input {
				postBody := PostBody{Name: targetGroupName}
				requestJson, err := json.Marshal(postBody)
				if err != nil {
					results[targetGroupName] = map[string]string{
						"error": fmt.Sprintf("Error creating JSON request: %v", err),
					}
					continue
				}

				req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/target_groups"), bytes.NewBuffer(requestJson))
				if err != nil {
					results[targetGroupName] = map[string]string{
						"error": fmt.Sprintf("Error creating request: %v", err),
					}
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := httpclient.MyHTTPClient.Do(req)
				if err != nil {
					results[targetGroupName] = map[string]string{
						"error": fmt.Sprintf("Error making request: %v", err),
					}
					continue
				}

				responseBody, err := io.ReadAll(resp.Body)
				resp.Body.Close()

				if err != nil {
					results[targetGroupName] = map[string]string{
						"error": fmt.Sprintf("Error reading response body: %v", err),
					}
					continue
				}

				var response ResponseBody
				err = json.Unmarshal(responseBody, &response)
				if err != nil {
					results[targetGroupName] = map[string]string{
						"error":        fmt.Sprintf("Error parsing response: %v", err),
						"raw_response": string(responseBody),
					}
				} else {
					results[targetGroupName] = response
				}
			}

			// Output only the JSON response
			jsonoutput.OutputJSON(results)
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no input provided"), "Error")
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
