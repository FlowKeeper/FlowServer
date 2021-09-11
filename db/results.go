package db

import (
	"context"
	"time"

	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddResults(Results []models.Result) error {
	var interfaceSlice []interface{} = make([]interface{}, len(Results))
	for i, d := range Results {
		interfaceSlice[i] = d
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := dbclient.Collection("results").InsertMany(ctx, interfaceSlice)

	if err != nil {
		logger.Error(loggingArea, "Couldn't insert results for agent", Results[0].HostID, ":", err)
	}

	return err
}

//GetResults returns all results without limiting the resultset
//Please be aware that this can be slow if you are working with a large set of data
func GetResults(AgentID primitive.ObjectID, ItemID primitive.ObjectID) (models.ResultSet, error) {
	return GetResultsWithLimit(AgentID, ItemID, 0)
}

func GetResultsWithLimit(AgentID primitive.ObjectID, ItemID primitive.ObjectID, Limit int64) (models.ResultSet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "capturedat", Value: -1}})
	findOptions.Limit = &Limit

	results, err := dbclient.Collection("results").Find(ctx, bson.M{"$and": []bson.M{
		{"itemid": ItemID},
		{"hostid": AgentID},
	}}, findOptions)

	if err != nil {
		logger.Error(loggingArea, "Couldn't construct result set:", err)
		return models.ResultSet{}, err
	}

	var resultSet models.ResultSet

	if err := results.All(ctx, &resultSet.Results); err != nil {
		logger.Error(loggingArea, "Couldn't decode resultset contents:", err)
		return models.ResultSet{}, err
	}

	logger.Debug(loggingArea, "Fetched ResultSet with", len(resultSet.Results), "results")

	return resultSet, nil
}
