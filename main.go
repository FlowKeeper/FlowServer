package main

import (
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/config"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/scheduler"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/webserver"
)

func main() {
	logger.Info("MAIN", "Starting up FlowKeeper FlowServer!")
	config.ReadConfig()
	if err := db.Connect(); err != nil {
		logger.Fatal("MAIN", "Couldn't connect to MongoDB:", err)
	}
	defer db.Disconnect()
	scheduler.Init()
	webserver.Init()
}
