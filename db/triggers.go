package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetTriggerAssignmentState(AgentID primitive.ObjectID, TriggerID primitive.ObjectID, Problematic bool) error {
	logger.Debug(loggingArea, "Trying to set problematic to", Problematic, "for TA", TriggerID)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := dbclient.Collection("agents").FindOneAndUpdate(ctx, bson.M{"_id": AgentID}, bson.M{"$set": bson.M{"triggers.$[elem].problematic": Problematic}},
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem.triggerid": TriggerID}},
		}))

	if result.Err() != nil {
		logger.Error(loggingArea, "Couldn't update triggerassignment state:", result.Err())
	}

	return result.Err()
}

func PersistTriggerError(AgentID primitive.ObjectID, TriggerID primitive.ObjectID, Error string) error {
	logger.Debug(loggingArea, "Trying to set error for trigger", TriggerID, "on agent", AgentID)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := dbclient.Collection("agents").FindOneAndUpdate(ctx, bson.M{"_id": AgentID}, bson.M{"$set": bson.M{"triggers.$[elem].error": Error}},
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem.triggerid": TriggerID}},
		}))

	if result.Err() != nil {
		logger.Error(loggingArea, "Couldn't update triggerassignment error:", result.Err())
	}

	return result.Err()
}

func ClearTriggerError(AgentID primitive.ObjectID, TriggerID primitive.ObjectID) error {
	return PersistTriggerError(AgentID, TriggerID, "")
}
