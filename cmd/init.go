package cmd

import (
	"os/exec"

	"cybervein.org/CyberveinDB/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization cybervein service",
	Long:  "Description:\n  Initialization cybervein service, init all basic file under chain directory.",
	Run:   initcybervein,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
	├── HomeDir
	│   ├── bin
	│   │   ├── cybervein
	│   ├── conf
	│   │   ├── redis.conf
	│   │   ├── configuration.yaml
	│   ├── chain
	│   │   ├── config
	│   │   │   ├── genesis.json
	│   │   │   ├── config.toml
	│   │   │   ├── ... ...
	│   │   ├── data
	│   │   │   ├── ... ...
*/

func initcybervein(cmd *cobra.Command, args []string) {
	initTendermint()
}


func initTendermint() {
	utils.DeleteFile("../chain")
	utils.DeleteFile("./tendermint.sock")
	cmd := exec.Command("tendermint", "init", "--home=../chain")
	cmd.Run()
}
