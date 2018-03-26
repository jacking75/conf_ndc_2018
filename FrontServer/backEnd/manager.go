package backEnd

import (
	"time"

	"../monitoring"
	def "../typeDefine"
	"../utils"

	"github.com/emirpasic/gods/lists/arraylist"
)

func ManagerStart(config *def.Config) {
	monitoring.ServerStatusCreateChannel()
	go processGoroutine()

	backendList := loadBackEndServerConnectInfo()
	backendCount := backendList.Size()

	for i := 0; i < backendCount; i++ {
		info, _ := backendList.Get(i)

		monitoring.ServerStatusCreateChannel()
		go connectGoroutine(info.(def.BackEndServerConnectInfo).Address, config)
	}

	utils.Logger.Info("wait all backend Server Connect ~")
	for {
		if backendCount != (int)(CurrentConnectCount()) {
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	utils.Logger.Info("CurrentConnectCount :", backendCount)
	utils.Logger.Info("backend Server complete!")
}

func loadBackEndServerConnectInfo() *arraylist.List {
	backEndList := arraylist.New()
	// backend := def.BackEndServerConnectInfo{}
	// backend.Address = "localhost:19999"

	// backEndList.Add(backend)
	// utils.Logger.Info("Load BackendServer Info:", backend.Address)

	return backEndList
}
