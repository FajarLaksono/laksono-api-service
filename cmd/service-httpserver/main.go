// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"fajarlaksono.github.io/laksono-api-service/app/service/api"
	"fajarlaksono.github.io/laksono-api-service/cmd/utils"
	"github.com/caarlos0/env"
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
	postgresClient, err := utils.InitPostgres(*conf)
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

	// Start the server
	service, err := api.NewService(conf, postgresClient)
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
	}()

	wg.Wait()
}
