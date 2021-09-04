package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func FindAgent(AgentID uuid.UUID) (models.Agent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := dbclient.Collection("agents").FindOne(ctx, bson.M{"agentid": AgentID})

	if result.Err() != nil {
		if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
			logger.Error(loggingArea, "Couldn't fetch agent from db:", result.Err())
		}

		return models.Agent{}, result.Err()
	}

	var agent models.Agent
	if err := result.Decode(&agent); err != nil {
		logger.Error(loggingArea, "Couldn't decode agent:", err)
		return models.Agent{}, err
	}

	//Fix if array is nil
	if agent.Items == nil {
		agent.Items = make([]primitive.ObjectID, 0)
	}
	if agent.ItemsResolved == nil {
		agent.ItemsResolved = make([]models.Item, 0)
	}

	if len(agent.Items) > 0 {
		var err error
		agent.ItemsResolved, err = GetItems(agent.Items)
		if err != nil {
			return agent, err
		}
	}
	return agent, nil
}
