package main

import (
	"fmt"

	"github.com/FlowKeeper/FlowServer/v2/config"
	"github.com/FlowKeeper/FlowServer/v2/db"
	"github.com/FlowKeeper/FlowServer/v2/scheduler"
	"github.com/FlowKeeper/FlowServer/v2/webserver"
	"github.com/FlowKeeper/FlowUtils/v2/flowutils"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
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
