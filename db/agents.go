package db

import (
	"context"
	"time"

	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAgent(Agent *models.Agent) error {
	if Agent.TemplateIDs == nil {
		Agent.TemplateIDs = make([]primitive.ObjectID, 0)
	}
	if Agent.TriggerMappings == nil {
		Agent.TriggerMappings = make([]models.TriggerAssignment, 0)
	}

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
