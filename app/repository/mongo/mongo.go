// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package mongo

/* import (
	"context"

	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	mongoDB       *mongo.Database
	mongoDBClient *mongo.Client
	locker        repository.Locker
}

func New(mongoDBURL, dbName string, locker repository.Locker) (*Client, func(), error) {
	mongoClient, disconnectMongoFunc, err := InitializeMongoClient(mongoDBURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "initialize mongo client")
	}

	return &Client{
		mongoDB:       mongoClient.Database(dbName),
		mongoDBClient: mongoClient,
		locker:        locker,
	}, disconnectMongoFunc, nil
}

func InitializeMongoClient(serverURI string) (*mongo.Client, func(), error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(serverURI)

	ctx := context.Background()

	mongoClient, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, nil, errors.Wrap(err, "initialize mongo client")
	}

	if err = mongoClient.Connect(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "mongo client connect")
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		return nil, nil, errors.Wrap(err, "mongo client ping")
	}

	disconnectFunc := func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			logrus.WithError(err).Error("unable to close mongodb client")

			return
		}

		logrus.Info("successfully close mongodb client connection")
	}

	return mongoClient, disconnectFunc, nil
}

func (c *Client) GetMongoClient() *mongo.Client {
	return c.mongoDBClient
}
*/
