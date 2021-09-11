package scheduler

import (
	"time"

	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//workloads stores all hosts managed by this server
var workloads map[primitive.ObjectID]models.Agent

//Init initializes all needed variables for the scheduler to work
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
