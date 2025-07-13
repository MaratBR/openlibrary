package session

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	_REDIS_PREFIX = "session:"
)

type redisStore struct {
	client *redis.Client
}

type redisSession struct {
	id      redactedString
	client  *redis.Client
	data    map[string]string
	touched bool
}

// Get implements Session.
func (s *redisSession) Get(ctx context.Context, key string) (string, bool) {
	value, ok := s.data[key]
	return value, ok
}

// Put implements Session.
func (s *redisSession) Put(ctx context.Context, key string, value string) {
	s.data[key] = value
	s.touched = true
}

// Save implements Session.
func (s *redisSession) Save(ctx context.Context) error {
	if !s.touched {
		return nil
	}

	if s.data == nil {
		return nil
	}

	_, err := s.client.HSet(ctx, _REDIS_PREFIX+string(s.id), s.data).Result()
	return err
}

func (s *redisSession) load(ctx context.Context) error {
	// TODO: encryption?
	data, err := s.client.HGetAll(ctx, _REDIS_PREFIX+string(s.id)).Result()
	if err != nil {
		return err
	}
	s.data = data
	s.touched = false
	return nil
}

// Get implements Store.
func (r *redisStore) Get(ctx context.Context, id string) (Session, error) {
	s := &redisSession{id: redactedString(id), client: r.client}
	err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewRedisStore(client *redis.Client) Store {
	return &redisStore{
		client: client,
	}
}
