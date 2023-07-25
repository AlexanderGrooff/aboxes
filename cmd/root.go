/*
Copyright © 2023 Alexander Grooff <alexandergrooff@gmail.com>
*/
package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set log level
		if Verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
	Use:   "boxes",
	Short: "Manage multiple machines",
	Long: `Run commands in a parallel manner on multiple machines.
This is typically useful for ad-hoc commands without too much fuzz.`,
}
var Verbose bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global verbose flag
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
}
