// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	postgresClient, err := utils.InitPostgres(conf)
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

	deps := &utils.DepsService{
		// Redis:    redisClient,
		Postgres: &postgresClient,
		// Mongo:    mongoClient,
	}

	defer disconnectFunc()

	// Start workers
	workerCtx, workerCtxCancel := context.WithCancel(context.Background())
	defer func() {
		log.Info("cancel all workers...")
		workerCtxCancel()
		log.Info("all workers has been shut down")
	}()

	go utils.RunWorkers(workerCtx, conf, deps)

	runHTTPServerFunc, stopHTTPServerFunc, err := utils.InitHTTPServer(conf.Port, conf.BasePath)
	if err != nil {
		log.WithError(err).Error("unable to initialize http server")

		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup, stopHTTPServerFunc func()) {
		defer wg.Done()

		<-sig
		log.Info("os quit signal received. stop all running process.")

		stopHTTPServerFunc()
	}(&wg, stopHTTPServerFunc)

	wg.Add(1)
	go func(wg *sync.WaitGroup, runHTTPServerFunc func()) {
		defer wg.Done()

		log.Info("running worker server")
		runHTTPServerFunc()
	}(&wg, runHTTPServerFunc)

	wg.Wait()
}
