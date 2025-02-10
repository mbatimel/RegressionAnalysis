package locker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=locker.go -destination locker_mock.go -package locker
type Locker interface {
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (Lock, error)
}
type locker struct {
	rs *redsync.Redsync
}

func (l *locker) AcquireLock(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	mutex := l.rs.NewMutex(key, redsync.WithExpiry(ttl))
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("could not lock mutex: %w", err)
	}
	return mutex, nil
}

type Lock interface {
	UnlockContext(context.Context) (bool, error)
}

func NewLocker(redisCli goredislib.UniversalClient) (Locker, error) {
	pool := goredis.NewPool(redisCli)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)
	return &locker{
		rs: rs,
	}, nil
}
