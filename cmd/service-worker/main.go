// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package main

import (
	"flag"
	"os"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"fajarlaksono.github.io/laksono-api-service/cmd/utils"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

var revisionID = "unknown"
var buildDate = "unknown"
var gitHash = "unknown"

const serviceNameLetters = "Laksono API Service Worker"
const serviceNameCode = "laksono-api-service-worker"

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

}
