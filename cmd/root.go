package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cybervein",
	Short: "Decentralized K-V database",
	Long:  "Description:\n  Decentralized K-V database based on Redis and Tendermint Blockchain.",
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}


