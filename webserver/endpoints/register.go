package endpoints

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	httphelper "gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpHelper"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/dbtemplate"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(w http.ResponseWriter, r *http.Request) {
	//Unfortunately we cannot use the right types directly as that results in empty uuid's if they arent specified in the request
	var request struct {
		AgentUUID string
		AgentOS   string
		AgentPort string
	}

	if httphelper.HasEmptyBody(w, r) {
		return
	}

	if err := httphelper.CastBodyToStruct(w, r, &request); err != nil {
		return
	}

	//Check if every needed field is present
	if stringHelper.IsEmpty(request.AgentUUID) {
		httpResponse.UserError(w, 400, "AgentUUID missing")
		return
	}

	if stringHelper.IsEmpty(request.AgentOS) {
		httpResponse.UserError(w, 400, "AgentOS missing")
		return
	}

	if stringHelper.IsEmpty(request.AgentPort) {
		httpResponse.UserError(w, 400, "AgentPort missing")
		return
	}

	//Check if the posted values are valid
	agentUUID, err := uuid.Parse(request.AgentUUID)
	if err != nil {
		httpResponse.UserError(w, 400, "AgentUUID is invalid")
		return
	}

	agentOS, err := models.AgentosFromString(request.AgentOS)
	if err != nil {
		httpResponse.UserError(w, 400, "AgentOS is invalid")
		return
	}

	agentPort, err := strconv.Atoi(request.AgentPort)
	if err != nil {
		httpResponse.UserError(w, 400, "AgentPort is invalid")
		return
	}

	//Find out if we now that agent already
	existingAgent, err := dbtemplate.GetAgentByUUID(db.Client(), agentUUID)
	if err == nil {
		httpResponse.Success(w, "OK", "Agent already registered:"+existingAgent.ID.Hex())
		return
	}

	//Error is not nil from this point on
	//Check if we got another error then not finding the specified agentUUID
	if !errors.Is(err, mongo.ErrNoDocuments) {
		httpResponse.InternalError(w, r, err)
		return
	}

	//If we get here it is a noDocuments error -> We need to create the agent

	//Split remote address so we cut out the source port
	agentURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", strings.Split(r.RemoteAddr, ":")[0], agentPort))

	//Register the new agent
	newAgent := models.Agent{
		AgentUUID: agentUUID,
		OS:        agentOS,
		Enabled:   true,
		LastSeen:  time.Now(),
		Endpoint:  agentURL,
	}

	if err := db.AddAgent(&newAgent); err != nil {
		httpResponse.InternalError(w, r, err)
		return
	}

	httpResponse.Success(w, "Added", "Agent successfully regiesterd")
}
