package main

import (
	"fmt"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/config"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/scheduler"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/webserver"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/flowutils"
)

func main() {
	logger.Info("MAIN", "Starting up FlowKeeper FlowServer!")

	utilsVersion := flowutils.Version()
	logger.Info("UTILS", fmt.Sprintf("Running FlowUtils Version: %d-%d-%s", utilsVersion.Major, utilsVersion.Minor, utilsVersion.Comment))
	config.ReadConfig()
	if err := db.Connect(); err != nil {
		logger.Fatal("MAIN", "Couldn't connect to MongoDB:", err)
	}
	defer db.Disconnect()
	scheduler.Init()
	webserver.Init()
}
