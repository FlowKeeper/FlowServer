package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func SetActiveTriggers(AgentID primitive.ObjectID, ActiveTriggers []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := dbclient.Collection("agents").UpdateByID(ctx, AgentID, bson.M{
		"$set": bson.M{
			"activetriggers": ActiveTriggers,
		},
	})

	if err != nil {
		logger.Error(loggingArea, "Couldn't update agents active triggers:", err)
	}

	return err
}
