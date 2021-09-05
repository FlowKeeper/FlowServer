package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTrigger(ID primitive.ObjectID) (models.Trigger, error) {
	triggers, err := GetTriggers([]primitive.ObjectID{ID})
	if err != nil {
		return models.Trigger{}, err
	}

	return triggers[0], nil
}

func GetTriggers(IDs []primitive.ObjectID) ([]models.Trigger, error) {
	triggers := make([]models.Trigger, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := dbclient.Collection("triggers").Find(ctx, bson.M{"_id": bson.M{"$in": IDs}})

	if err != nil {
		logger.Error(loggingArea, "Couldn't read items:", err)
		return triggers, err
	}

	if err := result.All(ctx, &triggers); err != nil {
		logger.Error(loggingArea, "Couldn't decode trigger array:", err)
	}

	return triggers, nil
}

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
