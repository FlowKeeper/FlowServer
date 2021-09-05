package trigger

import (
	"github.com/PaesslerAG/gval"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const loggingAreaEVAL = "Eval"
const loggingAreaTrigger = "Trigger"

func EvalutateTriggers(Agent models.Agent) {
	logger.Debug(loggingAreaTrigger, "Evaluating triggers for agent", Agent.ID.Hex())
	if len(Agent.TriggersResolved) == 0 {
		logger.Debug(loggingAreaTrigger, "Agent", Agent.ID.Hex(), "has no triggers assigned to it")
		return
	}

	itemFunctions := make(map[string]interface{})
	for _, k := range Agent.ItemsResolved {
		results, err := db.GetResults(Agent.ID, k.ID)
		if err != nil {
			continue
		}

		itemFunctions[k.Name] = results
	}

	activeTriggers := make([]primitive.ObjectID, 0)

	for _, k := range Agent.TriggersResolved {
		if !k.Enabled {
			continue
		}

		//ToDo: Dependencies
		value, err := gval.Evaluate(k.Expression, itemFunctions)
		if err != nil {
			logger.Error(loggingAreaEVAL, "Couldn't evaluate state of trigger", k.Name, ":", err)
			continue
		}

		logger.Debug(loggingAreaEVAL, "Evaluation of trigger", k.Name, "for agent", Agent.ID.Hex(), "returned:", value)

		expressionMatches, expressionIsBoolean := value.(bool)
		if !expressionIsBoolean {
			logger.Error(loggingAreaEVAL, "Expression for trigger", k.Name, "does not return true/false! Can't evaluate.")
			continue
		}

		if expressionMatches {
			activeTriggers = append(activeTriggers, k.ID)
		}
	}

	//Determine which triggers have switched to problematic
	for _, k := range activeTriggers {
		if !containsObjectId(Agent.ActiveTriggers, k) {
			//ToDo: Handle trigger going problematic
			logger.Info(loggingAreaTrigger, "Trigger", k.Hex(), "for agent", Agent.ID.Hex(), "is now active")
		}
	}

	//Determine which triggers have switched to ok
	for _, k := range Agent.ActiveTriggers {
		if !containsObjectId(activeTriggers, k) {
			//ToDo: Handle trigger going ok
			logger.Info(loggingAreaTrigger, "Trigger", k.Hex(), "for agent", Agent.ID.Hex(), "has recovered")
		}
	}

	db.SetActiveTriggers(Agent.ID, activeTriggers)
}

func containsObjectId(Slice []primitive.ObjectID, Id primitive.ObjectID) bool {
	for _, k := range Slice {
		if k == Id {
			return true
		}
	}

	return false
}
