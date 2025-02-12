// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package handler

import (
	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"github.com/emicklei/go-restful/v3"
	"github.com/segmentio/kafka-go"
)

type APIService struct {
	WebService  *restful.WebService
	DAOPostgres repository.PostgresDAO
	KafkaWriter *kafka.Writer
	// RedisDAO    redis.DAO
	// MongoDAO    mongo.DAO
}

func New(
	postgreDAO repository.PostgresDAO,
	kafkaWriter *kafka.Writer,
	// redisDAO redis.DAO,
	// mongoDAO mongo.DAO,
) *APIService {
	service := &APIService{
		DAOPostgres: postgreDAO,
		KafkaWriter: kafkaWriter,
		// RedisDAO:   redisDAO,
		// MongoDAO:   mongoDAO,
	}

	return service
}
