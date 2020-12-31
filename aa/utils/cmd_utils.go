package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	TDPID_FILE = "./.tendermint_pid"
	RDPID_FILE = "./.cybervein_pid"
	DBPID_FILE = "./.redis_pid"
)

func StartcyberveinDaemon() {
	cmd := exec.Command("./cybervein", "start", "-a")
	cmd.Start()
	fmt.Println("cybervein daemon process started")
	SavePID(cmd.Process.Pid, RDPID_FILE)
}

func StartTendermintDaemon() {
	DeleteFile("tendermint.sock")
	cmdStr := `nohup tendermint --home=../chain node --proxy_app=unix://tendermint.sock > ../log/tendermint.log 2>&1 &`
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Start()
	fmt.Println("Tendermint daemon process started")
	SavePID(cmd.Process.Pid, TDPID_FILE)
}

func StartRedisDaemon() {
	cmdStr := `redis-server ../conf/redis.conf`
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Start()
	fmt.Println("Redis daemon process started")
	//SavePID(cmd.Process.Pid, DBPID_FILE)
}


func SavePID(pid int, pidFile string) {

	file, err := os.Create(pidFile)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	file.Sync()
}

func StopPID(pid string) {
	cmd := exec.Command("kill", pid)
	cmd.Run()
	fmt.Println("Stop process ID is : ", pid)
}
