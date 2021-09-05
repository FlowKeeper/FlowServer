package config

import (
	"net"
	"os"
	"strconv"
	"strings"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
)

type SampleConfig struct {
	Debug      bool
	Hostname   string
	WebListen  string
	WebPort    int
	MongoDB    string
	AllowedIPs []net.IPNet
}

var Config SampleConfig

const loggingArea = "Config"

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

	allowedIPs := os.Getenv("FLOW_ALLOWED_IPS")
	if stringHelper.IsEmpty(allowedIPs) {
		logger.Fatal(loggingArea, "FLOW_ALLOWED_IPS is empty! Please add the CIDR you agents will be running in. Example: 192.168.1.0/24")
	}

	for _, k := range strings.Split(allowedIPs, ",") {
		_, net, err := net.ParseCIDR(k)
		if err != nil {
			logger.Fatal(loggingArea, "The following allowed IP entry was invalid:", k)
		}

		Config.AllowedIPs = append(Config.AllowedIPs, *net)
	}

	if len(Config.AllowedIPs) == 0 {
		logger.Fatal(loggingArea, "FLOW_ALLOWED_IPS was invalid! Please check your syntax. Example: FLOW_ALLOWED_IPS=\"192.168.1.0/24,192.168.2.0/24\"")
	}

	if Config.Debug {
		logger.EnableDebugLog()
		logger.Debug("CONFIG", "Enabled debug logging")
	}
}
