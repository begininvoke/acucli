/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package export

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

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Manage exports",
	Long:  `Commands for managing exports in Acunetix.`,
}

// ExportSource represents the source for an export
type ExportSource struct {
	ListType string   `json:"list_type"`
	IDList   []string `json:"id_list"`
}

// ExportRequest represents the request body for creating an export
type ExportRequest struct {
	ExportID string       `json:"export_id,omitempty"`
	Source   ExportSource `json:"source"`
}

// getExportTypesCmd represents the get_export_types command
var getExportTypesCmd = &cobra.Command{
	Use:   "get_export_types",
	Short: "Get available export types",
	Long:  `Get a list of available export types from Acunetix.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create the request
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/export_types", viper.GetString("URL")), nil)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error creating request")
			return
		}

		// Perform the request
		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error making request")
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error reading response body")
			return
		}

		// Parse the JSON response
		var result interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error parsing response")
			return
		}

		// Output the result as JSON
		jsonoutput.OutputJSON(result)
	},
}

// getExportCmd represents the get_export command
var getExportCmd = &cobra.Command{
	Use:   "get_export [export_id]",
	Short: "Get export details",
	Long:  `Get details of a specific export by its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exportID := args[0]

		// Create the request
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/exports/%s", viper.GetString("URL"), exportID), nil)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error creating request")
			return
		}

		// Perform the request
		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error making request")
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error reading response body")
			return
		}

		// Parse the JSON response
		var result interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error parsing response")
			return
		}

		// Output the result as JSON
		jsonoutput.OutputJSON(result)
	},
}

// createExportCmd represents the create command
var createExportCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new export",
	Long: `Create a new export for scans or other resources.
	
Takes input from stdin in the format of scan IDs or other resource IDs, one per line. Example:

echo "scan_id_here" | acucli export create
cat scan_ids.txt | acucli export create

By default, the command uses "scans" as the list type and a predefined export ID.
You can override these defaults with the --list-type and --export-id flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the list type from flags, default to "scans"
		listType, _ := cmd.Flags().GetString("list-type")
		if listType == "" {
			listType = "scans"
		}

		// Get the export ID from flags, default to predefined ID
		exportID, _ := cmd.Flags().GetString("export-id")
		if exportID == "" {
			exportID = "21111111-1111-1111-1111-111111111141"
		}

		// Read IDs from stdin
		idList := filehelper.ReadStdin()
		if idList == nil || len(idList) == 0 {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("no IDs provided"), "Error")
			return
		}

		// Create the export request
		exportRequest := ExportRequest{
			ExportID: exportID,
			Source: ExportSource{
				ListType: listType,
				IDList:   idList,
			},
		}

		// Marshal the request to JSON
		requestJSON, err := json.Marshal(exportRequest)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error creating JSON request")
			return
		}

		// Create the HTTP request
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/exports", viper.GetString("URL")), bytes.NewBuffer(requestJSON))
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error creating request")
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error making request")
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error reading response body")
			return
		}

		// Check if the response is valid JSON
		var responseBody interface{}
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error parsing JSON response")
			return
		}

		// Output only the JSON response
		jsonoutput.OutputRawJSON(body)
	},
}

func init() {
	// Add subcommands to the export command
	ExportCmd.AddCommand(getExportTypesCmd)
	ExportCmd.AddCommand(getExportCmd)
	ExportCmd.AddCommand(createExportCmd)

	// Add flags to the create command
	createExportCmd.Flags().String("list-type", "", "Type of list (e.g., 'scans', defaults to 'scans')")
	createExportCmd.Flags().String("export-id", "", "Optional export ID (defaults to '21111111-1111-1111-1111-111111111141')")
}
