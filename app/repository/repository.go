// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package repository

import (
	"context"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	postgresMigration "fajarlaksono.github.io/laksono-api-service/app/migration/postgres"
	modelapirequest "fajarlaksono.github.io/laksono-api-service/app/model/api/request"
	model "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"fajarlaksono.github.io/laksono-api-service/app/repository/postgres"
	"github.com/pkg/errors"
)

// type Locker interface {
// 	Lock(ctx context.Context, processID string, lockDuration time.Duration) (*redislock.Lock, error)
// 	Unlock(ctx context.Context, lock *redislock.Lock) error
// }

type PostgresDAO interface {
	Health(ctx context.Context) error

	CreateUser(ctx context.Context, data *model.User) (int64, error)
	GetUsers(ctx context.Context) (*model.Users, error)

	CreateProjects(ctx context.Context, data *model.Projects) (int64, error)
	GetProjects(ctx context.Context, projectNameFilter *string, isOverlappingFilter *bool,
		startDateFilter, endDateFilter *time.Time) (int64, *model.Projects, error)
	GetProjectByID(ctx context.Context, projectID string) (*model.Project, error)
	PatchProjects(ctx context.Context, input modelapirequest.UpdateProjectsRequest) (int64, error)
	DeleteProjects(ctx context.Context, input modelapirequest.DeleteProjectsByIDs) (int64, error)
	EvaluateNonOverlapProjects(ctx context.Context) (int64, error)
	EvaluateOverlapProjects(ctx context.Context) (int64, error)
}

func InitPostgres(configuration *config.Config) (PostgresDAO, error) {
	if configuration.PostgresIsInitMigrate {
		migConfig := postgresMigration.DBConfigToMigrationConfig(*configuration)
		pgMigration, err := postgresMigration.NewMigration(migConfig, configuration.PostgresMigrationPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to init database migrations")
		}
		err = pgMigration.Up()
		if err != nil {
			return nil, errors.Wrap(err, "unable to migrate database")
		}
	}

	dbConn, err := postgres.New("postgres", configuration)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database connection")
	}

	return dbConn, nil
}
