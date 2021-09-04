package endpoints

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/scheduler"
	"go.mongodb.org/mongo-driver/mongo"
)

func Config(w http.ResponseWriter, r *http.Request) {
	agentid := r.Header.Get("AgentID")

	if stringHelper.IsEmpty(agentid) {
		httpResponse.UserError(w, 400, "AgentID header missing")
		return
	}

	agentuuid, err := uuid.Parse(agentid)
	if err != nil {
		httpResponse.UserError(w, 400, "AgentID is not a valid uuid")
		return
	}

	agent, err := db.FindAgent(agentuuid)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			httpResponse.UserError(w, 404, "Specified agent isn't recognized")
		} else {
			httpResponse.InternalError(w, r, err)
		}
		return
	}

	httpResponse.SuccessWithPayload(w, "OK", agent)

	//Also ensure that we have a scheduler running to collect results
	go scheduler.StartScheduler(agent)
}
