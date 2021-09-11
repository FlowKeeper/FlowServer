package scheduler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/FlowKeeper/FlowServer/v2/db"
	"github.com/FlowKeeper/FlowServer/v2/trigger"
	"github.com/FlowKeeper/FlowUtils/v2/dbtemplate"
	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
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
	logger.Info(loggingArea, "Starting scheduler for agent", Agent.AgentUUID)

	for {
		var err error

		if Agent.ScrapeInterval <= 0 {
			logger.Debug(loggingArea, "Agent doesn't have a scrape interval set / its invalid. Using standart = 60 seconds")
			Agent.ScrapeInterval = 60
		} else {
			logger.Debug(loggingArea, "Using custom scrape interval for agent", Agent.Name, ":", Agent.ScrapeInterval)
		}
		time.Sleep(time.Second * time.Duration(Agent.ScrapeInterval))
		//Refresh agent
		Agent, err = dbtemplate.GetAgentByUUID(db.Client(), Agent.AgentUUID)
		if err != nil {
			logger.Error(loggingArea, "Couldn't refresh agent configuration:", err)
			continue
		}

		//Check if agent is still in our current workload set
		if _, found := workloads[Agent.ID]; !found {
			logger.Info(loggingArea, "Agent", Agent.AgentUUID, "is not our workload anymore -> Thread exiting")
			break
		}

		//Safety check if we still hold the lock
		if Agent.Scraper.UUID != db.InstanceConfig.InstanceID {
			logger.Warning(loggingArea, "Somehow lost lock for agent", Agent.AgentUUID, " -> Scheduler exiting")
			break
		}

		//Renew lock
		Agent.Scraper.Lock = time.Now()
		if err := db.UpdateLock(Agent); err != nil {
			logger.Error("Couldn't renew lock for agent", Agent.AgentUUID, ":", err)
			continue
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/api/v1/retrieve", Agent.Endpoint.Scheme, Agent.Endpoint.Host), nil)
		if err != nil {
			logger.Error(loggingArea, "Couldn't construct URL for agent", Agent.AgentUUID, ":", err, "-> Thread will exit")
			break
		}
		req.Header.Add("ScraperUUID", db.InstanceConfig.InstanceID.String())
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			logger.Error(loggingArea, "Couldn't retrieve results for agent", Agent.AgentUUID, ":", err)
			continue
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			logger.Error(loggingArea, "Couldn't retrieve results for agent", Agent.AgentUUID, "! Got:", string(bodyBytes))
			continue
		}

		var response struct {
			Status  string
			Payload []models.Result
		}

		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			logger.Error(loggingArea, "Couldn't decode results from agent", Agent.AgentUUID, ":", err)
			continue
		}

		//Populate the HostID field of the retrieved results
		for i := range response.Payload {
			response.Payload[i].HostID = Agent.ID
		}

		if len(response.Payload) > 0 {
			if err := db.AddResults(response.Payload); err != nil {
				logger.Error(loggingArea, "Database error occured")
				continue
			}
		}

		trigger.EvalutateTriggers(Agent)
	}

	delete(workloads, Agent.ID)
}
