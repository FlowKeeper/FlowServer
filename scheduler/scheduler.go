package scheduler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
)

const loggingArea = "Scheduler"

func StartScheduler(Agent models.Agent) {
	if _, found := workloads[Agent.ID]; found {
		logger.Debug(loggingArea, "Ignored request to start scheduler as its already in our current workload set")
		return
	}

	workloads[Agent.ID] = Agent
	schedulerThread(Agent)
}

func schedulerThread(Agent models.Agent) {
	logger.Info(loggingArea, "Starting scheduler for agent", Agent.AgentID)

	for {
		time.Sleep(time.Second * 60)
		//Check if agent is still in our current workload set
		if _, found := workloads[Agent.ID]; !found {
			logger.Info(loggingArea, "Agent", Agent.AgentID, "is not our workload anymore -> Thread exiting")
			break
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/api/v1/retrieve", Agent.Endpoint.Scheme, Agent.Endpoint.Host), nil)
		if err != nil {
			logger.Error(loggingArea, "Couldn't construct URL for agent", Agent.AgentID, ":", err, "-> Thread will exit")
			break
		}
		req.Header.Add("ScraperUUID", db.InstanceConfig.InstanceID.String())
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			logger.Error(loggingArea, "Couldn't retrieve results for agent", Agent.AgentID, ":", err, ", waiting an additional 5 minutes until next retrieval")
			timeoutRetrieval()
			continue
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			logger.Error(loggingArea, "Couldn't retrieve results for agent", Agent.AgentID, "! Got:", string(bodyBytes), ", waiting an additional 5 minutes until next retrieval")
			timeoutRetrieval()
			continue
		}

		var response struct {
			Status  string
			Payload []models.Result
		}

		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			logger.Error(loggingArea, "Couldn't decode results from agent", Agent.AgentID, ":", err, ", waiting an additional 5 minutes until next retrieval")
			timeoutRetrieval()
			continue
		}

		//Populate the HostID field of the retrieved results
		for i := range response.Payload {
			response.Payload[i].HostID = Agent.ID
		}

		if len(response.Payload) > 0 {
			if err := db.AddResults(response.Payload); err != nil {
				logger.Error(loggingArea, "Database error occured, waiting an additional 5 minutes until next retrieval")
				timeoutRetrieval()
				continue
			}
		}
	}

	delete(workloads, Agent.ID)
}

func timeoutRetrieval() {
	time.Sleep(5 * time.Minute)
}
