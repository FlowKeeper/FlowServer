package db

import (
	"context"
	"time"

	"github.com/FlowKeeper/FlowServer/v2/config"
	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var rawclient *mongo.Client
var dbclient *mongo.Database

const loggingArea = "DB"

//InstanceConfig gets loaded from the database upon calling Connect()
var InstanceConfig struct {
	Hostname   string
	InstanceID uuid.UUID
}

//Connect connects to the mongodb and initializes the client
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

	ensureIndexes()
	logger.Info(loggingArea, "Starting initializtaion")
	ensureCurrentScraper()
	logger.Info(loggingArea, "DB subsystem is ready")
	return nil
}

//Disconnect closes the connection to mongodb
func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rawclient.Disconnect(ctx)
}

//Client returns the current mongodb client with a preselected database
func Client() *mongo.Database {
	return dbclient
}
