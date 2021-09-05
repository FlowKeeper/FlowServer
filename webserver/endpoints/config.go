package endpoints

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/scheduler"
	"go.mongodb.org/mongo-driver/mongo"
)

func Config(w http.ResponseWriter, r *http.Request) {
	agentid := r.Header.Get("AgentUUID")

	if stringHelper.IsEmpty(agentid) {
		httpResponse.UserError(w, 400, "AgentUUID header missing")
		return
	}

	agentuuid, err := uuid.Parse(agentid)
	if err != nil {
		httpResponse.UserError(w, 400, "AgentUUID is not a valid uuid")
		return
	}

	agent, err := db.GetAgentByUUID(agentuuid)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			httpResponse.UserError(w, 404, "Specified agent isn't recognized")
		} else {
			httpResponse.InternalError(w, r, err)
		}
		return
	}

	//Only check lock lease if it is currently scheduled on another host than us
	//If it is scheduled on the current node, just pass it to the scheduler function as we know all our current workloads
	if agent.Scraper.UUID != db.InstanceConfig.InstanceID {
		//Check if the agent is currenty scheduled on a working node
		if time.Since(agent.Scraper.Lock) > time.Minute*3 {
			logger.Warning("Housekeeper", "A scraper seems to be overloaded or has failed as it hasn't scraped", agent.AgentUUID, "in 3 minutes -> Rescheduling")
			agent.Scraper.UUID = db.InstanceConfig.InstanceID
			agent.Scraper.Lock = time.Now()
			db.UpdateLock(agent)
		} else {
			//If lock is valid, don't start a new thread
			logger.Debug("Scheduler", "Ignored request to start scheduler for agent", agent.AgentUUID, "as its lock is valid")
			return
		}
	}

	httpResponse.SuccessWithPayload(w, "OK", agent)

	go scheduler.StartScheduler(agent)
}
