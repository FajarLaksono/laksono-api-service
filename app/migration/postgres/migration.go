// Copyright (c) 2019 Fajar Laksono. All Rights Reserved.

package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/aws_s3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	//	_ "github.com/golang-migrate/migrate/v4/source/google_cloud_storage"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host       string
	Username   string
	Password   string
	Name       string
	Schema     string
	Port       int
	SSLEnabled bool
	LogMode    bool
}

// Migration define migrations config
type Migration struct {
	DBDriver *database.Driver
	Source   string
	DBName   string
}

const (
	driverName          = "postgres"
	dbTimeoutConnection = 30 * time.Second
	dbConnectTimeOut    = 5
)

// NewMigration creates new migrations instance
func NewMigration(conf *Config, source string) (*Migration, error) {
	var dbDriver database.Driver
	var err error

	pgDB, err := sql.Open(driverName, connectionString(conf))
	if err != nil {
		return nil, err
	}

	// Close the connection for DB schema checker. This connection is used one time only for initializing DB.
	defer func() {
		err = pgDB.Close()
		if err != nil {
			logrus.WithError(err).Errorf("unable to close migration db connection [%v]", conf.Name)
		}
	}()

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = dbTimeoutConnection
	err = backoff.RetryNotify(func() error {
		dbDriver, err = postgres.WithInstance(pgDB, &postgres.Config{})
		if err != nil {
			return err
		}
		return nil
	}, b, func(err error, duration time.Duration) {
		logrus.Warnf("connection retries duration: %f sec , err: %s",
			duration.Seconds(), err.Error())
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: timeout: %w", err)
	}
	if pgDB == nil {
		return nil, fmt.Errorf("unable to initiate migrations: %w", err)
	}

	return &Migration{
		DBDriver: &dbDriver,
		Source:   source,
		DBName:   conf.Name,
	}, nil
}

// Migrate creates database schema
func (migration *Migration) Up() error {
	logrus.Infoln("Start schema migrations...")
	mig, err := migrate.NewWithDatabaseInstance(
		migration.Source,
		migration.DBName,
		*migration.DBDriver,
	)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = mig.Close()
	}()

	_, isDirty, err := mig.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		logrus.Infoln("Databases schema migrations... FAILED")
		return err
	}

	if isDirty {
		logrus.Infoln("Databases schema migrations... FAILED")
		logrus.Infoln("In order to protect your data integrity please fix you database schema using this tools.")
		logrus.Infoln("https://github.com/golang-migrate/migrate/releases")
	}
	err = mig.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logrus.Infoln("No changes, database schema is up to date")
		return nil
	}
	if err != nil {
		logrus.Infoln("Databases schema migrations... FAILED")
		return err
	}

	logrus.Infoln("Databases schema migrations... DONE")
	return nil
}

func (migration *Migration) Down() error {
	mig, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migration.Source),
		migration.DBName,
		*migration.DBDriver,
	)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = mig.Close()
	}()

	err = mig.Down()
	if errors.Is(err, migrate.ErrNoChange) {
		logrus.Println("Database schema is up to date")
		return nil
	}

	return err
}

func connectionString(config *Config) string {
	sslMode := "" // omit sslmode so it set to default (enabled)
	if !config.SSLEnabled {
		sslMode = "sslmode=disable"
	}

	searchPath := fmt.Sprintf("search_path=%s,public", config.Schema)

	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s connect_timeout=%v %s %s",
		config.Host, config.Port, config.Name, config.Username, config.Password, dbConnectTimeOut, searchPath, sslMode)
}

// DBConfigToMigrationConfig converts postgres config to migrations config
func DBConfigToMigrationConfig(conf config.Config) *Config {
	return &Config{
		Host:       conf.PostgresHost,
		Username:   conf.PostgresUsername,
		Password:   conf.PostgresPassword,
		Name:       conf.PostgresDBName,
		Schema:     conf.PostgresSchema,
		Port:       conf.PostgresPort,
		SSLEnabled: conf.PostgresSSLEnabled,
		LogMode:    conf.PostgresDebugMode,
	}
}
