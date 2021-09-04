package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetItems(IDs []primitive.ObjectID) ([]models.Item, error) {
	items := make([]models.Item, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := dbclient.Collection("items").Find(ctx, bson.M{"_id": bson.M{"$in": IDs}})

	if err != nil {
		logger.Error(loggingArea, "Couldn't read items:", err)
		return items, err
	}

	for result.Next(ctx) {
		var singleItem models.Item
		if err := result.Decode(&singleItem); err != nil {
			logger.Error(loggingArea, "Couldn't decode item:", err)
			return items, err
		}

		items = append(items, singleItem)
	}

	return items, nil
}
