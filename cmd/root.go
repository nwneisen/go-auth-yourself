package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/server"
	"os"
)

var (
	configFilePath string
)

// persistentPreRun initializes the configuration and other global settings
func persistentPreRun(cmd *cobra.Command, args []string) {
	// Initialize config
	err := config.InitConfig(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-auth-yourself",
	Short: "A simple proxy server",
	Long: `A longer description that explains what this proxy server does.
	
This application serves as a proxy server that handles various authentication flows
including OAuth, SAML, and configuration management.`,
	PreRun: persistentPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		server := server.NewServer()
		server.Start()
	},
}

func init() {
	rootCmd.Flags().StringVar(&configFilePath, "config", "", "Path to the config file")
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
