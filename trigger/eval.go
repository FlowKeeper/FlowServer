package trigger

import (
	"github.com/PaesslerAG/gval"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/db"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
)

const loggingAreaEVAL = "Eval"
const loggingAreaTrigger = "Trigger"

func EvalutateTriggers(Agent models.Agent) {
	logger.Debug(loggingAreaTrigger, "Evaluating triggers for agent", Agent.ID.Hex())
	if len(Agent.Triggers) == 0 {
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

	for _, k := range Agent.Triggers {
		if !k.Enabled || !k.Trigger.Enabled {
			continue
		}

		//ToDo: Dependencies
		value, err := gval.Evaluate(k.Trigger.Expression, itemFunctions)
		if err != nil {
			logger.Error(loggingAreaEVAL, "Couldn't evaluate state of trigger", k.Trigger.Name, ":", err)
			continue
		}

		logger.Debug(loggingAreaEVAL, "Evaluation of trigger", k.Trigger.Name, "for agent", Agent.ID.Hex(), "returned:", value)

		expressionMatches, expressionIsBoolean := value.(bool)
		if !expressionIsBoolean {
			logger.Error(loggingAreaEVAL, "Expression for trigger", k.Trigger.Name, "does not return true/false! Can't evaluate.")
			continue
		}

		if k.Problematic != expressionMatches {
			if err := db.SetTriggerAssignmentState(Agent.ID, k.TriggerID, expressionMatches); err != nil {
				continue
			}

			if expressionMatches {
				//ToDo: Handle trigger going problematic
				logger.Info(loggingAreaTrigger, "Trigger", k.Trigger.Name, "for agent", Agent.ID.Hex(), "is now active")
			} else {
				//ToDo: Handle trigger going ok
				logger.Info(loggingAreaTrigger, "Trigger", k.Trigger.Name, "for agent", Agent.ID.Hex(), "has recovered")
			}
		}
	}
}
