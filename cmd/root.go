package cmd

import (
	"fmt"
	"os"
	"nwneisen/go-proxy-yourself/internal/handlers"
	"nwneisen/go-proxy-yourself/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	port    int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-proxy-yourself",
	Short: "A simple proxy server",
	Long: `A longer description that explains what this proxy server does.
	
This application serves as a proxy server that handles various authentication flows
including OAuth, SAML, and configuration management.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize server with configuration
		server := server.NewServer()
		
		// Add handlers
		server.AddHandler("/", handlers.NewIndexHandler)
		server.AddHandler("/config", handlers.NewConfigHandler)
		server.AddHandler("/oauth", handlers.NewOAuthHandler)
		server.AddHandler("/saml", handlers.NewSamlHandler)
		server.AddHandler("/callback", handlers.NewCallbacksHandler)
		
		// Start server with configured port
		server.Start()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-proxy-yourself.yaml)")
	rootCmd.PersistentFlags().IntVar(&port, "port", 8080, "port to run the server on")
	
	// Bind flags to viper
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config file in home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-proxy-yourself")
	}
	
	// Read in config file
	err := viper.ReadInConfig()
	if err != nil {
		// If config file doesn't exist, create default one
		// For now, we'll just continue without config file
		fmt.Println("No config file found, using defaults")
	}
	
	// Set default values
	viper.SetDefault("port", 8080)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}