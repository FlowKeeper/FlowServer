package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ensureIndexes() {
	logger.Debug(loggingArea, "Ensuring all needed indexes are present")
	agentIndexes := []mongo.IndexModel{
		{Keys: bson.M{"agentid": 1}},
		{Keys: bson.M{"scraperid": 1}},
	}

	resultIndexes := []mongo.IndexModel{
		{Keys: bson.M{"itemid": 1}},
		{Keys: bson.M{"hostid": 1}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := dbclient.Collection("agents").Indexes().CreateMany(ctx, agentIndexes); err != nil {
		logger.Fatal(loggingArea, "Couldn't ensure indexes for agents collection:", err)
	}

	if _, err := dbclient.Collection("results").Indexes().CreateMany(ctx, resultIndexes); err != nil {
		logger.Fatal(loggingArea, "Couldn't ensure indexes for agents collection:", err)
	}
}
