/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

type configRequestBody struct {
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
	ClientCertificatePassword string   `json:"client_certificate_password"`
	ScanSpeed                 string   `json:"scan_speed"`
	CaseSensitive             string   `json:"case_sensitive"`
	Technologies              []string `json:"technologies"`
	CustomHeaders             []string `json:"custom_headers"`
	CustomCookies             []string `json:"custom_cookies"`
	ExcludedPaths             []string `json:"excluded_paths"`
	UserAgent                 string   `json:"user_agent"`
	Debug                     bool     `json:"debug"`
}

// setConfigCmd represents the setConfig command
var SetConfigCmd = &cobra.Command{
	Use:   "setConfig",
	Short: "Set scan config for target",
	Long: `Takes scan config variables from the config yaml file and the target from stdin. Example
	
	acucli targetGroup --id e3e5afcc-ee2e-431f-a8dc-9d894c93875d | cut -f2 | acucli target setConfig : Sets config for the targets in a target group
	`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input == nil || len(input) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no input provided"), "Error")
			return
		}

		results := make(map[string]interface{})
		for _, id := range input {
			statusCode, responseBody := setConfigRequest(id)
			results[id] = map[string]interface{}{
				"status_code": statusCode,
				"response":    responseBody,
			}
		}

		// Output only the JSON response
		jsonoutput.OutputJSON(results)
	},
}

func setConfigRequest(id string) (int, string) {
	configBody := defineConfig()
	requestJson, err := json.Marshal(configBody)
	if err != nil {
		return 500, fmt.Sprintf("Error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/targets/%s/configuration", viper.GetString("URL"), id), bytes.NewBuffer(requestJson))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return 500, fmt.Sprintf("Error creating request: %v", err)
	}

	// Perform the request using the custom client
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

func getConfigAsSlice(key string) []string {
	configValue := viper.GetString(key)
	if configValue == "" {
		return []string{}
	}
	return strings.Split(configValue, ",")
}

func defineConfig() configRequestBody {
	configBody := configRequestBody{
		Description:       viper.GetString("description"),
		LimitCrawlerScope: viper.GetBool("limit_crawler_scope"),
		Login: struct {
			Kind string `json:"kind"`
		}{Kind: viper.GetString("login.kind")},
		Sensor: viper.GetBool("sensor"),
		SSHCredentials: struct {
			Kind string `json:"kind"`
		}{Kind: viper.GetString("ssh_credentials.kind")},
		Proxy: struct {
			Enabled bool `json:"enabled"`
		}{Enabled: viper.GetBool("proxy.enabled")},
		Authentication: struct {
			Enabled bool `json:"enabled"`
		}{Enabled: viper.GetBool("authentication.enabled")},
		ClientCertificatePassword: viper.GetString("client_certificate_password"),
		ScanSpeed:                 viper.GetString("scan_speed"),
		CaseSensitive:             viper.GetString("case_sensitive"),
		UserAgent:                 viper.GetString("user_agent"),
		Debug:                     viper.GetBool("debug"),
	}

	configBody.Technologies = getConfigAsSlice("technologies")
	configBody.CustomHeaders = getConfigAsSlice("custom_headers")
	configBody.CustomCookies = getConfigAsSlice("custom_cookies")
	configBody.ExcludedPaths = getConfigAsSlice("excluded_paths")

	return configBody

}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
