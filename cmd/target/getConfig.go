/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
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

type ConfigResponseBody struct {
	Description       string `json:"description"`
	LimitCrawlerScope bool   `json:"limit_crawler_scope"`
	Login             struct {
		Kind string `json:"kind"`
	} `json:"login"`
	Sensor         bool `json:"sensor"`
	SSHCredentials struct {
		Kind string `json:"kind"`
	} `json:"ssh_credentials"`
	Proxy struct {
		Enabled bool `json:"enabled"`
	} `json:"proxy"`
	Authentication struct {
		Enabled bool `json:"enabled"`
	} `json:"authentication"`
	ClientCertificatePassword string `json:"client_certificate_password"`
	ScanSpeed                 string `json:"scan_speed"`
	CaseSensitive             string `json:"case_sensitive"`
	Technologies              string `json:"technologies"`
	CustomHeaders             string `json:"custom_headers"`
	CustomCookies             string `json:"custom_cookies"`
	ExcludedPaths             string `json:"excluded_paths"`
	UserAgent                 string `json:"user_agent"`
	Debug                     bool   `json:"debug"`
}

// getConfigCmd represents the getConfig command
var GetConfigCmd = &cobra.Command{
	Use:   "getConfig",
	Short: "Get scan configuration for target",
	Long: `Takes target ID from stdin example:

echo "5fac63fd-088c-4445-a2bf-a9f03f014832" | acucli target getConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			getConfigRequest(input[0])
		} else {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no input provided"), "Error")
		}
	},
}

func getConfigRequest(i string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/targets/%s/configuration", viper.GetString("URL"), i), nil)
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
	var respBody ConfigResponseBody
	err = json.Unmarshal(body, &respBody)
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
	// getConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
