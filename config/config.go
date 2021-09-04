package config

import (
	"os"
	"strconv"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
)

type SampleConfig struct {
	Debug     bool
	Hostname  string
	WebListen string
	WebPort   int
	MongoDB   string
}

var Config SampleConfig

func ReadConfig() {
	var err error

	hostname, _ := os.Hostname()
	Config.Hostname = os.Getenv("FLOW_HOSTNAME")
	if stringHelper.IsEmpty(Config.Hostname) {
		logger.Info("CONFIG", "No custom hostname specified. Using:", hostname)
		Config.Hostname = hostname
	}

	Config.WebListen = os.Getenv("FLOW_WEB_LISTEN")
	if stringHelper.IsEmpty(Config.WebListen) {
		logger.Info("CONFIG", "Using default listen address: 0.0.0.0")
		Config.WebListen = "0.0.0.0"
	}

	if Config.WebPort, err = strconv.Atoi(os.Getenv("FLOW_WEB_PORT")); err != nil {
		logger.Info("CONFIG", "Using default listen port: 5000")
		Config.WebPort = 5000
	}

	Config.MongoDB = os.Getenv("FLOW_MONGODB")
	if stringHelper.IsEmpty(Config.MongoDB) {
		logger.Info("CONFIG", "Using default mongodb address: mongodb://localhost:27017")
		Config.MongoDB = "mongodb://localhost:27017"
	}

	if Config.Debug, err = strconv.ParseBool(os.Getenv("FLOW_DEBUG")); err != nil {
		Config.Debug = false
	}

	if Config.Debug {
		logger.EnableDebugLog()
		logger.Debug("CONFIG", "Enabled debug logging")
	}
}
