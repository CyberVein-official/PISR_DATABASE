package cmd

import (
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart cybervein server",
	Long:  "Description:\n  Restart cybervein server, and includes tendermint server, redis server.",
	Run:   restart,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func restart(cmd *cobra.Command, args []string) {
	stop(cmd, args)
	daemon = true
	start(cmd, args)
}
