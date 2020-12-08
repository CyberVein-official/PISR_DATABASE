package core

import (
	"cybervein.org/CyberveinDB/database"
	"github.com/jasonlvhit/gocron"
)

var Schedulers *gocron.Scheduler

func InitAllJobs() {
	Schedulers = gocron.NewScheduler()
	Schedulers.Every(1).Seconds().Do(CheckRedisStatus)
}

func CheckRedisStatus() {
	isAlive := database.CheckAlive(3)
	if !isAlive {
		LogStoreApp.State.lock.Lock()
		err := AppService.RestoreLocalDatabase()
		if err != nil {
			return
		}
		LogStoreApp.State.UnLock()
	}
}


func StartAllJobs() {
	Schedulers.Start()
}

func StopAllJobs() {
	Schedulers.Clear()
}
