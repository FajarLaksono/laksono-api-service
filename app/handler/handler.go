// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package handler

import (
	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"github.com/emicklei/go-restful/v3"
)

type APIService struct {
	WebService  *restful.WebService
	DAOPostgres repository.PostgresDAO
	// RedisDAO    redis.DAO
	// MongoDAO    mongo.DAO
}

func New(
	postgreDAO repository.PostgresDAO,
	// redisDAO redis.DAO,
	// mongoDAO mongo.DAO,
) *APIService {
	service := &APIService{
		DAOPostgres: postgreDAO,
		// RedisDAO:   redisDAO,
		// MongoDAO:   mongoDAO,
	}

	return service
}
