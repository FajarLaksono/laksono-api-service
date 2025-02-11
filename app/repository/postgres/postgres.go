// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package postgres

import (
	"context"
	"fmt"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	postgresForeignKeyViolation = "23503"
	postgresUniqueKeyViolation  = "23505"
	maxRetryOperation           = 5
	dbTimeoutConnection         = 30 * time.Second
)

type PostgreClient struct {
	Client *gorm.DB
}

func New(dbDriver string, conf *config.Config) (*PostgreClient, error) {
	logrus.Info("Waiting client postgres...")

	sslMode := "" // omit sslmode so it set to default (enabled)
	if !conf.PostgresSSLEnabled {
		sslMode = "sslmode=disable"
	}

	connectionConfig := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s connect_timeout=5 %s",
		conf.PostgresHost, conf.PostgresPort, conf.PostgresDBName, conf.PostgresUsername, conf.PostgresPassword, sslMode)

	Logger := logger.Default.LogMode(logger.Warn)
	if conf.PostgresDebugMode {
		Logger = logger.Default.LogMode(logger.Info)
	}

	var (
		postgresORM *gorm.DB
		err         error
	)

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = dbTimeoutConnection

	err = backoff.RetryNotify(func() error {
		postgresORM, err = gorm.Open(postgres.Open(connectionConfig), &gorm.Config{Logger: Logger})
		if err != nil {
			return errors.Wrap(err, "unable to open gorm")
		}

		return nil
	}, b, func(err error, duration time.Duration) {
		logrus.Warnf("connection retries duration: %f sec , err: %s",
			duration.Seconds(), err.Error())
	})

	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres database: timeout")
	}

	if postgresORM == nil {
		return nil, errors.Wrapf(err, "unable to initiate postgres DAO")
	}

	sqlDB, err := postgresORM.DB()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get postgres sql instance")
	}

	sqlDB.SetMaxOpenConns(conf.PostgresMaxOpenConnections)
	sqlDB.SetMaxIdleConns(conf.PostgresMaxIdleConnections)
	sqlDB.SetConnMaxLifetime(conf.PostgresConnectionMaxLifeTime)

	return &PostgreClient{
		Client: postgresORM,
	}, nil
}

func (db *PostgreClient) Health(ctx context.Context) error {
	sqlDB, err := db.Client.DB()
	if err != nil {
		return errors.Wrapf(err, "unable to init postgres client")
	}
	return sqlDB.PingContext(ctx)
}
