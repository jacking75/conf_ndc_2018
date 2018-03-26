// FrontServer project main.go
package main

import (
	"./backEnd"
	"./clientGateway"
	"./innerMsg"
	"./monitoring"
	def "./typeDefine"
	"./utils"
)

func main() {
	defer utils.PrintPanicStack()

	mainLoop := make(chan struct{}, 1)

	config := mainInit(mainLoop)

	backEnd.ManagerStart(config)

	go startClientTCPServer(config)

	mainEnd(mainLoop)
}

func mainInit(mainloop chan<- struct{}) *def.Config {
	utils.SettingLog()
	utils.Logger.Info("Init. Front-Server")

	config := new(def.Config)
	config.Setting()

	utils.UniqueIdManagerInit(config.MachineId)

	go utils.Sighandler(mainloop)

	schedulStart()

	innerMsg.InnerMsgManagerInit(config.MaxServerInnerMsgChannelSize, config.MaxClientInnerMsgChannelSize)

	return config
}

func mainEnd(mainloop <-chan struct{}) {
	<-mainloop

	utils.Logger.Info("Server Terminate ....")

	clientGateway.TCPServerEnd()
	schedulEnd()

	utils.Logger.Info("Server Terminate !")
}

func startClientTCPServer(config *def.Config) {
	monitoring.ServerStatusCreateChannel()
	go clientGateway.TCPServerStart(config)
}
