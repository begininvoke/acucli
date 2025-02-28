/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package report

import (
	"github.com/spf13/cobra"
)

// ReportCmd represents the report command
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Commands for managing reports",
	Long:  `Commands for managing reports, including listing, getting, and generating reports.`,
}

func init() {
	// Add subcommands
	ReportCmd.AddCommand(ListCmd)
	ReportCmd.AddCommand(GenerateCmd)
	ReportCmd.AddCommand(RemoveCmd)
	ReportCmd.AddCommand(GetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ReportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ReportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
