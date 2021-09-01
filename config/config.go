package config

import (
	"os"
	"strconv"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
)

type SampleConfig struct {
	Hostname  string
	WebListen string
	WebPort   int
	MongoDB   string
}

var Config SampleConfig

func ReadConfig() {
	var err error

	hostname, _ := os.Hostname()
	Config.Hostname = os.Getenv("FLOW.HOSTNAME")
	if stringHelper.IsEmpty(Config.Hostname) {
		logger.Info("CONFIG", "No custom hostname specified. Using:", hostname)
		Config.Hostname = hostname
	}

	Config.WebListen = os.Getenv("FLOW.WEB_LISTEN")
	if stringHelper.IsEmpty(Config.WebListen) {
		logger.Info("CONFIG", "Using default listen address: 0.0.0.0")
		Config.WebListen = "0.0.0.0"
	}

	if Config.WebPort, err = strconv.Atoi(os.Getenv("FLOW.WEB_PORT")); err != nil {
		logger.Info("CONFIG", "Using default listen port: 5000")
		Config.WebPort = 5000
	}

	Config.MongoDB = os.Getenv("FLOW.MONGODB")
	if stringHelper.IsEmpty(Config.MongoDB) {
		logger.Info("CONFIG", "Using default mongodb address: mongodb://localhost:27017")
		Config.MongoDB = "mongodb://localhost:27017"
	}
}
