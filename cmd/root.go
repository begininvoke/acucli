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
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Println("acucli version 1.0.0")
			return nil
		}

		if autoMode {
			if targetURL == "" {
				return fmt.Errorf("target URL is required when using auto mode")
			}
			return auto.RunAutoCommand(targetURL, waitTimeout, outputPath, outputFormat)
		}

		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add commands
	RootCmd.AddCommand(target.TargetCmd)
	RootCmd.AddCommand(targetGroup.TargetGroupCmd)
	RootCmd.AddCommand(scanProfile.ScanProfileCmd)
	RootCmd.AddCommand(scan.ScanCmd)
	RootCmd.AddCommand(report.ReportCmd)
	RootCmd.AddCommand(export.ExportCmd)

	// Global flags
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.acucli.yaml)")

	// Auto mode flags
	RootCmd.Flags().BoolVarP(&autoMode, "auto", "a", false, "Run in auto mode")
	RootCmd.Flags().StringVarP(&targetURL, "u", "u", "", "Target URL to scan")
	RootCmd.Flags().IntVarP(&waitTimeout, "i", "i", 800, "Timeout in seconds for waiting operations")
	RootCmd.Flags().StringVarP(&outputPath, "o", "o", "", "Output path for downloaded report files")
	RootCmd.Flags().StringVarP(&outputFormat, "f", "f", "html", "Output format (csv or html)")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %v", err)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("acucli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	apiKey := viper.GetString("API")
	if apiKey == "" {
		return fmt.Errorf("API key not found in config file")
	}

	httpclient.CreateHttpClient(apiKey)
	return nil
}
