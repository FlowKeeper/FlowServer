package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAgent(Agent *models.Agent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := dbclient.Collection("agents").InsertOne(ctx, Agent)
	if err != nil {
		logger.Error(loggingArea, "Couldn't add agent:", err)
		return err
	}

	Agent.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func UpdateLock(Agent models.Agent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := dbclient.Collection("agents").UpdateByID(ctx, Agent.ID, bson.M{
		"$set": bson.M{
			"scraper": Agent.Scraper,
		},
	})

	if err != nil {
		logger.Error(loggingArea, "Couldn't update agent:", err)
		return err
	}

	return nil
}
