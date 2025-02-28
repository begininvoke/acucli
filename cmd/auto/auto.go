/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package auto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/tosbaa/acucli/helpers/jsonoutput"
)

// Target structure for adding a target
type Target struct {
	Address     string `json:"address"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Criticality int    `json:"criticality"`
}

// PostBody structure for adding a target
type PostBody struct {
	Targets []Target `json:"targets"`
	Groups  []string `json:"groups"`
}

// ScanSchedule structure for scan scheduling
type ScanSchedule struct {
	Disable       bool    `json:"disable"`
	TimeSensitive bool    `json:"time_sensitive"`
	StartDate     *string `json:"start_date"`
}

// ScanPostBody structure for starting a scan
type ScanPostBody struct {
	TargetID    string       `json:"target_id"`
	ProfileID   string       `json:"profile_id"`
	Schedule    ScanSchedule `json:"schedule"`
	Incremental bool         `json:"incremental"`
}

// ScanResponse structure for scan response
type ScanResponse struct {
	ScanID         string `json:"scan_id"`
	CurrentSession struct {
		Status string `json:"status"`
	} `json:"current_session"`
}

// ReportSource structure for report source
type ReportSource struct {
	Description string   `json:"description"`
	ListType    string   `json:"list_type"`
	IDList      []string `json:"id_list"`
}

// ReportRequest structure for generating a report
type ReportRequest struct {
	TemplateID string       `json:"template_id"`
	Source     ReportSource `json:"source"`
}

// ReportResponse structure for report response
type ReportResponse struct {
	ReportID string   `json:"report_id"`
	Status   string   `json:"status"`
	Download []string `json:"download"`
}

// RemoveReportRequest structure for removing reports
type RemoveReportRequest struct {
	ReportIDList []string `json:"report_id_list"`
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

// ExportResponse structure for export response
type ExportResponse struct {
	ReportID string `json:"report_id"`
}

var targetURL string
var scanProfileID string
var reportTemplateID string
var waitTimeout int
var outputPath string
var outputFormat string

// AutoCmd represents the auto command
var AutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automate the process of scanning and reporting",
	Long: `Automate the process of adding a target, scanning it, generating a report, and downloading the report files.
	
This command performs the following steps:
1. Add a target with the specified URL
2. Check if the target exists
3. Start a scan with the specified scan profile ID
4. Check if the scan exists
5. Wait for the scan to complete
6. Generate a report (HTML format) or create an export (CSV format) based on the format flag
7. Check if the report/export exists
8. Wait for the report/export to complete
9. Download the report/export files
`,
	Run: func(cmd *cobra.Command, args []string) {
		if targetURL == "" {
			jsonoutput.OutputErrorAsJSON(fmt.Errorf("target URL is required"), "Error")
			return
		}

		if scanProfileID == "" {
			// Use default scan profile ID if not provided
			scanProfileID = "11111111-1111-1111-1111-111111111111"
		}

		if reportTemplateID == "" && outputFormat != "csv" {
			// Use default report template ID if not provided and format is not CSV
			reportTemplateID = "11111111-1111-1111-1111-111111111126"
		}

		// Create output directory if specified and doesn't exist
		if outputPath != "" {
			// Extract directory part from the output path
			outputDir := filepath.Dir(outputPath)
			if outputDir != "." {
				err := os.MkdirAll(outputDir, 0755)
				if err != nil {
					jsonoutput.OutputErrorAsJSON(err, "Error creating output directory")
					return
				}
			}
		}

		// Step 1: Add target
		targetID, err := addTarget(targetURL)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error adding target")
			return
		}

		// Log progress
		progressLog := map[string]interface{}{
			"step":      "1. Add target",
			"target_id": targetID,
			"status":    "completed",
		}
		jsonoutput.OutputJSON(progressLog)

		// Step 2: Get target to check if it exists
		targetExists, err := checkTargetExists(targetID)
		if err != nil || !targetExists {
			jsonoutput.OutputErrorAsJSON(err, "Error checking target existence")
			// Clean up
			removeTarget(targetID)
			return
		}

		// Log progress
		progressLog = map[string]interface{}{
			"step":      "2. Check target exists",
			"target_id": targetID,
			"status":    "completed",
		}
		jsonoutput.OutputJSON(progressLog)

		// Step 3: Add scan with scan profile ID
		scanID, err := startScan(targetID, scanProfileID)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error starting scan")
			// Clean up in the specified order
			removeTarget(targetID)
			return
		}

		// Log progress
		progressLog = map[string]interface{}{
			"step":    "3. Start scan",
			"scan_id": scanID,
			"status":  "completed",
		}
		jsonoutput.OutputJSON(progressLog)

		// Step 4: Check scan ID
		scanExists, err := checkScanExists(scanID)
		if err != nil || !scanExists {
			jsonoutput.OutputErrorAsJSON(err, "Error checking scan existence")
			// Clean up in the specified order
			removeScan(scanID)
			removeTarget(targetID)
			return
		}

		// Log progress
		progressLog = map[string]interface{}{
			"step":    "4. Check scan exists",
			"scan_id": scanID,
			"status":  "completed",
		}
		jsonoutput.OutputJSON(progressLog)

		// Step 5: Wait for scan status to be completed
		scanCompleted, err := waitForScanCompletion(scanID, waitTimeout)
		if err != nil || !scanCompleted {
			jsonoutput.OutputErrorAsJSON(err, "Error waiting for scan completion")
			// Clean up in the specified order
			removeScan(scanID)
			removeTarget(targetID)
			return
		}

		// Log progress
		progressLog = map[string]interface{}{
			"step":    "5. Wait for scan completion",
			"scan_id": scanID,
			"status":  "completed",
		}
		jsonoutput.OutputJSON(progressLog)

		// Step 6: Generate report or create export based on format
		var reportID string
		var downloadLinks []string

		if strings.ToLower(outputFormat) == "csv" {
			// Create export for CSV format
			reportID, err = createExport("21111111-1111-1111-1111-111111111141", []string{scanID})
			if err != nil {
				jsonoutput.OutputErrorAsJSON(err, "Error creating export")
				// Clean up in the specified order
				removeScan(scanID)
				removeTarget(targetID)
				return
			}

			// Log progress
			progressLog = map[string]interface{}{
				"step":      "6. Create export",
				"report_id": reportID,
				"status":    "completed",
			}
			jsonoutput.OutputJSON(progressLog)

		} else {
			// Generate HTML report
			reportID, err = generateReport(reportTemplateID, "Auto-generated report", "scan_result", []string{scanID})
			if err != nil {
				jsonoutput.OutputErrorAsJSON(err, "Error generating report")
				// Clean up in the specified order
				removeScan(scanID)
				removeTarget(targetID)
				return
			}

			// Log progress
			progressLog = map[string]interface{}{
				"step":      "6. Generate report",
				"report_id": reportID,
				"status":    "completed",
			}
			jsonoutput.OutputJSON(progressLog)

		}
		// Step 7: Check if report exists
		reportExists, err := checkReportExists(reportID)
		if err != nil || !reportExists {
			jsonoutput.OutputErrorAsJSON(err, "Error checking report existence")
			// Clean up in the specified order
			removeReport(reportID)
			removeScan(scanID)
			removeTarget(targetID)
			return
		}

		// Log progress
		progressLog = map[string]interface{}{
			"step":      "7. Check report exists",
			"report_id": reportID,
			"status":    "completed",
		}
		jsonoutput.OutputJSON(progressLog)
		// Step 8: Wait for report completion and get download links
		downloadLinks, err = waitForReportCompletion(reportID, waitTimeout)
		if err != nil || len(downloadLinks) == 0 {
			jsonoutput.OutputErrorAsJSON(err, "Error waiting for report completion")
			// Clean up in the specified order
			removeReport(reportID)
			removeScan(scanID)
			removeTarget(targetID)
			return
		}
		// Log progress
		progressLog = map[string]interface{}{
			"step":           "8. Wait for completion",
			"report_id":      reportID,
			"download_links": downloadLinks,
			"status":         "completed",
		}
		jsonoutput.OutputJSON(progressLog)

		// Step 9: Download report/export files
		downloadedFiles, err := downloadReportFiles(downloadLinks, outputPath)
		if err != nil {
			jsonoutput.OutputErrorAsJSON(err, "Error downloading files")
			//// Clean up in the specified order
			//if strings.ToLower(outputFormat) == "csv" {
			//	removeExport(reportID)
			//} else {
			removeReport(reportID)
			//}
			removeScan(scanID)
			removeTarget(targetID)
			return
		} else {
			// Log progress
			progressLog = map[string]interface{}{
				"step":             "9. Download files",
				"downloaded_files": downloadedFiles,
				"status":           "completed",
			}
			jsonoutput.OutputJSON(progressLog)
		}

		// Final output with all IDs
		result := map[string]interface{}{
			"target_id": targetID,
			"scan_id":   scanID,
			"report_id": reportID,
			"files":     downloadedFiles,
			"status":    "success",
			"message":   "Auto process completed successfully",
		}
		jsonoutput.OutputJSON(result)

		// Clean up resources after successful download
		cleanupLog := map[string]interface{}{
			"step":    "Cleanup",
			"status":  "in_progress",
			"message": "Removing resources after successful download",
		}
		jsonoutput.OutputJSON(cleanupLog)

		// Remove in the specified order: report first, then scan, then target
		var cleanupErrors []string

		// Remove report/export
		//if strings.ToLower(outputFormat) == "csv" {
		//	if err := removeExport(reportID); err != nil {
		//		cleanupErrors = append(cleanupErrors, fmt.Sprintf("Failed to remove export: %v", err))
		//	} else {
		//		jsonoutput.OutputJSON(map[string]interface{}{
		//			"step":      "Cleanup - Remove export",
		//			"export_id": reportID,
		//			"status":    "completed",
		//		})
		//	}
		//} else {
		if err := removeReport(reportID); err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("Failed to remove report: %v", err))
		} else {
			jsonoutput.OutputJSON(map[string]interface{}{
				"step":      "Cleanup - Remove report",
				"report_id": reportID,
				"status":    "completed",
			})
		}
		//}

		// Remove scan
		if err := removeScan(scanID); err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("Failed to remove scan: %v", err))
		} else {
			jsonoutput.OutputJSON(map[string]interface{}{
				"step":    "Cleanup - Remove scan",
				"scan_id": scanID,
				"status":  "completed",
			})
		}

		// Remove target
		if err := removeTarget(targetID); err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("Failed to remove target: %v", err))
		} else {
			jsonoutput.OutputJSON(map[string]interface{}{
				"step":      "Cleanup - Remove target",
				"target_id": targetID,
				"status":    "completed",
			})
		}

		// Report any cleanup errors
		if len(cleanupErrors) > 0 {
			jsonoutput.OutputJSON(map[string]interface{}{
				"step":    "Cleanup",
				"status":  "warning",
				"message": "Some cleanup operations failed",
				"errors":  cleanupErrors,
			})
		} else {
			jsonoutput.OutputJSON(map[string]interface{}{
				"step":    "Cleanup",
				"status":  "completed",
				"message": "All resources successfully removed",
			})
		}
	},
}

// Add a target and return the target ID
func addTarget(targetURL string) (string, error) {
	targets := []Target{
		{
			Address:     targetURL,
			Description: "Auto-added target",
			Type:        "default",
			Criticality: 30,
		},
	}

	postBody := PostBody{
		Targets: targets,
		Groups:  []string{},
	}

	requestJson, err := json.Marshal(postBody)
	if err != nil {
		return "", fmt.Errorf("error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets/add"), bytes.NewBuffer(requestJson))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response to get the target ID
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Extract target ID from the response
	if targets, ok := response["targets"].([]interface{}); ok && len(targets) > 0 {
		if target, ok := targets[0].(map[string]interface{}); ok {
			if targetID, ok := target["target_id"].(string); ok {
				return targetID, nil
			}
		}
	}

	return "", fmt.Errorf("could not extract target ID from response")
}

// Check if a target exists
func checkTargetExists(targetID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/targets/%s", viper.GetString("URL"), targetID), nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("target does not exist, status code: %d", resp.StatusCode)
	}

	return true, nil
}

// Start a scan and return the scan ID
func startScan(targetID, scanProfileID string) (string, error) {
	postBody := ScanPostBody{
		TargetID:    targetID,
		ProfileID:   scanProfileID,
		Incremental: false,
		Schedule: ScanSchedule{
			Disable:       false,
			TimeSensitive: false,
			StartDate:     nil,
		},
	}

	requestJson, err := json.Marshal(postBody)
	if err != nil {
		return "", fmt.Errorf("error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/scans"), bytes.NewBuffer(requestJson))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response to get the scan ID
	var scanResponse ScanResponse
	err = json.Unmarshal(body, &scanResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return scanResponse.ScanID, nil
}

// Check if a scan exists
func checkScanExists(scanID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/scans/%s", viper.GetString("URL"), scanID), nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("scan does not exist, status code: %d", resp.StatusCode)
	}

	return true, nil
}

// Wait for scan completion
func waitForScanCompletion(scanID string, timeoutSeconds int) (bool, error) {
	startTime := time.Now()
	timeout := time.Duration(timeoutSeconds) * time.Second

	for {
		// Check if timeout has been reached
		if time.Since(startTime) > timeout {
			return false, fmt.Errorf("timeout waiting for scan completion")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/scans/%s", viper.GetString("URL"), scanID), nil)
		if err != nil {
			return false, fmt.Errorf("error creating request: %v", err)
		}

		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			return false, fmt.Errorf("error making request: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return false, fmt.Errorf("error reading response body: %v", err)
		}

		var scanResponse ScanResponse
		err = json.Unmarshal(body, &scanResponse)
		if err != nil {
			return false, fmt.Errorf("error parsing response: %v", err)
		}

		// Check if scan is completed
		if scanResponse.CurrentSession.Status == "completed" {
			return true, nil
		}

		// Wait for 5 seconds before checking again
		time.Sleep(5 * time.Second)
	}
}

// Generate a report and return the report ID
func generateReport(templateID, description, listType string, scanIDs []string) (string, error) {
	reportRequest := ReportRequest{
		TemplateID: templateID,
		Source: ReportSource{
			Description: description,
			ListType:    listType,
			IDList:      scanIDs,
		},
	}

	requestJson, err := json.Marshal(reportRequest)
	if err != nil {
		return "", fmt.Errorf("error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/reports"), bytes.NewBuffer(requestJson))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response to get the report ID
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Extract report ID from the response
	if reportID, ok := response["report_id"].(string); ok {
		return reportID, nil
	}

	return "", fmt.Errorf("could not extract report ID from response")
}

// Check if a report exists
func checkReportExists(reportID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/reports/%s", viper.GetString("URL"), reportID), nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("report does not exist, status code: %d", resp.StatusCode)
	}

	return true, nil
}

// Wait for report completion and get download links
func waitForReportCompletion(reportID string, timeoutSeconds int) ([]string, error) {
	startTime := time.Now()
	timeout := time.Duration(timeoutSeconds) * time.Second

	for {
		// Check if timeout has been reached
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("timeout waiting for report completion")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/reports/%s", viper.GetString("URL"), reportID), nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var reportResponse ReportResponse
		err = json.Unmarshal(body, &reportResponse)
		if err != nil {
			return nil, fmt.Errorf("error parsing response: %v", err)
		}

		// Check if report is completed
		if reportResponse.Status == "completed" {
			// Filter download links to only include .html files
			var htmlLinks []string
			for _, link := range reportResponse.Download {
				if strings.HasSuffix(link, ".csv") && outputFormat == "csv" {

				}
				if strings.HasSuffix(link, ".html") && outputFormat != "csv" {
					htmlLinks = append(htmlLinks, link)
				}
			}
			return htmlLinks, nil
		}

		// Wait for 5 seconds before checking again
		time.Sleep(5 * time.Second)
	}
}

// Download report files
func downloadReportFiles(downloadLinks []string, outputPath string) ([]string, error) {
	var downloadedFiles []string

	// If there are no download links, return empty result
	if len(downloadLinks) == 0 {
		return downloadedFiles, nil
	}

	// Determine if we're using a custom filename or just a directory
	var customFilename string
	var outputDir string

	if outputPath != "" {
		// Check if the output path ends with .html (or another extension)
		// to determine if it's a filename or just a directory
		if strings.HasSuffix(strings.ToLower(outputPath), ".html") ||
			strings.HasSuffix(strings.ToLower(outputPath), ".htm") {
			// It's a filename
			customFilename = filepath.Base(outputPath)
			outputDir = filepath.Dir(outputPath)
			if outputDir == "." {
				outputDir = ""
			}
		} else {
			// It's a directory
			outputDir = outputPath
		}
	}

	for i, link := range downloadLinks {
		// Extract filename from the download link
		defaultFilename := filepath.Base(link)

		// Determine which filename to use
		var filename string
		if customFilename != "" && i == 0 {
			// Use custom filename for the first file only
			filename = customFilename
		} else if customFilename != "" && i > 0 {
			// For additional files when a custom filename is provided,
			// use the custom name with an index suffix
			ext := filepath.Ext(customFilename)
			baseName := strings.TrimSuffix(customFilename, ext)
			filename = fmt.Sprintf("%s_%d%s", baseName, i, ext)
		} else {
			// Use default filename from the link
			filename = defaultFilename
		}

		// Create the request - Fix URL construction to avoid duplicate path segments
		baseURL := viper.GetString("URL")
		// Remove trailing "/api/v1" from baseURL if the link already starts with it
		if filepath.Base(baseURL) == "v1" && strings.HasPrefix(link, "/api/v1") {
			baseURL = strings.TrimSuffix(baseURL, "/api/v1")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", baseURL, link), nil)
		if err != nil {
			return downloadedFiles, fmt.Errorf("error creating request: %v", err)
		}

		// Perform the request
		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			return downloadedFiles, fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		// Determine file path
		var filePath string
		if outputDir != "" {
			filePath = filepath.Join(outputDir, filename)
		} else {
			filePath = filename
		}

		// Create the file
		out, err := os.Create(filePath)
		if err != nil {
			return downloadedFiles, fmt.Errorf("error creating file: %v", err)
		}
		defer out.Close()

		// Write the response body to the file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return downloadedFiles, fmt.Errorf("error writing to file: %v", err)
		}

		downloadedFiles = append(downloadedFiles, filePath)
	}

	return downloadedFiles, nil
}

// Remove a report
func removeReport(reportID string) error {
	request := RemoveReportRequest{
		ReportIDList: []string{reportID},
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/reports/delete"), bytes.NewBuffer(requestJson))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error removing report, status code: %d", resp.StatusCode)
	}

	return nil
}

// Remove a scan
func removeScan(scanID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/scans/%s", viper.GetString("URL"), scanID), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error removing scan, status code: %d", resp.StatusCode)
	}

	return nil
}

// Remove a target
func removeTarget(targetID string) error {
	type RemoveTargetRequest struct {
		TargetIDList []string `json:"target_id_list"`
	}

	request := RemoveTargetRequest{
		TargetIDList: []string{targetID},
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets/delete"), bytes.NewBuffer(requestJson))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error removing target, status code: %d", resp.StatusCode)
	}

	return nil
}

// Create an export and return the export ID
func createExport(exportID string, scanIDs []string) (string, error) {
	exportRequest := ExportRequest{
		ExportID: exportID,
		Source: ExportSource{
			ListType: "scans",
			IDList:   scanIDs,
		},
	}

	requestJson, err := json.Marshal(exportRequest)
	if err != nil {
		return "", fmt.Errorf("error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/exports"), bytes.NewBuffer(requestJson))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response to get the export ID
	var response ExportResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return response.ReportID, nil
}

// Check if an export exists
func checkExportExists(exportID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/exports/%s", viper.GetString("URL"), exportID), nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("export does not exist, status code: %d", resp.StatusCode)
	}

	return true, nil
}

// Wait for export completion and get download links
func waitForExportCompletion(exportID string, timeoutSeconds int) ([]string, error) {
	startTime := time.Now()
	timeout := time.Duration(timeoutSeconds) * time.Second

	for {
		// Check if timeout has been reached
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("timeout waiting for export completion")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/exports/%s", viper.GetString("URL"), exportID), nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var exportResponse map[string]interface{}
		err = json.Unmarshal(body, &exportResponse)
		if err != nil {
			return nil, fmt.Errorf("error parsing response: %v", err)
		}

		// Check if export is completed
		if status, ok := exportResponse["status"].(string); ok && status == "completed" {
			// Get download links
			if download, ok := exportResponse["download"].([]interface{}); ok {
				var links []string
				for _, link := range download {
					if linkStr, ok := link.(string); ok {
						links = append(links, linkStr)
					}
				}
				return links, nil
			}
			return nil, fmt.Errorf("no download links found in response")
		}

		// Wait for 5 seconds before checking again
		time.Sleep(5 * time.Second)
	}
}

// Remove an export
func removeExport(exportID string) error {
	// Create JSON request
	requestBody := map[string]interface{}{
		"export_id_list": []string{exportID},
	}
	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error creating JSON request: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/exports/delete"), bytes.NewBuffer(requestJson))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error removing export, status code: %d", resp.StatusCode)
	}

	return nil
}

func init() {
	// Define flags
	AutoCmd.Flags().StringVarP(&targetURL, "target", "t", "", "Target URL to scan (required)")
	AutoCmd.Flags().StringVarP(&scanProfileID, "scanProfileID", "s", "", "Scan profile ID to use")
	AutoCmd.Flags().StringVarP(&reportTemplateID, "reportTemplateID", "r", "", "Report template ID to use")
	AutoCmd.Flags().IntVarP(&waitTimeout, "timeout", "", 300, "Timeout in seconds for waiting operations")
	AutoCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for downloaded report files (directory or specific filename)")
	AutoCmd.Flags().StringVarP(&outputFormat, "format", "f", "html", "Output format for the report (csv or html)")

	// Mark required flags
	AutoCmd.MarkFlagRequired("target")
}
