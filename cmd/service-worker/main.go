// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package main

import (
	"context"
	"flag"
	"os"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	model "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	kafkaconsumerService "fajarlaksono.github.io/laksono-api-service/app/repository/kafkaconsumer"
	kafkaconsumer "fajarlaksono.github.io/laksono-api-service/app/repository/kafkaconsumer/kafka"
	"fajarlaksono.github.io/laksono-api-service/app/service/worker"
	"fajarlaksono.github.io/laksono-api-service/cmd/utils"
	"github.com/caarlos0/env"
	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var revisionID = "unknown"
var buildDate = "unknown"
var gitHash = "unknown"

const (
	serviceNameLetters = "Laksono API Service Worker"
	serviceNameCode    = "laksono-api-service-worker"

	kafkaTimeout     = 10 * time.Second
	kafkaMinSize     = 10e3 // 10KB
	kafkaMaxSize     = 10e6 // 10MB
	kafkaMaxWaitTime = 500 * time.Millisecond
)

var disconnectFunc func()

func main() {
	// Initial informations
	utils.PrintSplashInformation(serviceNameLetters, revisionID, buildDate, gitHash)

	// Load configuration
	conf := config.Config{}

	flag.Usage = config.FlagUsage(conf)
	flag.Parse()

	if err := env.Parse(&conf); err != nil {
		log.WithError(err).Error("unable to parse environment variables")
		return
	}

	// Set Log Level
	logLvl, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.WithError(err).Error("unable to parse log level from config")
		return
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(logLvl)

	// Initialize dependencies
	// // Redis connection configuration
	// redisClient, _, locker, err := utils.InitRedis(conf)
	// if err != nil {
	// 	log.WithError(err).Error("unable to initialize redis")
	// 	return
	// }

	// Postgres connection configuration
	postgresClient, err := utils.InitPostgres(conf, "worker")
	if err != nil {
		log.WithError(err).Error("unable to initialize postgres")
		return
	}

	// MongoDB connection configuration
	// mongoClient, disconnectFunc, err := utils.InitMongo(conf, locker)
	// if err != nil {
	// 	log.WithError(err).Error("unable to initialize mongo storage")
	// 	return
	// }

	if err := utils.CreateKafkaTopic(conf.KafkaBrokerList, conf.KafkaTopicProjects); err != nil {
		log.Errorf("unable to create Kafka topic (%s): %s", conf.KafkaTopicProjects, err)
		return
	}

	// Retry mechanism with exponential backoff for Kafka connection
	operation := func() error {
		log.Info("Attempting to connect to Kafka")
		_, err := kafka.DialLeader(context.Background(), "tcp", conf.KafkaBrokerList[0], conf.KafkaTopicProjects, 0)
		if err != nil {
			return err
		}
		return nil
	}

	notify := func(err error, duration time.Duration) {
		log.Warnf("Retrying in %v seconds due to error: %s", duration.Seconds(), err.Error())
	}

	err = backoff.RetryNotify(operation, backoff.NewExponentialBackOff(), notify)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka after retries: %v", err)
	}

	recordChan := make(chan *model.Record)
	recordFinishChan := make(chan *kafka.Message)
	recordRetryFinishChan := make(chan *kafka.Message)

	retryWriter := &kafka.Writer{
		Addr:  kafka.TCP(conf.KafkaBrokerList...),
		Topic: conf.KafkaTopicProjects,
	}

	kafkaDial := &kafka.Dialer{
		Timeout: kafkaTimeout,
	}

	if conf.KafkaEnableTLSConn {
		tlsConf, err := config.NewTLS(conf.TLSCertificatePath)
		if err != nil {
			logrus.Errorf("unable to generate TLS configuration from a given certificate (%s): %s",
				conf.TLSCertificatePath, err)
			return
		}

		retryWriter.Transport = &kafka.Transport{TLS: tlsConf}
		kafkaDial.TLS = tlsConf
	}

	var workers []worker.WorkerService
	for i := 0; i < conf.NumberOfWorker; i++ {
		workers = append(workers, worker.New(
			&conf,
			recordChan,
			recordFinishChan,
			make(chan bool),
			retryWriter,
			recordRetryFinishChan,
			postgresClient))
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               conf.KafkaBrokerList,
		GroupID:               conf.KafkaConsumerGroup,
		Topic:                 conf.KafkaTopicProjects,
		WatchPartitionChanges: true,
		MinBytes:              kafkaMinSize,
		MaxBytes:              kafkaMaxSize,
		MaxWait:               kafkaMaxWaitTime,
		Dialer:                kafkaDial,
	})
	defer func() {
		_ = reader.Close()
	}()

	retryReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               conf.KafkaBrokerList,
		GroupID:               conf.KafkaConsumerGroup,
		Topic:                 conf.KafkaTopicProjects,
		WatchPartitionChanges: true,
		MinBytes:              kafkaMinSize,
		MaxBytes:              kafkaMaxSize,
		MaxWait:               kafkaMaxWaitTime,
		Dialer:                kafkaDial,
	})
	defer func() {
		_ = retryReader.Close()
	}()

	ctx := context.Background()
	consumer := kafkaconsumer.NewConsumer(ctx, reader, recordChan, recordFinishChan, 0)
	retryConsumer := kafkaconsumer.NewConsumer(ctx, retryReader, recordChan, recordRetryFinishChan,
		conf.RetryWaitTimeMinute)

	service := kafkaconsumerService.KafkaConsumerService{
		Consumer:          consumer,
		NumberOfProcessor: conf.NumberOfWorker,
		RecordFinishChan:  recordFinishChan,
		Realm:             conf.Realm,
		Processors:        workers,
		RetryConsumer:     retryConsumer,
	}

	err = service.Start()
	logrus.Fatalf("service crashed. cause: %v", err)
}
