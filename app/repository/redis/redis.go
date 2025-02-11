// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package redis

/* import (
	"context"
	"fmt"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"github.com/bsm/redislock"
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
)

const (
	maxRetryOperation = 5
	DefaultOpsTimeout = 10 * time.Second
)

// Client functions
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewClient(cfg Config, options ...func(*goredis.Options)) (*goredis.Client, error) {
	redisOptions := &goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
	}

	for _, o := range options {
		o(redisOptions)
	}

	redisClient := goredis.NewClient(redisOptions)

	op := func() error {
		ctx, cancelFunc := context.WithTimeout(context.Background(), DefaultOpsTimeout)
		defer cancelFunc()

		err := redisClient.Ping(ctx).Err()

		return errors.Wrap(err, "redis ping")
	}

	err := backoff.Retry(op, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetryOperation))

	return redisClient, errors.Wrap(err, "initialize redis client")
}

func WithConnectionPoolSize(poolSize int) func(*goredis.Options) {
	return func(c *goredis.Options) {
		c.PoolSize = poolSize
	}
}

func WithMinIdleConnection(minIdleCons int) func(*goredis.Options) {
	return func(c *goredis.Options) {
		c.MinIdleConns = minIdleCons
	}
}

// Storage functions
type Storage struct {
	rdsClient *goredis.Client
}

func NewStorage(rdsClient *goredis.Client) *Storage {
	return &Storage{
		rdsClient: rdsClient,
	}
}

// Locker Functions
const (
	defaultLockRetryCount = 10
)

func NewLocker(client goredis.UniversalClient) repository.Locker {
	return &redisLocker{
		redisLockClient: redislock.New(client),
	}
}

type redisLocker struct {
	redisLockClient *redislock.Client
}

func (locker *redisLocker) Lock(ctx context.Context, processID string, lockDuration time.Duration) (*redislock.Lock, error) {
	lock, err := locker.redisLockClient.Obtain(ctx, processID, lockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), defaultLockRetryCount),
	})
	if err != nil {
		return nil, err
	}
	return lock, nil
}

func (locker *redisLocker) Unlock(ctx context.Context, lock *redislock.Lock) error {
	return lock.Release(ctx)
}
*/
