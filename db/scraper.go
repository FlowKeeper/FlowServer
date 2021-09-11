package db

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ensureCurrentScraper() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	hostname, _ := os.Hostname()

	result := dbclient.Collection("scrapers").FindOne(ctx, bson.M{"hostname": hostname})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			createScraper()
			return
		}

		logger.Fatal(loggingArea, "Couldn't look up scraper:", result.Err())
	}

	if err := result.Decode(&InstanceConfig); err != nil {
		logger.Fatal(loggingArea, "Couldn't decode scraper config:", err)
	}
}

func createScraper() {
	logger.Debug(loggingArea, "Scraper seems to be new, registering...")

	hostname, _ := os.Hostname()

	InstanceConfig = struct {
		Hostname   string
		InstanceID uuid.UUID
	}{
		Hostname:   hostname,
		InstanceID: uuid.New(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	dbclient.Collection("scrapers").InsertOne(ctx, InstanceConfig)
}
