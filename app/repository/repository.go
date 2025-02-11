// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package repository

import (
	"context"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	postgresMigration "fajarlaksono.github.io/laksono-api-service/app/migration/postgres"
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
