package trigger

import (
	"github.com/FlowKeeper/FlowServer/v2/db"
	"github.com/FlowKeeper/FlowUtils/v2/models"
	"github.com/PaesslerAG/gval"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
)

const loggingAreaEVAL = "Eval"
const loggingAreaTrigger = "Trigger"

//EvaluateTriggers updates all triggers to represent their respective current state
func EvaluateTriggers(Agent models.Agent) {
	logger.Debug(loggingAreaTrigger, "Evaluating triggers for agent", Agent.ID.Hex())
	if len(Agent.GetAllTriggers()) == 0 {
		logger.Debug(loggingAreaTrigger, "Agent", Agent.ID.Hex(), "has no triggers assigned to it")
		return
	}

	itemFunctions := make(map[string]interface{})
	for _, k := range Agent.GetAllItems() {
		results, err := db.GetResults(Agent.ID, k.ID)
		if err != nil {
			continue
		}

		itemFunctions[k.Name] = results
	}

	for _, k := range Agent.GetAllTriggers() {
		if !k.Enabled {
			continue
		}

		tm, err := Agent.GetTriggerMappingByTriggerID(k.ID)
		if err != nil {
			logger.Error(loggingAreaTrigger, "Trigger seems to be incosistent. Couldn't find TriggerMapping for trigger:", k.ID)
			continue
		}

		//ToDo: Dependencies
		value, err := gval.Evaluate(k.Expression, itemFunctions)
		if err != nil {
			logger.Error(loggingAreaEVAL, "Couldn't evaluate state of trigger", k.Name, ":", err)
			processTriggerError(Agent, tm, err.Error())
			continue
		}

		logger.Debug(loggingAreaEVAL, "Evaluation of trigger", k.Name, "for agent", Agent.ID.Hex(), "returned:", value)

		expressionMatches, expressionIsBoolean := value.(bool)
		if !expressionIsBoolean {
			logger.Error(loggingAreaEVAL, "Expression for trigger", k.Name, "does not return true/false! Can't evaluate.")
			processTriggerError(Agent, tm, "Trigger expression doesn't result in boolean")
			continue
		}

		if tm.Problematic != expressionMatches {
			if err := db.SetTriggerAssignmentState(Agent.ID, k.ID, expressionMatches); err != nil {
				processTriggerError(Agent, tm, err.Error())
				continue
			}

			if expressionMatches {
				//ToDo: Handle trigger going problematic
				logger.Info(loggingAreaTrigger, "Trigger", k.Name, "for agent", Agent.ID.Hex(), "is now active")
			} else {
				//ToDo: Handle trigger going ok
				logger.Info(loggingAreaTrigger, "Trigger", k.Name, "for agent", Agent.ID.Hex(), "has recovered")
			}
		}

		//If we get here, no continues were hit
		//That means the code ran successfully without any errors
		if tm.HasError() {
			if err := db.ClearTriggerError(Agent.ID, k.ID); err != nil {
				logger.Error(loggingAreaTrigger, "Couldn't clear trigger error:", err)
			}
		}
	}
}

func processTriggerError(Agent models.Agent, TriggerAssignment models.TriggerAssignment, Error string) {
	if !TriggerAssignment.HasError() {
		db.PersistTriggerError(Agent.ID, TriggerAssignment.TriggerID, Error)
	}
}
