// Copyright (c) 2023 FajarLaksono. All Rights Reserved.

package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"
)

type Config struct {
	// Service configurations
	Port            string   `env:"SERVER_PORT" envDefault:"8080" envDocs:"Port that the service will listen to"`
	BasePath        string   `env:"SERVER_BASE_PATH" envDefault:"/laksono" envDocs:"Base path of the service"`
	Realm           string   `env:"REALM" envDefault:"dev" envDocs:"Realm of the deployment"`
	ServiceRootPath string   `env:"SERVICE_ROOT_PATH" envDefault:"/srv" envDocs:"The root path of the service"`
	LogLevel        string   `env:"LOG_LEVEL" envDefault:"info" envDocs:"Log level"`
	AllowedOrigins  []string `env:"CORS_ALLOWED_ORIGINS" envDocs:"List of domains to be permitted by CORS"`
	AllowedMethods  []string `env:"CORS_ALLOWED_METHODS" envDocs:"List of methods to be permitted by CORS" envDefault:"GET,POST,PUT,PATCH,DELETE"`
	AllowedHeaders  []string `env:"CORS_ALLOWED_HEADERS" envDocs:"List of headers to be permitted by CORS" envDefault:"Access-Control-Allow-Origin,Access-Control-Allow-Methods,Authorization,Content-Type,Accept"`

	// Postgres configurations
	PostgresSSLEnabled            bool          `env:"POSTGRES_SSL_ENABLED" envDefault:"true" envDocs:"Enable SSL communication"`
	PostgresPort                  int           `env:"POSTGRES_PORT" envDefault:"5432" envDocs:"Port of the postgres"`
	PostgresMaxOpenConnections    int           `env:"POSTGRES_MAX_OPEN_CONNECTIONS" envDefault:"30" envDocs:"Postgres max open connection"`
	PostgresMaxIdleConnections    int           `env:"POSTGRES_MAX_IDLE_CONNECTIONS" envDefault:"10" envDocs:"Postgres max idle connection"`
	PostgresHost                  string        `env:"POSTGRES_HOST" envDocs:"Postgres Host"`
	PostgresDBName                string        `env:"POSTGRES_DB_NAME,required" envDocs:"Postgres Database name"`
	PostgresUsername              string        `env:"POSTGRES_USERNAME" envDocs:"Postgres Username"`
	PostgresPassword              string        `env:"POSTGRES_PASSWORD" envDocs:"Postgres Password"`
	PostgresSchema                string        `env:"POSTGRES_SCHEMA" envDefault:"public" envDocs:"Postgres schema name"`
	PostgresConnectionMaxLifeTime time.Duration `env:"POSTGRES_CONNECTION_MAX_LIFETIME" envDefault:"1h"`
	PostgresMigrationPath         string        `env:"POSTGRES_MIGRATION_PATH"  envDefault:"file:///srv/migrations/postgres" envDocs:"Postgres Migration file path"`
	PostgresIsInitMigrate         bool          `env:"POSTGRES_IS_INIT_MIGRATION" envDocs:"Init Database Migration" envDefault:"false"`
	PostgresDebugMode             bool          `env:"POSTGRES_DEBUG_MODE" envDocs:"show detailed database query logs" envDefault:"false"`

	// // Mongo configurations
	// MongoDBURL           string `env:"MONGO_DB_URL,required" envDocs:"MongoDB URI"`
	// MongoDBName          string `env:"MONGO_DB_NAME,required" envDocs:"MongoDB database name"`
	// MongoDBMigrationPath string `env:"MONGO_DB_MIGRATION_PATH" envDefault:"/srv/migrations/mongodb" envDocs:"MongoDB migration file path"`

	// // Redis configurations
	// RedisHost               string `env:"REDIS_HOST,required" envDocs:"Redis host"`
	// RedisPort               string `env:"REDIS_PORT,required" envDocs:"Redis port"`
	// RedisUsername           string `env:"REDIS_USERNAME" envDocs:"Redis username"`
	// RedisPassword           string `env:"REDIS_PASSWORD" envDocs:"Redis password"`
	// RedisConnectionPoolSize int    `env:"REDIS_CONNECTION_POOL_SIZE" envDefault:"20" envDocs:"Redis connection pool size"`
	// RedisMinIdleConnection  int    `env:"REDIS_MIN_IDLE_CONNECTIONS" envDefault:"0" envDocs:"redis min idle connection"`

	KafkaBrokerList    []string `env:"KAFKA_BROKER_LIST,required" envDocs:"Kafka endpoint"`
	KafkaTopicProjects string   `env:"KAFKA_TOPIC_PROJECTS" envDocs:"Kafka Topic name for Projects" envDefault:"laksono-projects"`
	KafkaEnableTLSConn bool     `env:"KAFKA_ENABLE_TLS_CONN" envDocs:"Enable TLS connection to Kafka Brokers" envDefault:"false"`
	TLSCertificatePath string   `env:"TLS_CERTIFICATE_PATH" envDocs:"TLS certificate path in string. If empty, will fallback to default TLS crt file for alpine linux at /etc/ssl/certs/ca-certificates.crt." envDefault:"/etc/ssl/certs/ca-certificates.crt"`
}

func (envVar Config) HelpDocs() []string {
	reflectEnvVar := reflect.TypeOf(envVar)
	doc := make([]string, 1+reflectEnvVar.NumField())
	doc[0] = "Service Config" + ":"
	for i := 0; i < reflectEnvVar.NumField(); i++ {
		field := reflectEnvVar.Field(i)
		envName := field.Tag.Get("env")
		envDefault := field.Tag.Get("envDefault")
		envDocs := field.Tag.Get("envDocs")
		doc[i+1] = fmt.Sprintf("  %v\t %v (default: %v)", envName, envDocs, envDefault)
	}

	return doc
}

type ConfigurationHelper interface {
	HelpDocs() []string
}

func FlagUsage(inp ConfigurationHelper) func() {
	return func() {
		flag.CommandLine.SetOutput(os.Stdout)
		for _, val := range inp.HelpDocs() {
			fmt.Println(val)
		}

		fmt.Println("")
		flag.PrintDefaults()
	}
}
