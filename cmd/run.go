/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the given command/scripts on the target hosts",
	Long:  `Run the given command/scripts on the target hosts and buffer the output to either stdout or a file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve all args
		targets, _ := cmd.Flags().GetStringSlice("target")
		commands, _ := cmd.Flags().GetStringSlice("command")
		outputFile, _ := cmd.Flags().GetString("output")

		// Execute the commands
		executeCommands(targets, commands, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceP("target", "t", []string{}, "Target host")
	runCmd.Flags().StringSliceP("command", "c", []string{}, "Command to run")
	runCmd.Flags().StringP("output", "o", "", "Output file")
}
