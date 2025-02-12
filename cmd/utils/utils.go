// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"github.com/cenkalti/backoff/v3"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type DepsService struct {
	// Redis    *redis.Storage
	Postgres *repository.PostgresDAO
	// Mongo    *mongo.Client
}

func PrintSplashInformation(serviceName, revisionID, buildDate, gitHash string) {
	log.Info(serviceName)
	log.Infof("RevisionID: %s, Build Date: %s, Git Hash: %s", revisionID, buildDate, gitHash)
}

// func InitRedis(cfg config.Config) (*redis.Storage, *goredis.Client, repository.Locker, error) {
// 	redisClient, err := redis.NewClient(
// 		redis.Config{
// 			Host:     cfg.RedisHost,
// 			Port:     cfg.RedisPort,
// 			Username: cfg.RedisUsername,
// 			Password: cfg.RedisPassword,
// 		},
// 		redis.WithConnectionPoolSize(cfg.RedisConnectionPoolSize),
// 		redis.WithMinIdleConnection(cfg.RedisMinIdleConnection))
// 	if err != nil {
// 		return nil, nil, nil, errors.Wrap(err, "unable to initialize redis")
// 	}

// 	locker := redis.NewLocker(redisClient)

// 	return redis.NewStorage(redisClient), redisClient, locker, nil
// }

func InitPostgres(cfg config.Config, serviceType string) (repository.PostgresDAO, error) {
	conf := cfg
	if serviceType == "worker" {
		conf.PostgresIsInitMigrate = false
	}

	postgresStorage, err := repository.InitPostgres(&conf)

	return postgresStorage, errors.Wrap(err, "unable to create new postgres connection")
}

// func InitMongo(cfg config.Config, locker repository.Locker) (*mongo.Client, func(), error) {
// 	mongoStorage, disconnectMongoFunc, err := mongo.New(cfg.MongoDBURL, cfg.MongoDBName, locker)

// 	return mongoStorage, disconnectMongoFunc, errors.Wrap(err, "unable to create mongodb client")
// }

func InitHTTPServer(port, basePath string) (runFunc, stopFunc func(), err error) {
	webService := new(restful.WebService)
	tags := []string{"Health"}

	webService.Path(basePath + "/healthz")
	webService.Route(webService.GET("/").
		To(func(req *restful.Request, resp *restful.Response) {
			resp.WriteHeader(http.StatusOK)
		}).
		Doc("Get health status").
		Notes(`Get healthiness status of this service.
This endpoint is also used by k8s to check whether the service is ready or not.`).
		Produces(restful.MIME_JSON).
		Metadata(restfulspec.KeyOpenAPITags, tags))

	goRestfulContainer := restful.NewContainer()
	goRestfulContainer.Add(webService)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to initialize http listener")
	}

	runFunc = func() {
		if err := http.Serve(listener, goRestfulContainer); err != nil {
			log.WithError(err).Error("unable to serve http server")
		}
	}

	stopFunc = func() {
		if err := listener.Close(); err != nil {
			log.WithError(err).Error("unable to close http listener")
		}
	}

	return runFunc, stopFunc, nil
}

func RunWorkers(ctx context.Context, cfg config.Config, deps *DepsService) {
	// runner := worker.Runner{}

	// runner.Run(ctx)
}

func CreateKafkaTopic(brokerAddress []string, topic string) error {
	var conn *kafka.Conn
	var err error

	operation := func() error {
		log.Infof("Connecting to Kafka broker at %s", brokerAddress[0])
		conn, err = kafka.Dial("tcp", brokerAddress[0])
		if err != nil {
			return errors.Wrap(err, "failed to dial Kafka")
		}
		defer conn.Close()

		// Check if the topic exists
		topics, err := conn.ReadPartitions()
		if err != nil {
			return errors.Wrap(err, "failed to read partitions")
		}

		topicExists := false
		for _, p := range topics {
			if p.Topic == topic {
				topicExists = true
				break
			}
		}

		if topicExists {
			fmt.Printf("Topic %s already exists\n", topic)
		} else {
			// Define the topic configuration
			topicConfig := kafka.TopicConfig{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			}

			// Create the topic
			err = conn.CreateTopics(topicConfig)
			if err != nil {
				return errors.Wrap(err, "failed to create topic")
			}

			fmt.Printf("Topic %s created successfully\n", topic)
		}

		return nil
	}

	notify := func(err error, duration time.Duration) {
		log.Warnf("Retrying in %v seconds due to error: %s", duration.Seconds(), err.Error())
	}

	err = backoff.RetryNotify(operation, backoff.NewExponentialBackOff(), notify)
	if err != nil {
		return errors.Wrap(err, "failed to create Kafka topic after retries")
	}

	return nil
}
