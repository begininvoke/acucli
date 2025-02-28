/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/cmd/auto"
	"github.com/tosbaa/acucli/cmd/export"
	"github.com/tosbaa/acucli/cmd/report"
	"github.com/tosbaa/acucli/cmd/scan"
	"github.com/tosbaa/acucli/cmd/scanProfile"
	"github.com/tosbaa/acucli/cmd/target"
	"github.com/tosbaa/acucli/cmd/targetGroup"
	"github.com/tosbaa/acucli/helpers/httpclient"
)

var (
	cfgFile string
	// Global flags
	targetURL    string
	waitTimeout  int
	outputPath   string
	outputFormat string
	autoMode     bool
	versionFlag  bool
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "acucli",
	Short: "A CLI tool for Acunetix",
	Long: `A command line interface tool for Acunetix that allows you to:
- Manage targets and target groups
- Configure and run scans
- Generate and manage reports
- Automate the entire scanning workflow`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Println("acucli version 1.0.0")
			return nil
		}

		if autoMode {
			// Run auto command with global flags
			return auto.RunAutoCommand(targetURL, waitTimeout, outputPath, outputFormat)
		}

		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(target.TargetCmd)
	RootCmd.AddCommand(targetGroup.TargetGroupCmd)
	RootCmd.AddCommand(scanProfile.ScanProfileCmd)
	RootCmd.AddCommand(scan.ScanCmd)
	RootCmd.AddCommand(report.ReportCmd)
	RootCmd.AddCommand(auto.AutoCmd)
	RootCmd.AddCommand(export.ExportCmd)
	cobra.OnInitialize(initConfig)

	// Global flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.acucli.yaml)")
	RootCmd.PersistentFlags().StringVarP(&targetURL, "u", "u", "", "Target URL to scan")
	RootCmd.PersistentFlags().IntVarP(&waitTimeout, "i", "i", 800, "Timeout in seconds for waiting operations")
	RootCmd.PersistentFlags().StringVarP(&outputPath, "o", "o", "", "Output path for downloaded report files")
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "f", "f", "html", "Output format (csv or html)")

	// Root-only flags
	RootCmd.Flags().BoolVarP(&autoMode, "auto", "a", false, "Run in auto mode")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Get current working directory
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory and current directory
		viper.AddConfigPath(home)
		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("acucli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		httpclient.CreateHttpClient(viper.GetString("API"))
	}
}
