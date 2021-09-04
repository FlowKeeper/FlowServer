package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowserver/v2/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var rawclient *mongo.Client
var dbclient *mongo.Database

const loggingArea = "DB"

type InstanceConfigSample struct {
	InstanceID uuid.UUID
}

var InstanceConfig InstanceConfigSample

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Config.MongoDB))

	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	rawclient = client
	dbclient = client.Database("flowkeeper")
	logger.Info(loggingArea, "Connected to MongoDB")

	initialize()
	return nil
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rawclient.Disconnect(ctx)
}

func initialize() {
	logger.Info(loggingArea, "Starting initializtaion")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := dbclient.Collection("config").FindOne(ctx, bson.M{})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			firstTimeSetup()
			return
		} else {
			logger.Fatal(loggingArea, "Couldn't read config from database:", result.Err())
		}
	}

	if err := result.Decode(&InstanceConfig); err != nil {
		logger.Fatal(loggingArea, "Couldn't parse config from databse:", err)
	}

	logger.Info(loggingArea, "Instance Config loaded")
}

func firstTimeSetup() {
	logger.Info(loggingArea, "Seems to be first time setup is needed!")
	InstanceConfig = InstanceConfigSample{
		InstanceID: uuid.New(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := dbclient.Collection("config").InsertOne(ctx, InstanceConfig)
	if err != nil {
		logger.Fatal(loggingArea, "Couldn't initialize config:", err)
	}

	logger.Info(loggingArea, "Created config with id", result.InsertedID)
}
