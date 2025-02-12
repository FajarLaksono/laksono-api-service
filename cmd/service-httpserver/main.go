// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"fajarlaksono.github.io/laksono-api-service/app/service/api"
	"fajarlaksono.github.io/laksono-api-service/cmd/utils"
	"github.com/caarlos0/env"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

var revisionID = "unknown"
var buildDate = "unknown"
var gitHash = "unknown"

const serviceNameLetters = "Laksono API Service"
const serviceNameCode = "laksono-api-service"

func main() {
	// Initial information
	utils.PrintSplashInformation(serviceNameLetters, revisionID, buildDate, gitHash)

	// Load configuration
	conf := &config.Config{}

	flag.Usage = config.FlagUsage(conf)
	flag.Parse()

	if err := env.Parse(conf); err != nil {
		log.WithError(err).Fatal("Failed to load environment variables")

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
	// redisClient, _, locker, err := utils.InitRedis(*conf)
	// if err != nil {
	// 	log.WithError(err).Error("unable to initialize redis")

	// 	return
	// }

	// Postgres connection configuration
	postgresClient, err := utils.InitPostgres(*conf, "api")
	if err != nil {
		log.WithError(err).Error("unable to initialize postgres")

		return
	}

	// // MongoDB connection configuration
	// mongoClient, disconnectFunc, err := utils.InitMongo(conf, locker)
	// if err != nil {
	// 	log.WithError(err).Error("unable to initialize mongo storage")

	// 	return
	// }

	// defer disconnectFunc()

	// Initial delay to give Kafka broker more time to initialize
	log.Info("Sleep for 60s to wait kafka ready (Homework to make to faster)")
	time.Sleep(60 * time.Second)
	log.Info("Running again")

	if err := utils.CreateKafkaTopic(conf.KafkaBrokerList, conf.KafkaTopicProjects); err != nil {
		log.Errorf("unable to create Kafka topic (%s): %s", conf.KafkaTopicProjects, err)
		return
	}

	// Kafka Writer Configuration
	kafkaWriter := &kafka.Writer{
		Addr:       kafka.TCP(conf.KafkaBrokerList...),
		Topic:      conf.KafkaTopicProjects,
		BatchSize:  100,
		BatchBytes: 200000000,
	}

	if conf.KafkaEnableTLSConn {
		tlsConf, err := config.NewTLS(conf.TLSCertificatePath)
		if err != nil {
			log.Errorf("unable to generate TLS configuration from a given certificate (%s): %s",
				conf.TLSCertificatePath, err)

			return
		}

		kafkaWriter.Transport = &kafka.Transport{TLS: tlsConf}
	}

	utils.CreateKafkaTopic(conf.KafkaBrokerList, conf.KafkaTopicProjects)

	// Start the server
	service, err := api.NewService(conf, postgresClient, kafkaWriter)
	if err != nil {
		log.WithError(err).Fatal("Failed to create HTTP server preparation")

		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Info("Quit signal received, Shutting down the Http server...")
		service.Stop()
	}()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info("Starting HTTP Server")

		if err := service.Start(); err != nil {
			log.WithError(err).Fatal("Failed to start HTTP server")
		}

		log.Info("The webservice is running now !!!")
	}()

	wg.Wait()
}
