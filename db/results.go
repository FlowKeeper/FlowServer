package db

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
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
