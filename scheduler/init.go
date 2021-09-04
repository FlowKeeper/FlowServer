package scheduler

import (
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//workloads stores all hosts managed by this server
var workloads map[primitive.ObjectID]models.Agent

func Init() {
	workloads = make(map[primitive.ObjectID]models.Agent)
	go debugThread()
	logger.Info(loggingArea, "Scheduler was initialized")
}

func debugThread() {
	for {
		logger.Debug("STATS", "Threads:", len(workloads))
		time.Sleep(time.Duration(10) * time.Second)
	}
}
