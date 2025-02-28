/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scanProfile

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

// removeCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the given scanProfile",
	Long: `Removes the given scanProfile. Takes the ids line by line from stdin. Example:

cat scanProfileids.txt | acucli scanProfile remove`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil && len(input) > 0 {
			makeDeleteRequest(input)
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no scan profile IDs provided"), "Error")
		}
	},
}

func makeDeleteRequest(ids []string) {
	results := make(map[string]interface{})

	for _, id := range ids {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", viper.GetString("URL"), "/scanning_profiles/", id), nil)
		if err != nil {
			results[id] = map[string]string{
				"error": fmt.Sprintf("Error creating request: %v", err),
			}
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			results[id] = map[string]string{
				"error": fmt.Sprintf("Error making request: %v", err),
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		result := map[string]interface{}{
			"status_code": resp.StatusCode,
			"status":      resp.Status,
		}

		// If there's a response body, include it
		if len(body) > 0 {
			result["response_body"] = string(body)
		}

		results[id] = result
	}

	// Output only the JSON response
	jsonoutput.OutputJSON(results)
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
