package cmd

import (
	"fmt"
	"cybervein.org/CyberveinDB/core"
	"cybervein.org/CyberveinDB/database"
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/network"
	"cybervein.org/CyberveinDB/utils"
	"github.com/spf13/cobra"
	abciserver "github.com/tendermint/tendermint/abci/server"
	tlog "github.com/tendermint/tendermint/libs/log"
	"os"
	"os/signal"
	"syscall"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cybervein server",
	Long:  "Description:\n  Start cybervein server, and includes tendermint server, redis server.",
	Run:   start,
}

var daemon bool
var alone bool

func init() {
	startCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "cybervein start mode")
	startCmd.Flags().BoolVarP(&alone, "alone", "a", false, "start cybervein server alone")

	rootCmd.AddCommand(startCmd)
}

func InitService() {
	utils.InitKey()
	utils.InitFiles()
	utils.InitConfig()
	core.InitClient()
	core.InitService()
	logger.InitLogger()
	database.InitRedisClient()
	database.InitBadgerDB()
	core.InitLogStoreApplication()

	core.InitAllJobs()
}

func start(cmd *cobra.Command, args []string) {
	if daemon {
		utils.StartRedisDaemon()
		utils.StartcyberveinDaemon()
		utils.StartTendermintDaemon()
		os.Exit(0)
	}
	if !alone {
		utils.StartRedisDaemon()
	}
	InitService()
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout))
	server := abciserver.NewSocketServer(core.SocketAddr, core.LogStoreApp)
	server.SetLogger(logger)
	if !alone {
		utils.StartTendermintDaemon()
	}

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()

	core.StartAllJobs()
	appServer := network.NewServer()
	appServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
