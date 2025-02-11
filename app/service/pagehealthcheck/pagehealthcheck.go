// Copyright (c) 2023 FajarLaksono. All Rights Reserved.

package pagehealthcheck

// import (
// 	"context"

// 	"fajarlaksono.github.io/laksono-api-service/app/repository/postgres"
// 	"github.com/sirupsen/logrus"
// )

// func NewPageHealthCheck(
// 	postgresClient *postgres.PostgreClient,
// 	// redisClient *redis.Storage,
// 	// mongoClient *mongo.Client,
// ) *PageHealthCheck {
// 	return &PageHealthCheck{
// 		Postgre: *postgresClient,
// 		// Redis:   *redisClient,
// 		// Mongo:   *mongoClient,
// 	}
// }

// type PageHealthCheck struct {
// 	Postgre postgres.PostgreClient
// 	// Redis   redis.Storage
// 	// Mongo   mongo.Client
// }

// func (s *PageHealthCheck) GetComponentsStatus(ctx context.Context) (bool, map[string]bool) {
// 	isHealthy := true
// 	components := make(map[string]bool)

// 	components["Postgres"] = true
// 	if !s.Postgre.Health() {
// 		logrus.Error(log.NewServiceUnhealthy("Postgres"))
// 		isHealthy = false
// 		components["Postgres"] = false
// 	}
// 	// components["Redis"] = true
// 	// if !s.Redis.HealthCheck() {
// 	// 	logrus.Error(log.NewServiceUnhealthy("Redis"))
// 	// 	isHealthy = false
// 	// 	components["Redis"] = false
// 	// }
// 	// components["Mongo"] = true
// 	// if !s.Mongo.HealthCheck() {
// 	// 	logrus.Error(log.NewServiceUnhealthy("Mongo"))
// 	// 	isHealthy = false
// 	// 	components["Mongo"] = false
// 	// }

// 	return isHealthy, components
// }
