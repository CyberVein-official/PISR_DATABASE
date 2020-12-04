package cmd

import (
	"cybervein.org/CyberveinDB/utils"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop cybervein server",
	Long:  "Description:\n  Stop cybervein server, and includes tendermint server, redis server.",
	Run:   stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, args []string) {
	pid := utils.ReadPIDFile(utils.DBPID_FILE)
	utils.StopPID(pid)
	utils.DeleteFile(utils.DBPID_FILE)

	pid = utils.ReadPIDFile(utils.TDPID_FILE)
	utils.StopPID(pid)
	utils.DeleteFile(utils.TDPID_FILE)

	pid = utils.ReadPIDFile(utils.RDPID_FILE)
	utils.StopPID(pid)
	utils.DeleteFile(utils.RDPID_FILE)
}
